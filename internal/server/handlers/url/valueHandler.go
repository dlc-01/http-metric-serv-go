package url

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValueHandler(gin *gin.Context) {
	types := gin.Param("types")
	key := gin.Param("name")

	var exist bool
	var metric metrics.Metric
	switch types {
	case metrics.CounterType:
		metric, exist = storage.GetCounter(key)

	case metrics.GaugeType:
		metric, exist = storage.GetGauge(key)

	default:
		logging.Info("cannot find metric type")
		gin.String(http.StatusNotFound, "Unsupported metric type")
		return
	}

	if !exist {
		logging.Info(fmt.Sprintf("cannot found metric %q", key))
		gin.String(http.StatusNotFound, fmt.Sprintf("Metric %q not found", key))
		return
	}
	switch types {
	case metrics.CounterType:
		gin.String(http.StatusOK, fmt.Sprintf("%v", *metric.Delta))

	case metrics.GaugeType:
		gin.String(http.StatusOK, fmt.Sprintf("%v", *metric.Value))
	}

}
