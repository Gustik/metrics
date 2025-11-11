package collector

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/Gustik/metrics/internal/model"
)

type MetricCollector interface {
	Collect() []model.Metric
}

type memStatsCollector struct {
	pollCount   int64
	randomValue float64
}

func NewCollector() MetricCollector {
	return &memStatsCollector{}
}

func (c *memStatsCollector) Collect() []model.Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	mestatsFields := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction",
		"GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
		"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse",
		"MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs",
		"StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	metrics := make([]model.Metric, len(mestatsFields))

	// Получаем значение структуры через рефлексию
	v := reflect.ValueOf(memStats)

	for _, fieldName := range mestatsFields {
		field := v.FieldByName(fieldName)

		var value float64
		switch field.Kind() {
		case reflect.Uint64:
			value = float64(field.Uint())
		case reflect.Float64:
			value = field.Float()
		case reflect.Uint32:
			value = float64(field.Uint())
		}

		metrics = append(metrics, model.Metric{ID: fieldName, MType: model.Gauge, Value: &value})
	}

	c.pollCount++
	metrics = append(metrics, model.Metric{ID: "pollCount", MType: model.Counter, Delta: &c.pollCount})

	c.randomValue = rand.Float64()
	metrics = append(metrics, model.Metric{ID: "randomValue", MType: model.Gauge, Value: &c.randomValue})

	return metrics
}
