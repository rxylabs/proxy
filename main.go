package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var config struct {
	googleToken string
	target      string
	port        string
}

func init() {
	flag.StringVar(&config.googleToken, "googleToken", "", "Google token")
	flag.StringVar(&config.target, "target", "http://www.apple.com", "Site to proxy")
	flag.StringVar(&config.port, "port", "5000", "Port to use")
	flag.Parse()

	if config.target == "" {
		log.Panic("target is required")
	}

	if config.port == "" {
		log.Panic("port is required")
	}
}

func main() {
	targetURL, err := url.Parse(config.target)
	if err != nil {
		log.Panic("Error parsing target: ", config.target, ": ", err)
	}

	jason := NewTransporter(targetURL)

	http.ListenAndServe(fmt.Sprintf(":%s", config.port), jason)
}
