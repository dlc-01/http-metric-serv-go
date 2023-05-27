package handlers

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func CreateResponse(key string, value interface{}) string {
	response := fmt.Sprintf("%s was adeed, value = %v", key, value)
	return response
}

type stor struct {
	storage.Storage
}

var ServerStor stor
