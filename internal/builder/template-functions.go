package builder

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/c8121/website-builder/internal/builder/converter"
)

type TemplateFunctions struct {
	Config      *BuildConfig
	FunctionMap map[string]any
}

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

func (f *TemplateFunctions) Invoke(functionName string, arg string) (string, error) { //TODO variable args

	slog.Info("** Invoke", "function", functionName, "arg", arg)

	function, exists := f.FunctionMap[functionName]
	if !exists {
		return "", errors.New("function not found: " + functionName)
	}

	fn, ok := function.(func(string) string)
	if !ok {
		return "", errors.New("Invalid function signature: " + functionName)
	}

	return fn(arg), nil
}
