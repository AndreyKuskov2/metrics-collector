package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

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

		var buf bytes.Buffer

		gz := gzip.NewWriter(&buf)
		defer gz.Close()

		_, err = gz.Write(jsonData)
		if err != nil {
			logger.Println("Error compressing metric data")
			continue
		}

		ro := grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type":     "application/json",
				// "Content-Encoding": "gzip",
			},
			// DisableCompression: false,
			JSON: jsonData,
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
