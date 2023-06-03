package limit

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Limit(limit int) gin.HandlerFunc {
	if limit <= 0 {
		logging.Fatalf("limit number is 0")
	}
	sema := make(chan struct{}, limit)
	return func(context *gin.Context) {
		var called, fulled bool
		defer func() {
			if !called && !fulled {
				<-sema
			}
		}()
		select {
		case sema <- struct{}{}:
			context.Next()
			called = true
			<-sema
		default:
			fulled = true
			context.AbortWithStatus(http.StatusBadGateway)
		}
	}
}
