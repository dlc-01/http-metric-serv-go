package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/all"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/jsonbatch"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func ExampleUpdateJSONHandler() {
	logging.InitLogger()
	router := gin.Default()
	storage.Init(context.Background(), &config.ServerConfig{})
	router.POST("/update/", json.UpdateJSONHandler)

	testCounter := `{"id":"TestCounter", "type":"counter", "delta":24}`

	req, err := http.NewRequest(http.MethodPost, `/update/`, strings.NewReader(testCounter))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)
}

func ExampleValueJSONHandler() {
	logging.InitLogger()
	router := gin.Default()
	storage.Init(context.Background(), &config.ServerConfig{})

	router.POST("/value/", json.ValueJSONHandler)

	testCounterGet := `{"id":"TestCounter", "type":"counter"}`

	get, err := http.NewRequest(http.MethodPost, `/value/`, strings.NewReader(testCounterGet))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(get)
}

func ExampleShowMetrics() {
	logging.InitLogger()
	router := gin.Default()
	storage.Init(context.Background(), &config.ServerConfig{})

	router.GET("/", all.ShowMetrics)

	get, err := http.NewRequest(http.MethodGet, `/`, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(get)
}

func ExampleUpdatesButchJSONHandler() {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.POST("/updates/", jsonbatch.UpdatesButchJSONHandler)

	testValue := 2022.02

	testDelta := int64(24)

	metricS := []metrics.Metric{{
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
		}}

	storage.Init(context.Background(), &config.ServerConfig{})
	jsons, err := metrics.ToJSON(metricS)
	if err != nil {
		logging.Fatalf("cannot generate request body: %w", err)
	}
	gzip, err := metrics.Gzipper(jsons)
	if err != nil {
		logging.Fatalf("cannot gzip body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, `/updates/`, gzip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(req)

}

func ExampleUpdateHandler() {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()

	router.POST("/update/:types/:name/:value", url.UpdateHandler)

	req, err := http.NewRequest(http.MethodPost, "/update/counter/test/666", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)

}

func ExampleValueHandler() {
	logging.InitLogger()
	storage.Init(context.Background(), &config.ServerConfig{})

	router := gin.Default()

	router.GET("/value/:types/:name", url.ValueHandler)

	req, err := http.NewRequest(http.MethodGet, "/value/counter/test/", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(req)

}
