package builder

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

var (
	placeholderRe     = regexp.MustCompile(`\{\{\s*([A-Za-z_][\w.\-| ]*)\s*\}\}`)
	phpBlockRe        = regexp.MustCompile(`(?s)<\?(?:php|=)?.*?\?>`)
	phpRestoreBlockRe = regexp.MustCompile(`(?s)<php>(.*?)</php>`)
)

type Placeholder struct {
	templateFunctions *TemplateFunctions
}

func NewPlaceholder(templateFunctions *TemplateFunctions) *Placeholder {
	return &Placeholder{
		templateFunctions: templateFunctions,
	}
}

// ReplacePlaceholders replaces
//   - {{ Name }}
//   - {{ Name | TemplateFunction }}
func (p *Placeholder) ReplacePlaceholders(s string, data map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		placeholder := placeholderRe.FindStringSubmatch(match)[1]
		fieldName, functionName := p.parsePlaceholder(placeholder)

		if val, ok := data[fieldName]; ok {

			if functionName != "" {
				var err error
				val, err = p.templateFunctions.Invoke(functionName, val)
				if err != nil {
					slog.Error("Failed to invoke", "field", fieldName, "function", functionName, "error", err)
					return err.Error()
				}
			}

			return val
		}
		return match
	})
}

func (p *Placeholder) parsePlaceholder(s string) (fieldName string, function string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ""
	}

	i := strings.Index(s, "|")
	if i == -1 {
		return s, ""
	}

	return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
}

// ReplacePhpBlocks finds embedded PHP-Code-Blocks and replaces them by <php>BASE64-encoded code</php>
// to protect code from being removed or altered by html.Parse
func (p *Placeholder) ReplacePhpBlocks(s string) string {
	return phpBlockRe.ReplaceAllStringFunc(s, func(match string) string {
		enc := base64.StdEncoding.EncodeToString([]byte(match))
		return fmt.Sprintf("<php>%s</php>", enc)
	})
}

// RestorePhpBlocks restores PHP-Code-Blocks that have been replaced by ReplacePhpBlocks
func (p *Placeholder) RestorePhpBlocks(s string) string {
	return phpRestoreBlockRe.ReplaceAllStringFunc(s, func(match string) string {
		php := phpRestoreBlockRe.FindStringSubmatch(match)[1]
		dec, err := base64.StdEncoding.DecodeString(php)
		if err != nil {
			return "<!-- " + err.Error() + " -->"
		}
		return string(dec)
	})
}
