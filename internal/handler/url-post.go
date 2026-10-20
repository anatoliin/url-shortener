package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/service"
)

type URLPostHandler struct {
	urlService *service.URLService
}

func NewURLPostHandler(service *service.URLService) *URLPostHandler {
	return &URLPostHandler{service}
}

func (h *URLPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 2048))
	rawURL := string(body)

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	shorUrl, err := h.urlService.GetShorURL(rawURL)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	result := (&url.URL{
		Scheme: "http",
		Host:   r.Host,
	}).JoinPath(shorUrl)
	_, _ = w.Write([]byte(result.String()))
}
