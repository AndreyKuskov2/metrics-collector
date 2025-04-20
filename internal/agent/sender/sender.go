package sender

import (
	"fmt"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/levigross/grequests"
)

func SendMetrics(address string, metrics map[string]models.Metric) error {
	for metricName, metricData := range metrics {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", address, metricData.Type, metricName, metricData.Value)
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
