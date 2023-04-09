package handlers

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValueHandler(gin *gin.Context) {
	types := gin.Param("types")
	key := gin.Param("name")
	switch types {
	case "counter":
		value, exist := storage.Ms.GetCounter(key)
		if !exist {
			gin.String(http.StatusNotFound, "Not a supported metric.")
			return
		}
		gin.String(http.StatusOK, fmt.Sprintf("%d", value))
	case "gauge":
		value, exist := storage.Ms.GetGauge(key)
		if !exist {
			gin.String(http.StatusNotFound, "Not a supported metric.")
			return
		}
		gin.String(http.StatusOK, fmt.Sprintf("%.3f", value))
	default:
		gin.String(http.StatusNotFound, "Not a supported metric.")
		return
	}
}
