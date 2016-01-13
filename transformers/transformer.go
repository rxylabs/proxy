package transformers

import (
	"golang.org/x/net/html"
)

type Transformer interface {
	Transform(*html.Node) (*html.Node, error)
}
