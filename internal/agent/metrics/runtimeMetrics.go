package metrics

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"math/rand"
	"runtime"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func Init() MemStorage {
	return MemStorage{make(map[string]float64), make(map[string]int64)}
}

func (metrics *MemStorage) Check() {
	var Runtime runtime.MemStats
	runtime.ReadMemStats(&Runtime)
	metrics.Gauge["Alloc"] = float64(Runtime.Alloc)
	metrics.Gauge["BuckHashSys"] = float64(Runtime.BuckHashSys)
	metrics.Gauge["Frees"] = float64(Runtime.Frees)
	metrics.Gauge["GCCPUFraction"] = float64(Runtime.GCCPUFraction)
	metrics.Gauge["GCSys"] = float64(Runtime.GCSys)
	metrics.Gauge["HeapAlloc"] = float64(Runtime.HeapAlloc)
	metrics.Gauge["HeapIdle"] = float64(Runtime.HeapIdle)
	metrics.Gauge["HeapInuse"] = float64(Runtime.HeapInuse)
	metrics.Gauge["HeapObjects"] = float64(Runtime.HeapObjects)
	metrics.Gauge["HeapReleased"] = float64(Runtime.HeapReleased)
	metrics.Gauge["HeapSys"] = float64(Runtime.HeapSys)
	metrics.Gauge["LastGC"] = float64(Runtime.LastGC)
	metrics.Gauge["Lookups"] = float64(Runtime.Lookups)
	metrics.Gauge["MCacheInuse"] = float64(Runtime.MCacheInuse)
	metrics.Gauge["MCacheSys"] = float64(Runtime.MCacheSys)
	metrics.Gauge["MSpanInuse"] = float64(Runtime.MSpanInuse)
	metrics.Gauge["MSpanSys"] = float64(Runtime.MSpanSys)
	metrics.Gauge["Mallocs"] = float64(Runtime.Mallocs)
	metrics.Gauge["NextGC"] = float64(Runtime.NextGC)
	metrics.Gauge["NumForcedGC"] = float64(Runtime.NumForcedGC)
	metrics.Gauge["NumGC"] = float64(Runtime.NumGC)
	metrics.Gauge["OtherSys"] = float64(Runtime.OtherSys)
	metrics.Gauge["PauseTotalNs"] = float64(Runtime.PauseTotalNs)
	metrics.Gauge["StackInuse"] = float64(Runtime.StackInuse)
	metrics.Gauge["StackSys"] = float64(Runtime.StackSys)
	metrics.Gauge["Sys"] = float64(Runtime.Sys)
	metrics.Gauge["TotalAlloc"] = float64(Runtime.TotalAlloc)
	metrics.Gauge["RandomValue"] = rand.Float64()
	metrics.Counter["PollCount"]++
}

func (metrics *MemStorage) GenerateURLMetrics(host string) []string {
	var urls []string

	for metric, value := range metrics.Gauge {
		generatedURL := fmt.Sprintf("http://%s/update/gauge/%s/%f", host, metric, value)
		urls = append(urls, generatedURL)
	}

	for metric, value := range metrics.Counter {
		generatedURL := fmt.Sprintf("http://%s/update/counter/%s/%d", host, metric, value)
		urls = append(urls, generatedURL)
	}
	return urls
}

func (metrics *MemStorage) GenerateStructMetrics() []interface{} {
	var requests []interface{}

	for metric, value := range metrics.Gauge {
		request := storage.Metrics{
			ID:    metric,
			MType: "Gauge",
			Value: &value,
		}
		requests = append(requests, request)
	}

	for metric, value := range metrics.Counter {
		request := storage.Metrics{
			ID:    metric,
			MType: "Counter",
			Delta: &value,
		}
		requests = append(requests, request)
	}
	return requests
}
