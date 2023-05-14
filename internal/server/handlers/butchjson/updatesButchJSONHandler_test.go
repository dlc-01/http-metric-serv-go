package butchjson

import (
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
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.POST("/updates/", UpdatesButchJSONHandler)
	storage.Init()

	testValue := 2022.02
	testDelta := int64(24)

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody []metrics.Metric
	}{
		{
			name:         `true gauge post`,
			expectedCode: http.StatusOK,
			url:          `/updates/`,
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
			name:         `wrong metric`,
			expectedCode: http.StatusNotImplemented,
			url:          `/updates/`,
			expectedBody: []metrics.Metric{
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
		jsons, err := metrics.ToJSONMetrics(tt.expectedBody)
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
			data := storage.GetMetrics()
			assert.Equal(t, tt.expectedBody, data)
		}
	}
}
