package storage

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

type StorageRepo interface {
	UpdateMetric(metricName string) error
	GetMetric(metricName string) (*models.Metrics, error)
	GetAllMetrics() (map[string]*models.Metrics, error)
}

type Storage struct {
	MemStorage map[string]*models.Metrics
}

func NewStorage() *Storage {
	return &Storage{
		MemStorage: make(map[string]*models.Metrics),
	}
}

func (s *Storage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return s.MemStorage, nil
}

func (s *Storage) GetMetric(metricName string) (*models.Metrics, bool) {
	if metric, ok := s.MemStorage[metricName]; ok {
		return metric, true
	}
	return nil, false
}

func (s *Storage) UpdateMetric(metric *models.Metrics) error {
	s.MemStorage[metric.ID] = metric
	return nil
}
