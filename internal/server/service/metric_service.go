package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gustik/metrics/internal/model"
	"github.com/Gustik/metrics/internal/server/repository"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

type MetricService interface {
	Update(ctx context.Context, metricType model.MType, name string, delta *int64, value *float64) error
	Get(ctx context.Context, metricType model.MType, name string) (*model.Metric, error)
}

type metricService struct {
	repo repository.MetricRepository
}

func NewMetricService(repo repository.MetricRepository) *metricService {
	return &metricService{
		repo: repo,
	}
}

func (m *metricService) Update(ctx context.Context, metricType model.MType, name string, delta *int64, value *float64) error {
	metric := model.Metric{
		ID:    name,
		MType: metricType,
		Delta: delta,
		Value: value,
	}

	err := m.repo.Save(ctx, metric)
	if errors.Is(err, repository.ErrUnknownMetric) {
		return ErrInvalidMetricType
	}

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func (m *metricService) Get(ctx context.Context, metricType model.MType, name string) (*model.Metric, error) {
	return m.repo.Get(ctx, metricType, name)
}
