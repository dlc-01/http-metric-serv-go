package storagesync

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestDumpRestoreFile(t *testing.T) {
	testValue1 := 2022.02
	testDelta1 := int64(24)
	testValue2 := 2003.03
	testDelta2 := int64(23)

	cfg := config.ServerConfig{FileStoragePath: "/tmp/test_save.json", StoreInterval: 0}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}
	os.Remove(cfg.FileStoragePath)
	s := storage.Init(context.Background(), &config.ServerConfig{})

	RunSync(&cfg, s)

	tests := []struct {
		name          string
		metricGauge   metrics.Metric
		metricCounter metrics.Metric
	}{
		{
			name:          "test1",
			metricCounter: metrics.Metric{ID: "firstTestCounter", MType: metrics.CounterType, Delta: &testDelta1},
			metricGauge:   metrics.Metric{ID: "firstTestGauge", MType: metrics.GaugeType, Value: &testValue1},
		},
		{
			name:          "test2",
			metricCounter: metrics.Metric{ID: "secondTestCounter", MType: metrics.CounterType, Delta: &testDelta2},
			metricGauge:   metrics.Metric{ID: "secondTestGauge", MType: metrics.GaugeType, Value: &testValue2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s.SetMetric(context.Background(), tt.metricGauge)
			s.SetMetric(context.Background(), tt.metricCounter)

			dumpFile()
			new := storage.Init(context.Background(), &config.ServerConfig{})

			restoreFile()

			gauge, err := new.GetMetric(context.Background(), tt.metricGauge)
			if err != nil {
				logging.Errorf("cannot get gauge metric: %s", err)
			}
			counter, err := new.GetMetric(context.Background(), tt.metricCounter)
			if err != nil {
				logging.Errorf("cannot get counter metric: %s", err)
			}

			assert.Equal(t, gauge, tt.metricGauge)
			assert.Equal(t, counter, tt.metricCounter)

			os.Remove(cfg.FileStoragePath)
		})
	}
}
