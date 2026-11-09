package builder

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
)

// TemplateParser takes a html-template and inserts content of pages (entity.Page)
type TemplateParser struct {
	doc    *goquery.Document
	Config *BuildConfig
}

func NewTemplateParser(doc *goquery.Document, config *BuildConfig) *TemplateParser {

	return &TemplateParser{
		doc:    doc,
		Config: config,
	}
}

func NewTemplateParserFromFile(path string, config *BuildConfig) (*TemplateParser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}
	return NewTemplateParser(doc, config), nil
}

func (p *TemplateParser) InsertContent(page *entity.Page) error {

	contents, err := entity.ListPageContents(page.Id)
	if err != nil {
		return err
	}

	parsedContent := make([]string, 0)

	for _, content := range contents {
		slog.Info("Insert content", "id", content.Id, "name", content.Name, "type", content.ContentType)

		_, contentTemplate := selectHtml(p.doc, ".tpl."+content.ContentType)

		if contentTemplate == "" {
			slog.Warn("No template found", "id", content.Id, "name", content.Name, "type", content.ContentType)
			parsedContent = append(parsedContent, "<pre>"+content.Content+"</pre>")
		} else {
			parsed := &bytes.Buffer{}
			err = p.Config.JsonBinder.Bind(contentTemplate, content, parsed)
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

func (p *TemplateParser) InsertMainNav() error {

	node, template := selectHtml(p.doc, ".main-nav")
	if node == nil || template == "" {
		return nil
	}

	binder := NewPageTreeBinder(p.Config)

	slog.Info("Insert main nav")
	node.Empty()

	parsed := &bytes.Buffer{}
	err := binder.Bind(template, parsed)
	if err != nil {
		slog.Warn("Failed to bind", "error", err)
	} else {
		node.AppendHtml("\n<!-- MainNav -->\n")
		node.AppendHtml(parsed.String())
	}

	return nil
}

func (p *TemplateParser) PostProcess() error {

	postProcessors := &[]PostProcessor{
		NewPostprocessorInsertAssets(p.Config),
		NewPostprocessorInsertPageLinks(p.Config),
		NewPostprocessorTranslateLinks(p.Config),
		NewPostprocessorCreateToc(p.Config),
	}

	for _, postProcessor := range *postProcessors {
		err := postProcessor.Process(p.doc)
		if err != nil {
			slog.Error("Failed to parse document", "error", err)
		}
	}

	return nil
}

func (p *TemplateParser) Save(path string) error {

	html, err := p.doc.Html()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(html), p.Config.FileMode)
	if err != nil {
		return err
	}

	return nil
}
