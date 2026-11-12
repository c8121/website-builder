package converter

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Markdown2Html struct {
	renderer *html.Renderer
}

func NewMarkdown2Html() *Markdown2Html {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return &Markdown2Html{
		renderer: renderer,
	}
}

func (c *Markdown2Html) Convert(md []byte) ([]byte, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	return markdown.Render(doc, c.renderer), nil
}
