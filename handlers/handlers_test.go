package handlers

import (
	"fmt"
	"github.com/Crushtain/shortenerURL/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type mockStorage struct {
	data     map[string]string
	mu       sync.RWMutex
	InMemory *storage.InMemory
}

func (m *mockStorage) Put(short, body string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[short] = body
}

func (m *mockStorage) Get(short string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[short]
}

func TestShortenHandler(t *testing.T) {
	// Создаем экземпляр URLHandler с использованием mockStorage
	storage := &mockStorage{
		data: make(map[string]string),
	}
	urlHandler := &URLHandler{
		storage: storage.InMemory,
	}

	// Создаем тестовый HTTP-запрос с телом запроса
	reqBody := strings.NewReader("http://example.com")
	req, err := http.NewRequest("POST", "/", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Вызываем обработчик Shorten
	handler := http.HandlerFunc(urlHandler.Shorten)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Ожидаемый статус-код %v, получен %v", http.StatusCreated, status)
	}

	// Проверяем заголовок Content-type
	expectedContentType := "text/plain"
	contentType := rr.Header().Get("Content-type")
	if contentType != expectedContentType {
		t.Errorf("Ожидаемый заголовок Content-type %q, получен %q", expectedContentType, contentType)
	}

	// Проверяем вывод в теле ответа
	expectedOutput := "http://localhost:8080/encoded_short"
	actualOutput := fmt.Sprint(rr.Body)
	if actualOutput != expectedOutput {
		t.Errorf("Ожидаемый вывод %q, получен %q", expectedOutput, actualOutput)
	}

	// Создаем тестовый HTTP-запрос без тела запроса (неверный запрос)
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Вызываем обработчик Shorten
	handler = http.HandlerFunc(urlHandler.Shorten)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидаемый статус-код %v, получен %v", http.StatusBadRequest, status)
	}
}

func TestOriginalHandler(t *testing.T) {
	// Создаем экземпляр URLHandler с использованием mockStorage
	storage := &mockStorage{
		data: make(map[string]string),
	}
	urlHandler := &URLHandler{
		storage: storage.InMemory,
	}

	// Добавляем данные в хранилище
	storage.Put("encoded_short", "http://example.com")

	// Создаем тестовый HTTP-запрос с корректным шорткодом
	req, err := http.NewRequest("GET", "/encoded_short", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Вызываем обработчик Original
	handler := http.HandlerFunc(urlHandler.Original)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Ожидаемый статус-код %v, получен %v", http.StatusTemporaryRedirect, status)
	}

	// Проверяем заголовок Location
	expectedLocation := "http://example.com"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("Ожидаемый заголовок Location %q, получен %q", expectedLocation, location)
	}

	// Проверяем вывод в консоль
	expectedOutput := "http://example.com"
	actualOutput := fmt.Sprint(rr.Body)
	if actualOutput != expectedOutput {
		t.Errorf("Ожидаемый вывод %q, получен %q", expectedOutput, actualOutput)
	}

	// Создаем тестовый HTTP-запрос с некорректным шорткодом
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Вызываем об работчик Original
	handler = http.HandlerFunc(urlHandler.Original)
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидаемый статус-код %v, получен %v", http.StatusBadRequest, status)
	}
}
