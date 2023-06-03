package checkinghash

import (
	"bytes"
	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CheckHash(key string) gin.HandlerFunc {
	//TODO я вроде как все сделал по тз а тест не проходит вроде как требуется чтоб с ключаом было хещирование а на тесте без хэша
	return func(gin *gin.Context) {
		if key != "" {
			body, err := io.ReadAll(gin.Request.Body)
			if err != nil {
				logging.Errorf("cannot read request body %s", err)
				return
			}
			gin.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			hash := gin.Request.Header.Get("HashSHA256")
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
