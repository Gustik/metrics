package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Gustik/metrics/internal/model"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
	ErrUnknownMetric  = errors.New("unknown metric type")
	ErrGaugeValue     = errors.New("gauge metric requires value")
	ErrCounterDelta   = errors.New("counter metric requires delta")
)

type MetricKey struct {
	ID    string      // имя метрики
	MType model.MType // тип метрики
}

type MetricRepository interface {
	Save(ctx context.Context, newMetric model.Metric) error
	Get(ctx context.Context, metricType model.MType, name string) (*model.Metric, error)
}

type inMemoryMetricRepository struct {
	mu      sync.RWMutex
	metrics map[MetricKey]*model.Metric
}

func NewInMemoryMetricRepository() *inMemoryMetricRepository {
	return &inMemoryMetricRepository{
		metrics: make(map[MetricKey]*model.Metric),
	}
}

func (r *inMemoryMetricRepository) Save(ctx context.Context, newMetric model.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := MetricKey{ID: newMetric.ID, MType: newMetric.MType}

	if newMetric.MType == model.Gauge {
		if newMetric.Value == nil {
			return ErrGaugeValue
		}

		r.metrics[key] = &newMetric

		return nil
	}

	if newMetric.MType == model.Counter {
		if newMetric.Delta == nil {
			return ErrCounterDelta
		}

		existsMetric := r.metrics[key]
		if existsMetric == nil {
			r.metrics[key] = &newMetric
			return nil
		}

		*existsMetric.Delta += *newMetric.Delta

		return nil
	}

	return ErrUnknownMetric
}

func (r *inMemoryMetricRepository) Get(ctx context.Context, metricType model.MType, name string) (*model.Metric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := MetricKey{ID: name, MType: metricType}
	metric := r.metrics[key]

	if metric == nil {
		return nil, ErrMetricNotFound
	}

	return metric, nil
}
