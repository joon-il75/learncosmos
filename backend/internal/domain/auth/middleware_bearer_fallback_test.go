package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newAuthTestContext(authHeader string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	c.Request = req
	return c
}

func TestExtractTokenDisablesBearerFallbackByDefaultInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "")

	_, err := extractToken(newAuthTestContext("Bearer test-token"))
	if err == nil {
		t.Fatal("extractToken() allowed bearer fallback in production without explicit opt-in")
	}
}

func TestExtractTokenAllowsBearerFallbackWhenExplicitlyEnabled(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "1")

	got, err := extractToken(newAuthTestContext("Bearer test-token"))
	if err != nil {
		t.Fatalf("extractToken() error = %v", err)
	}
	if got != "test-token" {
		t.Fatalf("extractToken() = %q, want test-token", got)
	}
}

func TestExtractTokenAllowsBearerFallbackOutsideProductionByDefault(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "")

	got, err := extractToken(newAuthTestContext("Bearer dev-token"))
	if err != nil {
		t.Fatalf("extractToken() error = %v", err)
	}
	if got != "dev-token" {
		t.Fatalf("extractToken() = %q, want dev-token", got)
	}
}

func TestExtractTokenCookieStillWorksWhenBearerFallbackDisabled(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "0")

	c := newAuthTestContext("Bearer ignored-token")
	c.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "cookie-token"})

	got, err := extractToken(c)
	if err != nil {
		t.Fatalf("extractToken() error = %v", err)
	}
	if got != "cookie-token" {
		t.Fatalf("extractToken() = %q, want cookie-token", got)
	}
}

func TestExtractSuperAdminTokenUsesOnlySuperAdminCookieWhenFallbackDisabled(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "0")

	regularCookieOnly := newAuthTestContext("")
	regularCookieOnly.Request.AddCookie(&http.Cookie{Name: "access_token", Value: "regular-token"})
	if _, err := extractSuperAdminToken(regularCookieOnly); err == nil {
		t.Fatal("extractSuperAdminToken() accepted regular access_token cookie")
	}

	superCookie := newAuthTestContext("")
	superCookie.Request.AddCookie(&http.Cookie{Name: superAdminAccessTokenCookieName, Value: "super-token"})
	got, err := extractSuperAdminToken(superCookie)
	if err != nil {
		t.Fatalf("extractSuperAdminToken() error = %v", err)
	}
	if got != "super-token" {
		t.Fatalf("extractSuperAdminToken() = %q, want super-token", got)
	}
}

func TestExtractTokenIgnoresSuperAdminCookie(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "0")

	c := newAuthTestContext("")
	c.Request.AddCookie(&http.Cookie{Name: superAdminAccessTokenCookieName, Value: "super-token"})
	if _, err := extractToken(c); err == nil {
		t.Fatal("extractToken() accepted super_admin_access_token cookie")
	}
}
