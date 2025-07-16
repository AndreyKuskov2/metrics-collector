package storage

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
)

// Базовая структура для хранения метрик в памяти
type MemStorage struct {
	memStorage map[string]*models.Metrics
}

// Создание структуры для хранения метрик
func NewMemStorage() *MemStorage {
	return &MemStorage{
		memStorage: make(map[string]*models.Metrics),
	}
}

func (s *MemStorage) UpdateBatchMetrics(metrics []models.Metrics) error {
	return nil
}

func (s *MemStorage) Ping() error {
	return nil
}

// Получение всех метрик из памяти
func (s *MemStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return s.memStorage, nil
}

// Получение конкретной метрики из памяти
func (s *MemStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	if metric, ok := s.memStorage[metricName]; ok {
		return metric, true
	}
	return nil, false
}

// Обновление метрики
func (s *MemStorage) UpdateMetric(metric *models.Metrics) error {
	s.memStorage[metric.ID] = metric
	return nil
}
