package main

import (
	"github.com/Crushtain/shortenerURL/handlers"
	"github.com/Crushtain/shortenerURL/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	memory := storage.NewInMemory()
	urlHandler := handlers.NewUrl(memory)
	router := chi.NewRouter()
	router.Post("/", urlHandler.Shorten)
	//router.Get("/:id", urlHandler.Original)

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		panic(err)
	}
}
