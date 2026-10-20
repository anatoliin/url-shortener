package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/service"
)

type UrlPostHandler struct {
	urlService *service.UrlService
}

func NewUrlPostHandler(service *service.UrlService) *UrlPostHandler {
	return &UrlPostHandler{service}
}

func (h *UrlPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	shorUrl, err := h.urlService.GetShorUrl(rawURL)
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
