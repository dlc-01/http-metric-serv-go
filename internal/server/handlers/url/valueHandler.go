package url

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

// ValueHandler — handler that searches for a metric by name and type taken from the request url.
func ValueHandler(gin *gin.Context) {
	metric := metrics.Metric{ID: gin.Param("name"), MType: gin.Param("types")}

	var err error

	switch metric.MType {
	case metrics.CounterType:
		metric, err = storage.GetMetric(gin, metric)
		if err != nil {
			logging.Info(fmt.Sprintf("cannot found metric %q", metric.ID))
			gin.String(http.StatusNotFound, fmt.Sprintf("Metric %q not found", metric.ID))
			return
		}

	case metrics.GaugeType:
		metric, err = storage.GetMetric(gin, metric)
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
