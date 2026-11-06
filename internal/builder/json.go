package builder

import (
	"context"
	"encoding/json"
	"io"

	"github.com/titpetric/vuego"
)

var (
	functionMap = vuego.FuncMap{
		"asset": func(url string) string {
			return "http://localhost:9998" + url
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
