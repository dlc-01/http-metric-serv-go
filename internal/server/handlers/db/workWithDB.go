package db

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

// PingDB — handler that checks if the database is connected.
func PingDB(gin *gin.Context) {

	if err := storage.PingStorage(gin); err != nil {
		logging.Errorf("error while try to ping db: %s", err)
		gin.Status(http.StatusInternalServerError)
	} else {
		gin.Status(http.StatusOK)
	}
}
