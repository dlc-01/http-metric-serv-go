package checkinghash

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
)

// CheckHash — middleware that use for checking header for http request "HashSHA256".
func CheckHash(key string) gin.HandlerFunc {
	return func(gin *gin.Context) {
		if key != "" {
			body, err := io.ReadAll(gin.Request.Body)
			if err != nil {
				logging.Errorf("cannot read request body %s", err)
				return
			}
			gin.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			hash := gin.Request.Header.Get("HashSHA256")
			if hash == "" {
				gin.Next()
			}
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
