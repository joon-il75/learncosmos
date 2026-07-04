package llm

import "testing"

func TestValidateLlamaBaseURLRejectsUnsafeEndpoint(t *testing.T) {
	for _, input := range []string{
		"",
		"http://api.example.com/openai/v1",
		"javascript:alert(1)",
		"https://user:pass@example.com/openai/v1",
		"https:///openai/v1",
	} {
		if err := validateLlamaBaseURL(input); err == nil {
			t.Fatalf("validateLlamaBaseURL(%q) error = nil, want error", input)
		}
	}
}

func TestValidateLlamaBaseURLAllowsHTTPS(t *testing.T) {
	if err := validateLlamaBaseURL("https://api.example.com/openai/v1"); err != nil {
		t.Fatalf("validateLlamaBaseURL() error = %v", err)
	}
}
