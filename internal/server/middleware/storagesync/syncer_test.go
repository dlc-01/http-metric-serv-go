package storagesync

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestDumpRestore(t *testing.T) {
	testValue1 := 2022.02
	testDelta1 := int64(24)
	testValue2 := 2003.03
	testDelta2 := int64(23)

	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}
	cfg.FileStoragePath = "/tmp/test_save.json"
	storage.Init()
	os.Remove(cfg.FileStoragePath)
	RunSync(cfg)

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
			storage.SetGauge(tt.metricGauge.ID, *tt.metricGauge.Value)
			storage.SetCounter(tt.metricCounter.ID, *tt.metricCounter.Delta)

			dump()

			storage.Init()

			restore("/tmp/test_save.json")

			gauge, _ := storage.GetGauge(tt.metricGauge.ID)
			counter, _ := storage.GetCounter(tt.metricCounter.ID)

			assert.Equal(t, gauge, tt.metricGauge)
			assert.Equal(t, counter, tt.metricCounter)

			os.Remove(cfg.FileStoragePath)
		})
	}
}
