package router

import (
	"net/http"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/handlers"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func GetRouter(cfg *config.ServerConfig, logger *logrus.Logger, h *handlers.MetricHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware(logger))
	middlewares.CheckHashMiddleware(cfg)
	r.Use(middleware.Compress(5, "text/html", "application/json"))

	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", h.UpdateMetricHandler)
	r.With(middlewares.GzipMiddleware).Post("/update/", h.UpdateMetricHandlerJSON)

	r.Get("/value/{metric_type}/{metric_name}", h.GetMetricHandler)
	r.With(middlewares.GzipMiddleware).Post("/value/", h.GetMetricHandlerJSON)

	r.With(middlewares.GzipMiddleware).Get("/", h.GetMetricsHandler)

	r.Get("/ping", h.Ping)

	r.With(middlewares.GzipMiddleware).With(middlewares.CheckHashMiddleware(cfg)).Post("/updates/", h.UpdateManyMetricsHandler)

	return r
}
