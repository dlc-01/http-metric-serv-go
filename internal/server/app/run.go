package app

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func Run(servAdress string, s storage.Storage) {

	handlers.ServerStor.Storage = s
	router := setupRouter()
	router.Run(servAdress)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	router.POST("/value/", handlers.ServerStor.ValueJSONHandler)
	router.GET("/value/:types/:name", handlers.ServerStor.ValueHandler)
	router.GET("/", handlers.ServerStor.ShowMetrics)
	router.GET("/ping", handlers.ServerStor.PingDB)
	updateRouterGroup := router.Group("/")
	updateRouterGroup.Use(storagesync.GetSyncMiddleware())

	{
		updateRouterGroup.POST("/update", handlers.ServerStor.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", handlers.ServerStor.UpdateHandler)
		updateRouterGroup.POST("/updates", handlers.ServerStor.UpdatesButchJSONHandler)

	}
	return router
}
