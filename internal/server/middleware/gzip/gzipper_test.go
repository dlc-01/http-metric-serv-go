package gzip

import (
	"compress/gzip"
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipWithUpdateJSONHandler(t *testing.T) {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()
	router.Use(Gzip(gzip.BestSpeed))
	router.POST("/update/", json.UpdateJSONHandler)
	storage.Init(context.Background(), &config.ServerConfig{})

	testValue := 2022.02
	testDelta := int64(24)

	tests := []struct {
		name         string
		url          string
		expectedCode int
		body         string
		expectedBody metrics.Metric
	}{
		{
			name:         `true gauge post`,
			expectedCode: http.StatusOK,
			url:          `/update/`,
			expectedBody: metrics.Metric{
				ID:    "TestGauge",
				MType: metrics.GaugeType,
				Delta: nil,
				Value: &testValue,
			},
		},
		{
			name:         `true counter post`,
			expectedCode: http.StatusOK,
			url:          `/update/`,
			expectedBody: metrics.Metric{
				ID:    "TestCounter",
				MType: metrics.CounterType,
				Delta: &testDelta,
				Value: nil,
			},
		},
		{
			name:         `wrong metric`,
			expectedCode: http.StatusNotImplemented,
			url:          `/update/`,
		},
		{
			name:         `bad metric value`,
			expectedCode: http.StatusNotImplemented,
			url:          `/update/`,
		},
	}

	for _, tt := range tests {
		jsons, err := tt.expectedBody.ToJSONWithGzip()
		if err != nil {
			logging.Fatalf("cannot generate request body: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, tt.url, jsons)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.expectedCode, w.Code)

		if tt.expectedCode == 200 {
			switch tt.expectedBody.MType {
			case metrics.GaugeType:
				value, _ := storage.GetMetric(context.TODO(), tt.expectedBody)
				assert.Equal(t, testValue, *value.Value)
			case metrics.CounterType:
				delta, _ := storage.GetMetric(context.TODO(), tt.expectedBody)
				assert.Equal(t, testDelta, *delta.Delta)
			}
		}
	}

}
