package urlsafe

import "testing"

func TestNormalizeHTTPURLAllowsHTTPAndHTTPS(t *testing.T) {
	got, ok := NormalizeHTTPURL(" HTTPS://Example.COM/path?q=1#frag ")
	if !ok {
		t.Fatal("NormalizeHTTPURL() rejected valid HTTPS URL")
	}
	want := "https://example.com/path?q=1#frag"
	if got != want {
		t.Fatalf("NormalizeHTTPURL() = %q, want %q", got, want)
	}
}

func TestNormalizeHTTPURLRejectsUnsafeURLs(t *testing.T) {
	for _, input := range []string{
		"",
		"javascript:alert(1)",
		"data:text/html,hi",
		"//example.com/path",
		"https://user:pass@example.com/path",
		"https://example.com/a b",
		"https://example.com/\r\nSet-Cookie:x=y",
		"http:example.com",
	} {
		if got, ok := NormalizeHTTPURL(input); ok {
			t.Fatalf("NormalizeHTTPURL(%q) = %q, want rejected", input, got)
		}
	}
}
