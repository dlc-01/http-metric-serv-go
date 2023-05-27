package checkinghash

import (
	"bytes"
	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ChekHash(key string) gin.HandlerFunc {
	return func(gin *gin.Context) {
		if key != "" {
			hash := gin.Request.Header.Get("HashSHA256")
			body, err := io.ReadAll(gin.Request.Body)
			if err != nil {
				logging.Errorf("cannot read request body %s", err)
				return
			}
			gin.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			check, err := hashing.CheckingHash(hash, key, body)
			if err != nil {
				logging.Errorf("cannot check hash %s", err)
				gin.AbortWithStatus(http.StatusBadRequest)

			}
			if !check {
				gin.AbortWithStatus(http.StatusBadRequest)
			}
		}

		gin.Next()
	}
}