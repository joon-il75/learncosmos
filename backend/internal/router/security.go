package router

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var fallbackSameOriginHosts = map[string]struct{}{
	"dev.learnweavr.com": {},
	"www.learnweavr.com": {},
	"learnweavr.com":     {},
}

func unsafeRequestOriginGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isUnsafeMethod(c.Request.Method) {
			c.Next()
			return
		}
		if isAllowedBrowserOrigin(c.Request) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":      "request origin is not allowed",
			"error_code": "origin_not_allowed",
		})
	}
}

func isUnsafeMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func isAllowedBrowserOrigin(req *http.Request) bool {
	origin := strings.TrimSpace(req.Header.Get("Origin"))
	if origin != "" {
		return isAllowedOriginHost(origin, req.Host)
	}

	referer := strings.TrimSpace(req.Header.Get("Referer"))
	if referer != "" {
		return isAllowedOriginHost(referer, req.Host)
	}

	return true
}

func isAllowedOriginHost(rawURL, requestHost string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed == nil || parsed.Host == "" {
		return false
	}
	host := normalizedRequestHost(parsed.Host)
	if host == "" {
		return false
	}
	if host == normalizedRequestHost(requestHost) {
		return true
	}
	if frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL")); frontendURL != "" {
		if parsedFrontend, err := url.Parse(frontendURL); err == nil && parsedFrontend != nil && host == normalizedRequestHost(parsedFrontend.Host) {
			return true
		}
	}
	_, ok := fallbackSameOriginHosts[host]
	return ok
}

func normalizedRequestHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" || strings.ContainsAny(host, "\\/\r\n\t @") {
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
	if host == "" || strings.ContainsAny(host, "\\/\r\n\t @") {
		return ""
	}
	return host
}
