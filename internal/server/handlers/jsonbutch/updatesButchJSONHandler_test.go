package jsonbutch

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdatesButchJSONHandler(t *testing.T) {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.POST("/updates/", UpdatesButchJSONHandler)

	testValue := 2022.02
	testValueOther := 2022.01
	testDelta := int64(24)
	testDeltaSum := int64(48)

	tests := []struct {
		name         string
		url          string
		expectedCode int
		responseBody []metrics.Metric
		expectedBody []metrics.Metric
	}{
		{
			name:         `true gauge post`,
			expectedCode: http.StatusOK,
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
			name:         `true gauge post`,
			expectedCode: http.StatusOK,
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
		{
			name:         `wrong metric`,
			expectedCode: http.StatusNotImplemented,
			url:          `/updates/`,
			responseBody: []metrics.Metric{
				{
					ID:    "TestWrongMetric",
					MType: "qwert",
					Delta: &testDelta,
					Value: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage.Init(context.Background(), &config.ServerConfig{})
			jsons, err := metrics.ToJSONs(tt.responseBody)
			if err != nil {
				logging.Fatalf("cannot generate request body: %w", err)
			}
			gzip, err := metrics.Gzipper(jsons)
			if err != nil {
				logging.Fatalf("cannot gzip body: %w", err)
			}

			req, err := http.NewRequest(http.MethodPost, tt.url, gzip)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")

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
