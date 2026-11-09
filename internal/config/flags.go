package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
)

const defaultBaseURL = "http://localhost:8080"

type baseURL url.URL

func (b *baseURL) String() string {
	return (*url.URL)(b).String()
}

func (b *baseURL) Set(flagValue string) error {
	u, err := url.Parse(flagValue)
	if err != nil {
		return err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("invalid base URL: %q", flagValue)
	}
	*b = baseURL(*u)
	return nil
}

func (b *baseURL) UnmarshalText(text []byte) error {
	return b.Set(string(text))
}

func (b *baseURL) URL() *url.URL {
	return (*url.URL)(b)
}

type Options struct {
	RunAddress      string  `env:"SERVER_ADDRESS"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"`
	RedirectPrefix  baseURL `env:"BASE_URL"`
}

func InitOptions() Options {
	options := new(Options)
	flag.StringVar(&options.RunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&options.FileStoragePath, "f", "data.json", "file storage path")

	if err := options.RedirectPrefix.Set(defaultBaseURL); err != nil {
		panic(err)
	}
	flag.Var(&options.RedirectPrefix, "b", "base address for short URLs")
	flag.Parse()

	_ = env.Parse(options)

	return *options
}
