package builder

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/c8121/website-builder/internal/entity"
	"github.com/c8121/website-builder/shared/config"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GetPageDirectoryPath(page entity.Page) (string, error) {

	path, err := entity.LoadPagePath(page.Id)
	if err != nil {
		return "", err
	}

	directoryPath := config.AppConfig.OutputDir
	for _, item := range path {
		directoryPath = filepath.Join(directoryPath, GetPageDirectoryName(item))
	}

	return directoryPath, nil
}

func GetPageDirectoryName(page entity.Page) string {
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
