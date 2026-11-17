package builder

import (
	"log/slog"
	"os"
)

type TemplatePreprocessor struct {
	Config    *BuildConfig
	Templates map[string][]byte
}

func NewTemplatePreprocessor(config *BuildConfig) (*TemplatePreprocessor, error) {

	return &TemplatePreprocessor{
		Config:    config,
		Templates: make(map[string][]byte),
	}, nil
}

func (pt *TemplatePreprocessor) Preprocess(path string) ([]byte, error) {

	template, exist := pt.Templates[path]
	if exist {
		slog.Info("Found preprocessed template", "path", path)
		return template, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	processed := ReplacePlaceholders(string(content), pt.Config.SiteConfig)
	processed = ReplacePhpBlocks(processed)
	processedBytes := []byte(processed)

	pt.Templates[path] = processedBytes

	return processedBytes, nil
}
