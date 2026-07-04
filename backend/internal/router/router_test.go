package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestConfigureTrustedProxiesAllowsOnlyLoopbackProxyHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := configureTrustedProxies(r); err != nil {
		t.Fatalf("configureTrustedProxies() error = %v", err)
	}
	r.GET("/client-ip", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})

	loopbackReq := httptest.NewRequest(http.MethodGet, "/client-ip", nil)
	loopbackReq.RemoteAddr = "127.0.0.1:12345"
	loopbackReq.Header.Set("X-Forwarded-For", "203.0.113.10")
	loopbackResp := httptest.NewRecorder()
	r.ServeHTTP(loopbackResp, loopbackReq)
	if got := loopbackResp.Body.String(); got != "203.0.113.10" {
		t.Fatalf("loopback proxy ClientIP() = %q, want forwarded client IP", got)
	}

	directReq := httptest.NewRequest(http.MethodGet, "/client-ip", nil)
	directReq.RemoteAddr = "198.51.100.20:12345"
	directReq.Header.Set("X-Forwarded-For", "203.0.113.99")
	directResp := httptest.NewRecorder()
	r.ServeHTTP(directResp, directReq)
	if got := directResp.Body.String(); got != "198.51.100.20" {
		t.Fatalf("untrusted direct ClientIP() = %q, want remote address without spoofed X-Forwarded-For", got)
	}
}

func TestUnsafeRequestOriginGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		method     string
		origin     string
		referer    string
		host       string
		wantStatus int
	}{
		{name: "same origin post allowed", method: http.MethodPost, origin: "https://dev.learnweavr.com", host: "dev.learnweavr.com", wantStatus: http.StatusNoContent},
		{name: "fallback frontend origin allowed", method: http.MethodPatch, origin: "https://www.learnweavr.com", host: "dev.learnweavr.com", wantStatus: http.StatusNoContent},
		{name: "cross origin post blocked", method: http.MethodPost, origin: "https://attacker.example", host: "dev.learnweavr.com", wantStatus: http.StatusForbidden},
		{name: "cross origin referer blocked", method: http.MethodDelete, referer: "https://attacker.example/path", host: "dev.learnweavr.com", wantStatus: http.StatusForbidden},
		{name: "non browser unsafe request allowed", method: http.MethodPost, host: "dev.learnweavr.com", wantStatus: http.StatusNoContent},
		{name: "safe method allowed", method: http.MethodGet, origin: "https://attacker.example", host: "dev.learnweavr.com", wantStatus: http.StatusNoContent},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(unsafeRequestOriginGuard())
			r.Handle(tc.method, "/probe", func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})

			req := httptest.NewRequest(tc.method, "/probe", nil)
			req.Host = tc.host
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			if tc.referer != "" {
				req.Header.Set("Referer", tc.referer)
			}
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			if resp.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tc.wantStatus)
			}
		})
	}
}
