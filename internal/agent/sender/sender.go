package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

// Отправка метрик
func SendMetrics(address string, metrics map[string]models.Metrics, logger *logrus.Logger) error {
	for metricName, metricData := range metrics {
		var url string
		if metricData.Value == nil {
			url = fmt.Sprintf("http://%s/update/%s/%s/%v", address, metricData.MType, metricName, *metricData.Delta)
		} else {
			url = fmt.Sprintf("http://%s/update/%s/%s/%v", address, metricData.MType, metricName, *metricData.Value)
		}
		ro := grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
		resp, err := grequests.Post(url, &ro)
		if err != nil {
			logger.Printf("Failed to send metric %s: %v\n", metricName, err)
			continue
		}

		if resp.StatusCode != 200 {
			logger.Printf("Failed to send metric %s: status code %d\n", metricName, resp.StatusCode)
		}
	}
	return nil
}

// Отправка метрик в формате JSON
func SendMetricsJSON(address string, metrics map[string]models.Metrics, logger *logrus.Logger) error {
	for metricName, metricData := range metrics {
		url := fmt.Sprintf("http://%s/update/", address)

		jsonData, err := json.Marshal(metricData)
		if err != nil {
			logger.Printf("Failed to marshal metric %s: %v\n", metricName, err)
			continue
		}

		ro := grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type":     "application/json",
				"Content-Encoding": "gzip",
			},
			DisableCompression: false,
			JSON:               jsonData,
		}
		resp, err := grequests.Post(url, &ro)
		if err != nil {
			logger.Printf("Failed to send metric %s: %v\n", metricName, err)
			continue
		}

		if resp.StatusCode != 200 {
			logger.Printf("Failed to send metric %s: status code %d\n", metricName, resp.StatusCode)
		}
	}
	return nil
}

var gzipNewWriter = func(w io.Writer) *gzip.Writer {
	return gzip.NewWriter(w)
}

// Отправка метрик пачкой в формате JSON
func SendMetricsBatch(cfg *config.AgentConfig, metricsData models.AllMetrics, logger *logrus.Logger) error {
	url := fmt.Sprintf("http://%s/updates/", cfg.Address)

	var requestBody []models.Metrics
	for _, metric := range metricsData.RuntimeMetrics {
		requestBody = append(requestBody, metric)
	}

	for _, metric := range metricsData.AdditionalMetrics {
		requestBody = append(requestBody, metric)
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		logger.Println("Error json Encode metric data")
		return err
	}

	var buf bytes.Buffer

	gz := gzipNewWriter(&buf)

	_, err = gz.Write(body)
	if err != nil {
		logger.Println("Error compressing metric data")
		return err
	}

	err = gz.Close()
	if err != nil {
		logger.Println("Error close gzip compressor")
		return err
	}

	var response *http.Response
	for trying := 0; trying <= cfg.MaxRetries; trying++ {
		var hash string
		if cfg.SecretKey != "" {
			hash = calculateHash(body, []byte(cfg.SecretKey))
		}
		req, _ := http.NewRequest("POST", url, &buf)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("HashSHA256", hash)

		httpClient := &http.Client{
			Timeout: 100 * time.Millisecond,
		}

		response, err = httpClient.Do(req)
		if err != nil {
			if trying < cfg.MaxRetries && strings.Contains(err.Error(), "connection refused") {
				logger.Printf("Bad %v trying sending metric: %v. BODY: %v\n", trying+1, err, requestBody)
				time.Sleep(cfg.RetryDelay)
				continue
			}
			logger.Printf("Error sending metric: %v. BODY: %v\n", err, requestBody)
			return err
		}
		response.Body.Close()
		break
	}

	if response != nil {
		if response.StatusCode == http.StatusOK {
			logger.Printf("Successfully sent metric: %v\n", requestBody)
		} else {
			logger.Printf("Failed to send metric: %v, status code: %d\n", requestBody, response.StatusCode)
		}
	}
	return nil
}

// Повторная отправка метрик
func sendWithRetry(cfg *config.AgentConfig, ro grequests.RequestOptions, url string, logger *logrus.Logger) error {
	delay := cfg.RetryDelay
	for i := 0; i < cfg.MaxRetries; i++ {
		resp, err := grequests.Post(url, &ro)
		if err != nil {
			logger.Infof("Failed to send request: %v", err)
		} else if resp.StatusCode == 200 {
			return nil
		} else {
			logger.Infof("Failed to send request: status code %d", resp.StatusCode)
		}

		time.Sleep(delay)
		delay += 2 * time.Second
	}
	return fmt.Errorf("failed to send request after %d attempts", cfg.MaxRetries)
}
