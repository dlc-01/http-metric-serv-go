package url

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateHandler(gin *gin.Context) {

	types := gin.Param("types")
	key := gin.Param("name")
	values := gin.Param("value")

	switch types {
	case metrics.CounterType:
		value, err := strconv.ParseInt(values, 10, 64)
		if err != nil {
			logging.Errorf("cannot parse counter: %s", err)
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}

		storage.SetCounter(key, value)
		metric, _ := storage.GetCounter(key)

		gin.String(http.StatusOK, handlers.CreateResponse(key, *metric.Delta))

	case metrics.GaugeType:
		value, err := strconv.ParseFloat(values, 64)
		if err != nil {
			logging.Errorf("cannot parse gauge: %s", err)
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}

		storage.SetGauge(key, value)
		gin.String(http.StatusOK, handlers.CreateResponse(key, value))

	default:
		logging.Info("cannot find metric type")
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}
}
