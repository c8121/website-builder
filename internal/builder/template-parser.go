package builder

import (
	"bytes"
	"log/slog"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/c8121/website-builder/internal/entity"
	"golang.org/x/net/html"
)

const (
	SlotCssClassPrefix = "slot-"
	DefaultSlotName    = "default"

	PageContentSelector           = ".page-content"
	MainNavSelector               = ".main-nav"
	ContentTemplateSelectorPrefix = ".tpl."
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

func NewTemplateParserFromFile(template []byte, config *BuildConfig) (*TemplateParser, error) {

	root, err := html.Parse(bytes.NewReader(template))
	if err != nil {
		return nil, err
	}

	return NewTemplateParser(
		goquery.NewDocumentFromNode(root),
		config,
	)
}

func (p *TemplateParser) InsertContent(page *entity.Page) error {

	contents, err := entity.ListPageContents(entity.NewSession(), page.Id)
	if err != nil {
		return err
	}

	parsedContent := make(map[string][]string)

	for _, content := range contents {

		slot := strings.TrimSpace(content.Slot)
		if slot == "" {
			slot = DefaultSlotName
		}

		slog.Info("Create content", "id", content.Id, "name", content.Name, "type", content.ContentType, "slot", slot)

		_, contentTemplate := selectHtml(p.doc, ContentTemplateSelectorPrefix+content.ContentType)

		if contentTemplate == "" {
			slog.Warn("No template found", "id", content.Id, "name", content.Name, "type", content.ContentType)
			parsedContent[slot] = append(parsedContent[slot], "<pre>"+content.Content+"</pre>")
		} else {
			parsed := &bytes.Buffer{}
			err = p.Config.JsonBinder.Bind(contentTemplate, page, &content, parsed)
			if err != nil {
				slog.Warn("Failed to bind", "error", err, "id", content.Id, "name", content.Name, "type", content.ContentType)
				parsedContent[slot] = append(parsedContent[slot], "<pre>"+err.Error()+"</pre><pre>"+content.Content+"</pre>")
			} else {
				parsedContent[slot] = append(parsedContent[slot], parsed.String())
			}
		}
	}

	slog.Info("Insert", "page", page.Id)

	p.doc.Find(PageContentSelector).Each(func(i int, s *goquery.Selection) {

		s.Empty()

		css, _ := s.Attr("class")
		slot := p.extractSlotName(css)
		if slot == "" {
			slot = DefaultSlotName
		}

		list, ok := parsedContent[slot]
		if !ok {
			slog.Warn("No content items.", "slot", slot)
			return
		}

		slog.Info("Insert content", "slot", slot, "items", len(list))

		for _, content := range list {
			s.AppendHtml("\n<!-- PageContent -->\n")
			s.AppendHtml(content)
		}
	})

	return nil
}

func (p *TemplateParser) extractSlotName(cssClasses string) string {

	i := strings.Index(cssClasses, SlotCssClassPrefix)
	if i < 0 {
		return ""
	}

	slot := cssClasses[i+len(SlotCssClassPrefix):]
	i = strings.IndexAny(slot, " '\"\t")
	if i > 0 {
		slot = slot[:i]
	}
	return slot
}

func (p *TemplateParser) InsertMainNav(page *entity.Page) error {

	p.doc.Find(MainNavSelector).Each(func(i int, s *goquery.Selection) {

		template, err := s.Html()
		if err == nil {
			slog.Info("Insert main-nav")

			s.Empty()

			pageTreeBinder := NewPageTreeBinder(page, p.Config)

			parsed := &bytes.Buffer{}
			err := pageTreeBinder.Bind(template, parsed, entity.VisibilityDefault)
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

	htmlContent = p.Config.Placeholder.RestorePhpBlocks(htmlContent)

	err = os.WriteFile(path, []byte(htmlContent), p.Config.FileMode)
	if err != nil {
		return err
	}

	return nil
}
