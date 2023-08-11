package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	// Создаем новый тестовый запрос
	body := []byte("example-url")
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseWriter для записи ответа
	recorder := httptest.NewRecorder()

	// Вызываем хендлер
	h := &URLHandler{}
	h.Shorten(recorder, req)

	// Проверяем статус код ответа
	if recorder.Code != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, получен %d",
			http.StatusCreated, recorder.Code)
	}

	// Проверяем тип данных контента
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Ожидался тип данных 'text/plain', получен '%s'",
			contentType)
	}

	// Проверяем тело ответа
	responseBody, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedResponse := "http://localhost:8080/" + string(body)
	actualResponse := strings.TrimSpace(string(responseBody))
	if expectedResponse != actualResponse {
		t.Errorf("Ожидался ответ '%s', получен '%s'",
			expectedResponse, actualResponse)
	}
}
