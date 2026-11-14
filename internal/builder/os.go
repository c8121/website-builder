package builder

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/c8121/website-builder/internal/entity"
	"github.com/c8121/website-builder/shared/config"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	hyphensRegEx = regexp.MustCompile(`-{2,}`)
)

func GetPageDirectoryPath(page *entity.Page) (string, error) {

	path, err := entity.LoadPagePath(entity.NewSession(), page.Id)
	if err != nil {
		return "", err
	}

	directoryPath := config.AppConfig.OutputDir
	for _, item := range path {
		directoryPath = filepath.Join(directoryPath, GetPageDirectoryName(item))
	}

	return directoryPath + "/", nil
}

func GetPageDirectoryName(page entity.Page) string {
	return NormalizeAsFileName(page.Name, fmt.Sprintf("page-%d", page.Id))
}

func NormalizeAsFileName(name string, fallback string) string {

	s := name
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), "-")
	s = replaceUmlautsGerman(s)
	s = replaceIllegalCharacters(s)
	s = strings.Trim(s, "-")
	s = strings.ToLower(s)

	if s == "" {
		slog.Warn("Generated empty name", "name", name, "fallback", fallback)
		s = fallback
	}

	return s
}

func replaceUmlautsGerman(input string) string {
	r := strings.NewReplacer(
		"ä", "ae", "ö", "oe", "ü", "ue",
		"Ä", "Ae", "Ö", "Oe", "Ü", "Ue",
		"ß", "ss",
	)
	return r.Replace(input)
}

func replaceIllegalCharacters(input string) string {

	// Characters explicitly forbidden in Windows, Linux, and macOS file paths
	// Windows: < > : " / \ | ? *
	// Linux/macOS: / and NUL (handled via control characters check)
	invalidPathChars := `<>:"/\|?*`

	// Replace unwanted punctuation with a hyphen
	replaceInvalidPathRunes := runes.Map(func(r rune) rune {
		// 1. Keep valid
		if r == '-' || r == '.' {
			return r
		}

		// 2. Check for ASCII control characters (0-31 and 127)
		if r < 32 || r == 127 {
			return '-'
		}

		// 3. Check for explicitly forbidden OS path characters
		if strings.ContainsRune(invalidPathChars, r) {
			return '-'
		}

		// 4. Fallback: Replace other unwanted punctuation
		if unicode.IsPunct(r) {
			return '-'
		}

		return r
	})

	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)), // Remove accents/diacritics
		replaceInvalidPathRunes,
		norm.NFC,
	)

	output, _, err := transform.String(t, input)
	if err != nil {
		slog.Error("Transform failed", "error", err)
		return input
	}

	output = hyphensRegEx.ReplaceAllString(output, "-")
	return output
}
