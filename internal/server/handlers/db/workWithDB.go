package db

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConnectDB(gin *gin.Context) {
	if database.ConnectDB() {
		gin.Status(http.StatusOK)
	} else {
		gin.Status(http.StatusInternalServerError)
	}
}
