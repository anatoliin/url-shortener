package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/service"
)

type URLPostHandler struct {
	urlService  *service.URLService
	redirectURL *url.URL
}

func NewURLPostHandler(service *service.URLService, redirectAddress *url.URL) *URLPostHandler {
	return &URLPostHandler{service, redirectAddress}
}

func (h *URLPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil || len(body) > maxBodySize {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rawURL := string(body)

	result, err := shortenURL(h.urlService, h.redirectURL, rawURL)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(result))
}
