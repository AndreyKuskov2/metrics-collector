package collector

import (
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	tests := []struct {
		name      string
		pollCount int64
	}{
		{"Test with pollCount 0", 0},
		{"Test with pollCount 1", 1},
		{"Test with pollCount 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := CollectMetrics(tt.pollCount)
			if len(metrics) == 0 {
				t.Errorf("CollectMetrics() returned empty slice")
			}

			for metricName, metricData := range metrics {
				switch metricName {
				case "PollCount":
					if v := metricData.Delta; v != nil {
						if *v != tt.pollCount {
							t.Errorf("Expected PollCount to be %v, got %v", tt.pollCount, *v)
						}
					} else {
						t.Errorf("PollCount metric is nil")
					}
				case "RandomValue":
					if v := metricData.Value; v != nil {
						if *v < 0 || *v > 1 {
							t.Errorf("Expected RandomValue to be between 0 and 1, got %v", *v)
						}
					} else {
						t.Errorf("RandomValue metric is nil")
					}
				}
			}
		})
	}
}
