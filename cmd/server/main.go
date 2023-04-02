package main

import (
	"handlers"
	"net/http"
	"storage"
)

func main() {
	mux := http.NewServeMux()
	storage.Ms.Init()

	mux.HandleFunc("/update/gauge/", handlers.UpdateGauge)
	mux.HandleFunc("/update/counter/", handlers.UpdateCounter)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
