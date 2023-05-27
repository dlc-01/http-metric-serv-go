package db

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PingDB(gin *gin.Context) {

	if err := storage.ServerStorage.PingStorage(gin); err != nil {
		logging.Errorf("error while try to ping db: %s", err)
		gin.Status(http.StatusInternalServerError)
	} else {
		gin.Status(http.StatusOK)
	}
}
