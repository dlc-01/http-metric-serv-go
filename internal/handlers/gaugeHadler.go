package handlers

import (
	"net/http"
	"storage"
	"strconv"
)

func UpdateGauge(writer http.ResponseWriter, request *http.Request) {
	check(writer, request)
	key, valueStr := parsUrl(request.URL.Path, writer)

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(writer, "Unsupported value", http.StatusBadRequest)
		return
	}
	storage.Ms.SetGauge(key, value)

	writer.Write([]byte(createResponse(key, storage.Ms.GetGauge(key))))
	writer.WriteHeader(http.StatusOK)

}
