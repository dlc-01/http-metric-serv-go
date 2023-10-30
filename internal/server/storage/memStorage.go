package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

type memStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var memS Storage = &memStorage{}

var mux sync.RWMutex

func (m *memStorage) Сreate(ctx context.Context, cfg *config.ServerConfig) Storage {
	m.Gauges = make(map[string]float64)
	m.Counters = make(map[string]int64)
	return memS
}

func (m *memStorage) SetMetric(ctx context.Context, metric metrics.Metric) error {
	switch metric.MType {
	case metrics.CounterType:
		m.setCounter(metric.ID, *metric.Delta)
		return nil
	case metrics.GaugeType:
		m.setGauge(metric.ID, *metric.Value)
		return nil
	default:
		return fmt.Errorf("unsupported metric type")
	}
}

func (m *memStorage) setGauge(k string, v float64) {
	mux.Lock()
	defer mux.Unlock()

	m.Gauges[k] = v
}

func (m *memStorage) setCounter(k string, v int64) {
	mux.Lock()
	defer mux.Unlock()

	if _, ok := m.Counters[k]; !ok {
		m.Counters[k] = 0
	}
	m.Counters[k] += v

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
	mux.RLock()
	defer mux.RUnlock()

	v, exist := m.Counters[metric.ID]
	if !exist {
		return metric, fmt.Errorf("cannot find countert")
	}
	metric.Delta = &v
	return metric, nil
}

func (m *memStorage) getGauge(metric metrics.Metric) (metrics.Metric, error) {
	mux.RLock()
	defer mux.RUnlock()

	v, exist := m.Gauges[metric.ID]
	if !exist {
		return metric, fmt.Errorf("cannot find countert")
	}
	metric.Value = &v
	return metric, nil
}

func (m *memStorage) GetAllMetrics(ctx context.Context) ([]metrics.Metric, error) {
	mux.RLock()
	defer mux.RUnlock()

	res := []metrics.Metric{}
	for name := range m.Counters {
		v := m.Counters[name]
		res = append(res, metrics.Metric{
			ID:    name,
			MType: metrics.CounterType,
			Delta: &v,
		})
	}
	for name := range m.Gauges {
		v := m.Gauges[name]
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

func (m *memStorage) GetAllStrings(ctx context.Context) ([]string, error) {
	mux.RLock()
	defer mux.RUnlock()
	names := make([]string, 0)
	for cm := range m.Counters {
		names = append(names, cm)
	}
	for gm := range m.Gauges {
		names = append(names, gm)
	}
	return names, nil
}

func (m *memStorage) Close(ctx context.Context) error {
	return nil
}
