package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/html"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/params"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(logging.Logger(), gzip.Gzip(gzip.BestSpeed))
	router.POST("/update/:types/:name/:value", url.UpdateHandler)
	router.POST("/update/", json.UpdateJSONHandler)
	router.POST("/value/", json.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	router.GET("/", html.ShowMetrics)
	return router
}

func main() {
	logging.InitLogger()

	params.ParseFlagsOs()

	router := setupRouter()

	storage.Init()

	if params.Restore {
		storage.Restore()
	}
	go func() {
		for {
			storage.Save()
			time.Sleep(time.Duration(params.StoreInterval) * time.Second)
		}

	}()
	router.Run(params.ServerAddress)
	defer storage.Save()
}
