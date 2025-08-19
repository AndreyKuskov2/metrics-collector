package handlers

import (
	"context"

	"github.com/AndreyKuskov2/metrics-collector/pkg/converter"
	pb "github.com/AndreyKuskov2/metrics-collector/proto/metrics"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	services MetricServicer
	logger   *logrus.Logger
	pb.UnimplementedMetricsServiceServer
}

func NewGRPCHandler(services MetricServicer, logger *logrus.Logger) *GRPCHandler {
	return &GRPCHandler{
		services: services,
		logger:   logger,
	}
}

func (s *GRPCHandler) GetAllMetrics(_ context.Context, in *pb.GetAllMetricsRequest) (*pb.GetAllMetricsResponse, error) {
	metrics, err := s.services.GetAllMetrics()
	if err != nil {
		s.logger.Infof("all metrics handler error: %s", err)
		return nil, status.Errorf(codes.Internal, "all metrics handler error")
	}

	return &pb.GetAllMetricsResponse{
		Metrics: converter.HTTPMetricsMapToGRPC(metrics),
	}, nil
}

func (s *GRPCHandler) GetMetricByNameAndType(_ context.Context, in *pb.GetMetricByNameAndTypeRequest) (*pb.GetMetricByNameAndTypeResponse, error) {
	if responseMetric, ok := s.services.GetMetric(in.Id); ok {
		s.logger.Infof("response metric: %v", responseMetric)
		return &pb.GetMetricByNameAndTypeResponse{
			Metric: converter.HTTPMetricToGRPC(responseMetric),
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "metric not found by name")
}

func (s *GRPCHandler) UpdateMetric(_ context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	if in.Metric == nil {
		return nil, status.Errorf(codes.NotFound, "metric cannot be null")
	}

	metric, err := s.services.UpdateMetric(converter.GRPCMetricToHTTP(in.Metric))
	if err != nil {
		s.logger.Infof("cannot update metric: %v", err)
		return nil, status.Errorf(codes.Internal, "cannot update metric")
	}

	if responseMetric, ok := s.services.GetMetric(metric.ID); ok {
		return &pb.UpdateMetricResponse{
			Metric: converter.HTTPMetricToGRPC(responseMetric),
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "metric not found")
}

func (s *GRPCHandler) UpdateBatchMetrics(_ context.Context, in *pb.UpdateBatchMetricsRequest) (*pb.UpdateBatchMetricsResponse, error) {
	if err := s.services.UpdateBatchMetricsServ(converter.GRPCMetricsListToHTTPList(in.Metric), nil); err != nil {
		s.logger.Infof("Failed to update metrics: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update metrics")
	}

	return &pb.UpdateBatchMetricsResponse{
		Status: "OK",
	}, nil
}

func (s *GRPCHandler) Ping(_ context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	if err := s.services.Ping(); err != nil {
		s.logger.Infof("ping error: %s", err)
		return nil, status.Errorf(codes.Internal, "ping error")
	}
	return &pb.PingResponse{
		Status: "pong",
	}, nil
}
