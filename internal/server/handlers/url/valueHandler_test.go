package url

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func TestValueHandlerGauge(t *testing.T) {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})
	router := gin.Default()
	router.POST("/update/:types/:name/:value", UpdateHandler)
	router.GET("/value/:types/:name", ValueHandler)

	testsGauge := []struct {
		name               string
		urlPost            string
		urlGet             string
		expectedCode       int
		expectedValueFloat string
	}{
		{
			"/update/gauge/test/666.35",
			"/update/gauge/test/666.35",
			"/value/gauge/test",
			http.StatusOK,
			"666.35",
		},
		{
			"/value/gauge/tests/",
			"",
			"/value/gauge/tests",
			http.StatusNotFound,
			"",
		},
	}

	for _, tt := range testsGauge {
		reqPost, _ := http.NewRequest(http.MethodPost, tt.urlPost, nil)
		wPost := httptest.NewRecorder()
		router.ServeHTTP(wPost, reqPost)
		reqGet, _ := http.NewRequest(http.MethodGet, tt.urlGet, nil)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)
		assert.Equal(t, tt.expectedCode, wGet.Code)
		if wGet.Code == http.StatusOK {

			assert.Equal(t, tt.expectedValueFloat, wGet.Body.String())

		}
	}

}
func TestValueHandlerCounter(t *testing.T) {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})
	router := gin.Default()

	router.POST("/update/:types/:name/:value", UpdateHandler)
	router.GET("/value/:types/:name", ValueHandler)

	testsCounter := []struct {
		name             string
		urlPost          string
		urlGet           string
		expectedCode     int
		expectedValueInt string
	}{
		{
			"/update/counter/test/666",
			"/update/counter/test/666",
			"/value/counter/test",
			http.StatusOK,
			"666",
		},
		{
			"/value/counter/tests",
			"",
			"/value/counter/tests",
			http.StatusNotFound,
			"",
		},
	}

	for _, tt := range testsCounter {
		reqPost, _ := http.NewRequest(http.MethodPost, tt.urlPost, nil)
		wPost := httptest.NewRecorder()
		router.ServeHTTP(wPost, reqPost)

		reqGet, _ := http.NewRequest(http.MethodGet, tt.urlGet, nil)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)

		assert.Equal(t, tt.expectedCode, wGet.Code)
		if wGet.Code == http.StatusOK {
			assert.Equal(t, tt.expectedValueInt, wGet.Body.String())
		}
	}

}
