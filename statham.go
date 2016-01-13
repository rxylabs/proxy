package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Transporter struct {
	target *url.URL
	rp     *httputil.ReverseProxy
}

func NewTransporter(target *url.URL) *Transporter {
	jason := &Transporter{
		target: target,
		rp:     httputil.NewSingleHostReverseProxy(target),
	}

	return jason
}

func (t *Transporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.fixHeaders(r)
	t.rp.ServeHTTP(w, r)
}

func (t *Transporter) fixHeaders(r *http.Request) {
	r.URL.Scheme = t.target.Scheme
	r.URL.Host = t.target.Host
	r.Host = t.target.Host

	req, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Panic("Error dumping request:", err)
	}
	log.Print(string(req))
}
