package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func UpdateMetricHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			metricType := chi.URLParam(r, "metric_type")
			if metricType == "" || (metricType != "counter" && metricType != "gauge") {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}
			metricName := chi.URLParam(r, "metric_name")
			if metricName == "" {
				render.Status(r, http.StatusNotFound)
				render.PlainText(w, r, "")
				return
			}
			metricValue := chi.URLParam(r, "metric_value")
			if metricValue == "" {
				render.Status(r, http.StatusNotFound)
				render.PlainText(w, r, "")
				return
			}

			// TODO: Выносим это в слой сервиса
			switch metricType {
			case "counter":
				value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					render.Status(r, http.StatusBadRequest)
					render.PlainText(w, r, "")
					return
				}
				oldMetric, ok := s.GetMetric(metricName)
				if !ok {
					s.UpdateMetric(metricType, metricName, value)
				} else {
					s.UpdateMetric(metricType, metricName, oldMetric.Value.(int64)+value)
				}
			case "gauge":
				value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					render.Status(r, http.StatusBadRequest)
					render.PlainText(w, r, "")
					return
				}
				s.UpdateMetric(metricType, metricName, value)
			}

			render.Status(r, http.StatusOK)
			render.PlainText(w, r, "")
		} else {
			render.Status(r, http.StatusMethodNotAllowed)
			render.PlainText(w, r, "")
		}
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

		value := fmt.Sprintf("%v", metric.Value)
		render.PlainText(w, r, value)
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

func main() {
	c := config.NewConfig()

	s := storage.NewStorage()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", UpdateMetricHandler(s))
	r.Get("/value/{metric_type}/{metric_name}", GetMetricHandler(s))
	r.Get("/", GetMetricsHandler(s))

	log.Printf("Start web-server on %s", c.Address)
	if err := http.ListenAndServe(c.Address, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
