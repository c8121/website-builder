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

func NewTemplateParser(doc *goquery.Document, config *BuildConfig) (*TemplateParser, error) {

	return &TemplateParser{
		doc:    doc,
		Config: config,
	}, nil
}

func NewTemplateParserFromFile(path string, config *BuildConfig) (*TemplateParser, error) {

	root, err := loadAndPreprocessFile(path, config)
	if err != nil {
		return nil, err
	}

	return NewTemplateParser(
		goquery.NewDocumentFromNode(root),
		config,
	)
}

func loadAndPreprocessFile(path string, config *BuildConfig) (*html.Node, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content, err = preprocessTemplate(content, config)

	root, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return root, nil
}

func preprocessTemplate(content []byte, config *BuildConfig) ([]byte, error) {

	siteConfig, err := entity.ListSiteConfig(entity.NewSession())
	if err != nil {
		return nil, err
	}
	configMap := entity.SiteConfigAsMap(siteConfig)

	processed := ReplacePlaceholders(string(content), configMap)
	processed = ReplacePhpBlocks(processed)

	return []byte(processed), nil
}

func (p *TemplateParser) InsertContent(page *entity.Page) error {

	contents, err := entity.ListPageContents(entity.NewSession(), page.Id)
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
			err = p.Config.JsonBinder.Bind(contentTemplate, page, &content, parsed)
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

func (p *TemplateParser) InsertMainNav(page *entity.Page) error {

	p.doc.Find(".main-nav").Each(func(i int, s *goquery.Selection) {

		template, err := s.Html()
		if err == nil {
			slog.Info("Insert main-nav")

			s.Empty()

			pageTreeBinder := NewPageTreeBinder(page, p.Config)

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

func (p *TemplateParser) PostProcess(page *entity.Page) error {

	postProcessors := &[]PostProcessor{
		NewPostprocessorInsertAssets(p.Config),
		NewPostprocessorInsertPageLinks(p.Config),
		NewPostprocessorTranslateLinks(p.Config),
		NewPostprocessorCreateToc(p.Config),
	}

	for _, postProcessor := range *postProcessors {
		err := postProcessor.Process(page, p.doc)
		if err != nil {
			slog.Error("Failed to parse document", "error", err)
		}
	}

	return nil
}

func (p *TemplateParser) Save(path string) error {

	htmlContent, err := p.doc.Html()
	if err != nil {
		return err
	}

	htmlContent = RestorePhpBlocks(htmlContent)

	err = os.WriteFile(path, []byte(htmlContent), p.Config.FileMode)
	if err != nil {
		return err
	}

	return nil
}
