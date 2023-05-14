package storagesync

import (
	"context"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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
	storage.Init()

	RunSync(context.Background(), &cfg)

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

			dumpFile()

			storage.Init()

			restoreFile()

			gauge, _ := storage.GetGauge(tt.metricGauge.ID)
			counter, _ := storage.GetCounter(tt.metricCounter.ID)

			assert.Equal(t, gauge, tt.metricGauge)
			assert.Equal(t, counter, tt.metricCounter)

			os.Remove(cfg.FileStoragePath)
		})
	}
}

func TestGetSyncMiddlewareFile(t *testing.T) {
	testValue1 := 2022.02
	testDelta1 := int64(24)
	testValue2 := 2003.03
	testDelta2 := int64(23)

	cfg := config.ServerConfig{FileStoragePath: "/tmp/test_save.json", StoreInterval: 1}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}
	os.Remove(cfg.FileStoragePath)

	storage.Init()

	RunSync(context.Background(), &cfg)

	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed), GetSyncMiddleware())
	router.POST("/update/", json.UpdateJSONHandler)

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
		jsonsCounter, err := tt.metricCounter.ToJSON()
		if err != nil {
			logging.Fatalf("cannot generate request body: %wCounter", err)
		}

		reqCounter, err := http.NewRequest(http.MethodPost, "/update/", jsonsCounter)
		if err != nil {
			t.Fatal(err)
		}

		reqCounter.Header.Set("Content-Type", "application/json")
		reqCounter.Header.Set("Content-Encoding", "gzip")

		wCounter := httptest.NewRecorder()
		router.ServeHTTP(wCounter, reqCounter)
		jsonsGauge, err := tt.metricGauge.ToJSON()
		if err != nil {
			logging.Fatalf("cannot generate request body: %wCounter", err)
		}

		reqGauge, err := http.NewRequest(http.MethodPost, "/update/", jsonsGauge)
		if err != nil {
			t.Fatal(err)
		}

		reqGauge.Header.Set("Content-Type", "application/json")
		reqGauge.Header.Set("Content-Encoding", "gzip")

		wGauge := httptest.NewRecorder()
		router.ServeHTTP(wGauge, reqGauge)

		time.Sleep(1 * time.Second)

		t.Run(tt.name, func(t *testing.T) {
			new := storage.GetStorage()
			fmt.Println(new)
			storage.Init()
			restoreFile()
			new = storage.GetStorage()
			fmt.Println(new)
			gauge, _ := storage.GetGauge(tt.metricGauge.ID)
			counter, _ := storage.GetCounter(tt.metricCounter.ID)

			assert.Equal(t, gauge, tt.metricGauge)
			assert.Equal(t, counter, tt.metricCounter)

			os.Remove(cfg.FileStoragePath)
		})
	}

}
