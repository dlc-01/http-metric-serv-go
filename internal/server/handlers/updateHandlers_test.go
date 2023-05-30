package handlers

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandlerGauge(t *testing.T) {
	logging.InitLogger()
	s := storage.Init(context.Background(), &config.ServerConfig{})
	ServerStorage.Storage = s
	router := gin.Default()

	router.POST("/update/:types/:name/:value", ServerStorage.UpdateHandler)

	testsGauge := []struct {
		name               string
		url                string
		expectedCode       int
		nameValue          string
		expectedValueFloat float64
	}{
		{
			"/update/gauge/test/666.35",
			"/update/gauge/test/666.35",
			http.StatusOK,
			"test",
			666.35,
		},
		{
			"/update/gauge",
			"/update/gauge",
			http.StatusNotFound,
			"nil",
			0,
		},
	}

	for _, tt := range testsGauge {
		req, _ := http.NewRequest(http.MethodPost, tt.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, tt.expectedCode, w.Code)
		if w.Code == http.StatusOK {
			val, _ := s.GetMetric(context.TODO(), metrics.Metric{ID: tt.nameValue, MType: metrics.GaugeType})
			assert.Equal(t, tt.expectedValueFloat, *val.Value)

		}
	}

}

func TestUpdateHandlerCounter(t *testing.T) {
	s := storage.Init(context.Background(), &config.ServerConfig{})
	ServerStorage.Storage = s
	router := gin.Default()

	router.POST("/update/:types/:name/:value", ServerStorage.UpdateHandler)

	testsCounter := []struct {
		name             string
		url              string
		expectedCode     int
		nameValue        string
		expectedValueInt int64
	}{
		{
			"/update/counter/test/666",
			"/update/counter/test/666",
			http.StatusOK,
			"test",
			666,
		},
		{
			"/update/counter/",
			"/update/counter",
			http.StatusNotFound,
			"nil",
			0,
		},
	}

	for _, tt := range testsCounter {
		req, _ := http.NewRequest(http.MethodPost, tt.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.expectedCode, w.Code)
		if w.Code == http.StatusOK {
			val, _ := s.GetMetric(context.TODO(), metrics.Metric{ID: tt.nameValue, MType: metrics.CounterType})
			assert.Equal(t, tt.expectedValueInt, *val.Delta)
		}
	}

}
