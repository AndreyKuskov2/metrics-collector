package sender

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

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

func SendMetricsBatch(cfg *config.AgentConfig, metricsData map[string]models.Metrics, logger *logrus.Logger) error {
	url := fmt.Sprintf("http://%s/updates/", cfg.Address)

	var requestBody []models.Metrics
	for _, metric := range metricsData {
		requestBody = append(requestBody, metric)
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Infof("Failed to marshal metrics: %v", err)
		return err
	}

	ro := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		},
		DisableCompression: false,
		JSON:               jsonData,
	}
	if err := sendWithRetry(cfg, ro, url, logger); err != nil {
		logger.Infof("Failed to send metrics: %v", err)
	}
	return nil
}

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
