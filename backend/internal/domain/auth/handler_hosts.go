package auth

import (
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultFrontendURL = "https://dev.learnweavr.com"

var allowedFallbackFrontendHosts = map[string]struct{}{
	"dev.learnweavr.com": {},
	"www.learnweavr.com": {},
	"learnweavr.com":     {},
}

func (h *Handler) frontendURL(c *gin.Context) string {
	if configured := configuredFrontendURL(); configured != "" {
		return configured
	}

	host := normalizedHost(c.Request.Host)
	if _, ok := allowedFallbackFrontendHosts[host]; ok {
		return "https://" + host
	}

	return defaultFrontendURL
}

func configuredFrontendURL() string {
	configured := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_URL")), "/")
	if configured == "" {
		return ""
	}
	parsed, err := url.Parse(configured)
	if err != nil || parsed == nil || parsed.Host == "" {
		return defaultFrontendURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return defaultFrontendURL
	}
	return configured
}

const invalidHostChars = "\\/\r\n\t @"

func normalizedHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" || strings.ContainsAny(host, invalidHostChars) {
		return ""
	}

	if splitHost, _, err := net.SplitHostPort(host); err == nil {
		host = splitHost
	} else if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.TrimPrefix(strings.TrimSuffix(host, "]"), "[")
	} else if strings.Contains(host, ":") {
		return ""
	}

	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	if host == "" || strings.ContainsAny(host, invalidHostChars) {
		return ""
	}
	return host
}
