package json

import (
	"bytes"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateJSONHandler(gin *gin.Context) {
	var metric metrics.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(gin.Request.Body)
	if err != nil {
		logging.Errorf("cannot read postRequest body: %s", err)
		gin.String(http.StatusBadRequest, "Unsupported postRequest body")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		logging.Errorf("cannot unmarshal json: %s", err)
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		return
	}

	check := storage.SetMetric(metric.ID, metric.MType, metric.Value, metric.Delta)
	if !check {
		logging.Info("cannot find metric type")
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}

	gin.SecureJSON(http.StatusOK, metric)

}
