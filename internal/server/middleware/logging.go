package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"time"
)

var sugarLog zap.SugaredLogger

func InitLogging() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	defer logger.Sync()
	sugarLog = *logger.Sugar()
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		sugarLog.Infoln(
			"type", "request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", latency,
		)
		sugarLog.Infoln(
			"type", "response",
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}

}
