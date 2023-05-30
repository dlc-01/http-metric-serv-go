package storagesync

import (
	"context"
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

	storage.Init(context.Background(), &config.ServerConfig{})

	RunSync(&cfg)

	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed), GetSyncMiddleware())
	router.POST("/update/", json.UpdateJSONHandler)
	router.POST("/value/", json.ValueJSONHandler)

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
		jsonsCounter, err := tt.metricCounter.ToJSONWithGzip()
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
		jsonsGauge, err := tt.metricGauge.ToJSONWithGzip()
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
		time.Sleep(time.Second)
		t.Run(tt.name, func(t *testing.T) {
			dumpFile()

			storage.Init(context.TODO(), conf)
			RunSync(conf)

			gauge, err := storage.GetMetric(context.Background(), tt.metricGauge)
			if err != nil {
				logging.Errorf("cannot get gauge metric: %s", err)
			}
			counter, err := storage.GetMetric(context.Background(), tt.metricCounter)
			if err != nil {
				logging.Errorf("cannot get counter metric: %s", err)
			}

			assert.Equal(t, *gauge.Value, *tt.metricGauge.Value)
			assert.Equal(t, *counter.Delta, *tt.metricCounter.Delta)

			os.Remove(cfg.FileStoragePath)
		})
	}

}
