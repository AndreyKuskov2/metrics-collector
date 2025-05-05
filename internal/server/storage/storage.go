package storage

import (
	"context"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/sirupsen/logrus"
)

type Storager interface {
	GetAllMetrics() (map[string]*models.Metrics, error)
	GetMetric(metricName string) (*models.Metrics, bool)
	UpdateMetric(metric *models.Metrics) error
	Ping(ctx context.Context) error
}

//  && cfg.DatabaseDSN == ""
// else if cfg.DatabaseDSN != "" {
// 	logger.Info("Selected storage: DB")
// 	DB, err := NewDBStorage(cfg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return DB, nil
// }

func NewStorage(cfg *config.ServerConfig, logger *logrus.Logger) (Storager, error) {
	if cfg.FileStoragePath == "" {
		logger.Info("No storage selected using default: MemoryStorage")
		return NewMemStorage(cfg), nil
	} else {
		logger.Info("Selected storage: File")
		storage := NewFileMemStorage(cfg)
		StartFileStorageLogic(cfg, storage, logger)
		return storage, nil
	}
}
