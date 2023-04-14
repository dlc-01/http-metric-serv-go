package handlers

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	GaugeTypeName   = "gauge"
	CounterTypeName = "counter"
)

func ValueHandler(gin *gin.Context) {
	types := gin.Param("types")
	key := gin.Param("name")

	switch types {
	case CounterTypeName:
		value, exist := storage.GetCounter(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Counter %q not found", key))
			return
		}

		gin.String(http.StatusOK, fmt.Sprintf("%v", value))

	case GaugeTypeName:
		value, exist := storage.GetGauge(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Gauge %q not found", key))
			return
		}
		gin.String(http.StatusOK, fmt.Sprintf("%v", value))

	default:
		gin.String(http.StatusNotFound, "Unsupported metric type")
		return
	}
}
