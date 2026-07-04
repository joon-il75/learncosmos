package safety

import (
	"regexp"
	"strings"
	"unicode"
)

var whitespacePattern = regexp.MustCompile(`\s+`)

func normalizeText(input string) string {
	value := strings.TrimSpace(strings.ToLower(input))
	if value == "" {
		return ""
	}
	value = strings.Map(func(r rune) rune {
		switch {
		case r == '\u200b' || r == '\u200c' || r == '\u200d' || r == '\ufeff':
			return -1
		case unicode.IsControl(r) && r != '\n' && r != '\t':
			return -1
		default:
			return r
		}
	}, value)
	value = whitespacePattern.ReplaceAllString(value, " ")
	return strings.TrimSpace(value)
}

func textHash(input string) string {
	return sha256Hex(strings.TrimSpace(input))
}
