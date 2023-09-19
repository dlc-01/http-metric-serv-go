package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

func BenchmarkSetMetric(b *testing.B) {
	conf := config.ServerConfig{}
	Init(context.Background(), &conf)
	b.Run("setMetrics benchmarks", func(b *testing.B) {
		v := 50.1001
		metric := metrics.Metric{
			ID:    "new",
			MType: "gauge",
			Value: &v,
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := SetMetric(context.Background(), metric)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkSetMetricsBatch(b *testing.B) {
	conf := config.ServerConfig{}
	Init(context.Background(), &conf)
	b.Run("setMetrics benchmarks", func(b *testing.B) {
		vF := 50.1001
		var vI int64 = 50
		metric := []metrics.Metric{
			{
				ID:    "new",
				MType: "counter",
				Delta: &vI,
			},
			{ID: "newG",
				MType: "gauge",
				Value: &vF},
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			err := SetMetricsBatch(context.Background(), metric)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkGetMetric(b *testing.B) {
	conf := config.ServerConfig{}
	Init(context.Background(), &conf)
	b.Run("setMetrics benchmarks", func(b *testing.B) {
		vF := 50.1001
		var vI int64 = 50
		metrics := []metrics.Metric{
			{
				ID:    "newC",
				MType: "counter",
				Delta: &vI,
			},
			{ID: "newG",
				MType: "gauge",
				Value: &vF,
			},
		}

		err := SetMetricsBatch(context.Background(), metrics)
		assert.NoError(b, err)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_, err := GetAll(context.Background())
			assert.NoError(b, err)

		}
	})
}
