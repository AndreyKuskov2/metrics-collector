package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/AndreyKuskov2/metrics-collector/internal/middlewares"
	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/storage"
	"github.com/AndreyKuskov2/metrics-collector/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func UpdateMetricHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metric_type")
		if metricType == "" || (metricType != "counter" && metricType != "gauge") {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "")
			return
		}
		metricName := chi.URLParam(r, "metric_name")
		metricValue := chi.URLParam(r, "metric_value")

		if metricName == "" || metricValue == "" {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
			return
		}

		// TODO: Выносим это в слой сервиса
		var metric *models.Metrics
		switch metricType {
		case "counter":
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}
			oldMetric, ok := s.GetMetric(metricName)
			if ok {
				totalDelta := *oldMetric.Delta + value
				metric = &models.Metrics{
					ID:    metricName,
					MType: metricType,
					Delta: &totalDelta,
				}
			} else {
				metric = &models.Metrics{
					ID:    metricName,
					MType: metricType,
					Delta: &value,
				}
			}
		case "gauge":
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}
			metric = &models.Metrics{
				ID:    metricName,
				MType: metricType,
				Value: &value,
			}
		}

		if err := s.UpdateMetric(metric); err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "")
	}
}

func GetMetricHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricName := chi.URLParam(r, "metric_name")
		if metricName == "" {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
			return
		}

		metric, ok := s.GetMetric(metricName)
		if !ok {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
			return
		}

		switch metric.MType {
		case "counter":
			value := fmt.Sprintf("%v", *metric.Delta)
			render.PlainText(w, r, value)
		case "gauge":
			value := fmt.Sprintf("%v", *metric.Value)
			render.PlainText(w, r, value)
		}
	}
}

func GetMetricsHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := s.GetAllMetrics()
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
			return
		}
		tmpl, err := template.ParseFiles("./web/template/index.html")
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "")
			return
		}
		tmpl.Execute(w, metrics)
	}
}

func UpdateMetricHandlerJSON(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetric models.Metrics

		if err := render.Bind(r, &requestMetric); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		var metric *models.Metrics
		switch requestMetric.MType {
		case "counter":
			oldMetric, ok := s.GetMetric(requestMetric.ID)
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
		case "gauge":
			metric = &models.Metrics{
				ID:    requestMetric.ID,
				MType: requestMetric.MType,
				Value: requestMetric.Value,
			}
		}

		if err := s.UpdateMetric(metric); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		if responseMetric, ok := s.GetMetric(metric.ID); ok {
			render.Status(r, http.StatusOK)
			render.JSON(w, r, responseMetric)
			return
		}
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
	}
}

func GetMetricHandlerJSON(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics

		if err := render.Bind(r, &metric); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}

		if responseMetric, ok := s.GetMetric(metric.ID); ok {
			render.Status(r, http.StatusOK)
			render.JSON(w, r, responseMetric)
			return
		}

		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
	}
}

func main() {
	c := config.NewConfig()
	logger := logger.NewLogger("./logs/server.log")

	s := storage.NewStorage()

	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware(logger))

	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", UpdateMetricHandler(s))
	r.Post("/update/", UpdateMetricHandlerJSON(s))

	r.Get("/value/{metric_type}/{metric_name}", GetMetricHandler(s))
	r.Post("/value/", GetMetricHandlerJSON(s))

	r.Get("/", GetMetricsHandler(s))

	logger.Printf("Start web-server on %s", c.Address)
	if err := http.ListenAndServe(c.Address, r); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
