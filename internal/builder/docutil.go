package builder

import (
	"bytes"
	"log/slog"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func selectHtml(doc *goquery.Document, selector string) (*goquery.Selection, string) {

	var selection *goquery.Selection = nil
	htmlContent := ""
	var err error
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		selection = s
		htmlContent, err = s.Html()
	})

	if err != nil {
		slog.Error("Failed to get html from dom", "selector", selector, "error", err)
		return nil, ""
	}

	return selection, htmlContent
}

func textContent(n *html.Node) string {
	var buf bytes.Buffer

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return buf.String()
}
