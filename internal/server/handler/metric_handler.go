package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gustik/metrics/internal/model"
	"github.com/Gustik/metrics/internal/server/service"
)

type MetricHandler struct {
	service service.MetricService
}

func NewMetricHandler(metricService service.MetricService) *MetricHandler {
	return &MetricHandler{
		service: metricService,
	}
}

func (h *MetricHandler) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	delta, err := strconv.ParseInt(r.PathValue("value"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), model.Counter, metric, &delta, nil)
	if err != nil {
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	metric := r.PathValue("metric")
	value, err := strconv.ParseFloat(r.PathValue("value"), 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	err = h.service.Update(r.Context(), model.Gauge, metric, nil, &value)
	if err != nil {
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) GetCounter(w http.ResponseWriter, r *http.Request) {
	metricName := r.PathValue("metric")

	metric, err := h.service.Get(r.Context(), model.Counter, metricName)
	if err != nil {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", *metric.Delta)
}

func (h *MetricHandler) GetGauge(w http.ResponseWriter, r *http.Request) {
	metricName := r.PathValue("metric")

	metric, err := h.service.Get(r.Context(), model.Gauge, metricName)
	if err != nil {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%f", *metric.Value)
}
