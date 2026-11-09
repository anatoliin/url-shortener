package repository

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/anatoliin/url-shortener/internal/model"
)

type URLFileRepository struct {
	*URLMemoryRepository
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
		storage[item.ShortURL] = item.OriginalURL
	}

	return &URLFileRepository{
		URLMemoryRepository: NewURLMemoryRepository(storage),
		filename:            filename,
	}, nil
}

func (h *URLFileRepository) Close() error {
	h.URLMemoryRepository.mu.RLock()
	defer h.URLMemoryRepository.mu.RUnlock()

	dataItems := make([]model.URL, 0, len(h.URLMemoryRepository.storage))
	i := 1
	for shortURL, originalURL := range h.URLMemoryRepository.storage {
		dataItems = append(dataItems, model.URL{
			ID:          strconv.Itoa(i),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		})
		i++
	}

	data, err := json.Marshal(dataItems)
	if err != nil {
		return err
	}
	return os.WriteFile(h.filename, data, 0600)
}
