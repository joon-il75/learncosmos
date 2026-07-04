package abuseguard

import (
	"testing"
	"time"
)

func TestFixedWindowLimiterAllowsThenBlocksUntilWindowResets(t *testing.T) {
	now := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	limiter := NewFixedWindowLimiter(2, time.Minute)
	limiter.now = func() time.Time { return now }

	if ok, _ := limiter.Allow("client"); !ok {
		t.Fatal("first request blocked")
	}
	if ok, _ := limiter.Allow("client"); !ok {
		t.Fatal("second request blocked")
	}
	if ok, retryAfter := limiter.Allow("client"); ok || retryAfter <= 0 {
		t.Fatalf("third request ok=%v retryAfter=%v, want blocked with retryAfter", ok, retryAfter)
	}

	now = now.Add(time.Minute)
	if ok, _ := limiter.Allow("client"); !ok {
		t.Fatal("request after reset blocked")
	}
}

func TestFixedWindowLimiterScopesKeysIndependently(t *testing.T) {
	limiter := NewFixedWindowLimiter(1, time.Minute)
	if ok, _ := limiter.Allow("a"); !ok {
		t.Fatal("first key blocked")
	}
	if ok, _ := limiter.Allow("b"); !ok {
		t.Fatal("second key should have independent budget")
	}
	if ok, _ := limiter.Allow("a"); ok {
		t.Fatal("first key should be exhausted")
	}
}
