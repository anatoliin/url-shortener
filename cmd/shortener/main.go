package main

import (
	"net/http"
	"os"

	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/router"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	options := config.InitOptions(logger)
	r, repo, err := router.NewRouter(options.RedirectPrefix.URL(), options.FileStoragePath, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create router")
	}
	defer repo.Close()

	err = http.ListenAndServe(options.RunAddress, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start server")
	}
}
