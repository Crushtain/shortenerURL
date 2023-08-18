package storage

import (
	"sync"
)

type InMemory struct {
	mu      sync.RWMutex
	storage map[string]string
}

func NewInMemory() *InMemory {
	return &InMemory{
		storage: make(map[string]string),
	}
}

func (m *InMemory) Put(short string, body string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[short] = body
}

func (m *InMemory) Get(short string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	original, ok := m.storage[short]

	if !ok {
		return "No key"
	}
	return original
}
