package jsonbatch

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

// UpdatesButchJSONHandler — handler that saves metrics in butches in json format.
func UpdatesButchJSONHandler(gin *gin.Context) {
	var data []metrics.Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(gin.Request.Body)
	if err != nil {
		logging.Errorf("cannot read request body: %s", err)
		gin.String(http.StatusBadRequest, "Unsupported request body")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		logging.Errorf("cannot unmarshal json: %s", err)
		gin.String(http.StatusBadRequest, "Unsupported type JSON")
		return
	}

	if err = storage.SetMetricsBatch(gin, data); err != nil {
		logging.Errorf("cannot save metric type: %s", err)
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}

	gin.String(http.StatusOK, "Saved metrics")
}
