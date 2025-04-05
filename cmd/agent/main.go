package main

import (
	// "fmt"

	"fmt"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/collector"
	"github.com/AndreyKuskov2/metrics-collector/internal/agent/sender"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

var (
	pollCount int64
	metrics   map[string]models.Metric
)

func main() {
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(10) * time.Second

	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	for {
		select {
		case <-tickerPoll.C:
			pollCount++
			metrics = collector.CollectMetrics(pollCount)
		case <-tickerReport.C:
			sender.SendMetrics(metrics)
			fmt.Println("Sent metrics")
		}
	}
}
