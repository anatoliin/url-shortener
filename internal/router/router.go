package router

import (
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/handler"
	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func NewRouter(redirectURL *url.URL) http.Handler {
	repo := repository.NewURLMemoryRepository(make(map[string]string))
	urlService := service.NewURLService(repo)
	r := chi.NewRouter()
	r.Use(Logger(log.Logger), GzipMiddleware)
	r.Method(http.MethodPost, "/", handler.NewURLPostHandler(urlService, redirectURL))
	r.Method(http.MethodPost, "/api/shorten", handler.NewURLShortenHandle(urlService, redirectURL))
	r.Method(http.MethodGet, "/{id}", handler.NewURLGetHandler(urlService))
	return r
}
