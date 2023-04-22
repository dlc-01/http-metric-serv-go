package json

import (
	"bytes"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJSONHandler(gin *gin.Context) {
	var metrics storage.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(gin.Request.Body)
	if err != nil {
		gin.String(http.StatusBadRequest, "Unsupported request body")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		return
	}

	types := metrics.MType
	key := metrics.ID

	switch types {
	case url.CounterTypeName:
		value := *metrics.Delta

		storage.SetCounter(key, value)
		value, _ = storage.GetCounter(key)

		result := storage.Metrics{
			ID:    metrics.ID,
			MType: metrics.MType,
			Delta: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	case url.GaugeTypeName:
		value := *metrics.Value

		storage.SetGauge(key, value)

		result := storage.Metrics{
			ID:    metrics.ID,
			MType: metrics.MType,
			Value: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	default:
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}

}
