package main

import (
	"flag"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware"
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
	router.Use(middleware.Logger())
	router.POST("/update/:types/:name/:value", handlers.UpdateHandler)
	router.GET("/value/:types/:name", handlers.ValueHandler)
	return router
}

func main() {
	middleware.InitLogging()

	parseFlagsOs()

	router := setupRouter()

	storage.Init()

	router.Run(serverAddress)

}
