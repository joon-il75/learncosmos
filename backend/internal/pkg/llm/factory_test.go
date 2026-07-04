package llm

import (
	"testing"
)

func TestNewClient_AllProviders(t *testing.T) {
	providers := []struct {
		name     string
		provider string
	}{
		{"openai", "openai"},
		{"google", "google"},
		{"grok", "grok"},
		{"anthropic", "anthropic"},
		{"solar", "solar"},
		{"hyperclova", "hyperclova"},
		{"llama", "llama"},
		{"exaone", "exaone"},
	}

	for _, tc := range providers {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.provider, "test-key", "")
			if err != nil {
				t.Errorf("NewClient(%q) error = %v", tc.provider, err)
			}
			if client == nil {
				t.Errorf("NewClient(%q) returned nil client", tc.provider)
			}
		})
	}
}

func TestNewClient_UnknownProvider(t *testing.T) {
	_, err := NewClient("unknown-provider", "test-key", "")
	if err == nil {
		t.Error("expected error for unknown provider, got nil")
	}
}

func TestNewClient_LlamaWithEndpoint(t *testing.T) {
	client, err := NewClient("llama", "test-key", "", "https://custom.endpoint/v1")
	if err != nil {
		t.Fatalf("NewClient(llama) error = %v", err)
	}
	llama, ok := client.(*LlamaClient)
	if !ok {
		t.Fatal("expected *LlamaClient")
	}
	if llama.baseURL != "https://custom.endpoint/v1" {
		t.Errorf("expected custom endpoint, got %q", llama.baseURL)
	}
}

func TestNewClient_LlamaDefaultEndpoint(t *testing.T) {
	client, err := NewClient("llama", "test-key", "")
	if err != nil {
		t.Fatalf("NewClient(llama) error = %v", err)
	}
	llama, ok := client.(*LlamaClient)
	if !ok {
		t.Fatal("expected *LlamaClient")
	}
	if llama.baseURL != llamaDefaultBaseURL {
		t.Errorf("expected default endpoint %q, got %q", llamaDefaultBaseURL, llama.baseURL)
	}
}
