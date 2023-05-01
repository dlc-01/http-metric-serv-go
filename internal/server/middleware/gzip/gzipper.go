package gzip

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
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

			gin.Writer.Header().Set("Content-Encoding", "gzip")

			r, err := gzip.NewReader(gin.Request.Body)
			if err != nil {
				gin.AbortWithError(http.StatusBadRequest, err)
				return
			}
			gin.Request.Body = r
			r.Close()
			gin.Next()
		}
		if !shouldCompress(gin.Request) {
			return
		}
		gz, err := gzip.NewWriterLevel(gin.Writer, level)
		if err != nil {
			return
		}

		gin.Header("Content-Encoding", "gzip")
		gin.Header("Vary", "Accept-Encoding")
		gin.Writer = &gzipWriter{gin.Writer, gz}
		defer func() {
			gin.Header("Content-Length", "0")
			gz.Close()
		}()
		gin.Next()
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

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
