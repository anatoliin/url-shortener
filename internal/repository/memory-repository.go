package repository

import (
	"sync"
)

type URLMemoryRepository struct {
	mu      sync.RWMutex
	storage map[string]string
}

func NewURLMemoryRepository(storage map[string]string) *URLMemoryRepository {
	return &URLMemoryRepository{storage: storage}
}

func (h *URLMemoryRepository) FindByID(id string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if val, ok := h.storage[id]; ok {
		return val, nil
	}
	return "", ErrNotFound
}

func (h *URLMemoryRepository) Save(id, url string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.storage[id]; exists {
		return ErrCollision
	}
	h.storage[id] = url
	return nil
}
