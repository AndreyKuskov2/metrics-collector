package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/handlers"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/router"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/services"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/storage"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/AndreyKuskov2/metrics-collector/proto/metrics"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	logger := logger.NewLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Info("failed to get config")
		return
	}

	stor, err := storage.NewStorage(context.Background(), cfg, logger)
	if err != nil {
		logger.Fatalf("failed to create repository: %v", err)
	}
	service := services.NewMetricService(stor, logger)
	handler := handlers.NewMetricHandler(service, logger)
	grpcHandler := handlers.NewGRPCHandler(service, logger)

	metricRouter := router.GetRouter(cfg, logger, handler)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if cfg.CryptoKey != "" {
			certFile := "certs/cert.crt"

			logger.Infof("Start HTTPS web-server on %s", cfg.Address)
			if err := http.ListenAndServeTLS(cfg.Address, certFile, cfg.CryptoKey, metricRouter); err != nil {
				logger.Fatalf("Failed to start server: %v", err)
			}
		} else {
			logger.Infof("Start web-server on %s", cfg.Address)
			if err := http.ListenAndServe(cfg.Address, metricRouter); err != nil {
				logger.Fatalf("Failed to start server: %v", err)
			}
		}
	}()

	go func() {
		listen, err := net.Listen("tcp", cfg.GRPCAddress)
		if err != nil {
			logger.Fatalf("Failed to listen connection: %v", err)
		}

		// Настройка базового логгирования
		const component = "grpc-example"
		grpclogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
		rpcLogger := grpclogger.With("service", "gRPC/server", "component", component)

		s := grpc.NewServer(grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(rpcLogger)),
		))

		pb.RegisterMetricsServiceServer(s, grpcHandler)

		reflection.Register(s)

		logger.Infof("Start GRPC server on %s", cfg.GRPCAddress)
		if err := s.Serve(listen); err != nil {
			logger.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	<-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
}
