package urlsafe

import (
	"net/url"
	"strings"
	"unicode"
)

func NormalizeHTTPURL(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.HasPrefix(trimmed, "//") || containsSpaceOrControl(trimmed) {
		return "", false
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || parsed == nil {
		return "", false
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", false
	}
	if parsed.Host == "" || parsed.Hostname() == "" || parsed.User != nil || parsed.Opaque != "" {
		return "", false
	}
	if containsSpaceOrControl(parsed.Host) {
		return "", false
	}

	parsed.Scheme = scheme
	parsed.Host = strings.ToLower(parsed.Host)
	return parsed.String(), true
}

func containsSpaceOrControl(value string) bool {
	return strings.IndexFunc(value, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsControl(r)
	}) >= 0
}
