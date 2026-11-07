package builder

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

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
	}

	vueRenderer := vuego.New(
		vuego.WithFuncs(functionMap))

	return &JsonBinder{
		assetLoader: assetLoader,
		vueRenderer: vueRenderer,
	}

}

func (b *JsonBinder) Bind(html string, jsonData string, writer io.Writer) error {

	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return err
	}

	err = b.vueRenderer.New().Fill(data).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
