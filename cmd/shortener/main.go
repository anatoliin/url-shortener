package main

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/router"
)

func main() {
	options := config.InitOptions()
	r := router.NewRouter(options.RedirectPrefix.URL())

	err := http.ListenAndServe(options.RunAddress, r)
	if err != nil {
		panic(err)
	}
}
