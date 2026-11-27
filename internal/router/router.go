package router

import (
	"net/http"
	"net/url"

	"github.com/anatoliin/url-shortener/internal/handler"
	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func NewRouter(redirectURL *url.URL, path string, conn *pgx.Conn, logger zerolog.Logger) (http.Handler, error) {
	repo, err := repository.NewURLFileRepository(path)
	if err != nil {
		return nil, err
	}
	urlService := service.NewURLService(repo)
	r := chi.NewRouter()
	r.Use(LoggingMiddleware(logger), GzipMiddleware)
	r.Method(http.MethodPost, "/", handler.NewURLPostHandler(urlService, redirectURL))
	r.Method(http.MethodPost, "/api/shorten", handler.NewURLShortenHandle(urlService, redirectURL))
	r.Method(http.MethodGet, "/{id}", handler.NewURLGetHandler(urlService))
	r.Method(http.MethodGet, "/ping", handler.NewPingHandler(conn))
	return r, nil
}
