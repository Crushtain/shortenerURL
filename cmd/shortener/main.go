package main

import (
	"encoding/json"
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

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp.ShortURL))
	//fmt.Printf("%s\n", resp.ShortURL)
}

func ShortenURL(url string) string {
	shortUrlID := "TestShortURL_id"
	URLStore[shortUrlID] = url
	return shortUrlID
}

func GetFunc(w http.ResponseWriter, r *http.Request) {

	shortUrlId := r.URL.Path[len("/"):]
	originalURL := GetOrigianlURL(shortUrlId)
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
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
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, PostFunc)
	mux.HandleFunc(`/{id}`, GetFunc)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
