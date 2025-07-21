package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type mockMetricService struct {
	UpdateMetricFunc           func(requestMetric *models.Metrics) (*models.Metrics, error)
	GetMetricFunc              func(metricName string) (*models.Metrics, bool)
	GetAllMetricsFunc          func() (map[string]*models.Metrics, error)
	PingFunc                   func() error
	UpdateBatchMetricsServFunc func(metrics []models.Metrics, r *http.Request) error
}

func (m *mockMetricService) UpdateMetric(requestMetric *models.Metrics) (*models.Metrics, error) {
	return m.UpdateMetricFunc(requestMetric)
}
func (m *mockMetricService) GetMetric(metricName string) (*models.Metrics, bool) {
	return m.GetMetricFunc(metricName)
}
func (m *mockMetricService) GetAllMetrics() (map[string]*models.Metrics, error) {
	return m.GetAllMetricsFunc()
}
func (m *mockMetricService) Ping() error {
	return m.PingFunc()
}
func (m *mockMetricService) UpdateBatchMetricsServ(metrics []models.Metrics, r *http.Request) error {
	return m.UpdateBatchMetricsServFunc(metrics, r)
}

func TestUpdateMetricHandler_Counter_Success(t *testing.T) {
	mockService := &mockMetricService{
		UpdateMetricFunc: func(requestMetric *models.Metrics) (*models.Metrics, error) {
			return requestMetric, nil
		},
	}
	logger := logrus.New()
	mh := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("POST", "/update/counter/testCounter/42", nil)
	req = muxSetURLParams(req, map[string]string{"metric_type": utils.COUNTER, "metric_name": "testCounter", "metric_value": "42"})
	rw := httptest.NewRecorder()

	mh.UpdateMetricHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}
}

func TestUpdateMetricHandler_Counter_BadValue(t *testing.T) {
	mockService := &mockMetricService{}
	logger := logrus.New()
	mh := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("POST", "/update/counter/testCounter/notanumber", nil)
	req = muxSetURLParams(req, map[string]string{"metric_type": utils.COUNTER, "metric_name": "testCounter", "metric_value": "notanumber"})
	rw := httptest.NewRecorder()

	mh.UpdateMetricHandler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rw.Code)
	}
}

func TestGetMetricHandler_Found_Counter(t *testing.T) {
	val := int64(123)
	mockService := &mockMetricService{
		GetMetricFunc: func(metricName string) (*models.Metrics, bool) {
			return &models.Metrics{ID: metricName, MType: utils.COUNTER, Delta: &val}, true
		},
	}
	logger := logrus.New()
	mh := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/value/counter/testCounter", nil)
	req = muxSetURLParams(req, map[string]string{"metric_name": "testCounter"})
	rw := httptest.NewRecorder()

	mh.GetMetricHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}
	if rw.Body.String() != "123" {
		t.Errorf("Expected body '123', got '%s'", rw.Body.String())
	}
}

func TestGetMetricHandler_NotFound(t *testing.T) {
	mockService := &mockMetricService{
		GetMetricFunc: func(metricName string) (*models.Metrics, bool) {
			return nil, false
		},
	}
	logger := logrus.New()
	mh := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/value/counter/unknown", nil)
	req = muxSetURLParams(req, map[string]string{"metric_name": "unknown"})
	rw := httptest.NewRecorder()

	mh.GetMetricHandler(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rw.Code)
	}
}

// Helper to set chi URL params in tests
func muxSetURLParams(r *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}

func ExampleMetricHandler_UpdateMetricHandler() {
	mockService := &mockMetricService{
		UpdateMetricFunc: func(requestMetric *models.Metrics) (*models.Metrics, error) {
			return requestMetric, nil
		},
	}
	logger := logrus.New()
	handler := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("POST", "/update/counter/myCounter/42", nil)
	req = muxSetURLParams(req, map[string]string{"metric_type": "counter", "metric_name": "myCounter", "metric_value": "42"})
	rw := httptest.NewRecorder()

	handler.UpdateMetricHandler(rw, req)
	fmt.Println(rw.Code)
	// Output: 200
}

func ExampleMetricHandler_GetMetricHandler() {
	val := int64(123)
	mockService := &mockMetricService{
		GetMetricFunc: func(metricName string) (*models.Metrics, bool) {
			return &models.Metrics{ID: metricName, MType: "counter", Delta: &val}, true
		},
	}
	logger := logrus.New()
	handler := NewMetricHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/value/counter/myCounter", nil)
	req = muxSetURLParams(req, map[string]string{"metric_name": "myCounter"})
	rw := httptest.NewRecorder()

	handler.GetMetricHandler(rw, req)
	fmt.Println(rw.Code)
	fmt.Println(rw.Body.String())
	// Output:
	// 200
	// 123
}
