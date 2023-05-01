package gzip

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func DefaultDecompressHandle(gin *gin.Context) {
	if gin.Request.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewWriterLevel(gin.Writer, gzip.BestSpeed)
		if err != nil {
			gin.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		gin.Writer = &gzipWriter{gin.Writer, gz}
		gz.Close()

		r, err := gzip.NewReader(gin.Request.Body)
		if err != nil {
			gin.AbortWithError(http.StatusBadRequest, err)
			return
		}
		gin.Request.Body = r
		r.Close()

	} else if strings.Contains(gin.Request.Header.Get("Accept-Encoding"), "gzip") && !strings.Contains(gin.Request.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewWriterLevel(gin.Writer, gzip.BestSpeed)
		if err != nil {
			gin.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		gin.Header("Content-Encoding", "gzip")
		gin.Writer = &gzipWriter{gin.Writer, gz}
		gz.Close()
	}

}
