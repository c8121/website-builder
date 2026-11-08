package builder

import (
	"context"
	"io"

	"github.com/titpetric/vuego"
)

type PageTreeBinder struct {
	Config      *BuildConfig
	vueRenderer vuego.Template
}

func NewPageTreeBinder(config *BuildConfig) *PageTreeBinder {
	vueRenderer := vuego.New()

	return &PageTreeBinder{
		Config:      config,
		vueRenderer: vueRenderer,
	}
}

func (p *PageTreeBinder) Bind(html string, writer io.Writer) error {

	pageTree, err := NewPageTree(p.Config).CreatePageTree()
	if err != nil {
		return err
	}

	err = p.vueRenderer.New().Fill(pageTree).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
