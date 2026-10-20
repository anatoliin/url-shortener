package handler

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/service"
)

type UrlGetHandler struct {
	urlService *service.UrlService
}

func NewUrlGetHandler(service *service.UrlService) *UrlGetHandler {
	return &UrlGetHandler{service}
}

func (h *UrlGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	originalUrl, err := h.urlService.GetOriginalUrl(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", originalUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
