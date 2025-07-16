package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"
	"github.com/sirupsen/logrus"
)

type mockStorage struct {
	UpdateMetricFunc       func(metric *models.Metrics) error
	GetMetricFunc          func(metricName string) (*models.Metrics, bool)
	GetAllMetricsFunc      func() (map[string]*models.Metrics, error)
	PingFunc               func() error
	UpdateBatchMetricsFunc func(metrics []models.Metrics) error
}

func (m *mockStorage) UpdateMetric(metric *models.Metrics) error {
	return m.UpdateMetricFunc(metric)
}
func (m *mockStorage) GetMetric(metricName string) (*models.Metrics, bool) {
	return m.GetMetricFunc(metricName)
}
func (m *mockStorage) GetAllMetrics() (map[string]*models.Metrics, error) {
	return m.GetAllMetricsFunc()
}
func (m *mockStorage) Ping() error {
	return m.PingFunc()
}
func (m *mockStorage) UpdateBatchMetrics(metrics []models.Metrics) error {
	return m.UpdateBatchMetricsFunc(metrics)
}

func TestUpdateMetric_Success(t *testing.T) {
	mock := &mockStorage{
		UpdateMetricFunc: func(metric *models.Metrics) error { return nil },
		GetMetricFunc: func(name string) (*models.Metrics, bool) {
			return nil, false
		},
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	m := &models.Metrics{ID: "foo", MType: utils.COUNTER, Delta: ptrInt64(1)}
	res, err := svc.UpdateMetric(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "foo" || res.MType != utils.COUNTER {
		t.Errorf("unexpected metric: %+v", res)
	}
}

func TestUpdateMetric_Error(t *testing.T) {
	mock := &mockStorage{
		UpdateMetricFunc: func(metric *models.Metrics) error { return errors.New("fail") },
		GetMetricFunc:    func(name string) (*models.Metrics, bool) { return nil, false },
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	m := &models.Metrics{ID: "foo", MType: utils.COUNTER, Delta: ptrInt64(1)}
	_, err := svc.UpdateMetric(m)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetMetric_Found(t *testing.T) {
	mock := &mockStorage{
		GetMetricFunc: func(name string) (*models.Metrics, bool) {
			return &models.Metrics{ID: name, MType: utils.GAUGE, Value: ptrFloat64(3.14)}, true
		},
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	m, ok := svc.GetMetric("bar")
	if !ok || m.ID != "bar" {
		t.Errorf("expected found metric 'bar', got %+v", m)
	}
}

func TestGetMetric_NotFound(t *testing.T) {
	mock := &mockStorage{
		GetMetricFunc: func(name string) (*models.Metrics, bool) { return nil, false },
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	m, ok := svc.GetMetric("baz")
	if ok || m != nil {
		t.Errorf("expected not found, got %+v", m)
	}
}

func TestGetAllMetrics_Success(t *testing.T) {
	mock := &mockStorage{
		GetAllMetricsFunc: func() (map[string]*models.Metrics, error) {
			return map[string]*models.Metrics{"foo": {ID: "foo"}}, nil
		},
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	metrics, err := svc.GetAllMetrics()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := metrics["foo"]; !ok {
		t.Errorf("expected key 'foo' in metrics")
	}
}

func TestGetAllMetrics_Error(t *testing.T) {
	mock := &mockStorage{
		GetAllMetricsFunc: func() (map[string]*models.Metrics, error) { return nil, errors.New("fail") },
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	_, err := svc.GetAllMetrics()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPing(t *testing.T) {
	mock := &mockStorage{
		PingFunc: func() error { return nil },
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	if err := svc.Ping(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateBatchMetricsServ_Empty(t *testing.T) {
	mock := &mockStorage{}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	err := svc.UpdateBatchMetricsServ([]models.Metrics{}, &http.Request{})
	if err == nil {
		t.Error("expected error for empty metrics, got nil")
	}
}

func TestUpdateBatchMetricsServ_ValidationError(t *testing.T) {
	mock := &mockStorage{}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	invalid := models.Metrics{ID: "", MType: ""}
	err := svc.UpdateBatchMetricsServ([]models.Metrics{invalid}, &http.Request{})
	if err == nil {
		t.Error("expected validation error, got nil")
	}
}

func TestUpdateBatchMetricsServ_StorageError(t *testing.T) {
	mock := &mockStorage{
		UpdateBatchMetricsFunc: func(metrics []models.Metrics) error { return errors.New("fail") },
	}
	logger := logrus.New()
	svc := NewMetricService(mock, logger)
	valid := models.Metrics{ID: "foo", MType: utils.COUNTER, Delta: ptrInt64(1)}
	err := svc.UpdateBatchMetricsServ([]models.Metrics{valid}, &http.Request{})
	if err == nil {
		t.Error("expected storage error, got nil")
	}
}

// helpers
func ptrInt64(v int64) *int64       { return &v }
func ptrFloat64(v float64) *float64 { return &v }
