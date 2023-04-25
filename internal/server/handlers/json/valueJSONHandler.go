package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValueJSONHandler(gin *gin.Context) {
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
		value, exist := storage.GetCounter(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Counter %q not found", key))
			return
		}

		result := storage.Metrics{
			ID:    metrics.ID,
			MType: metrics.MType,
			Delta: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	case url.GaugeTypeName:
		value, exist := storage.GetGauge(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Gauge %q not found", key))
			return
		}
		result := storage.Metrics{
			ID:    metrics.ID,
			MType: metrics.MType,
			Value: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	default:
		gin.String(http.StatusNotFound, "Unsupported metric type")
		return
	}
}
