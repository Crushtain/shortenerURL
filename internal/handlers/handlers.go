package handlers

import (
	"fmt"
	"github.com/Crushtain/shortenerURL/internal/encode"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/Crushtain/shortenerURL/internal/config"
	"github.com/Crushtain/shortenerURL/internal/storage"
)

type Config struct {
	mu      sync.RWMutex
	storage *storage.InMemory
	url     *config.Config
}

func NewURL(store *storage.InMemory, config *config.Config) *Config {
	return &Config{storage: store, url: config}
}

func (c *Config) Save(short string, body string) {

	c.storage.Put(short, body)

}

func (c *Config) Shorten(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	short := encode.Encode(body)

	c.Save(string(short), string(body))

	response := fmt.Sprintf(c.url.ResultURL+"/%s", short)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, _ = io.WriteString(w, response)

}

func (c *Config) Original(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	lastSlashIndex := strings.LastIndex(urlPath, "/")
	short := urlPath[lastSlashIndex+1:]

	if short == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	out := c.storage.Get(string(short))

	w.Header().Set("Location", string(out))
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println(out)
}
