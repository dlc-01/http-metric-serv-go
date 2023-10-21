package decryptor

import (
	"bytes"
	"github.com/dlc-01/http-metric-serv-go/internal/general/encryption"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func DecryptMiddleware() gin.HandlerFunc {
	return func(gin *gin.Context) {
		if encryption.MetDecryptor != nil {
			buf, _ := io.ReadAll(gin.Request.Body)
			message, err := encryption.MetDecryptor.Decrypt(buf)
			if err != nil {
				gin.AbortWithStatus(http.StatusBadRequest)
				logging.Errorf("error while decryption request body: %w", err)
				return
			}

			gin.Request.Body = io.NopCloser(bytes.NewBuffer(message))
		}
	}
}
