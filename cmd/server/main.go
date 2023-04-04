package main

import (
	"handlers"
	"net/http"
	"storage"
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
