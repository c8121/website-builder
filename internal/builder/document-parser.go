package builder

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// DocumentParser analyses the whole DOM and parses contents (finding assets, parsing paths...)
type DocumentParser struct {
	doc    *goquery.Document
	Config *TemplateConfig
}

func NewDocumentParser(document *goquery.Document, config *TemplateConfig) *DocumentParser {
	return &DocumentParser{
		doc:    document,
		Config: config,
	}
}

func (p *DocumentParser) Parse() error {

	if p.doc == nil {
		return errors.New("document is nil")
	}
	if len(p.doc.Nodes) == 0 {
		return errors.New("document is empty")
	}

	return p.traverseNodes(p.doc.Nodes[0])
}

func (p *DocumentParser) traverseNodes(n *html.Node) error {

	err := p.parseNode(n)
	if err != nil {
		return err
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = p.traverseNodes(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DocumentParser) parseNode(n *html.Node) error {

	for i, attr := range n.Attr {

		key := strings.ToLower(attr.Key)
		switch key {
		case "src":
			if p.Config.AssetLoader.IsAssetUrl(attr.Val) {
				assetSrc, err := p.Config.AssetLoader.CopyAsset(attr.Val)
				if err == nil {
					n.Attr[i].Val = assetSrc
				} else {
					slog.Warn("Failed to asset from URL", "src", attr.Val, "err", err)
				}
			}

			n.Attr[i].Val = p.parseUrl(n.Attr[i].Val)
		}
	}

	return nil
}

func (p *DocumentParser) parseUrl(url string) string {
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, p.Config.OutputBaseDir) {
		return url[len(p.Config.OutputBaseDir):]
	}
	return url
}
