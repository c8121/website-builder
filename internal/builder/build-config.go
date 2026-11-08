package builder

import (
	"os"

	"github.com/c8121/website-builder/internal/entity"
)

type BuildConfig struct {
	FileMode                os.FileMode
	OutputBaseDir           string
	TemplateDir             string
	DefaultTemplateFileName string
	AssetLoader             *AssetLoader
	JsonBinder              *JsonBinder

	HomePage *entity.Page
}

func (c *BuildConfig) IsHomePage(page *entity.Page) bool {
	// TODO this just selects the fist page
	if page.Id > 0 && c.HomePage == nil {
		c.HomePage = page
		return true
	}
	return c.HomePage.Id == page.Id
}
