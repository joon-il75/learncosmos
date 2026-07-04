package safehttp

import (
	"context"
	"errors"
	"testing"
)

func TestValidateHTTPURLBlocksInternalTargets(t *testing.T) {
	tests := []string{
		"http://localhost:8080/health",
		"http://127.0.0.1:8080/health",
		"http://10.0.0.1/",
		"http://172.16.0.1/",
		"http://192.168.0.1/",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]/",
		"http://[fe80::1]/",
	}

	for _, rawURL := range tests {
		_, err := ValidateHTTPURL(context.Background(), rawURL)
		if !errors.Is(err, ErrBlockedHost) {
			t.Fatalf("ValidateHTTPURL(%q) error = %v, want ErrBlockedHost", rawURL, err)
		}
	}
}

func TestValidateHTTPURLAllowsPublicHTTPSTarget(t *testing.T) {
	_, err := ValidateHTTPURL(context.Background(), "https://8.8.8.8/dns-query")
	if err != nil {
		t.Fatalf("ValidateHTTPURL() error = %v, want nil", err)
	}
}

func TestValidateHTTPURLRejectsUnsupportedSchemeAndUserInfo(t *testing.T) {
	if _, err := ValidateHTTPURL(context.Background(), "file:///etc/passwd"); !errors.Is(err, ErrInvalidURL) && !errors.Is(err, ErrUnsupported) {
		t.Fatalf("file URL error = %v, want invalid or unsupported", err)
	}
	if _, err := ValidateHTTPURL(context.Background(), "https://user:pass@example.com/"); !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("userinfo URL error = %v, want ErrInvalidURL", err)
	}
}
