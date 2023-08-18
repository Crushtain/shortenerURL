package storage

import "testing"

func TestInMemory(t *testing.T) {
	inMemory := NewInMemory()

	// Проверяем что значение не найдено
	result := inMemory.Get("nonexistent")
	if result != "No key" {
		t.Errorf("Expected 'No key' but got %q", result)
	}

	// Проверяем добавление и получение значения
	inMemory.Put("key", "value")
	result = inMemory.Get("key")
	if result != "value" {
		t.Errorf("Expected 'value' but got %q", result)
	}
}
