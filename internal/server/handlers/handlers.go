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
func parsURL(url string, writer http.ResponseWriter) (string, string, string) {
	urlContainer := strings.Split(url, "/")
	if len(urlContainer) != 5 {
		http.Error(writer, "Unsupported URL.", http.StatusNotFound)
		return "", "", ""
	}
	return urlContainer[2], urlContainer[3], urlContainer[4]
}
func check(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "HTTP method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
