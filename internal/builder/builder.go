package builder

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
	"github.com/otiai10/copy"
)

type Builder struct {
	FileMode                os.FileMode
	DefaultTemplateFileName string

	Parser *Parser

	templateDir string
	outputDir   string

	template string
}

func NewBuilder(templateDir string, outputDir string) (*Builder, error) {

	if !util.IsExistingDirectory(templateDir, slog.LevelError) {
		return nil, os.ErrNotExist
	}

	if !util.IsExistingDirectory(outputDir, slog.LevelError) {
		return nil, os.ErrNotExist
	}

	builder := &Builder{
		FileMode:                0740,
		DefaultTemplateFileName: "index.html",

		templateDir: templateDir,
		outputDir:   outputDir,
	}

	builder.template = filepath.Join(templateDir, builder.DefaultTemplateFileName)
	parser, err := NewParserFromFile(builder.template)
	if err != nil {
		return nil, err
	}
	builder.Parser = parser

	return builder, nil
}

func (b *Builder) Build() error {

	util.LogError(entity.WalkPages(func(page entity.Page) {
		path, err := GetPageDirectoryPath(page)
		if err != nil {
			slog.Error("Failed to get path", "path", path)
			return
		}

		slog.Info("Create Page", "id", page.Id, "path", path)
		err = os.MkdirAll(path, b.FileMode)
		if err != nil {
			slog.Error("Failed to get path", "path", path)
			return
		}

		template := filepath.Join(path, b.DefaultTemplateFileName)
		err = copy.Copy(b.template, template)
		if err != nil {
			slog.Error("Failed to copy page template", "path", path, "err", err)
			return
		}

		contents, err := entity.ListPageContents(page.Id)
		if err != nil {
			slog.Error("Failed to list contents", "path", path, "err", err)
			return
		}

		err = b.ParseTemplateFile(template, page, contents)
		if err != nil {
			slog.Error("Failed to parse template", "path", path, "err", err)
			return
		}
	}))

	return nil
}

func (b *Builder) ParseTemplateFile(file string, page entity.Page, contents []entity.PageContent) error {
	slog.Info("Parse Template", "path", file, "page", page.Id, "contents", len(contents))

	err := b.Parser.Parse(page, contents)
	if err != nil {
		return err
	}

	err = b.Parser.Save(file, b.FileMode)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) CopyTemplate() error {

	opts := copy.Options{
		// No chmod
		PermissionControl: copy.DoNothing,
	}

	err := copy.Copy(b.templateDir, b.outputDir, opts)
	if err != nil {
		slog.Error("Failed to copy template", "error", err)
		return err
	}

	return nil
}
