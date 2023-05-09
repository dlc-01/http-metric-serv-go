package json

import (
	"bytes"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJSONHandler(gin *gin.Context) {
	var metric metrics.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(gin.Request.Body)
	if err != nil {
		gin.String(http.StatusBadRequest, "Unsupported postRequest body")
		logging.Errorf("cannot read postRequest body: %s", err)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		logging.Errorf("cannot unmarshal json: %s", err)
		return
	}

	types := metric.MType
	key := metric.ID

	switch types {
	case url.CounterTypeName:
		value := *metric.Delta

		storage.SetCounter(key, value)
		value, _ = storage.GetCounter(key)

		result := metrics.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	case url.GaugeTypeName:
		value := *metric.Value

		storage.SetGauge(key, value)

		result := metrics.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Value: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	default:
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		logging.Info("cannot find metric type")
		return
	}

}
