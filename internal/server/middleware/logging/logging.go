package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"time"
)

var SLog zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	defer logger.Sync()
	SLog = *logger.Sugar()
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		SLog.Infoln(
			"type", "request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", latency,
		)
		SLog.Infoln(
			"type", "response",
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}

}
