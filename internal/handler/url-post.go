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

const maxBodySize = 2048

func (h *URLPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil || len(body) > maxBodySize {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rawURL := string(body)

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	shorURL, err := h.urlService.GetShorURL(rawURL)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	result := h.redirectURL.JoinPath(shorURL)
	_, _ = w.Write([]byte(result.String()))
}
