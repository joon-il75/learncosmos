package logsafe

import (
	"strings"
	"testing"
)

func TestURLRedactsSensitiveParts(t *testing.T) {
	got := URL("https://user:pass@example.com/path/file?token=secret#frag")
	want := "https://example.com/path/file"
	if got != want {
		t.Fatalf("URL() = %q, want %q", got, want)
	}
}

func TestTextRedactsEmbeddedURLs(t *testing.T) {
	got := Text("failed https://example.com/callback?code=secret&state=x and http://user:pass@host.test/a#b")
	if got != "failed https://example.com/callback and http://host.test/a" {
		t.Fatalf("Text() = %q", got)
	}
}

func TestSummaryOmitsRawInput(t *testing.T) {
	secret := "프라이빗 목표: 토큰 secret=abc123"
	got := Summary("  " + secret + "  ")
	if got == "" || got == "empty" {
		t.Fatalf("Summary() = %q", got)
	}
	if strings.Contains(got, secret) || strings.Contains(got, "abc123") || strings.Contains(got, "프라이빗") {
		t.Fatalf("Summary() leaked raw input: %q", got)
	}
	if got != Summary(secret) {
		t.Fatalf("Summary() should be stable after trimming, got %q", got)
	}
}

func TestSummaryEmptyInput(t *testing.T) {
	if got := Summary(" "); got != "empty" {
		t.Fatalf("Summary(empty) = %q", got)
	}
}

func TestPersistedErrorOmitsProviderSecretMaterialAndIsIdempotent(t *testing.T) {
	secret := "provider failed https://api.example.com/v1/chat?api_key=secret-token Authorization Bearer sk-secret"
	got := PersistedError(secret)
	if got == "" || got == "empty" {
		t.Fatalf("PersistedError() = %q", got)
	}
	for _, leaked := range []string{"secret-token", "api_key=secret-token", "sk-secret"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("PersistedError() leaked %q in %q", leaked, got)
		}
	}
	if got != PersistedError(got) {
		t.Fatalf("PersistedError() should be idempotent, got %q then %q", got, PersistedError(got))
	}
}

func TestPersistedErrorEmptyInput(t *testing.T) {
	if got := PersistedError(" "); got != "" {
		t.Fatalf("PersistedError(empty) = %q", got)
	}
}

func TestPersistedErrorClassifiesProviderErrors(t *testing.T) {
	got := PersistedError(`openai: API error 401: invalid api key sk-provider-secret`)
	if got != "provider_unauthorized" {
		t.Fatalf("PersistedError() = %q, want provider_unauthorized", got)
	}
}

func TestPersistedErrorSummarizesURLsAndSecretAssignments(t *testing.T) {
	got := PersistedError(`failed callback https://user:pass@example.com/path?token=secret Authorization Bearer sk-secret api_key=abc123`)
	for _, forbidden := range []string{"user:pass", "token=secret", "sk-secret", "abc123", "api_key=abc123", "https://example.com/path"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("PersistedError() leaked %q in %q", forbidden, got)
		}
	}
	if !strings.HasPrefix(got, "persisted_error len=") || !strings.Contains(got, " sha256=") {
		t.Fatalf("PersistedError() did not return summary shape: %q", got)
	}
}

func TestProviderErrorCodeFallsBackToGenericProviderError(t *testing.T) {
	if got := ProviderErrorCode("provider echoed unexpected opaque body"); got != "provider_error" {
		t.Fatalf("ProviderErrorCode() = %q, want provider_error", got)
	}
}

func TestPersistedErrorClassifiesProviderFailures(t *testing.T) {
	if got := PersistedError("openai API error 401 invalid api key"); got != "provider_unauthorized" {
		t.Fatalf("PersistedError(401) = %q, want provider_unauthorized", got)
	}
	if got := PersistedError("request timeout while calling provider"); got != "provider_timeout" {
		t.Fatalf("PersistedError(timeout) = %q, want provider_timeout", got)
	}
}
