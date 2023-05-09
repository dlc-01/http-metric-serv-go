package url

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValueHandler(gin *gin.Context) {
	types := gin.Param("types")
	key := gin.Param("name")

	switch types {
	case metrics.CounterType:
		value, exist := storage.GetCounter(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Counter %q not found", key))
			//logging.Info(fmt.Sprintf("cannot found caounter %q", key))
			return
		}

		gin.String(http.StatusOK, fmt.Sprintf("%v", value))

	case metrics.GaugeType:
		value, exist := storage.GetGauge(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Gauge %q not found", key))
			//logging.Info(fmt.Sprintf("cannot found gauge %q", key))
			return
		}
		gin.String(http.StatusOK, fmt.Sprintf("%v", value))

	default:
		gin.String(http.StatusNotFound, "Unsupported metric type")
		//logging.Info("cannot find metric type")
		return
	}
}
