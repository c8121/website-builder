package builder

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/website-builder/internal/entity"
	"golang.org/x/net/html"
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

	root, err := loadAndPreprocessFile(path, config)
	if err != nil {
		return nil, err
	}

	p := &TemplateParser{
		doc:    goquery.NewDocumentFromNode(root),
		Config: config,
	}
	return p, nil
}

func loadAndPreprocessFile(path string, config *BuildConfig) (*html.Node, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content, err = insertSiteConfig(content, config)

	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return root, nil
}

// insertSiteConfig is a preparse step currently: could not figure how to parse goquery.Document or html.Node with vuego
func insertSiteConfig(content []byte, config *BuildConfig) ([]byte, error) {

	siteConfig, err := entity.ListSiteConfig()
	if err != nil {
		return nil, err
	}
	configMap := entity.SiteConfigAsMap(siteConfig)

	return []byte(ReplacePlaceholders(string(content), configMap)), nil
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

	p.doc.Find(".main-nav").Each(func(i int, s *goquery.Selection) {

		template, err := s.Html()
		if err == nil {
			slog.Info("Insert main-nav")

			s.Empty()

			pageTreeBinder := NewPageTreeBinder(p.Config)

			parsed := &bytes.Buffer{}
			err := pageTreeBinder.Bind(template, parsed)
			if err != nil {
				slog.Warn("Failed to bind", "error", err)
			} else {
				s.AppendHtml("\n<!-- MainNav -->\n")
				s.AppendHtml(parsed.String())
			}

		} else {
			slog.Warn("Failed to get main-nav template", "error", err)
		}

	})

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
