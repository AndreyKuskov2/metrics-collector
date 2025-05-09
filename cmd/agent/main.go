package main

import (
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/collector"
	"github.com/AndreyKuskov2/metrics-collector/internal/agent/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/agent/sender"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
)

var (
	pollCount int64
	metrics   map[string]models.Metrics
)

func main() {
	c := config.NewConfig()
	logger := logger.NewLogger("./logs/agent.log")

	tickerPoll := time.NewTicker(time.Duration(c.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(c.ReportInterval) * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			pollCount++
			metrics = collector.CollectMetrics(pollCount)
		case <-tickerReport.C:
			sender.SendMetrics(c.Address, metrics, logger)
			sender.SendMetricsJSON(c.Address, metrics, logger)
			logger.Info("Sent metrics")
		}
	}
}
