package builder

import (
	"context"
	"io"

	"github.com/c8121/website-builder/internal/entity"
	"github.com/titpetric/vuego"
)

type PageTreeBinder struct {
	Config      *BuildConfig
	Page        *entity.Page
	vueRenderer vuego.Template
}

func NewPageTreeBinder(page *entity.Page, config *BuildConfig) *PageTreeBinder {
	vueRenderer := vuego.New()

	return &PageTreeBinder{
		Config:      config,
		Page:        page,
		vueRenderer: vueRenderer,
	}
}

func (p *PageTreeBinder) Bind(html string, writer io.Writer) error {

	pageTree, err := NewPageTree(p.Page, p.Config).CreatePageTree()
	if err != nil {
		return err
	}

	err = p.vueRenderer.New().Fill(pageTree).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
