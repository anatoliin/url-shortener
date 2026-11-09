package main

import (
	"net/http"

	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/router"
)

func main() {
	options := config.InitOptions()
	r, err := router.NewRouter(options.RedirectPrefix.URL(), options.FileStoragePath)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(options.RunAddress, r)
	if err != nil {
		panic(err)
	}
}
