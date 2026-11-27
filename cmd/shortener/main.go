package main

import (
	"context"
	"net/http"
	"os"

	"github.com/anatoliin/url-shortener/internal/config"
	"github.com/anatoliin/url-shortener/internal/config/db"
	"github.com/anatoliin/url-shortener/internal/router"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	options := config.InitOptions(logger)

	conn, err := db.NewDBConnection(options.DatabaseDSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer conn.Close(context.Background())

	r, err := router.NewRouter(options.RedirectPrefix.URL(), options.FileStoragePath, conn, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create router")
	}

	err = http.ListenAndServe(options.RunAddress, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to start server")
	}
}
