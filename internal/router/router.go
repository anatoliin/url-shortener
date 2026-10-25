package router

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/handler"
	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func StartWebServer() {
	repo := repository.NewURLMemoryRepository(make(map[string]string))
	urlService := service.NewURLService(repo)
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/", handler.NewURLPostHandler(urlService))
	r.Method(http.MethodGet, "/{id}", handler.NewURLGetHandler(urlService))

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
