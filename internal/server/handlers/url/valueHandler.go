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
	metric := metrics.Metric{ID: gin.Param("name"), MType: gin.Param("types")}

	var err error

	switch metric.MType {
	case metrics.CounterType:
		metric, err = storage.ServerStorage.GetMetric(gin, metric)
		if err != nil {
			logging.Info(fmt.Sprintf("cannot found metric %q", metric.ID))
			gin.String(http.StatusNotFound, fmt.Sprintf("Metric %q not found", metric.ID))
			return
		}

	case metrics.GaugeType:
		metric, err = storage.ServerStorage.GetMetric(gin, metric)
		if err != nil {
			logging.Info(fmt.Sprintf("cannot found metric %q", metric.ID))
			gin.String(http.StatusNotFound, fmt.Sprintf("Metric %q not found", metric.ID))
			return
		}

	default:
		logging.Info("cannot find metric type")
		gin.String(http.StatusNotFound, "Unsupported metric type")
	}
	switch metric.MType {
	case metrics.CounterType:
		gin.String(http.StatusOK, fmt.Sprintf("%v", *metric.Delta))

	case metrics.GaugeType:
		gin.String(http.StatusOK, fmt.Sprintf("%v", *metric.Value))
	}
}
