package storage

import (
	"context"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
)

type MemStorage struct {
	memStorage map[string]*models.Metrics
	DB         *DBStorage
}

func NewMemStorage(cfg *config.ServerConfig) *MemStorage {
	db, err := NewDBStorage(cfg)
	if err != nil {
		return nil
	}
	return &MemStorage{
		memStorage: make(map[string]*models.Metrics),
		DB:         db,
	}
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return s.DB.Ping(ctx)
	// return nil
}

func (s *MemStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return s.memStorage, nil
}

func (s *MemStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	if metric, ok := s.memStorage[metricName]; ok {
		return metric, true
	}
	return nil, false
}

func (s *MemStorage) UpdateMetric(metric *models.Metrics) error {
	s.memStorage[metric.ID] = metric
	return nil
}
