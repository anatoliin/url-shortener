package config

import (
	"flag"
	"fmt"
	"net/url"
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

func (b *baseURL) URL() *url.URL {
	return (*url.URL)(b)
}

type Options struct {
	RunAddress     string
	RedirectPrefix baseURL
}

func InitOptions() Options {
	options := new(Options)
	flag.StringVar(&options.RunAddress, "a", ":8080", "address and port to run server")

	if err := options.RedirectPrefix.Set(defaultBaseURL); err != nil {
		panic(err)
	}
	flag.Var(&options.RedirectPrefix, "b", "base address for short URLs")
	flag.Parse()
	return *options
}
