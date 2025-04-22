package main

import (
	"net/http"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/handlers"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/router"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/services"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/storage"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
)

func main() {
	logger := logger.NewLogger("./logs/server.log")
	c, err := config.NewConfig()
	if err != nil {
		logger.Info("failed to get config")
		return
	}

	// var actualStorage services.MetricStorage
	// if с.StoreInterval == 0 {
	// 	actualStorage = storage.NewAutoDump(memStorage, dump)
	// } else {
	// 	actualStorage = memStorage
	// 	go func() {
	// 		ticker := time.NewTicker(cfg.StoreInterval)
	// 		defer ticker.Stop()

	// 		for range ticker.C {
	// 			if err := dump.SaveMetricToFile(); err != nil {
	// 				logger.Log.Info("error save to file", zap.Any("error", err))
	// 			} else {
	// 				logger.Log.Info("metrics saved to file successfully")
	// 			}
	// 		}
	// 	}()
	// }

	// Загрузка данных из файла
	if c.Restore {
		// err := dump.LoadMetricsFromFile()
		// if err != nil {
		// 	logger.Log.Info("error load from file", zap.Any("error", err))
		// }
	}

	storage := storage.NewStorage()
	service := services.NewMetricService(storage)
	handler := handlers.NewMetricHandler(service)

	metricRouter := router.GetRouter(logger, handler)

	logger.Printf("Start web-server on %s", c.Address)
	if err := http.ListenAndServe(c.Address, metricRouter); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
