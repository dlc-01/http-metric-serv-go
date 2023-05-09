package gzip

import (
	"compress/gzip"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int) gin.HandlerFunc {
	return func(gin *gin.Context) {
		if gin.Request.Header.Get("Content-Encoding") == "gzip" {
			newCompressReader(gin)
			newCompressWriter(gin, level)
		} else if strings.Contains(gin.Request.Header.Get("Accept-Encoding"), "gzip") && gin.Request.Header.Get("Content-Encoding") != "gzip" {
			newCompressWriter(gin, level)
		}

	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func newCompressReader(gin *gin.Context) {
	r, err := gzip.NewReader(gin.Request.Body)
	if err != nil {
		gin.AbortWithError(http.StatusBadRequest, err)
		logging.Errorf("cannot uncompressed request body: %s", err)
		return
	}
	gin.Request.Body = r
	defer r.Close()
	gin.Next()
}

func newCompressWriter(gin *gin.Context, level int) {
	gz, err := gzip.NewWriterLevel(gin.Writer, level)
	if err != nil {
		logging.Errorf("cannot compress request body: %s", err)
		return
	}

	gin.Writer = &gzipWriter{gin.Writer, gz}
	defer gz.Close()
	gin.Header("Content-Encoding", "gzip")
	gin.Next()
}
