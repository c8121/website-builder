package builder

import (
	"fmt"
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

	p.doc.Find(".page-content").Each(func(i int, s *goquery.Selection) {

		slog.Info("Insert", "page", page.Id)

		s.Empty()
		s.AppendHtml(fmt.Sprintf("<p>PAGE %d, %d content items</p>", page.Id, len(contents)))
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
