package repository

import (
	"errors"
	"sync"
)

type UrlRepository struct {
	mu      sync.RWMutex
	storage map[string]string
}

func NewUrlRepository() *UrlRepository {
	return &UrlRepository{storage: make(map[string]string)}
}

func (h *UrlRepository) FindById(id string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if val, ok := h.storage[id]; ok {
		return val, nil
	}
	return "", errors.New("url not found")
}

func (h *UrlRepository) Save(id, url string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.storage[id]; exists {
		return errors.New("code collision")
	}
	h.storage[id] = url
	return nil
}
