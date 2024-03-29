package storagesync

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
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
	storage.Init(context.Background(), &config.ServerConfig{})

	RunSync(&cfg)

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

			storage.SetMetric(context.Background(), tt.metricGauge)
			storage.SetMetric(context.Background(), tt.metricCounter)

			dumpFile()
			storage.Init(context.Background(), &config.ServerConfig{})

			err := restoreFile()
			if err != nil {
				logging.Fatalf("cannot restore %s", err)
			}

			gauge, err := storage.GetMetric(context.Background(), tt.metricGauge)
			if err != nil {
				logging.Errorf("cannot get gauge metric: %s", err)
			}
			counter, err := storage.GetMetric(context.Background(), tt.metricCounter)
			if err != nil {
				logging.Errorf("cannot get counter metric: %s", err)
			}

			assert.Equal(t, gauge, tt.metricGauge)
			assert.Equal(t, counter, tt.metricCounter)

			os.Remove(cfg.FileStoragePath)
		})
	}
}
