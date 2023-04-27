package gzip

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func DefaultDecompressHandle(gin *gin.Context) {
	if gin.Request.Header.Get("Content-Encoding") != "gzip" {
		gin.Next()
		return
	}
	if gin.Request.Body == nil {
		gin.Next()
		return
	}
	r, err := gzip.NewReader(gin.Request.Body)
	if err != nil {
		gin.AbortWithError(http.StatusBadRequest, err)
		return
	}
	gin.Request.Body = r
}
func DefaultCompressionResponse(gin *gin.Context) {
	if gin.Request.Header.Get("Accept-Encoding") != "gzip" {
		gin.Next()
		return
	}
	gz, err := gzip.NewWriterLevel(gin.Writer, gzip.BestSpeed)
	defer gz.Close()
	if err != nil {
		gin.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gin.Header("Content-Encoding", "gzip")

	gin.Writer = &gzipWriter{gin.Writer, gz}

}
