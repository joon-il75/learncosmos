package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNormalizeRedirectAfterAllowsInternalPaths(t *testing.T) {
	tests := map[string]string{
		"/dashboard":                         "/dashboard",
		"/dashboard/course-drafts/abc?x=1":   "/dashboard/course-drafts/abc?x=1",
		"/en/agreements?redirect_after=%2Fx": "/en/agreements?redirect_after=%2Fx",
	}
	for input, want := range tests {
		if got := normalizeRedirectAfter(input); got != want {
			t.Fatalf("normalizeRedirectAfter(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeRedirectAfterRejectsExternalOrAmbiguousTargets(t *testing.T) {
	tests := []string{
		"",
		"https://evil.example/dashboard",
		"http://evil.example/dashboard",
		"//evil.example/dashboard",
		"javascript:alert(1)",
		"/\\evil",
		"/dashboard\r\nSet-Cookie:x=y",
	}
	for _, input := range tests {
		if got := normalizeRedirectAfter(input); got != "" {
			t.Fatalf("normalizeRedirectAfter(%q) = %q, want empty", input, got)
		}
	}
}

func TestFrontendURLIgnoresSpoofedForwardedHost(t *testing.T) {
	t.Setenv("FRONTEND_URL", "")
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/auth/google/callback", nil)
	ctx.Request.Host = "dev.learnweavr.com"
	ctx.Request.Header.Set("X-Forwarded-Host", "evil.example")
	ctx.Request.Header.Set("X-Forwarded-Proto", "http")

	h := &Handler{}
	if got := h.frontendURL(ctx); got != "https://dev.learnweavr.com" {
		t.Fatalf("frontendURL() = %q, want default dev host without spoofed forwarded host", got)
	}
}

func TestFrontendURLFallsBackWhenHostIsNotAllowed(t *testing.T) {
	t.Setenv("FRONTEND_URL", "")
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/auth/google/callback", nil)
	ctx.Request.Host = "evil.example"

	h := &Handler{}
	if got := h.frontendURL(ctx); got != "https://dev.learnweavr.com" {
		t.Fatalf("frontendURL() = %q, want default dev host", got)
	}
}

func TestFrontendURLPrefersConfiguredFrontendURL(t *testing.T) {
	t.Setenv("FRONTEND_URL", "https://www.learnweavr.com/")
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/api/v1/auth/google/callback", nil)
	ctx.Request.Host = "dev.learnweavr.com"
	ctx.Request.Header.Set("X-Forwarded-Host", "evil.example")

	h := &Handler{}
	if got := h.frontendURL(ctx); got != "https://www.learnweavr.com" {
		t.Fatalf("frontendURL() = %q, want configured frontend URL", got)
	}
}

func TestDemoLoginHostGateIgnoresSpoofedForwardedHost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/demo/login", nil)
	ctx.Request.Host = "evil.example"
	ctx.Request.Header.Set("X-Forwarded-Host", "dev.learnweavr.com")

	h := &Handler{demoLogin: DemoLoginConfig{AllowedHosts: []string{"dev.learnweavr.com"}}}
	if h.isDemoLoginHostAllowed(ctx) {
		t.Fatal("isDemoLoginHostAllowed() = true, want false for spoofed X-Forwarded-Host")
	}

	ctx.Request.Host = "dev.learnweavr.com:443"
	if !h.isDemoLoginHostAllowed(ctx) {
		t.Fatal("isDemoLoginHostAllowed() = false, want true for allowed Host")
	}
}
