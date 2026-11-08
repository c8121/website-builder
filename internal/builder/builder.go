package builder

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
	"github.com/otiai10/copy"
)

// Builder creates the website
type Builder struct {
	Config *BuildConfig
}

func NewBuilder(templateDir string, outputDir string) (*Builder, error) {

	if !util.IsExistingDirectory(templateDir, slog.LevelError) {
		return nil, os.ErrNotExist
	}

	if !util.IsExistingDirectory(outputDir, slog.LevelError) {
		return nil, os.ErrNotExist
	}

	config := &BuildConfig{
		FileMode:                0740,
		OutputBaseDir:           outputDir,
		TemplateDir:             templateDir,
		DefaultTemplateFileName: "index.html",
	}

	assetsOutputDir := filepath.Join(outputDir, "assets")
	err := util.CreateDirIfNotExists(assetsOutputDir, config.FileMode)
	if err != nil {
		return nil, err
	}

	config.AssetLoader = NewAssetLoader(
		"http://localhost:9998", //TODO
		assetsOutputDir,
	)

	config.JsonBinder = NewJsonBinder(config.AssetLoader)

	builder := &Builder{
		Config: config,
	}

	return builder, nil
}

func (b *Builder) SelectTemplate(page *entity.Page) (string, error) {
	return filepath.Join(b.Config.TemplateDir, b.Config.DefaultTemplateFileName), nil
}

func (b *Builder) CreateTemplateParser(page *entity.Page) (*TemplateParser, string, error) {
	templateFile, err := b.SelectTemplate(page)
	if err != nil {
		return nil, "", err
	}

	parser, err := NewTemplateParserFromFile(templateFile, b.Config)
	return parser, templateFile, err
}

func (b *Builder) Build() {

	util.LogError(entity.WalkPages(func(page *entity.Page) {

		var outputDir string
		var err error

		if b.Config.IsHomePage(page) {
			outputDir = b.Config.OutputBaseDir
		} else {
			outputDir, err = GetPageDirectoryPath(page)
			if err != nil {
				slog.Error("Failed to get path", "path", outputDir, "err", err)
				return
			}
		}

		slog.Info("Create Page", "id", page.Id, "path", outputDir)
		err = os.MkdirAll(outputDir, b.Config.FileMode)
		if err != nil {
			slog.Error("Failed to get path", "path", outputDir, "err", err)
			return
		}

		templateParser, templateFile, err := b.CreateTemplateParser(page)
		if err != nil {
			slog.Error("Failed to create parser", "path", outputDir, "err", err)
			return
		}

		err = templateParser.InsertContent(page)
		if err != nil {
			slog.Error("Failed to parse", "path", outputDir, "err", err)
			return
		}

		err = templateParser.InsertMainNav()
		if err != nil {
			slog.Error("Failed to insert main nav", "path", outputDir, "err", err)
			return
		}

		err = templateParser.ParseDocument()
		if err != nil {
			slog.Error("Failed to parse document", "path", outputDir, "err", err)
			return
		}

		outputFile := filepath.Join(outputDir, filepath.Base(templateFile))
		err = templateParser.Save(outputFile)
		if err != nil {
			slog.Error("Failed to save", "path", outputDir, "err", err)
			return
		}
	}))
}

// CopyTemplate copies the whole TemplateDir to OutputBaseDir
func (b *Builder) CopyTemplate() error {

	opts := copy.Options{
		// No chmod
		PermissionControl: copy.DoNothing,
	}

	err := copy.Copy(b.Config.TemplateDir, b.Config.OutputBaseDir, opts)
	if err != nil {
		slog.Error("Failed to copy template", "error", err)
		return err
	}

	return nil
}
