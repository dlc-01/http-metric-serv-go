package json

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func TestUpdateJSONHandler(t *testing.T) {
	logging.InitLogger()
	router := gin.Default()
	storage.Init(context.Background(), &config.ServerConfig{})
	router.POST("/update/", UpdateJSONHandler)

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

			data, _ := storage.GetMetric(context.TODO(), tt.expectedBody)
			assert.Equal(t, tt.expectedBody, data)
		}
	}

}
