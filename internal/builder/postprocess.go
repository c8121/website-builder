package builder

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type PostProcessor interface {
	Process(doc *goquery.Document) error
}

func traverseDocumentNodes(doc *goquery.Document, parseNode func(n *html.Node) error) error {

	if doc == nil {
		return errors.New("document is nil")
	}
	if len(doc.Nodes) == 0 {
		return errors.New("document is empty")
	}

	return traverseNodes(doc.Nodes[0], parseNode)
}

func traverseNodes(n *html.Node, parseNode func(n *html.Node) error) error {

	err := parseNode(n)
	if err != nil {
		return err
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = traverseNodes(c, parseNode)
		if err != nil {
			return err
		}
	}

	return nil
}
