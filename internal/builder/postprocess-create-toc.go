package builder

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/website-builder/internal/entity"
	"golang.org/x/net/html"
)

// PostprocessorCreateToc finds headings (h1, h2,...) and creates TOC
type PostprocessorCreateToc struct {
	Config *BuildConfig
}

func NewPostprocessorCreateToc(config *BuildConfig) *PostprocessorCreateToc {
	return &PostprocessorCreateToc{
		Config: config,
	}
}

type tocItem struct {
	Page *entity.Page

	Number  int
	Level   int
	Name    string
	Caption string

	Parent   *tocItem
	Children []*tocItem
}

func newTocItem(page *entity.Page, n *html.Node) *tocItem {

	name := n.Data
	if len(name) != 2 || (!strings.HasPrefix(name, "h") && !strings.HasPrefix(name, "H")) {
		return nil
	}

	level, err := strconv.Atoi(strings.TrimSpace(name[1:]))
	if err != nil || level <= 0 {
		//Not a h-tag
		return nil
	}

	text := textContent(n)
	nodeName := NormalizeAsFileName(text, "")

	return &tocItem{
		Page:     page,
		Level:    level,
		Name:     nodeName,
		Caption:  text,
		Children: []*tocItem{},
	}
}

func (p *PostprocessorCreateToc) Process(page *entity.Page, doc *goquery.Document) error {

	node, template := selectHtml(doc, ".toc")
	if node == nil || template == "" {
		return nil
	}

	root := &tocItem{
		Page:     page,
		Children: []*tocItem{},
	}
	lastItem := root
	num := 0

	err := traverseDocumentNodes(doc, func(n *html.Node) error {

		if n.Type == html.ElementNode {
			item := newTocItem(page, n)
			if item != nil {

				if item.Level < lastItem.Level {
					if lastItem.Parent != nil {
						lastItem = lastItem.Parent
					}
				}

				if item.Level > lastItem.Level {
					item.Parent = lastItem
					lastItem.Children = append(lastItem.Children, item)
				} else if lastItem.Parent != nil {
					item.Parent = lastItem.Parent
					lastItem.Parent.Children = append(lastItem.Parent.Children, item)
				} else {
					//invalid structure
					slog.Warn("Invalid structure")
				}

				num = num + 1
				item.Number = num

				if item.Name == "" {
					item.Name = fmt.Sprintf("toc-%d", num)
				}
				n.Attr = append(n.Attr, html.Attribute{
					Key: "id",
					Val: item.Name,
				})

				lastItem = item
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	parsed := &bytes.Buffer{}
	err = p.Config.JsonBinder.BindData(template, root, parsed)
	if err != nil {
		slog.Warn("Failed to bind toc", "error", err)
	} else {
		node.Empty()
		node.AppendHtml("\n<!-- TOC -->\n")
		node.AppendHtml(parsed.String())
	}
	return nil
}
