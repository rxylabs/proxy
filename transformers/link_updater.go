package transformers

import (
	"log"

	"golang.org/x/net/html"
)

type LinkUpdater struct {
}

func (lu *LinkUpdater) Transform(doc *html.Node) (*html.Node, error) {

	log.Println("Starting crawl")
	crawl(doc)

	return doc, nil
}

func crawl(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for i, a := range n.Attr {
			if a.Key == "href" {
				n.Attr[i].Val = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		crawl(c)
	}
}
