package builder

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"

	"github.com/c8121/website-builder/internal/entity"
	"github.com/titpetric/vuego"
)

// JsonBinder takes html and json (each as string) and inserts json into html (using vuego)
type JsonBinder struct {
	assetLoader *AssetLoader
	vueRenderer vuego.Template
}

func NewJsonBinder(assetLoader *AssetLoader) *JsonBinder {

	functionMap := vuego.FuncMap{

		//not needed for src-attributes, because these will be parsed by DocumentParser
		"asset": func(url string) string {
			path, err := assetLoader.CopyAsset(url)
			if err != nil {
				slog.Warn("Asset download failed", "error", err)
				return url
			}
			return path
		},

		"replace": func(s string, find string, replace string) string {
			return strings.Replace(s, find, replace, -1)
		},
	}

	vueRenderer := vuego.New(
		vuego.WithFuncs(functionMap))

	return &JsonBinder{
		assetLoader: assetLoader,
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
