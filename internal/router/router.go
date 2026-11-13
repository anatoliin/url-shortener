package router

import (
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/handler"
	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Closable interface {
	Close() error
}

func NewRouter(redirectURL *url.URL, path string, logger zerolog.Logger) (http.Handler, Closable, error) {
	repo, err := repository.NewURLFileRepository(path)
	if err != nil {
		return nil, nil, err
	}
	urlService := service.NewURLService(repo)
	r := chi.NewRouter()
	r.Use(LoggingMiddleware(logger), GzipMiddleware)
	r.Method(http.MethodPost, "/", handler.NewURLPostHandler(urlService, redirectURL))
	r.Method(http.MethodPost, "/api/shorten", handler.NewURLShortenHandle(urlService, redirectURL))
	r.Method(http.MethodGet, "/{id}", handler.NewURLGetHandler(urlService))
	return r, repo, nil
}
