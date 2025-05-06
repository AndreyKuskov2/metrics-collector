package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/AndreyKuskov2/metrics-collector/internal/models"
	"github.com/AndreyKuskov2/metrics-collector/internal/server/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type IMetricHandler interface {
	UpdateMetric(requestMetric *models.Metrics) (*models.Metrics, error)
	GetMetric(metricName string) (*models.Metrics, bool)
	GetAllMetrics() (map[string]*models.Metrics, error)
	Ping() error
}

type MetricHandler struct {
	services IMetricHandler
	logger   *logrus.Logger
}

func NewMetricHandler(services IMetricHandler, logger *logrus.Logger) *MetricHandler {
	return &MetricHandler{
		services: services,
		logger:   logger,
	}
}

func (mh *MetricHandler) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metric_type")
	if metricType == "" || (metricType != utils.COUNTER && metricType != utils.GAUGE) {
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

	var requestMetric *models.Metrics
	switch metricType {
	case utils.COUNTER:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "")
			return
		}
		requestMetric = &models.Metrics{
			ID:    metricName,
			MType: metricType,
			Delta: &value,
		}
	case utils.GAUGE:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, "")
			return
		}
		requestMetric = &models.Metrics{
			ID:    metricName,
			MType: metricType,
			Value: &value,
		}
	}

	_, err := mh.services.UpdateMetric(requestMetric)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "")
}

func (mh *MetricHandler) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricName := chi.URLParam(r, "metric_name")
	if metricName == "" {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
		return
	}

	metric, ok := mh.services.GetMetric(metricName)
	if !ok {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
		return
	}

	switch metric.MType {
	case utils.COUNTER:
		value := fmt.Sprintf("%v", *metric.Delta)
		render.PlainText(w, r, value)
	case utils.GAUGE:
		value := fmt.Sprintf("%v", *metric.Value)
		render.PlainText(w, r, value)
	}
}

func (mh *MetricHandler) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics, err := mh.services.GetAllMetrics()
	if err != nil {
		mh.logger.Infof("all metrics handler error: %s", err)
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
		return
	}
	tmpl, err := template.ParseFiles("./web/template/index.html")
	if err != nil {
		mh.logger.Errorf("cannot parse template file: %s", err)
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "")
		return
	}
	tmpl.Execute(w, metrics)
}

func (mh *MetricHandler) UpdateMetricHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var requestMetric models.Metrics

	if err := render.Bind(r, &requestMetric); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	metric, err := mh.services.UpdateMetric(&requestMetric)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	if responseMetric, ok := mh.services.GetMetric(metric.ID); ok {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, responseMetric)
		return
	}
	render.Status(r, http.StatusNotFound)
	render.PlainText(w, r, "")
}

func (mh *MetricHandler) GetMetricHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	if err := render.Bind(r, &metric); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	if responseMetric, ok := mh.services.GetMetric(metric.ID); ok {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, responseMetric)
		return
	}

	render.Status(r, http.StatusNotFound)
	render.PlainText(w, r, "")
}

func (mh *MetricHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := mh.services.Ping(); err != nil {
		mh.logger.Infof("ping error: %s", err)
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "")
}
