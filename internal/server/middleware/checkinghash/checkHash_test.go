package checkinghash

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/jsonbatch"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func TestMiddleware(t *testing.T) {
	key := "secret_key"
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.Use(CheckHash(key))
	router.POST("/updates/", jsonbatch.UpdatesButchJSONHandler)

	testValue := 2022.02
	testValueOther := 2022.01
	testDelta := int64(24)
	testDeltaSum := int64(48)

	tests := []struct {
		name         string
		url          string
		encodeKey    string
		expectedCode int
		responseBody []metrics.Metric
		expectedBody []metrics.Metric
	}{
		{
			name:         `true key`,
			expectedCode: http.StatusOK,
			url:          `/updates/`,
			encodeKey:    key,
			responseBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
			expectedBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
		},
		{
			name:         `false key`,
			encodeKey:    "key",
			expectedCode: http.StatusBadRequest,
			url:          `/updates/`,
			responseBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValueOther,
				},
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
			expectedBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &(testDeltaSum),
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage.Init(context.Background(), &config.ServerConfig{})

			jsons, err := metrics.ToJSON(tt.responseBody)
			if err != nil {
				logging.Fatalf("cannot generate request body: %s", err)
			}
			hash := hashing.HashingData(tt.encodeKey, jsons)
			gzip, err := metrics.Gzipper(jsons)
			if err != nil {
				logging.Fatalf("cannot gzip body: %s", err)
			}
			req, err := http.NewRequest(http.MethodPost, tt.url, gzip)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("HashSHA256", hash)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == 200 {
				data, _ := storage.GetAllMetrics(context.Background())
				assert.Equal(t, tt.expectedBody, data)
			}
		})
	}
}
