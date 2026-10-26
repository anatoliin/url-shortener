package config

import (
	"flag"
)

var Options struct {
	RunAddress     string
	RedirectPrefix string
}

func InitOptions() {
	flag.StringVar(&Options.RunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.RedirectPrefix, "b", "http://localhost:8080", "address and port to redirect to")
	flag.Parse()
}
