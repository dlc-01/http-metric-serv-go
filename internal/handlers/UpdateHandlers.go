package handlers

import (
	"net/http"
	"storage"
	"strconv"
)

func UpdateHandlers(writer http.ResponseWriter, request *http.Request) {
	check(writer, request)

	types, key, valueStr := parsUrl(request.URL.String(), writer)
	switch types {
	case "counter":
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(writer, "Unsupported value", http.StatusBadRequest)
			return
		}
		storage.Ms.SetCounter(key, value)
		writer.Write([]byte(createResponse(key, storage.Ms.GetCounter(key))))
		writer.WriteHeader(http.StatusOK)
	case "gauge":
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(writer, "Unsupported value", http.StatusBadRequest)
			return
		}
		storage.Ms.SetGauge(key, value)
		writer.Write([]byte(createResponse(key, storage.Ms.GetGauge(key))))
		writer.WriteHeader(http.StatusOK)
	default:
		http.Error(writer, "Not a supported metric.", http.StatusNotImplemented)
		return
	}

}
