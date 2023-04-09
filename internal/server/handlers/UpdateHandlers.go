package handlers

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func UpdateHandler(gin *gin.Context) {

	types := gin.Param("types")
	keyValue := gin.Param("name")
	url := strings.Split(keyValue, "/")
	if len(url) != 3 {
		gin.String(http.StatusNotFound, "Unsupported URL.")
		return
	}
	key := url[1]

	switch types {
	case "counter":
		value, err := strconv.ParseInt(url[2], 10, 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported value")
			return
		}
		storage.Ms.SetCounter(key, value)
		value, _ = storage.Ms.GetCounter(key)
		gin.String(http.StatusOK, createResponse(key, value))
	case "gauge":
		value, err := strconv.ParseFloat(url[3], 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported value")
			return
		}
		storage.Ms.SetGauge(key, value)
		value, _ = storage.Ms.GetGauge(key)
		gin.String(http.StatusOK, createResponse(key, value))
	default:
		gin.String(http.StatusNotImplemented, "Not a supported metric.")
		return
	}
}
