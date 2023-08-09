package main

import (
	"encoding/json"
	"fmt"
	"github.com/itchyny/base58-go"
	"io"
	"net/http"
)

type OriginalURL struct {
	URL string `json:"url"`
}
type ShortURL struct {
	ShortURL string `json:"short-url"`
}

var URLStore map[string]string

func PostFunc(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//resp := ShortURL{
	//	ShortURL: shortURL,
	//}
	//out, _ := json.Marshal(&resp)
	encoding := base58.FlickrEncoding
	encode, err := encoding.Encode(body)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Println(string(encode))

}

func ShortenURL(url string) string {
	shortURLID := "TestShortURL_id"
	URLStore[shortURLID] = url
	return shortURLID
}

func GetFunc(w http.ResponseWriter, r *http.Request) {

	shortURLId := r.URL.Path[len("/"):]
	originalURL := GetOrigianlURL(shortURLId)
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	out, _ := json.Marshal(&originalURL)

	w.Header().Set("Location", string(out))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func GetOrigianlURL(id string) string {
	originalURL, ok := URLStore[id]
	if !ok {
		return ""
	}
	return originalURL
}

func main() {

	http.HandleFunc("/", PostFunc)
	http.HandleFunc("/{id}", GetFunc)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
