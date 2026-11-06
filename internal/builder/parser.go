package builder

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
)

type Parser struct {
	doc *goquery.Document
}

func NewParser(doc *goquery.Document) *Parser {
	return &Parser{
		doc: doc,
	}
}

func NewParserFromFile(path string) (*Parser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}
	return NewParser(doc), nil
}

func (p *Parser) Parse(page entity.Page, contents []entity.PageContent) error {

	contents, err := entity.ListPageContents(page.Id)
	if err != nil {
		return err
	}

	parsedContent := make([]string, 0)

	for _, content := range contents {
		slog.Info("Insert content", "id", content.Id, "name", content.Name, "type", content.ContentType)

		contentTemplateSelector := ".tpl." + content.ContentType
		contentTemplate := ""
		var err error
		p.doc.Find(contentTemplateSelector).Each(func(i int, s *goquery.Selection) {
			contentTemplate, err = s.Html()
		})
		if err != nil {
			return err
		}

		if contentTemplate == "" {
			slog.Warn("No template found", "id", content.Id, "name", content.Name, "type", content.ContentType, "selector", contentTemplateSelector)
			parsedContent = append(parsedContent, "<pre>"+content.Content+"</pre>")
		} else {
			parsed := &bytes.Buffer{}
			err = Bind(contentTemplate, content.Content, parsed)
			if err != nil {
				slog.Warn("Failed to bind", "error", err, "id", content.Id, "name", content.Name, "type", content.ContentType)
				parsedContent = append(parsedContent, "<pre>"+err.Error()+"</pre><pre>"+content.Content+"</pre>")
			} else {
				parsedContent = append(parsedContent, parsed.String())
			}
		}
	}

	p.doc.Find(".page-content").Each(func(i int, s *goquery.Selection) {

		slog.Info("Insert", "page", page.Id)
		s.Empty()

		for _, content := range parsedContent {
			s.AppendHtml("\n<!-- PageContent -->\n")
			s.AppendHtml(content)
		}
	})

	return nil
}

func (p *Parser) Save(path string, perm os.FileMode) error {

	html, err := p.doc.Html()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(html), perm)
	if err != nil {
		return err
	}

	return nil
}
