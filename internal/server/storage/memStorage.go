package storage

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

type memStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var defaultStorage memStorage

func GetStorage() memStorage {
	return defaultStorage
}

func SetStorage(new memStorage) {
	defaultStorage = new
}

func Init() {
	defaultStorage.Gauges = make(map[string]float64)
	defaultStorage.Counters = make(map[string]int64)
}

func SetMetric(id string, t string, f *float64, i *int64) bool {
	switch t {
	case metrics.CounterType:
		SetCounter(id, *i)
		return true
	case metrics.GaugeType:
		SetGauge(id, *f)
		return true
	default:
		return false
	}
}

func GetMetric(k string, t string) (metrics.Metric, bool, bool) {
	switch t {
	case metrics.CounterType:
		v, e := GetCounter(k)
		return v, e, true
	case metrics.GaugeType:
		v, e := GetGauge(k)
		return v, e, true
	default:
		return metrics.Metric{}, false, false
	}
}

func SetGauge(k string, v float64) {
	defaultStorage.Gauges[k] = v
}

func GetGauge(k string) (metrics.Metric, bool) {
	v, exist := defaultStorage.Gauges[k]
	return metrics.Metric{
		ID:    k,
		MType: metrics.GaugeType,
		Value: &v,
	}, exist
}

func SetCounter(k string, v int64) {
	if _, ok := defaultStorage.Counters[k]; !ok {
		defaultStorage.Counters[k] = 0
	}
	defaultStorage.Counters[k] += v

}

func GetCounter(k string) (metrics.Metric, bool) {
	v, exist := defaultStorage.Counters[k]
	return metrics.Metric{
		ID:    k,
		MType: metrics.CounterType,
		Delta: &v,
	}, exist
}

func GetAll() []string {
	names := make([]string, 0)
	for cm := range defaultStorage.Counters {
		names = append(names, cm)
	}
	for gm := range defaultStorage.Gauges {
		names = append(names, gm)
	}
	return names
}

func GetMetrics() []metrics.Metric {
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
	return res
}

func SetMetrics(res []metrics.Metric) error {
	for _, m := range res {
		if !SetMetric(m.ID, m.MType, m.Value, m.Delta) {
			return fmt.Errorf("Unsuported metric type")
		}
	}
	return nil
}
