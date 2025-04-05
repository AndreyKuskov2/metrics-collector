package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndreyKuskov2/metrics-collector/internal/storage"
)

func UpdateMetricHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			metricType := r.PathValue("metric_type")
			if metricType == "" || (metricType != "counter" && metricType != "gauge") {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-type", "text/plain; charset=utf-8")
				return
			}
			metricName := r.PathValue("metric_name")
			if metricName == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-type", "text/plain; charset=utf-8")
				return
			}
			metricValue := r.PathValue("metric_value")
			if metricValue == "" {
				// TODO: Добавить обработку ошибок что ли
				w.Header().Set("Content-type", "text/plain; charset=utf-8")
				return
			}

			// TODO: Выносим это в слой сервиса
			switch metricType {
			case "counter":
				value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					// TODO: Добавить обработку, если не удалось распарсить значения
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
					// TODO: Добавить обработку, если не удалось распарсить значения
					return
				}
				s.UpdateMetric(metricType, metricName, value)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "text/plain; charset=utf-8")
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-type", "text/plain; charset=utf-8")
		}
	}
}
