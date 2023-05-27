package storage

import (
	"context"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

type memStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var memS storage = &memStorage{}

var defaultStorage memStorage

func (m *memStorage) Сreate(ctx context.Context, cfg *config.ServerConfig) storage {
	defaultStorage.Gauges = make(map[string]float64)
	defaultStorage.Counters = make(map[string]int64)
	return memS
}

func (m *memStorage) SetMetric(ctx context.Context, metric metrics.Metric) error {
	switch metric.MType {
	case metrics.CounterType:
		setCounter(metric.ID, *metric.Delta)
		return nil
	case metrics.GaugeType:
		setGauge(metric.ID, *metric.Value)
		return nil
	default:
		return fmt.Errorf("usupported metric type")
	}
}

func (m *memStorage) SetMetricsBatch(ctx context.Context, metric []metrics.Metric) error {
	for _, metr := range metric {
		if err := m.SetMetric(ctx, metr); err != nil {
			return fmt.Errorf("cannot set butch metrics %w", err)
		}
	}
	return nil
}

func (m *memStorage) GetMetric(ctx context.Context, metric metrics.Metric) (metrics.Metric, error) {
	switch metric.MType {
	case metrics.CounterType:
		return m.getCounter(metric)
	case metrics.GaugeType:
		return m.getGauge(metric)
	default:
		return metric, fmt.Errorf("cannot find type metic")
	}
}

func (m *memStorage) getCounter(metric metrics.Metric) (metrics.Metric, error) {
	v, exist := defaultStorage.Counters[metric.ID]
	if !exist {
		return metric, fmt.Errorf("cannot find countert")
	}
	metric.Delta = &v
	return metric, nil
}

func (m *memStorage) getGauge(metric metrics.Metric) (metrics.Metric, error) {
	v, exist := defaultStorage.Gauges[metric.ID]
	if !exist {
		return metric, fmt.Errorf("cannot find countert")
	}
	metric.Value = &v
	return metric, nil
}

func (m *memStorage) GetAllMetrics(ctx context.Context) ([]metrics.Metric, error) {
	res := []metrics.Metric{}
	for name := range defaultStorage.Counters {
		v := defaultStorage.Counters[name]
		res = append(res, metrics.Metric{
			ID:    name,
			MType: metrics.CounterType,
			Delta: &v,
		})
	}
	for name := range defaultStorage.Gauges {
		v := defaultStorage.Gauges[name]
		res = append(res, metrics.Metric{
			ID:    name,
			MType: metrics.GaugeType,
			Value: &v,
		})
	}
	return res, nil
}

func (m *memStorage) PingStorage(ctx context.Context) error {
	return fmt.Errorf("databse not connected")
}

func (m *memStorage) GetAll(ctx context.Context) ([]string, error) {
	names := make([]string, 0)
	for cm := range defaultStorage.Counters {
		names = append(names, cm)
	}
	for gm := range defaultStorage.Gauges {
		names = append(names, gm)
	}
	return names, nil
}

func (m *memStorage) Close(ctx context.Context) {

}

func setGauge(k string, v float64) {
	defaultStorage.Gauges[k] = v
}

func setCounter(k string, v int64) {
	if _, ok := defaultStorage.Counters[k]; !ok {
		defaultStorage.Counters[k] = 0
	}
	defaultStorage.Counters[k] += v

}
