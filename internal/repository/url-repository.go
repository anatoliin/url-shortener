package repository

import (
	"errors"
	"sync"
)

type URLRepository struct {
	mu      sync.RWMutex
	storage map[string]string
}

func NewURLRepository() *URLRepository {
	return &URLRepository{storage: make(map[string]string)}
}

func (h *URLRepository) FindByID(id string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if val, ok := h.storage[id]; ok {
		return val, nil
	}
	return "", errors.New("url not found")
}

func (h *URLRepository) Save(id, url string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.storage[id]; exists {
		return errors.New("code collision")
	}
	h.storage[id] = url
	return nil
}
