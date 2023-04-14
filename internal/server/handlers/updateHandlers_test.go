package handlers

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler1(t *testing.T) {
	router := gin.Default()
	router.POST("/update/:types/:name/:value", UpdateHandler)
	storage.Init()

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
			"/update/gauge/",
			"/update/gauge/",
			http.StatusTemporaryRedirect,
			"nil",
			0,
		},
	}

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
			"/update/counter/",
			http.StatusTemporaryRedirect,
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
			val, _ := storage.GetGauge(tt.nameValue)
			assert.Equal(t, tt.expectedValueFloat, val)

		}
	}
	for _, tt := range testsCounter {
		req, _ := http.NewRequest(http.MethodPost, tt.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, tt.expectedCode, w.Code)
		if w.Code == http.StatusOK {
			val, _ := storage.GetCounter(tt.nameValue)
			assert.Equal(t, tt.expectedValueInt, val)
		}
	}
}
