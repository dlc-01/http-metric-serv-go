package storage

import (
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
func SetGauge(k string, v float64) {
	defaultStorage.Gauges[k] = v

}

func GetGauge(k string) (float64, bool) {
	v, exist := defaultStorage.Gauges[k]
	return v, exist
}

func SetCounter(k string, v int64) {
	if _, ok := defaultStorage.Counters[k]; !ok {
		defaultStorage.Counters[k] = 0
	}

	defaultStorage.Counters[k] += v

}

func GetCounter(k string) (int64, bool) {
	v, exist := defaultStorage.Counters[k]
	return v, exist
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
			MType: metrics.CounterType,
			Delta: &v,
		})
	}
	for name := range defaultStorage.Gauges {
		v := defaultStorage.Gauges[name]
		res = append(res, metrics.Metric{
			MType: metrics.GaugeType,
			Value: &v,
		})
	}
	return res
}
