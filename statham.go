package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/rxylabs/proxy/transformers"

	"golang.org/x/net/html"
)

type Transporter struct {
	Target       *url.URL
	Transformers []transformers.Transformer
	rp           *httputil.ReverseProxy
}

func NewTransporter(target *url.URL) *Transporter {
	jason := &Transporter{
		Target: target,
		rp:     httputil.NewSingleHostReverseProxy(target),
		Transformers: []transformers.Transformer{
			&transformers.LinkUpdater{},
		},
	}

	jason.rp.Transport = jason

	return jason
}

func (t *Transporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.fixHeaders(r)
	t.rp.ServeHTTP(w, r)
}

func (t *Transporter) fixHeaders(r *http.Request) {
	r.URL.Scheme = t.Target.Scheme
	r.URL.Host = t.Target.Host
	r.Host = t.Target.Host

	r.Header.Del("Accept-Encoding")
	r.Header.Del("Accept-Language")

	if false {
		req, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.Panic("Error dumping request:", err)
		}
		log.Print(string(req))
	}
}

func (t *Transporter) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		log.Println("Error fetching site:", err)
	}

	if htmlResponse(*response) {
		doc, err := html.Parse(response.Body)
		if err != nil {
			log.Println("Error reading response body", err)
			return nil, err
		}

		for _, transformer := range t.Transformers {
			doc, err = transformer.Transform(doc)
			if err != nil {
				log.Println("Error with transformer:", reflect.TypeOf(transformer), ": ", err)
				return nil, err
			}
		}

		finalBuffer := bytes.NewBuffer(nil)
		err = html.Render(finalBuffer, doc)
		if err != nil {
			log.Println("Error rendering html:", err)
			return nil, err
		}

		final := finalBuffer.Bytes()

		body := ioutil.NopCloser(bytes.NewReader(final))
		response.Body = body
		response.ContentLength = int64(len(final))
		response.Header.Set("Content-Length", strconv.Itoa(len(final)))
	}

	return response, nil
}

func htmlResponse(r http.Response) bool {
	if !strings.Contains(r.Header.Get("Content-Type"), "html") {
		return false
	}

	return true
}
