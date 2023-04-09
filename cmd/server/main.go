package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	storage.Ms.Init()

	router.POST("/update/:types/*name", handlers.UpdateHandler)

	router.GET("/value/:types/:name", handlers.ValueHandler)

	router.Run(":8080")

}
