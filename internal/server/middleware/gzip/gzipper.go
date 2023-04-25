package gzip

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DefaultDecompressHandle(gin *gin.Context) {
	if gin.Request.Header.Get("Content-Encoding") != "gzip" {
		return
	}
	if gin.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(gin.Request.Body)
	if err != nil {
		_ = gin.AbortWithError(http.StatusBadRequest, err)
		return
	}
	gin.Request.Body = r

}
