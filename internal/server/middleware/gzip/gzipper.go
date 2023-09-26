package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// Gzip — middleware for gin that use for compress and decompress gzip.
func Gzip(level int) gin.HandlerFunc {
	return func(gin *gin.Context) {
		headerContentGzip := strings.Contains(gin.Request.Header.Get("Content-Encoding"), "gzip")
		headerAcceptGzip := strings.Contains(gin.Request.Header.Get("Accept-Encoding"), "gzip")

		if headerContentGzip {
			newCompressReader(gin)
			newCompressWriter(gin, level)
		} else if headerAcceptGzip && !headerContentGzip {
			newCompressWriter(gin, level)
		}
	}
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
