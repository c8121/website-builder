package builder

import "os"

type TemplateConfig struct {
	FileMode                os.FileMode
	OutputBaseDir           string
	TemplateDir             string
	DefaultTemplateFileName string
	AssetLoader             *AssetLoader
	JsonBinder              *JsonBinder
}
