package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/htmlh"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/jsonh"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/paramss"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(logging.Logger(), gzip.Gzip(gzip.BestSpeed))
	router.POST("/update/:types/:name/:value", url.UpdateHandler)
	router.POST("/update/", jsonh.UpdateJSONHandler)
	router.POST("/value/", jsonh.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	router.GET("/", htmlh.ShowMetrics)
	return router
}

func main() {
	logging.InitLogger()

	paramss.ParseFlagsOs()

	router := setupRouter()

	storage.Init()

	if paramss.Restore {
		if err := storage.Restore(paramss.FileStoragePath); err != nil {
			logging.SLog.Error(err, "restore")
		}
	}

	go func() {
		for {
			if err := storage.Save(paramss.FileStoragePath); err != nil {
				logging.SLog.Error(err, "saved")
			}
			time.Sleep(time.Duration(paramss.StoreInterval) * time.Second)
		}
	}()

	router.Run(paramss.ServerAddress)

}
