package collector

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"math/rand"
	"runtime"
)

func CollectMetrics() {
	var Runtime runtime.MemStats
	runtime.ReadMemStats(&Runtime)
	storage.SetGauge("Alloc", float64(Runtime.Alloc))
	storage.SetGauge("BuckHashSys", float64(Runtime.BuckHashSys))
	storage.SetGauge("Frees", float64(Runtime.Frees))
	storage.SetGauge("GCCPUFraction", float64(Runtime.GCCPUFraction))
	storage.SetGauge("GCSys", float64(Runtime.GCSys))
	storage.SetGauge("HeapAlloc", float64(Runtime.HeapAlloc))
	storage.SetGauge("HeapIdle", float64(Runtime.HeapIdle))
	storage.SetGauge("HeapInuse", float64(Runtime.HeapInuse))
	storage.SetGauge("HeapObjects", float64(Runtime.HeapObjects))
	storage.SetGauge("HeapReleased", float64(Runtime.HeapReleased))
	storage.SetGauge("HeapSys", float64(Runtime.HeapSys))
	storage.SetGauge("LastGC", float64(Runtime.LastGC))
	storage.SetGauge("Lookups", float64(Runtime.Lookups))
	storage.SetGauge("MCacheInuse", float64(Runtime.MCacheInuse))
	storage.SetGauge("MCacheSys", float64(Runtime.MCacheSys))
	storage.SetGauge("MSpanInuse", float64(Runtime.MSpanInuse))
	storage.SetGauge("MSpanSys", float64(Runtime.MSpanSys))
	storage.SetGauge("Mallocs", float64(Runtime.Mallocs))
	storage.SetGauge("NextGC", float64(Runtime.NextGC))
	storage.SetGauge("NumForcedGC", float64(Runtime.NumForcedGC))
	storage.SetGauge("NumGC", float64(Runtime.NumGC))
	storage.SetGauge("OtherSys", float64(Runtime.OtherSys))
	storage.SetGauge("PauseTotalNs", float64(Runtime.PauseTotalNs))
	storage.SetGauge("StackInuse", float64(Runtime.StackInuse))
	storage.SetGauge("StackSys", float64(Runtime.StackSys))
	storage.SetGauge("Sys", float64(Runtime.Sys))
	storage.SetGauge("TotalAlloc", float64(Runtime.TotalAlloc))
	storage.SetGauge("RandomValue", rand.Float64())
	storage.SetCounter("PollCount", 1)
}
