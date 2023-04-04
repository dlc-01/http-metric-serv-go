package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	storage.Ms.Init()

	mux.HandleFunc("/update/", handlers.UpdateHandlers)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
