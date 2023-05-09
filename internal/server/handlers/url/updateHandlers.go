package url

import (
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
	case CounterTypeName:
		value, err := strconv.ParseInt(values, 10, 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported values")
			//logging.Errorf("cannot parse counter: %s", err)
			return
		}
		storage.SetCounter(key, value)
		value, _ = storage.GetCounter(key)

		gin.String(http.StatusOK, handlers.CreateResponse(key, value))

	case GaugeTypeName:
		value, err := strconv.ParseFloat(values, 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported values")
			//logging.Errorf("cannot parse gauge: %s", err)
			return
		}
		storage.SetGauge(key, value)
		gin.String(http.StatusOK, handlers.CreateResponse(key, value))

	default:
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		//logging.Info("cannot find metric type")
		return
	}
}
