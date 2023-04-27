package gzip

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	gin.Writer.Header().Add("Content-Encoding", "gzip")
}
