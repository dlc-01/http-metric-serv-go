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
		logging.Errorf("cannot read postRequest body", err)
		gin.String(http.StatusBadRequest, "Unsupported postRequest body")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		logging.Errorf("cannot unmarshal JSON", err)
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		return
	}

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
		return
	}

	gin.SecureJSON(http.StatusOK, metric)
}
