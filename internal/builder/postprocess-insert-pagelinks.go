package builder

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/website-builder/internal/entity"
	"golang.org/x/net/html"
)

// PostprocessorInsertPageLinks insert Page-Links: page://<ID>
type PostprocessorInsertPageLinks struct {
	Config *BuildConfig
}

func NewPostprocessorInsertPageLinks(config *BuildConfig) *PostprocessorInsertPageLinks {
	return &PostprocessorInsertPageLinks{
		Config: config,
	}
}

func (p *PostprocessorInsertPageLinks) Process(page *entity.Page, doc *goquery.Document) error {
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

func (p *PostprocessorInsertPageLinks) parseUrl(url string) string {
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
func (p *PostprocessorInsertPageLinks) createPageLink(url string) string {
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
