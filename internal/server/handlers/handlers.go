package handlers

import (
	"fmt"
)

func createResponse(key string, value interface{}) string {
	response := fmt.Sprintf("%s was adeed, value = %v", key, value)
	return response
}
