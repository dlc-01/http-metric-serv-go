package collector

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"runtime"
	"time"
)

func CollectMetrics(ctx context.Context, result chan<- []metrics.Metric, duration time.Duration) {
	for {
		time.Sleep(duration)
		v, _ := mem.VirtualMemory()
		countCPU, _ := cpu.Counts(true)
		var Runtime runtime.MemStats
		runtime.ReadMemStats(&Runtime)
		alloc := float64(Runtime.Alloc)
		storage.SetMetric(ctx, metrics.Metric{ID: "Alloc", MType: metrics.GaugeType, Value: &alloc})
		buckHashSys := float64(Runtime.BuckHashSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "BuckHashSys", MType: metrics.GaugeType, Value: &buckHashSys})
		frees := float64(Runtime.Frees)
		storage.SetMetric(ctx, metrics.Metric{ID: "Frees", MType: metrics.GaugeType, Value: &frees})
		gpuFraction := float64(Runtime.GCCPUFraction)
		storage.SetMetric(ctx, metrics.Metric{ID: "GCCPUFraction", MType: metrics.GaugeType, Value: &gpuFraction})
		gcSys := float64(Runtime.GCSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "GCSys", MType: metrics.GaugeType, Value: &gcSys})
		heapAlloc := float64(Runtime.HeapAlloc)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapAlloc", MType: metrics.GaugeType, Value: &heapAlloc})
		heapIdle := float64(Runtime.HeapIdle)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapIdle", MType: metrics.GaugeType, Value: &heapIdle})
		heapInuse := float64(Runtime.HeapInuse)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapInuse", MType: metrics.GaugeType, Value: &heapInuse})
		heapObjects := float64(Runtime.HeapObjects)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapObjects", MType: metrics.GaugeType, Value: &heapObjects})
		heapReleased := float64(Runtime.HeapReleased)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapReleased", MType: metrics.GaugeType, Value: &heapReleased})
		heapSys := float64(Runtime.HeapSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "HeapSys", MType: metrics.GaugeType, Value: &heapSys})
		lastGC := float64(Runtime.LastGC)
		storage.SetMetric(ctx, metrics.Metric{ID: "LastGC", MType: metrics.GaugeType, Value: &lastGC})
		lookups := float64(Runtime.Lookups)
		storage.SetMetric(ctx, metrics.Metric{ID: "Lookups", MType: metrics.GaugeType, Value: &lookups})
		mCacheInuse := float64(Runtime.MCacheInuse)
		storage.SetMetric(ctx, metrics.Metric{ID: "MCacheInuse", MType: metrics.GaugeType, Value: &mCacheInuse})
		mSpanInuse := float64(Runtime.MSpanInuse)
		storage.SetMetric(ctx, metrics.Metric{ID: "MSpanInuse", MType: metrics.GaugeType, Value: &mSpanInuse})
		mSpanSys := float64(Runtime.MSpanSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "MSpanSys", MType: metrics.GaugeType, Value: &mSpanSys})
		mCacheSys := float64(Runtime.MCacheSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "MCacheSys", MType: metrics.GaugeType, Value: &mCacheSys})
		mallocs := float64(Runtime.Mallocs)
		storage.SetMetric(ctx, metrics.Metric{ID: "Mallocs", MType: metrics.GaugeType, Value: &mallocs})
		nextGC := float64(Runtime.NextGC)
		storage.SetMetric(ctx, metrics.Metric{ID: "NextGC", MType: metrics.GaugeType, Value: &nextGC})
		numForcedGC := float64(Runtime.NumForcedGC)
		storage.SetMetric(ctx, metrics.Metric{ID: "NumForcedGC", MType: metrics.GaugeType, Value: &numForcedGC})
		numGC := float64(Runtime.NumGC)
		storage.SetMetric(ctx, metrics.Metric{ID: "NumGC", MType: metrics.GaugeType, Value: &numGC})
		otherSys := float64(Runtime.OtherSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "OtherSys", MType: metrics.GaugeType, Value: &otherSys})
		pauseTotalNs := float64(Runtime.PauseTotalNs)
		storage.SetMetric(ctx, metrics.Metric{ID: "PauseTotalNs", MType: metrics.GaugeType, Value: &pauseTotalNs})
		stackInuse := float64(Runtime.StackInuse)
		storage.SetMetric(ctx, metrics.Metric{ID: "StackInuse", MType: metrics.GaugeType, Value: &stackInuse})
		stackSys := float64(Runtime.StackSys)
		storage.SetMetric(ctx, metrics.Metric{ID: "StackSys", MType: metrics.GaugeType, Value: &stackSys})
		sys := float64(Runtime.Sys)
		storage.SetMetric(ctx, metrics.Metric{ID: "Sys", MType: metrics.GaugeType, Value: &sys})
		totalAlloc := float64(Runtime.TotalAlloc)
		storage.SetMetric(ctx, metrics.Metric{ID: "TotalAlloc", MType: metrics.GaugeType, Value: &totalAlloc})
		randomValue := rand.Float64()
		storage.SetMetric(ctx, metrics.Metric{ID: "RandomValue", MType: metrics.GaugeType, Value: &randomValue})
		totalMemory := float64(v.Total)
		storage.SetMetric(ctx, metrics.Metric{ID: "TotalMemory", MType: metrics.GaugeType, Value: &totalMemory})
		freeMemory := float64(v.Free)
		storage.SetMetric(ctx, metrics.Metric{ID: "FreeMemory", MType: metrics.GaugeType, Value: &freeMemory})
		cPUutilization1 := float64(countCPU)
		storage.SetMetric(ctx, metrics.Metric{ID: "CPUutilization1", MType: metrics.GaugeType, Value: &cPUutilization1})
		var count int64 = 1
		storage.SetMetric(ctx, metrics.Metric{ID: "PollCount", MType: metrics.CounterType, Delta: &count})
		metric, err := storage.GetAllMetrics(context.Background())
		if err != nil {
			logging.Errorf("cannot get metrics :%w", err)
		}
		result <- metric
	}
}
