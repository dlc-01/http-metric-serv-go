package json

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

// UpdateJSONHandler — handler that save metric that in json format.
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

	err = storage.SetMetric(gin, metric)
	if err != nil {
		logging.Errorf("cannot save metric: %s", err)
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}

	gin.SecureJSON(http.StatusOK, metric)

}
