package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
)

var config struct {
	googleToken string
	target      string
}

func init() {
	flag.StringVar(&config.googleToken, "googleToken", "", "Google token")
	flag.StringVar(&config.target, "target", "http://www.apple.com", "Site to proxy")
	flag.Parse()

	if config.googleToken == nil {
		log.Panic("googleToken is required")
	}

	if config.target == nil {
		log.Panic("target is required")
	}

}

func main() {

}
