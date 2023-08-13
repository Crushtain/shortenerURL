package handlers

import (
	"fmt"
	"github.com/Crushtain/shortenerURL/config"
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
	url     *config.Config
}

func NewURL(store *storage.InMemory, config *config.Config) *URLHandler {
	return &URLHandler{storage: store, url: config}
}

func (h *URLHandler) Save(short string, body string) {

	h.storage.Put(short, body)

}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	address := r.Host
	if strings.Contains(address, ":") {
		// Удаление порта из адреса
		address = strings.Split(address, ":")[0]
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	short := base62.Encode(body)        //кодирую
	h.Save(string(short), string(body)) //сохраняю
	//out := "http://localhost:8080/" + string(short)
	//flag.StringVar("b", out, "fgwewe")
	response := fmt.Sprintf(h.url.ResultURL+"/%s", short)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	//fmt.Fprint(w, string(out))
	_, err = io.WriteString(w, response)
	if err != nil {
		return
	}

}

func (h *URLHandler) Original(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	lastSlashIndex := strings.LastIndex(urlPath, "/")
	short := urlPath[lastSlashIndex+1:]

	if short == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	out := h.storage.Get(string(short))

	w.Header().Set("Location", string(out))
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println(out)
}
