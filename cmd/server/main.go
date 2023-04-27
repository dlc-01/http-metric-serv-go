package main

import (
	"flag"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	serverAddress string
)

func parseFlagsOs() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "server address")
	flag.Parse()
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(logging.Logger(), gzip.DefaultDecompressHandle, gzip.DefaultCompressionResponse)
	router.POST("/update/:types/:name/:value", url.UpdateHandler)
	router.POST("/update/", json.UpdateJSONHandler)
	router.POST("/value/", json.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	return router
}

func main() {
	logging.InitLogger()

	parseFlagsOs()

	router := setupRouter()

	storage.Init()

	router.Run(serverAddress)

}
