package handlers

import (
	"fmt"
)

func CreateResponse(key string, value interface{}) string {
	response := fmt.Sprintf("%s was adeed, value = %v", key, value)
	return response
}
