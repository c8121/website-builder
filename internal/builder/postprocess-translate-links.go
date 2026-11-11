package builder

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/website-builder/internal/entity"
	"golang.org/x/net/html"
)

// PostprocessorTranslateLinks translate absolute path to paths within output-base
type PostprocessorTranslateLinks struct {
	Config *BuildConfig
}

func NewPostprocessorTranslateLinks(config *BuildConfig) *PostprocessorTranslateLinks {
	return &PostprocessorTranslateLinks{
		Config: config,
	}
}

func (p *PostprocessorTranslateLinks) Process(page *entity.Page, doc *goquery.Document) error {
	return traverseDocumentNodes(doc, func(n *html.Node) error {
		for i, attr := range n.Attr {

			key := strings.ToLower(attr.Key)
			switch key {
			case "src", "href":
				n.Attr[i].Val = p.parseUrl(n.Attr[i].Val)
			}
		}
		return nil
	})
}

func (p *PostprocessorTranslateLinks) parseUrl(url string) string {
	if url == "" {
		return ""
	}

	parsed := url

	if strings.HasPrefix(parsed, p.Config.OutputBaseDir) {
		parsed = parsed[len(p.Config.OutputBaseDir):]
	}

	return parsed
}
