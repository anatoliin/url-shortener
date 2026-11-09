package repository

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/anatoliin/url-shortener/internal/model"
)

type URLFileRepository struct {
	mu       sync.RWMutex
	storage  map[string]string
	filename string
}

func NewURLFileRepository(filename string) (*URLFileRepository, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var dataItems []model.URL
	if len(data) > 0 {
		if err := json.Unmarshal(data, &dataItems); err != nil {
			return nil, err
		}
	}

	storage := make(map[string]string, len(dataItems))
	for _, item := range dataItems {
		storage[item.ShortUrl] = item.OriginalURL
	}

	return &URLFileRepository{storage: storage, filename: filename}, nil
}

func (h *URLFileRepository) FindByID(id string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if val, ok := h.storage[id]; ok {
		return val, nil
	}
	return "", ErrNotFound
}

func (h *URLFileRepository) Save(id, url string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.storage[id]; exists {
		return ErrCollision
	}
	h.storage[id] = url
	return nil
}

func (h *URLFileRepository) Close() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	dataItems := make([]model.URL, 0, len(h.storage))
	for id, originalURL := range h.storage {
		dataItems = append(dataItems, model.URL{
			ShortUrl:    id,
			OriginalURL: originalURL,
		})
	}
	for i := range dataItems {
		dataItems[i].ID = strconv.Itoa(i + 1)
	}

	data, err := json.Marshal(dataItems)
	if err != nil {
		return err
	}
	return os.WriteFile(h.filename, data, 0600)
}
