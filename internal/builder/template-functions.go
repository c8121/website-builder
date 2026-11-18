package builder

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/c8121/website-builder/internal/builder/converter"
)

type TemplateFunctions struct {
	Config      *BuildConfig
	FunctionMap map[string]any
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func NewTemplateFunctions(config *BuildConfig) *TemplateFunctions {

	return &TemplateFunctions{

		Config: config,

		FunctionMap: map[string]any{
			//not needed for src-attributes, because these will be parsed by DocumentParser
			"asset": func(url string) string {
				path, err := config.AssetLoader.CopyAsset(url)
				if err != nil {
					slog.Warn("Asset download failed", "error", err)
					return url
				}
				if strings.HasPrefix(path, config.OutputBaseDir) {
					path = path[len(config.OutputBaseDir):]
				}
				return path
			},

			"replace": func(s string, find string, replace string) string {
				return strings.Replace(s, find, replace, -1)
			},

			"md2html": func(s string) string {
				md2html := converter.NewMarkdown2Html()
				html, err := md2html.Convert([]byte(s))
				if err != nil {
					slog.Warn("Convert failed", "error", err)
				}
				return string(html)
			},
		},
	}
}

func (f *TemplateFunctions) Invoke(expression string, subject any) (string, error) {

	fn, err := NewFunctionExpression(expression, subject)
	if err != nil {
		return "", err
	}
	result, err := fn.InvokeFunctionExpression(f.FunctionMap)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", result), nil
}
