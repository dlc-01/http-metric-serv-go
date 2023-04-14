package handlers

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateHandler(gin *gin.Context) {
	types := gin.Param("types")
	key := gin.Param("name")
	values := gin.Param("value")

	//url := strings.Split(keyValue, "/")
	//if len(url) != 3 {
	//	gin.String(http.StatusNotFound, "Unsupported URL.")
	//	return
	//}
	//
	//key := url[1]

	switch types {
	case "counter":
		value, err := strconv.ParseInt(values, 10, 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}
		storage.SetCounter(key, value)

		gin.String(http.StatusOK, createResponse(key, value))

	case "gauge":
		value, err := strconv.ParseFloat(values, 64)
		if err != nil {
			gin.String(http.StatusBadRequest, "Unsupported values")
			return
		}

		storage.SetGauge(key, value)

		gin.String(http.StatusOK, createResponse(key, value))

	default:
		gin.String(http.StatusNotImplemented, "Unsupported metric type")
		return
	}
}
