package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValueJSONHandler(gin *gin.Context) {
	var metric metrics.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(gin.Request.Body)
	if err != nil {
		gin.String(http.StatusBadRequest, "Unsupported postRequest body")
		logging.Errorf("cannot read postRequest body", err)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		logging.Errorf("cannot unmarshal JSON", err)
		return
	}

	types := metric.MType
	key := metric.ID

	switch types {
	case metrics.CounterType:
		value, exist := storage.GetCounter(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Counter %q not found", key))
			logging.Info(fmt.Sprintf("cannot found counter %q ", key))
			return
		}

		result := metrics.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	case metrics.GaugeType:
		value, exist := storage.GetGauge(key)
		if !exist {
			gin.String(http.StatusNotFound, fmt.Sprintf("Gauge %q not found", key))
			logging.Info(fmt.Sprintf("cannot found gauge %q", key))
			return
		}
		result := metrics.Metric{
			ID:    metric.ID,
			MType: metric.MType,
			Value: &value,
		}
		gin.SecureJSON(http.StatusOK, result)

	default:
		gin.String(http.StatusNotFound, "Unsupported metric type")
		logging.Info("cannot find metric type")
		return
	}
}
