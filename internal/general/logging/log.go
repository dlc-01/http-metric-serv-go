package logging

import (
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var sLog zap.SugaredLogger

// InitLogger — initialization function of zap logger.
func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("cannot init logger: %w", err)
	}
	sLog = *logger.Sugar()
	return nil
}

// Fatalf — function that equals Fatalf error in the logger.
func Fatalf(format string, opts ...any) {
	sLog.Fatalf(format, opts)
}

// Errorf — function that equals Errorf error in the logger.
func Errorf(format string, opts ...any) {
	sLog.Errorf(format, opts)
}

// Infof — function that equals Infof error in the logger.
func Infof(format string, opts ...any) {
	sLog.Infof(format, opts)
}

// Warnf — function that equals Warnf error in the logger.
func Warnf(format string, opts ...any) {
	sLog.Warnf(format, opts)
}

// Panicf — function that equals Panicf error in the logger.
func Panicf(format string, opts ...any) {
	sLog.Panicf(format, opts)
}

// Info — function that equals Info error in the logger.
func Info(msg string) {
	sLog.Info(msg)
}

// GetMiddlewareLogger — middleware for gin that input request and response logging.
func GetMiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		sLog.Infoln(
			"type", "request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"duration", latency,
		)
		sLog.Infoln(
			"type", "response",
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}
}
