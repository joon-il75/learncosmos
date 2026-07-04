package logsafe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	urlPattern          = regexp.MustCompile(`https?://[^\s"'<>]+`)
	bearerPattern       = regexp.MustCompile(`(?i)\b(bearer)\s+[^\s"'<>]+`)
	secretAssignPattern = regexp.MustCompile(`(?i)\b(api[_-]?key|x-api-key|token|key|secret|password|authorization)=([^\s"'<>]+)`)
)

func URL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed == nil || parsed.Host == "" {
		return "[invalid-url]"
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed.String()
}

func Text(value string) string {
	if strings.TrimSpace(value) == "" {
		return value
	}
	return urlPattern.ReplaceAllStringFunc(value, URL)
}

func Error(err error) string {
	if err == nil {
		return ""
	}
	return Text(err.Error())
}

func PersistedError(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "persisted_error ") {
		return trimmed
	}
	if code, ok := ClassifyProviderError(trimmed); ok {
		return code
	}
	redacted := redactSecretLikeText(Text(trimmed))
	return "persisted_error " + Summary(redacted)
}

func ProviderErrorCode(value string) string {
	if code, ok := ClassifyProviderError(value); ok {
		return code
	}
	return "provider_error"
}

func ClassifyProviderError(value string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", false
	}
	switch {
	case strings.Contains(normalized, "context canceled"):
		return "provider_request_canceled", true
	case strings.Contains(normalized, "deadline exceeded") || strings.Contains(normalized, "timeout"):
		return "provider_timeout", true
	case strings.Contains(normalized, "401") || strings.Contains(normalized, "unauthorized") || strings.Contains(normalized, "invalid api key"):
		return "provider_unauthorized", true
	case strings.Contains(normalized, "403") || strings.Contains(normalized, "forbidden") || strings.Contains(normalized, "permission"):
		return "provider_forbidden", true
	case strings.Contains(normalized, "429") || strings.Contains(normalized, "rate limit") || strings.Contains(normalized, "quota"):
		return "provider_rate_limited", true
	case strings.Contains(normalized, "500") || strings.Contains(normalized, "502") || strings.Contains(normalized, "503") || strings.Contains(normalized, "504"):
		return "provider_unavailable", true
	case strings.Contains(normalized, "api error"):
		return "provider_api_error", true
	default:
		return "", false
	}
}

func redactSecretLikeText(value string) string {
	redacted := bearerPattern.ReplaceAllString(value, "$1 [redacted]")
	redacted = secretAssignPattern.ReplaceAllString(redacted, "$1=[redacted]")
	return redacted
}

func Summary(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "empty"
	}
	sum := sha256.Sum256([]byte(trimmed))
	return fmt.Sprintf("len=%d sha256=%s", utf8.RuneCountInString(trimmed), hex.EncodeToString(sum[:])[:12])
}
