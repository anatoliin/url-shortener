package router

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/handler"
	"github.com/anatoliin/url-shortener/internal/repository"
	"github.com/anatoliin/url-shortener/internal/service"
)

func StartWebServer() {
	repo := repository.NewUrlRepository()
	urlService := service.NewUrlService(repo)
	mux := http.NewServeMux()
	mux.Handle(`POST /`, handler.NewUrlPostHandler(urlService))
	mux.Handle(`GET /{id}`, handler.NewUrlGetHandler(urlService))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
