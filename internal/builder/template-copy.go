package builder

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/go_vue_spa/util"
)

type TemplateCopy struct {
	src    string
	dest   string
	Config *BuildConfig
}

func NewTemplateCopy(src string, dest string, config *BuildConfig) *TemplateCopy {
	return &TemplateCopy{
		src:    filepath.Clean(src),
		dest:   filepath.Clean(dest),
		Config: config,
	}
}

func (c *TemplateCopy) Copy() error {

	slog.Info("Copy template", "src", c.src, "dst", c.dest)

	return filepath.Walk(c.src, func(path string, info os.FileInfo, err error) error {
		return c.copyEntry(path, info, err)
	})

}

func (c *TemplateCopy) copyEntry(path string, info os.FileInfo, err error) error {

	dest := filepath.Join(c.dest, path[len(c.src):])
	//slog.Info("Copy entry", "path", path, "dest", dest)

	if info.IsDir() {
		_, err := os.Stat(dest)
		if errors.Is(err, os.ErrNotExist) {
			slog.Info("Create directory", "dest", dest)
			if err := os.MkdirAll(dest, c.Config.FileMode); err != nil {
				return err
			}
		}
	} else {
		isTemplate, err := c.isTemplateFile(path)
		if err != nil {
			return err
		}
		if !isTemplate {
			err = c.copyFile(path, dest)
		}
	}

	return nil
}

// isTemplateFile determines if the given file is a template (currently only by name, if starting with "index.")
func (c *TemplateCopy) isTemplateFile(path string) (bool, error) {
	name := filepath.Base(path)
	if strings.HasPrefix(name, "index.") {
		return true, nil
	}
	return false, nil
}

func (c *TemplateCopy) copyFile(src string, dest string) error {

	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer util.CloseOrLog(reader)

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer util.CloseOrLog(f)

	var buf []byte = nil
	var w io.Writer = f
	var r io.Reader = reader

	if _, err = io.CopyBuffer(w, r, buf); err != nil {
		return err
	}

	return f.Sync()
}
