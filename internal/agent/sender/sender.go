package sender

import (
	"encoding/json"
	"fmt"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/levigross/grequests"
)

func SendMetrics(address string, metrics map[string]models.Metrics) error {
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
			fmt.Printf("Failed to send metric %s: %v\n", metricName, err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Failed to send metric %s: status code %d\n", metricName, resp.StatusCode)
		}
	}
	return nil
}

func SendMetricsJSON(address string, metrics map[string]models.Metrics) error {
	for metricName, metricData := range metrics {
		url := fmt.Sprintf("http://%s/update", address)

		jsonData, err := json.Marshal(metricData)
		ro := grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			JSON: jsonData,
		}
		resp, err := grequests.Post(url, &ro)
		if err != nil {
			fmt.Printf("Failed to send metric %s: %v\n", metricName, err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Failed to send metric %s: status code %d\n", metricName, resp.StatusCode)
		}
	}
	return nil
}
