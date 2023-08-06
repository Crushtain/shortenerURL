package main

import (
	"encoding/json"
	"fmt"
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

	var original OriginalURL
	err := json.NewDecoder(r.Body).Decode(&original)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := ShortenURL(original.URL)

	resp := ShortURL{
		ShortURL: shortURL,
	}
	out, _ := json.Marshal(&resp)
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Println(string(out))

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
