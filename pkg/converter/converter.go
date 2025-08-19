package converter

import (
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	pb "github.com/AndreyKuskov2/metrics-collector/proto/metrics"
)

func GRPCMetricToHTTP(metric *pb.Metrics) *models.Metrics {
	return &models.Metrics{
		ID:    metric.Id,
		MType: metric.Type,
		Value: &metric.Value,
		Delta: &metric.Delta,
	}
}

func HTTPMetricToGRPC(metric *models.Metrics) *pb.Metrics {
	var value float64
	if metric.Value != nil {
		value = *metric.Value
	}

	var delta int64
	if metric.Delta != nil {
		delta = *metric.Delta
	}

	return &pb.Metrics{
		Id:    metric.ID,
		Type:  metric.MType,
		Value: value,
		Delta: delta,
	}
}

func GRPCMetricsListToHTTPList(metrics []*pb.Metrics) []models.Metrics {
	metricsList := make([]models.Metrics, len(metrics))

	for _, metric := range metrics {
		metricsList = append(metricsList, models.Metrics{
			ID:    metric.Id,
			MType: metric.Type,
			Value: &metric.Value,
			Delta: &metric.Delta,
		})
	}
	return metricsList
}

func HTTPMetricsListToGRPCList(metrics []models.Metrics) []*pb.Metrics {
	grpcMetrics := make([]*pb.Metrics, len(metrics))

	for _, metric := range metrics {
		grpcMetrics = append(grpcMetrics, HTTPMetricToGRPC(&metric))
	}
	return grpcMetrics
}

func HTTPMetricsMapToGRPC(metrics map[string]*models.Metrics) map[string]*pb.Metrics {
	grpcMetrics := make(map[string]*pb.Metrics, len(metrics))

	for key, metric := range metrics {
		grpcMetrics[key] = HTTPMetricToGRPC(metric)
	}
	return grpcMetrics
}
