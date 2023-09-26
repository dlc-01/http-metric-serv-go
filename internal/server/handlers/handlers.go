package handlers

import (
	"fmt"
)

// CreateResponse — function for generating string for  http response.
func CreateResponse(key string, value interface{}) string {
	response := fmt.Sprintf("%s was adeed, value = %v", key, value)
	return response
}
