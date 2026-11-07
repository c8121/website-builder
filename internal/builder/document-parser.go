package builder

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/website-builder/internal/entity"
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
		case "src", "href":
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

	parsed := url

	if strings.HasPrefix(url, "page://") {
		parsed = p.createPageLink(parsed)
	}

	if strings.HasPrefix(parsed, p.Config.OutputBaseDir) {
		parsed = parsed[len(p.Config.OutputBaseDir):]
	}

	return parsed
}

// createPageLink transforms page://<ID> into a valid href
func (p *DocumentParser) createPageLink(url string) string {
	pos := strings.LastIndex(url, "/")
	if pos < 0 {
		slog.Warn("Failed to create page link", "url", url)
		return url
	}

	id, err := strconv.Atoi(url[pos+1:])
	if err != nil {
		slog.Warn("Failed to create page link", "url", url)
		return url
	}

	page := &entity.Page{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(int64(id)))
	if err != nil {
		slog.Warn("Failed to load page", "url", url, "err", err)
		return url
	}

	path, err := GetPageDirectoryPath(page)
	if err != nil {
		slog.Warn("Failed to create page link", "url", url, "err", err)
		return url
	}

	return path
}
