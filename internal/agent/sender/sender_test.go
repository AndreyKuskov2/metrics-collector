package sender

import (
	"errors"
	"testing"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/agent/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/jarcoal/httpmock"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

type mockLogger struct{}

func (m *mockLogger) Printf(string, ...interface{}) {}
func (m *mockLogger) Println(...interface{})        {}
func (m *mockLogger) Infof(string, ...interface{})  {}
func (m *mockLogger) Info(...interface{})           {}
func (m *mockLogger) Errorf(string, ...interface{}) {}
func (m *mockLogger) Fatalf(string, ...interface{}) {}
func (m *mockLogger) Warnf(string, ...interface{})  {}

func TestSendMetrics_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/update/counter/foo/1",
		httpmock.NewStringResponder(200, ""))

	metrics := map[string]models.Metrics{
		"foo": {MType: "counter", Delta: ptrInt64(1)},
	}
	err := SendMetrics("localhost:8080", metrics, &logrus.Logger{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendMetrics_HTTPError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(httpmock.NewErrorResponder(errors.New("fail")))

	metrics := map[string]models.Metrics{
		"foo": {MType: "counter", Delta: ptrInt64(1)},
	}
	_ = SendMetrics("localhost:8080", metrics, &logrus.Logger{})
}

func TestSendMetrics_Non200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/update/counter/foo/1",
		httpmock.NewStringResponder(500, "fail"))

	metrics := map[string]models.Metrics{
		"foo": {MType: "counter", Delta: ptrInt64(1)},
	}
	_ = SendMetrics("localhost:8080", metrics, &logrus.Logger{})
}

func TestSendMetricsJSON_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/update/",
		httpmock.NewStringResponder(200, ""))

	metrics := map[string]models.Metrics{
		"foo": {MType: "counter", Delta: ptrInt64(1)},
	}
	err := SendMetricsJSON("localhost:8080", metrics, &logrus.Logger{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendWithRetry_Exhausted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterNoResponder(httpmock.NewErrorResponder(errors.New("fail")))

	cfg := &config.AgentConfig{MaxRetries: 2, RetryDelay: 1 * time.Millisecond}
	err := sendWithRetry(cfg, grequests.RequestOptions{}, "http://localhost:8080", &logrus.Logger{})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// helpers
func ptrInt64(v int64) *int64 { return &v }
