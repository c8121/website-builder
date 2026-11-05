package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/otiai10/copy"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/c8121/website-builder/internal/entity"
	boostrap "github.com/c8121/website-builder/shared/bootstrap"
	"github.com/c8121/website-builder/shared/config"
)

const (
	FileMode os.FileMode = 0750
)

func main() {

	boostrap.InitWebsiteBuiderApp()
	defer boostrap.FinalizeWebsiteBuiderApp()

	if config.AppConfig.TemplateDir == "" {
		slog.Error("No template selected")
		return
	}

	err := CopyTemplate(config.AppConfig.TemplateDir, config.AppConfig.OutputDir)
	if err != nil {
		slog.Error("Exit", "error", err)
		return
	}

	entity.WalkPages(func(page entity.Page) {
		path, err := GetDirectoryPath(page)
		if err != nil {
			slog.Error("Failed to get path", "path", path)
			return
		}

		slog.Info("Create Page", "id", page.Id, "path", path)
		err = os.MkdirAll(path, FileMode)
		if err != nil {
			slog.Error("Failed to get path", "path", path)
			return
		}
	})
}

func CopyTemplate(templatePath string, destPath string) error {

	opts := copy.Options{
		// No chmod
		PermissionControl: copy.DoNothing,
	}

	err := copy.Copy(templatePath, destPath, opts)
	if err != nil {
		slog.Error("Failed to copy template", "error", err)
		return err
	}

	return nil
}

func GetDirectoryPath(page entity.Page) (string, error) {

	path, err := entity.LoadPagePath(page.Id)
	if err != nil {
		return "", err
	}

	directoryPath := config.AppConfig.OutputDir
	for _, item := range path {
		directoryPath = filepath.Join(directoryPath, GetDirectoryName(item))
	}

	return directoryPath, nil
}

func GetDirectoryName(page entity.Page) string {
	s := page.Name
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = ReplaceUmlautsGerman(s)
	s = RemoveIllegalCharacters(s)
	s = strings.ToLower(s)

	if s == "" {
		slog.Warn("Generated empty dir name", "path", s)
		s = fmt.Sprintf("page-%d", page.Id)
	}

	return s
}

func ReplaceUmlautsGerman(input string) string {
	r := strings.NewReplacer(
		"ä", "ae", "ö", "oe", "ü", "ue",
		"Ä", "Ae", "Ö", "Oe", "Ü", "Ue",
		"ß", "ss",
	)
	return r.Replace(input)
}

func RemoveIllegalCharacters(input string) string {

	removePunctExceptHyphen := runes.Predicate(func(r rune) bool {
		return unicode.IsPunct(r) && r != '-'
	})

	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		runes.Remove(removePunctExceptHyphen),
		norm.NFC,
	)

	output, _, err := transform.String(t, input)
	if err != nil {
		slog.Error("Transform failed", "error", err)
		return input
	}
	return output
}
