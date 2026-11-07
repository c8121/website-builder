package builder

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/titpetric/vuego"
)

var (
	functionMap = vuego.FuncMap{
		"asset": func(url string) string {
			path, err := InsertAsset(url)
			if err != nil {
				slog.Warn("Asset download failed", "error", err)
				return url
			}
			return path
		},
	}

	renderer = vuego.New(vuego.WithFuncs(functionMap))
)

func Bind(html string, jsonData string, writer io.Writer) error {

	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return err
	}

	err = renderer.New().Fill(data).RenderString(context.Background(), writer, html)
	if err != nil {
		return err
	}

	return nil
}
