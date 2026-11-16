package builder

import (
	"log/slog"
	"os"

	"github.com/c8121/website-builder/internal/entity"
)

type TemplatePreprocessor struct {
	Config    *BuildConfig
	Templates map[string][]byte

	siteConfig map[string]string
}

func NewTemplatePreprocessor(config *BuildConfig) (*TemplatePreprocessor, error) {

	siteConfig, err := entity.ListSiteConfig(entity.NewSession())
	if err != nil {
		return nil, err
	}
	configMap := entity.SiteConfigAsMap(siteConfig)

	return &TemplatePreprocessor{
		Config:     config,
		Templates:  make(map[string][]byte),
		siteConfig: configMap,
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

	processed := ReplacePlaceholders(string(content), pt.siteConfig)
	processed = ReplacePhpBlocks(processed)
	processedBytes := []byte(processed)

	pt.Templates[path] = processedBytes

	return processedBytes, nil
}
