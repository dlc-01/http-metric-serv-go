package json

import (
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValueJSONHandler(t *testing.T) {
	logging.InitLogger()
	router := gin.Default()
	router.POST("/value/", ValueJSONHandler)
	router.POST("/update/", UpdateJSONHandler)
	storage.Init()

	testGaugePost := `{"id":"TestGauge", "type":"gauge", "value":2022.02}`
	testCounterPost := `{"id":"TestCounter", "type":"counter", "delta":24}`
	testWrongPost := `{"id":"TestWrong", "type":"qwert", "delta":24}`
	testGaugeGet := `{"id":"TestGauge", "type":"gauge"}`
	testCounterGet := `{"id":"TestCounter", "type":"counter"}`
	testWrongMetricGet := `{"id":"TestWrong", "type":"qwert"}`
	testValue := 2022.02
	testDelta := int64(24)

	tests := []struct {
		name         string
		url          string
		postRequest  string
		expectedCode int
		getRequest   string
		expectedBody metrics.Metric
	}{
		{
			name:         `true gauge post`,
			expectedCode: http.StatusOK,
			postRequest:  testGaugePost,
			url:          `/value/`,
			getRequest:   testGaugeGet,
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
			url:          `/value/`,
			postRequest:  testCounterPost,
			getRequest:   testCounterGet,
			expectedBody: metrics.Metric{
				ID:    "TestCounter",
				MType: metrics.CounterType,
				Delta: &testDelta,
				Value: nil,
			},
		},
		{
			name:         `wrong metric`,
			expectedCode: http.StatusNotFound,
			url:          `/value/`,
			getRequest:   testWrongMetricGet,
			postRequest:  testWrongPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, err := http.NewRequest(http.MethodPost, "/update/", strings.NewReader(tt.postRequest))
			if err != nil {
				t.Fatal(err)
			}
			wPost := httptest.NewRecorder()
			router.ServeHTTP(wPost, set)

			get, err := http.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.getRequest))
			if err != nil {
				t.Fatal(err)
			}
			get.Header.Set("Content-Type", "application/json")
			wGet := httptest.NewRecorder()
			router.ServeHTTP(wGet, get)

			assert.Equal(t, tt.expectedCode, wGet.Code)
			if tt.expectedCode == 200 {
				var data metrics.Metric
				json.Unmarshal(wGet.Body.Bytes(), &data)
				assert.Equal(t, tt.expectedBody, data)
			}
		})
	}
}
