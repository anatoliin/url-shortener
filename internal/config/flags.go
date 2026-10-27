package config

import (
	"flag"
	"fmt"
	"net/url"
)

type Options struct {
	RunAddress     string
	RedirectPrefix url.URL
}

func (o *Options) String() string {
	return fmt.Sprint(o.RedirectPrefix.String())
}

func (o *Options) Set(flagValue string) error {
	u, err := url.Parse(flagValue)
	if err != nil {
		return err
	}
	o.RedirectPrefix = *u
	return nil
}

func InitOptions() Options {
	options := new(Options)
	flag.StringVar(&options.RunAddress, "a", ":8080", "address and port to run server")
	flag.Var(options, "b", "address and port to redirect to")
	flag.Parse()
	return *options
}
