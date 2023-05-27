package handlers

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s stor) PingDB(gin *gin.Context) {

	if err := s.PingStorage(gin); err != nil {
		logging.Errorf("error while try to ping db: %s", err)
		gin.Status(http.StatusInternalServerError)
	} else {
		gin.Status(http.StatusOK)
	}
}
