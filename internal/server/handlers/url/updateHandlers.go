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
	metric := metrics.Metric{ID: gin.Param("name"), MType: gin.Param("types")}
	values := gin.Param("value")

	switch metric.MType {
	case metrics.CounterType:
		value, err := strconv.ParseInt(values, 10, 64)
		if err != nil {
			logging.Errorf("cannot parse counter: %s", err)
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}
		metric.Delta = &value

		err = storage.ServerStorage.SetMetric(gin, metric)
		if err != nil {
			logging.Errorf("cannot save metric: %s", err)
		}

		metric, err = storage.ServerStorage.GetMetric(gin, metric)
		if err != nil {
			logging.Errorf("cannot get metric: %s", err)
		}

		gin.String(http.StatusOK, handlers.CreateResponse(metric.ID, *metric.Delta))

	case metrics.GaugeType:
		value, err := strconv.ParseFloat(values, 64)
		if err != nil {
			logging.Errorf("cannot parse gauge: %s", err)
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}
		metric.Value = &value

		err = storage.ServerStorage.SetMetric(gin, metric)
		if err != nil {
			logging.Errorf("cannot save metric: %s", err)
		}

		gin.String(http.StatusOK, handlers.CreateResponse(metric.ID, metric.Value))

	default:
		logging.Info("cannot find metric type")
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}
}
