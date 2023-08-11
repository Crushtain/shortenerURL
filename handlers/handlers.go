package handlers

import (
	"fmt"
	storage "github.com/Crushtain/shortenerURL/storage"
	"github.com/jxskiss/base62"
	"io"
	"net/http"
	"strings"
	"sync"
)

type URLHandler struct {
	mu      sync.RWMutex
	storage *storage.InMemory
}

func NewURL(store *storage.InMemory) *URLHandler {
	return &URLHandler{storage: store}
}

func (h *URLHandler) Save(short string, body string) {

	h.storage.Put(short, body)

}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
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

func (h *URLHandler) Original(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	lastSlashIndex := strings.LastIndex(urlPath, "/")
	short := urlPath[lastSlashIndex+1:]
	//short := strings.TrimPrefix(r.URL.Path, "/")
	if short == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	out := h.storage.Get(string(short))

	w.Header().Set("Location", string(out))
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println(out)
}
