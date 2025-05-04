package services

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"
)

type IMetricService interface {
	UpdateMetric(metric *models.Metrics) error
	GetMetric(metricName string) (*models.Metrics, bool)
	GetAllMetrics() (map[string]*models.Metrics, error)
}

type MetricService struct {
	storageRepo IMetricService
}

func NewMetricService(storageRepo IMetricService) *MetricService {
	return &MetricService{
		storageRepo: storageRepo,
	}
}

func (s *MetricService) UpdateMetric(requestMetric *models.Metrics) (*models.Metrics, error) {
	var metric *models.Metrics
	switch requestMetric.MType {
	case utils.COUNTER:
		oldMetric, ok := s.storageRepo.GetMetric(requestMetric.ID)
		if ok {
			totalDelta := *oldMetric.Delta + *requestMetric.Delta
			metric = &models.Metrics{
				ID:    requestMetric.ID,
				MType: requestMetric.MType,
				Delta: &totalDelta,
			}
		} else {
			metric = &models.Metrics{
				ID:    requestMetric.ID,
				MType: requestMetric.MType,
				Delta: requestMetric.Delta,
			}
		}
	case utils.GAUGE:
		metric = &models.Metrics{
			ID:    requestMetric.ID,
			MType: requestMetric.MType,
			Value: requestMetric.Value,
		}
	}

	if err := s.storageRepo.UpdateMetric(metric); err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *MetricService) GetMetric(metricName string) (*models.Metrics, bool) {
	responseMetric, ok := s.storageRepo.GetMetric(metricName)
	if ok {
		return responseMetric, ok
	}
	return nil, ok
}

func (s *MetricService) GetAllMetrics() (map[string]*models.Metrics, error) {
	metrics, err := s.storageRepo.GetAllMetrics()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
