package metrics

import (
	"fmt"
	"math/rand"
	"runtime"
)

type MemMetrics struct {
	gauge   map[string]float64
	counter map[string]int64
}

var gaugeMetrics = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue"}

func (metrics *MemMetrics) Init() {
	metrics.gauge = make(map[string]float64)
	metrics.counter = make(map[string]int64)

}
func (metrics *MemMetrics) Check() {
	var Runtime runtime.MemStats
	runtime.ReadMemStats(&Runtime)
	metrics.gauge["Alloc"] = float64(Runtime.Alloc)
	metrics.gauge["BuckHashSys"] = float64(Runtime.BuckHashSys)
	metrics.gauge["Frees"] = float64(Runtime.Frees)
	metrics.gauge["GCCPUFraction"] = float64(Runtime.GCCPUFraction)
	metrics.gauge["GCSys"] = float64(Runtime.GCSys)
	metrics.gauge["HeapAlloc"] = float64(Runtime.HeapAlloc)
	metrics.gauge["HeapIdle"] = float64(Runtime.HeapIdle)
	metrics.gauge["HeapInuse"] = float64(Runtime.HeapInuse)
	metrics.gauge["HeapObjects"] = float64(Runtime.HeapObjects)
	metrics.gauge["HeapReleased"] = float64(Runtime.HeapReleased)
	metrics.gauge["HeapSys"] = float64(Runtime.HeapSys)
	metrics.gauge["LastGC"] = float64(Runtime.LastGC)
	metrics.gauge["Lookups"] = float64(Runtime.Lookups)
	metrics.gauge["MCacheInuse"] = float64(Runtime.MCacheInuse)
	metrics.gauge["MCacheSys"] = float64(Runtime.MCacheSys)
	metrics.gauge["MSpanInuse"] = float64(Runtime.MSpanInuse)
	metrics.gauge["MSpanSys"] = float64(Runtime.MSpanSys)
	metrics.gauge["Mallocs"] = float64(Runtime.Mallocs)
	metrics.gauge["NextGC"] = float64(Runtime.NextGC)
	metrics.gauge["NumForcedGC"] = float64(Runtime.NumForcedGC)
	metrics.gauge["NumGC"] = float64(Runtime.NumGC)
	metrics.gauge["OtherSys"] = float64(Runtime.OtherSys)
	metrics.gauge["PauseTotalNs"] = float64(Runtime.PauseTotalNs)
	metrics.gauge["StackInuse"] = float64(Runtime.StackInuse)
	metrics.gauge["StackSys"] = float64(Runtime.StackSys)
	metrics.gauge["Sys"] = float64(Runtime.Sys)
	metrics.gauge["TotalAlloc"] = float64(Runtime.TotalAlloc)
	metrics.gauge["RandomValue"] = rand.Float64()
	metrics.counter["PollCount"]++
}
func (metrics *MemMetrics) GenerateUrlMetrics(host string) []string {
	var urls []string
	for metric, value := range metrics.gauge {
		generatedUrl := fmt.Sprintf("%s/update/gauge/%s/%f", host, metric, value)
		urls = append(urls, generatedUrl)
	}
	for metric, value := range metrics.counter {
		generatedUrl := fmt.Sprintf("%s/update/counter/%s/%d", host, metric, value)
		urls = append(urls, generatedUrl)
	}
	return urls
}
