package handlers

import (
	"net/http"
	"storage"
	"strconv"
)

func UpdateCounter(writer http.ResponseWriter, request *http.Request) {
	check(writer, request)

	key, valueStr := parsUrl(request.URL.Path, writer)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		http.Error(writer, "Unsupported value", http.StatusBadRequest)
		return
	}
	storage.Ms.SetCounter(key, value)

	writer.Write([]byte(createResponse(key, storage.Ms.GetCounter(key))))
	writer.WriteHeader(http.StatusOK)

}
