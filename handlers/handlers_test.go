package handlers

import (
	"bytes"
	"fmt"
	"github.com/Crushtain/shortenerURL/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	// Создаем новый тестовый запрос
	body := []byte("example-url")
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
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

type mockStorage struct {
	data     map[string]string
	InMemory *storage.InMemory
}

func (m *mockStorage) Put(short, body string) {
	m.data[short] = body
}

func (m *mockStorage) Get(short string) string {
	return m.data[short]
}

func TestOriginalHandler(t *testing.T) {
	// Создаем экземпляр URLHandler с использованием mockStorage
	storage := &mockStorage{
		data: make(map[string]string),
	}
	urlHandler := &URLHandler{
		storage: storage.InMemory,
	}

	// Создаем тестовый HTTP-запрос с корректным шорткодом
	req, err := http.NewRequest("GET", "/abc123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Вызываем обработчик Original
	handler := http.HandlerFunc(urlHandler.Original)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %v but got %v", http.StatusTemporaryRedirect, status)
	}

	// Проверяем заголовок Location
	expectedLocation := "value"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Expected Location header %q but got %q", expectedLocation, location)
	}

	// Проверяем вывод в консоль
	expectedOutput := "value"
	actualOutput := fmt.Sprint(rr.Body)
	if actualOutput != expectedOutput {
		t.Errorf("Expected output %q but got %q", expectedOutput, actualOutput)
	}

	// Создаем тестовый HTTP-запрос с некорректным шорткодом
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Вызываем обработчик Original
	handler = http.HandlerFunc(urlHandler.Original)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %v but got %v", http.StatusBadRequest, status)
	}
}
