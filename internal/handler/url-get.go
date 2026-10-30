package handler

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type URLGetHandler struct {
	urlService *service.URLService
}

func NewURLGetHandler(service *service.URLService) *URLGetHandler {
	return &URLGetHandler{service}
}

func (h *URLGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, err := h.urlService.GetOriginalURL(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
