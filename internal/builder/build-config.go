package builder

import "os"

type BuildConfig struct {
	FileMode                os.FileMode
	OutputBaseDir           string
	TemplateDir             string
	DefaultTemplateFileName string
	AssetLoader             *AssetLoader
	JsonBinder              *JsonBinder
}
