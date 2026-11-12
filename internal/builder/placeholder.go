package builder

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"regexp"
)

var (
	placeholderRe     = regexp.MustCompile(`\{\{\s*([A-Za-z_][\w.\-]*)\s*\}\}`)
	phpBlockRe        = regexp.MustCompile(`(?s)<\?(?:php|=)?.*?\?>`)
	phpRestoreBlockRe = regexp.MustCompile(`(?s)<php>(.*?)</php>`)
)

// ReplacePlaceholders replaces {{ Name }}...
func ReplacePlaceholders(s string, data map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		key := placeholderRe.FindStringSubmatch(match)[1]
		if val, ok := data[key]; ok {
			return val
		}
		return match
	})
}

func ReplacePhpBlocks(s string) string {
	return phpBlockRe.ReplaceAllStringFunc(s, func(match string) string {
		enc := base64.StdEncoding.EncodeToString([]byte(match))
		return fmt.Sprintf("<php>%s</php>", enc)
	})
}

func RestorePhpBlocks(s string) string {
	return phpRestoreBlockRe.ReplaceAllStringFunc(s, func(match string) string {
		slog.Info("php", "php", match)
		php := phpRestoreBlockRe.FindStringSubmatch(match)[1]
		dec, err := base64.StdEncoding.DecodeString(php)
		if err != nil {
			return "<!-- " + err.Error() + " -->"
		}
		return string(dec)
	})
}
