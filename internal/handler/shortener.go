package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
)

type URLShortenHandler struct {
	urlService  *service.URLService
	redirectURL *url.URL
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func NewURLShortenHandle(service *service.URLService, redirectAddress *url.URL) *URLShortenHandler {
	return &URLShortenHandler{service, redirectAddress}
}

const maxBodySize = 2048

func shortenURL(urlService *service.URLService, redirectURL *url.URL, rawURL string) (string, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", ErrInvalidURL
	}

	shortURL, err := urlService.GetShortURL(rawURL)
	if err != nil {
		return "", err
	}

	return redirectURL.JoinPath(shortURL).String(), nil
}

func (h *URLShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil || len(body) > maxBodySize {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	result, err := shortenURL(h.urlService, h.redirectURL, req.URL)
	if err != nil {
		if errors.Is(err, repository.ErrCollision) || errors.Is(err, ErrInvalidURL) {
			http.Error(w, "", http.StatusBadRequest)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	respObject := ShortenResponse{
		Result: result,
	}

	resp, err := json.Marshal(respObject)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}
