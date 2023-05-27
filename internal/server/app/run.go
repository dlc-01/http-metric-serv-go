package app

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/checkinghash"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.ServerConfig, s storage.Storage) {

	handlers.ServerStor.Storage = s
	router := setupRouter(cfg)
	router.Run(cfg.ServerAddress)
}

func setupRouter(cfg *config.ServerConfig) *gin.Engine {
	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	if cfg.HashKey != "" {
		router.Use(checkinghash.ChekHash(cfg.HashKey))
	}
	router.POST("/value/", handlers.ServerStor.ValueJSONHandler)
	router.GET("/value/:types/:name", handlers.ServerStor.ValueHandler)
	router.GET("/", handlers.ServerStor.ShowMetrics)
	router.GET("/ping", handlers.ServerStor.PingDB)
	updateRouterGroup := router.Group("/")
	if cfg.DatabaseAddress != "" {
		updateRouterGroup.Use(storagesync.GetSyncMiddleware(cfg.DatabaseAddress))
	}

	{
		updateRouterGroup.POST("/update", handlers.ServerStor.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", handlers.ServerStor.UpdateHandler)
		updateRouterGroup.POST("/updates", handlers.ServerStor.UpdatesButchJSONHandler)

	}
	return router
}
