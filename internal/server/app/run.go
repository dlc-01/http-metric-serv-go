package app

import (
	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/all"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/db"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/jsonbatch"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/checkinghash"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
)

// Run — function that set up and run sever.
func Run(cfg *config.ServerConfig) {
	router := setupRouter(cfg)

	router.Run(cfg.ServerAddress)

}

func setupRouter(cfg *config.ServerConfig) *gin.Engine {
	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	if cfg.HashKey != "" {
		router.Use(checkinghash.CheckHash(cfg.HashKey))
	}
	router.POST("/value/", json.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	router.GET("/", all.ShowMetrics)
	router.GET("/ping", db.PingDB)
	updateRouterGroup := router.Group("/")
	if cfg.DatabaseAddress == "" {
		updateRouterGroup.Use(storagesync.GetSyncMiddleware())
	}

	{
		updateRouterGroup.POST("/update", json.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", url.UpdateHandler)
		updateRouterGroup.POST("/updates", jsonbatch.UpdatesButchJSONHandler)

	}
	return router
}
