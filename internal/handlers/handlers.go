package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func createResponse(key string, value interface{}) string {
	response := fmt.Sprintf("%s was adeed, value = %v", key, value)
	return response
}
func parsUrl(url string, writer http.ResponseWriter) (string, string) {
	urlContainer := strings.Split(url, "/")
	if len(urlContainer) != 5 {
		http.Error(writer, "Unsupported URL.", http.StatusBadRequest)
		return "", ""
	}
	return urlContainer[3], urlContainer[4]
}
func check(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "HTTP method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if request.Header.Get("Content-Type") != "text/plain" {
		http.Error(writer, "Unsupported content-type.", http.StatusUnsupportedMediaType)
		return
	}
}
