package content

import (
	"os"
	"strings"
	"testing"
)

func TestBuildCanonicalURL(t *testing.T) {
	raw := "HTTPS://Blog.Naver.com/sample/12345/?from=abc#section"
	got := buildCanonicalURL(&raw)
	if got == nil {
		t.Fatal("expected canonical url")
	}
	want := "https://blog.naver.com/sample/12345/?from=abc"
	if *got != want {
		t.Fatalf("expected %s, got %s", want, *got)
	}
}

func TestBuildCanonicalURLRejectsUnsafeURL(t *testing.T) {
	for _, input := range []string{
		"javascript:alert(1)",
		"//example.com/path",
		"https://user:pass@example.com/path",
		"https://example.com/a b",
	} {
		if got := buildCanonicalURL(&input); got != nil {
			t.Fatalf("buildCanonicalURL(%q) = %q, want nil", input, *got)
		}
	}
}

func TestNormalizeCreateContentRequestURLsRejectsExternalUnsafeURL(t *testing.T) {
	badURL := "javascript:alert(1)"
	req := CreateContentRequest{ContentType: ContentTypeArticle, URL: &badURL, Title: "bad"}
	if err := normalizeCreateContentRequestURLs(&req); err == nil {
		t.Fatal("normalizeCreateContentRequestURLs() error = nil, want error")
	}
}

func TestSaveEmbeddingFailureSanitizesPersistedErrorMessage(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatalf("read repository.go: %v", err)
	}
	text := string(source)
	if !strings.Contains(text, "errMessage = logsafe.PersistedError(embedErr.Error())") {
		t.Fatal("SaveEmbeddingFailure must sanitize persisted embedding error messages")
	}
	if strings.Contains(text, "errMessage = truncateString(embedErr.Error()") {
		t.Fatal("SaveEmbeddingFailure must not persist raw embedding error messages")
	}
}
