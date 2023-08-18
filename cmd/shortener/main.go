package main

import (
	"flag"
	"net/http"

	"github.com/Crushtain/shortenerURL/internal/config"
	"github.com/Crushtain/shortenerURL/internal/handlers"
	"github.com/Crushtain/shortenerURL/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	conf := config.New()
	flag.Parse()

	memory := storage.NewInMemory()
	urlHandler := handlers.NewURL(memory, conf)

	router := chi.NewRouter()
	router.Post("/", urlHandler.Shorten)
	router.Get("/{id}", urlHandler.Original)

	err := http.ListenAndServe(conf.Host, router)
	if err != nil {
		panic(err)
	}
}
