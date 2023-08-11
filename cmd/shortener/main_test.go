package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFoo(t *testing.T) {
	// Создаем новый маршрутизатор
	router := chi.NewRouter()

	// Заглушка для хендлера "Shorten"
	shortenHandlerCalls := 0
	shortenHandler := func(w http.ResponseWriter, r *http.Request) {
		shortenHandlerCalls++
		w.WriteHeader(http.StatusCreated)
	}

	// Заглушка для хендлера "Original"
	originalHandlerCalls := 0
	originalHandler := func(w http.ResponseWriter, r *http.Request) {
		originalHandlerCalls++
	}

	// Регистрируем заглушки хендлеров на соответствующие маршруты
	router.Post("/", shortenHandler)
	router.Get("/{id}", originalHandler)

	// Запускаем тестовый сервер
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Отправляем POST запрос на "/"
	resp, err := http.Post(ts.URL+"/", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверяем статус код ответа POST запроса
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, получен %d",
			http.StatusCreated, resp.StatusCode)
	}

	// Отправляем GET запрос на "/{id}"
	resp, err = http.Get(ts.URL + "/abc123")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Проверяем статус код ответа GET запроса
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получен %d",
			http.StatusOK, resp.StatusCode)
	}

	// Проверяем, что хендлеры были вызваны
	if shortenHandlerCalls != 1 {
		t.Errorf("Ожидался один вызов хендлера Shorten, получено %d",
			shortenHandlerCalls)
	}
	if originalHandlerCalls != 1 {
		t.Errorf("Ожидался один вызов хендлера Original, получено %d",
			originalHandlerCalls)
	}
}
