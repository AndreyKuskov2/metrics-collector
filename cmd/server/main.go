package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/handlers"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/router"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/services"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/storage"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
)

func main() {
	logger := logger.NewLogger("./logs/server.log")
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Info("failed to get config")
		return
	}

	stor, err := storage.NewStorage(context.Background(), cfg, logger)
	if err != nil {
		logger.Fatalf("failed to create repository: %v", err)
	}
	service := services.NewMetricService(stor)
	handler := handlers.NewMetricHandler(service, logger)

	metricRouter := router.GetRouter(logger, handler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Start web-server on %s", cfg.Address)
		if err := http.ListenAndServe(cfg.Address, metricRouter); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
}
