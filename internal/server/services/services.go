package services

import (
	"fmt"
	"net/http"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"
	"github.com/sirupsen/logrus"
)

type IMetricService interface {
	UpdateMetric(metric *models.Metrics) error
	GetMetric(metricName string) (*models.Metrics, bool)
	GetAllMetrics() (map[string]*models.Metrics, error)
	Ping() error
	UpdateBatchMetrics(metrics []models.Metrics) error
}

type MetricService struct {
	storageRepo IMetricService
	logger      *logrus.Logger
}

func NewMetricService(storageRepo IMetricService, logger *logrus.Logger) *MetricService {
	return &MetricService{
		storageRepo: storageRepo,
		logger:      logger,
	}
}

func (s *MetricService) Ping() error {
	return s.storageRepo.Ping()
}

func (s *MetricService) localUpdateMetric(requestMetric *models.Metrics) (*models.Metrics, error) {
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
	return metric, nil
}

func (s *MetricService) UpdateMetric(requestMetric *models.Metrics) (*models.Metrics, error) {
	metric, _ := s.localUpdateMetric(requestMetric)

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

func (s *MetricService) UpdateBatchMetricsServ(metrics []models.Metrics, r *http.Request) error {
	if len(metrics) == 0 {
		return fmt.Errorf("empty metrics")
	}

	// Валидация данных
	for _, metric := range metrics {
		if err := metric.Bind(r); err != nil {
			return err
		}
	}

	if err := s.storageRepo.UpdateBatchMetrics(metrics); err != nil {
		return fmt.Errorf("failed to update received metrics: %s", err)
	}

	return nil
}
