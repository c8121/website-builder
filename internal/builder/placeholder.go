package builder

import "regexp"

var (
	placeholderRe = regexp.MustCompile(`\{\{\s*([A-Za-z_][\w.\-]*)\s*\}\}`)
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
