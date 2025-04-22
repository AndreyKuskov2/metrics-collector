package storage

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

type StorageRepo interface {
	UpdateMetric(metricName string) error
	GetMetric(metricName string) (*models.Metrics, bool)
	GetAllMetrics() (map[string]*models.Metrics, error)
}

type Storage struct {
	memStorage map[string]*models.Metrics
}

func NewStorage() *Storage {
	return &Storage{
		memStorage: make(map[string]*models.Metrics),
	}
}

func (s *Storage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return s.memStorage, nil
}

func (s *Storage) GetMetric(metricName string) (*models.Metrics, bool) {
	if metric, ok := s.memStorage[metricName]; ok {
		return metric, true
	}
	return nil, false
}

func (s *Storage) UpdateMetric(metric *models.Metrics) error {
	s.memStorage[metric.ID] = metric
	return nil
}
