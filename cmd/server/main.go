package main

import (
	"flag"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
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
	router.POST("/update/:types/*name", handlers.UpdateHandler)
	router.GET("/value/:types/:name", handlers.ValueHandler)
	return router
}

func main() {
	parseFlagsOs()
	router := setupRouter()
	storage.Ms.Init()
	router.Run(serverAddress)
	fmt.Println("Server address:" + serverAddress)
}
