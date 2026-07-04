package extsearch

import "testing"

func TestNormalizeYouTubeRelevanceLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		want     string
	}{
		{name: "english", language: "en", want: "en"},
		{name: "english spaced", language: " EN ", want: "en"},
		{name: "korean", language: "ko", want: "ko"},
		{name: "default", language: "", want: "ko"},
		{name: "unsupported", language: "ja", want: "ko"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeYouTubeRelevanceLanguage(tt.language); got != tt.want {
				t.Fatalf("normalizeYouTubeRelevanceLanguage(%q)=%q want %q", tt.language, got, tt.want)
			}
		})
	}
}

func TestNewYouTubeProviderWithLanguage(t *testing.T) {
	provider := NewYouTubeProviderWithLanguage("key", "en")
	if provider.relevanceLanguage != "en" {
		t.Fatalf("relevanceLanguage=%q want en", provider.relevanceLanguage)
	}

	defaultProvider := NewYouTubeProvider("key")
	if defaultProvider.relevanceLanguage != "ko" {
		t.Fatalf("default relevanceLanguage=%q want ko", defaultProvider.relevanceLanguage)
	}
}
