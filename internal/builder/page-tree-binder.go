package builder

import (
	"context"
	"io"

	"github.com/titpetric/vuego"
)

type PageTreeBinder struct {
	vueRenderer vuego.Template
}

func NewPageTreeBinder() *PageTreeBinder {
	vueRenderer := vuego.New()

	return &PageTreeBinder{
		vueRenderer: vueRenderer,
	}
}

func (p *PageTreeBinder) Bind(html string, writer io.Writer) error {

	pageTree, err := CreatePageTree()
	if err != nil {
		return err
	}

	err = p.vueRenderer.New().Fill(pageTree).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
