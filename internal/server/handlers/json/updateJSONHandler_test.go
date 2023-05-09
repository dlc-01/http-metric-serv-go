package json

import (
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateJSONHandler(t *testing.T) {
	logging.InitLogger()
	router := gin.Default()
	router.POST("/update/", UpdateJSONHandler)
	storage.Init()

	testGauge := `{"id":"TestGauge", "type":"gauge", "value":2022.02}`
	testCounter := `{"id":"TestCounter", "type":"counter", "delta":24}`
	testWrongMetric := `{"id":"TestWrongMetric", "type":"qwert", "delta":24}`
	testWrongValue := `{"id":"TestWrongValue", "type":"counter", "delta":2022.02}`
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
			body:         testGauge,
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
			body:         testCounter,
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
			body:         testWrongMetric,
		},
		{
			name:         `bad metric value`,
			expectedCode: http.StatusBadRequest,
			url:          `/update/`,
			body:         testWrongValue,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, tt.expectedCode, w.Code)
		if tt.expectedCode == 200 {
			var data metrics.Metric
			json.Unmarshal(w.Body.Bytes(), &data)
			assert.Equal(t, tt.expectedBody, data)
		}
	}

}
func TestUpdateJSONHandlerWithGzip(t *testing.T) {
	logging.InitLogger()
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.POST("/update/", UpdateJSONHandler)
	storage.Init()

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
		jsons, err := tt.expectedBody.ToJson()
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
			var data metrics.Metric
			json.Unmarshal(w.Body.Bytes(), &data)
			assert.Equal(t, tt.expectedBody, data)
		}
	}

}
