package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/auth"
)

const adminAccessBoundaryJWTSecret = "12345678901234567890123456789012"

func newAdminBoundaryRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handlers := append([]gin.HandlerFunc{}, middlewares...)
	handlers = append(handlers, func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	r.GET("/", handlers...)
	return r
}

func adminBoundaryContext(userID string, role auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, userID)
		c.Set(auth.ContextKeyRole, role)
		c.Next()
	}
}

func TestSuperAdminMiddlewareChainRequiresSuperAdminCookieAndRole(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_BEARER_FALLBACK_ENABLED", "0")

	svc := auth.NewService(nil, nil, adminAccessBoundaryJWTSecret, auth.OAuthConfig{})
	superToken, err := svc.IssueSuperAdminToken(auth.User{ID: "super_admin", Role: auth.RoleSuperAdmin})
	if err != nil {
		t.Fatalf("issue super admin token: %v", err)
	}
	adminToken, err := svc.IssueAdminTOTPToken("admin-1")
	if err != nil {
		t.Fatalf("issue admin token: %v", err)
	}

	r := newAdminBoundaryRouter(svc.JWTMiddlewareHeaderOnly(), TOTPMiddleware())

	t.Run("allows super admin cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: superAdminAccessTokenCookieName, Value: superToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("rejects admin role token in super admin cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: superAdminAccessTokenCookieName, Value: adminToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("ignores regular access token cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: adminToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestRequireAdminTOTPSessionBindsTokenToAdminUserAndRole(t *testing.T) {
	svc := auth.NewService(nil, nil, adminAccessBoundaryJWTSecret, auth.OAuthConfig{})
	matchingAdminToken, err := svc.IssueAdminTOTPToken("admin-1")
	if err != nil {
		t.Fatalf("issue matching admin token: %v", err)
	}
	otherAdminToken, err := svc.IssueAdminTOTPToken("admin-2")
	if err != nil {
		t.Fatalf("issue other admin token: %v", err)
	}
	superToken, err := svc.IssueSuperAdminToken(auth.User{ID: "super_admin", Role: auth.RoleSuperAdmin})
	if err != nil {
		t.Fatalf("issue super admin token: %v", err)
	}

	h := &AdminHandler{authSvc: svc}

	t.Run("allows matching admin totp session", func(t *testing.T) {
		r := newAdminBoundaryRouter(adminBoundaryContext("admin-1", auth.RoleAdmin), h.requireAdminTOTPSession())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: adminTOTPSessionCookieName, Value: matchingAdminToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("rejects other admin totp session", func(t *testing.T) {
		r := newAdminBoundaryRouter(adminBoundaryContext("admin-1", auth.RoleAdmin), h.requireAdminTOTPSession())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: adminTOTPSessionCookieName, Value: otherAdminToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("rejects super admin token as admin totp session", func(t *testing.T) {
		r := newAdminBoundaryRouter(adminBoundaryContext("admin-1", auth.RoleAdmin), h.requireAdminTOTPSession())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: adminTOTPSessionCookieName, Value: superToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("rejects non admin context", func(t *testing.T) {
		r := newAdminBoundaryRouter(adminBoundaryContext("admin-1", auth.RoleLearner), h.requireAdminTOTPSession())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: adminTOTPSessionCookieName, Value: matchingAdminToken})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})
}
