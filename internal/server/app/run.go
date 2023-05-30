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

	handlers.ServerStorage.Storage = s
	router := setupRouter(cfg)
	router.Run(cfg.ServerAddress)
}

func setupRouter(cfg *config.ServerConfig) *gin.Engine {
	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	if cfg.HashKey != "" {
		router.Use(checkinghash.CheckHash(cfg.HashKey))
	}
	router.POST("/value/", handlers.ServerStorage.ValueJSONHandler)
	router.GET("/value/:types/:name", handlers.ServerStorage.ValueHandler)
	router.GET("/", handlers.ServerStorage.ShowMetrics)
	router.GET("/ping", handlers.ServerStorage.PingDB)
	updateRouterGroup := router.Group("/")
	if cfg.DatabaseAddress != "" {
		updateRouterGroup.Use(storagesync.GetSyncMiddleware(cfg.DatabaseAddress))
	}

	{
		updateRouterGroup.POST("/update", handlers.ServerStorage.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", handlers.ServerStorage.UpdateHandler)
		updateRouterGroup.POST("/updates", handlers.ServerStorage.UpdatesButchJSONHandler)

	}
	return router
}
