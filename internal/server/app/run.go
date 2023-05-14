package app

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/butchjson"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/db"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/html"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
	"github.com/gin-gonic/gin"
)

func Run(serverAddress string) {
	router := setupRouter()
	router.Run(serverAddress)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(logging.GetMiddlewareLogger(), gzip.Gzip(gzip.BestSpeed))
	router.POST("/value/", json.ValueJSONHandler)
	router.GET("/value/:types/:name", url.ValueHandler)
	router.GET("/", html.ShowMetrics)
	router.GET("/ping", db.ConnectDB)
	updateRouterGroup := router.Group("/")
	updateRouterGroup.Use(storagesync.GetSyncMiddleware())
	{
		updateRouterGroup.POST("/update", json.UpdateJSONHandler)
		updateRouterGroup.POST("/update/:types/:name/:value", url.UpdateHandler)
		updateRouterGroup.POST("/updates", butchjson.UpdatesButchJSONHandler)

	}

	return router
}
