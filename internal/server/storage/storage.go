package storage

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

type StorageRepo interface {
	UpdateMetric(metricName string) error
	GetMetric(metricName string) (*models.Metric, error)
	GetAllMetrics() (map[string]*models.Metric, error)
}

type Storage struct {
	MemStorage map[string]*models.Metric
}

func NewStorage() *Storage {
	return &Storage{
		MemStorage: make(map[string]*models.Metric),
	}
}

func (s *Storage) GetAllMetrics() (map[string]*models.Metric, error) {
	return s.MemStorage, nil
}

func (s *Storage) GetMetric(metricName string) (*models.Metric, bool) {
	metric, ok := s.MemStorage[metricName]
	if !ok {
		return &models.Metric{}, false
	}
	return metric, true
}

func (s *Storage) UpdateMetric(metricType, metricName string, metricValue interface{}) error {
	metric, ok := s.GetMetric(metricName)
	if !ok {
		s.MemStorage[metricName] = &models.Metric{
			Type:  metricType,
			Value: metricValue,
		}
	} else {
		metric.Value = metricValue
	}
	return nil
}
