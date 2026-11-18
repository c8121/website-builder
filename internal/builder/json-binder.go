package builder

import (
	"context"
	"encoding/json"
	"io"

	"github.com/c8121/website-builder/internal/entity"
	"github.com/titpetric/vuego"
)

// JsonBinder takes html and json (each as string) and inserts json into html (using vuego)
type JsonBinder struct {
	vueRenderer vuego.Template
}

func NewJsonBinder(templateFunctions *TemplateFunctions) *JsonBinder {

	vueRenderer := vuego.New(
		vuego.WithFuncs(templateFunctions.FunctionMap))

	return &JsonBinder{
		vueRenderer: vueRenderer,
	}

}

type BindContent struct {
	Page    *entity.Page
	Meta    *entity.PageContent
	Content map[string]any
}

func (b *JsonBinder) Bind(html string, page *entity.Page, content *entity.PageContent, writer io.Writer) error {

	var data map[string]any
	if content.Content != "" {
		err := json.Unmarshal([]byte(content.Content), &data)
		if err != nil {
			return err
		}
	}

	bindContent := &BindContent{
		Page:    page,
		Meta:    content,
		Content: data,
	}

	err := b.vueRenderer.New().Fill(bindContent).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}

func (b *JsonBinder) BindData(html string, data any, writer io.Writer) error {

	err := b.vueRenderer.New().Fill(data).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
