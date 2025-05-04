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
	c, err := config.NewConfig()
	if err != nil {
		logger.Info("failed to get config")
		return
	}

	stor := storage.NewStorage()
	service := services.NewMetricService(stor)
	handler := handlers.NewMetricHandler(service)

	metricRouter := router.GetRouter(logger, handler)

	storage.StartFileStorageLogic(c, stor, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("Start web-server on %s", c.Address)
		if err := http.ListenAndServe(c.Address, metricRouter); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop

	// Создание контекста с тайм-аутом для завершения работы сервера
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stor.SaveMemStorageToFile()
	stor.CloseFile()
  
	// Логирование завершения работы сервера
	logger.Info("Shutting down server...")
}
