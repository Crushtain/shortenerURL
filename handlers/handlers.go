package handlers

import (
	"fmt"
	storage "github.com/Crushtain/shortenerURL/storage"
	"github.com/jxskiss/base62"
	"io"
	"net/http"
	"sync"
)

type UrlHandler struct {
	mu      sync.RWMutex
	storage *storage.InMemory
}

func NewUrl(store *storage.InMemory) *UrlHandler {
	return &UrlHandler{storage: store}
}

func (h *UrlHandler) Save(short string, body string) {

	h.storage.Put(short, body)

}

func (h *UrlHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	short := base62.Encode(body)        //кодирую
	h.Save(string(short), string(body)) //сохраняю

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Println(string(short))

}

func (h *UrlHandler) Original(w http.ResponseWriter, r *http.Request) {
	short, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	out := h.storage.Get(string(short))

	w.Header().Set("Location", string(out))
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println(out)
}
