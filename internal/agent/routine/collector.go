package routine

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"math/rand"
	"runtime"
)

func (s stor) CollectMetrics(ctx context.Context) {
	var Runtime runtime.MemStats
	runtime.ReadMemStats(&Runtime)
	alloc := float64(Runtime.Alloc)
	s.SetMetric(ctx, metrics.Metric{ID: "Alloc", MType: metrics.GaugeType, Value: &alloc})
	buckHashSys := float64(Runtime.BuckHashSys)
	s.SetMetric(ctx, metrics.Metric{ID: "BuckHashSys", MType: metrics.GaugeType, Value: &buckHashSys})
	frees := float64(Runtime.Frees)
	s.SetMetric(ctx, metrics.Metric{ID: "Frees", MType: metrics.GaugeType, Value: &frees})
	gpuFraction := float64(Runtime.GCCPUFraction)
	s.SetMetric(ctx, metrics.Metric{ID: "GCCPUFraction", MType: metrics.GaugeType, Value: &gpuFraction})
	gcSys := float64(Runtime.GCSys)
	s.SetMetric(ctx, metrics.Metric{ID: "GCSys", MType: metrics.GaugeType, Value: &gcSys})
	heapAlloc := float64(Runtime.HeapAlloc)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapAlloc", MType: metrics.GaugeType, Value: &heapAlloc})
	heapIdle := float64(Runtime.HeapIdle)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapIdle", MType: metrics.GaugeType, Value: &heapIdle})
	heapInuse := float64(Runtime.HeapInuse)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapInuse", MType: metrics.GaugeType, Value: &heapInuse})
	heapObjects := float64(Runtime.HeapObjects)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapObjects", MType: metrics.GaugeType, Value: &heapObjects})
	heapReleased := float64(Runtime.HeapReleased)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapReleased", MType: metrics.GaugeType, Value: &heapReleased})
	heapSys := float64(Runtime.HeapSys)
	s.SetMetric(ctx, metrics.Metric{ID: "HeapSys", MType: metrics.GaugeType, Value: &heapSys})
	lastGC := float64(Runtime.LastGC)
	s.SetMetric(ctx, metrics.Metric{ID: "LastGC", MType: metrics.GaugeType, Value: &lastGC})
	lookups := float64(Runtime.Lookups)
	s.SetMetric(ctx, metrics.Metric{ID: "Lookups", MType: metrics.GaugeType, Value: &lookups})
	mCacheInuse := float64(Runtime.MCacheInuse)
	s.SetMetric(ctx, metrics.Metric{ID: "MCacheInuse", MType: metrics.GaugeType, Value: &mCacheInuse})
	mSpanInuse := float64(Runtime.MSpanInuse)
	s.SetMetric(ctx, metrics.Metric{ID: "MSpanInuse", MType: metrics.GaugeType, Value: &mSpanInuse})
	mSpanSys := float64(Runtime.MSpanSys)
	s.SetMetric(ctx, metrics.Metric{ID: "MSpanSys", MType: metrics.GaugeType, Value: &mSpanSys})
	mCacheSys := float64(Runtime.MCacheSys)
	s.SetMetric(ctx, metrics.Metric{ID: "MCacheSys", MType: metrics.GaugeType, Value: &mCacheSys})
	mallocs := float64(Runtime.Mallocs)
	s.SetMetric(ctx, metrics.Metric{ID: "Mallocs", MType: metrics.GaugeType, Value: &mallocs})
	nextGC := float64(Runtime.NextGC)
	s.SetMetric(ctx, metrics.Metric{ID: "NextGC", MType: metrics.GaugeType, Value: &nextGC})
	numForcedGC := float64(Runtime.NumForcedGC)
	s.SetMetric(ctx, metrics.Metric{ID: "NumForcedGC", MType: metrics.GaugeType, Value: &numForcedGC})
	numGC := float64(Runtime.NumGC)
	s.SetMetric(ctx, metrics.Metric{ID: "NumGC", MType: metrics.GaugeType, Value: &numGC})
	otherSys := float64(Runtime.OtherSys)
	s.SetMetric(ctx, metrics.Metric{ID: "OtherSys", MType: metrics.GaugeType, Value: &otherSys})
	pauseTotalNs := float64(Runtime.PauseTotalNs)
	s.SetMetric(ctx, metrics.Metric{ID: "PauseTotalNs", MType: metrics.GaugeType, Value: &pauseTotalNs})
	stackInuse := float64(Runtime.StackInuse)
	s.SetMetric(ctx, metrics.Metric{ID: "StackInuse", MType: metrics.GaugeType, Value: &stackInuse})
	stackSys := float64(Runtime.StackSys)
	s.SetMetric(ctx, metrics.Metric{ID: "StackSys", MType: metrics.GaugeType, Value: &stackSys})
	sys := float64(Runtime.Sys)
	s.SetMetric(ctx, metrics.Metric{ID: "Sys", MType: metrics.GaugeType, Value: &sys})
	totalAlloc := float64(Runtime.TotalAlloc)
	s.SetMetric(ctx, metrics.Metric{ID: "TotalAlloc", MType: metrics.GaugeType, Value: &totalAlloc})
	randomValue := rand.Float64()
	s.SetMetric(ctx, metrics.Metric{ID: "RandomValue", MType: metrics.GaugeType, Value: &randomValue})
	totalMemory := rand.Float64()
	s.SetMetric(ctx, metrics.Metric{ID: "TotalMemory", MType: metrics.GaugeType, Value: &totalMemory})
	freeMemory := rand.Float64()
	s.SetMetric(ctx, metrics.Metric{ID: "FreeMemory", MType: metrics.GaugeType, Value: &freeMemory})
	cPUutilization1 := rand.Float64()
	s.SetMetric(ctx, metrics.Metric{ID: "CPUutilization1", MType: metrics.GaugeType, Value: &cPUutilization1})
	var count int64 = 1
	s.SetMetric(ctx, metrics.Metric{ID: "PollCount", MType: metrics.CounterType, Delta: &count})

}
