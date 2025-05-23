package main

import (
	"maps"
	"sync"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/collector"
	"github.com/AndreyKuskov2/metrics-collector/internal/agent/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/agent/sender"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	pollCount    int64
	metrics      map[string]models.Metrics
	metricsMutex sync.Mutex
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

	if cfg.RateLimit == 0 {
		for {
			select {
			case <-tickerPoll.C:
				pollCount++
				metrics = collector.CollectMetrics(pollCount)
			case <-tickerReport.C:
				sender.SendMetrics(cfg.Address, metrics, logger)
				sender.SendMetricsJSON(cfg.Address, metrics, logger)
				sender.SendMetricsBatch(cfg, models.AllMetrics{RuntimeMetrics: metrics}, logger)
				logger.Info("Sent metrics")
			}
		}
	} else {
		metricsChan := make(chan models.AllMetrics, cfg.RateLimit)
		var wg sync.WaitGroup

		// Запускаем воркеры
		for i := 0; i < cfg.RateLimit; i++ {
			wg.Add(1)
			go worker(metricsChan, &wg, cfg, logger)
		}

		go func() {
			for range tickerPoll.C {
				pollCount++
				metricsMutex.Lock()
				runtimeMetrics := collector.CollectMetrics(pollCount)
				metricsMutex.Unlock()

				metricsChan <- models.AllMetrics{RuntimeMetrics: runtimeMetrics}
			}
		}()

		// Горутина для сбора дополнительных метрик
		go func() {
			for range tickerPoll.C {
				metricsMutex.Lock()
				additionalMetrics := collector.CollectAdditionMetrics()
				metricsMutex.Unlock()

				metricsChan <- models.AllMetrics{AdditionalMetrics: additionalMetrics}
			}
		}()

		// Горутина для отправки метрик на сервер
		go func() {
			for range tickerReport.C {
				metricsMutex.Lock()
				var combinedMetrics models.AllMetrics
				for i := 0; i < cfg.RateLimit; i++ {
					metrics := <-metricsChan
					maps.Copy(combinedMetrics.RuntimeMetrics, metrics.RuntimeMetrics)
					maps.Copy(combinedMetrics.AdditionalMetrics, metrics.AdditionalMetrics)
				}
				metricsMutex.Unlock()

				maps.Copy(combinedMetrics.RuntimeMetrics, combinedMetrics.AdditionalMetrics)
				sender.SendMetricsBatch(cfg, combinedMetrics, logger)
			}
		}()

		wg.Wait()
	}
}

func worker(metricsChan chan models.AllMetrics, wg *sync.WaitGroup, config *config.AgentConfig, logger *logrus.Logger) {
	defer wg.Done()
	for metrics := range metricsChan {
		maps.Copy(metrics.RuntimeMetrics, metrics.AdditionalMetrics)
		sender.SendMetricsBatch(config, metrics, logger)
	}
}
