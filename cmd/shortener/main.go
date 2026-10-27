package main

import (
	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/router"
)

func main() {
	options := config.InitOptions()
	router.StartWebServer(options.RunAddress, options.RedirectPrefix.URL())
}
