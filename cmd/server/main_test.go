package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitRouter(t *testing.T) {
	router := setupRouter()
	storage.Ms.Init()
	wPost := httptest.NewRecorder()
	reqPost, _ := http.NewRequest(http.MethodPost, "/update/gauge/test/666", nil)
	router.ServeHTTP(wPost, reqPost)
	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest(http.MethodGet, "/value/gauge/test", nil)
	router.ServeHTTP(wGet, reqGet)
	assert.Equal(t, 200, wPost.Code)
	assert.Equal(t, "666", wGet.Body.String())
}
