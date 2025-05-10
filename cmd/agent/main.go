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
	logger := logger.NewLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Info("failed to get config")
		return
	}

	tickerPoll := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-tickerPoll.C:
			pollCount++
			metrics = collector.CollectMetrics(pollCount)
		case <-tickerReport.C:
			// sender.SendMetrics(cfg.Address, metrics, logger)
			// sender.SendMetricsJSON(cfg.Address, metrics, logger)
			sender.SendMetricsBatch(cfg, metrics, logger)
			logger.Info("Sent metrics")
		}
	}
}
