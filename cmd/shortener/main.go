package main

import (
	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/router"
)

func main() {
	config.InitOptions()
	router.StartWebServer(config.Options.RunAddress, config.Options.RedirectPrefix)
}
