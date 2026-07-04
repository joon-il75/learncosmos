package llm

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestProviderClientsDoNotReturnRawAPIErrorBodies(t *testing.T) {
	cases := []struct {
		name      string
		provider  string
		client    Client
		status    int
		body      string
		sensitive []string
	}{
		{
			name:     "openai",
			provider: "openai",
			client:   NewOpenAIClient("sk-provider-secret", "gpt-4o-mini"),
			status:   http.StatusUnauthorized,
			body:     `{"error":"invalid key sk-provider-secret"}`,
			sensitive: []string{
				"sk-provider-secret",
				"invalid key",
			},
		},
		{
			name:     "google",
			provider: "google",
			client:   NewGoogleClient("AIza-provider-secret", "gemini-2.5-flash"),
			status:   http.StatusTooManyRequests,
			body:     `{"error":"quota exceeded for key AIza-provider-secret"}`,
			sensitive: []string{
				"AIza-provider-secret",
				"quota exceeded for key",
			},
		},
		{
			name:     "llama",
			provider: "llama",
			client:   NewLlamaClient("gsk-provider-secret", "llama-test", "https://llama.example.invalid/openai/v1"),
			status:   http.StatusUnauthorized,
			body:     `{"error":"Authorization Bearer gsk-provider-secret rejected"}`,
			sensitive: []string{
				"gsk-provider-secret",
				"Authorization Bearer",
			},
		},
		{
			name:     "grok",
			provider: "grok",
			client:   NewGrokClient("xai-provider-secret", "grok-beta"),
			status:   http.StatusForbidden,
			body:     `{"error":"forbidden bearer xai-provider-secret"}`,
			sensitive: []string{
				"xai-provider-secret",
				"forbidden bearer",
			},
		},
		{
			name:     "anthropic",
			provider: "anthropic",
			client:   NewAnthropicClient("ant-provider-secret", "claude-test"),
			status:   http.StatusUnauthorized,
			body:     `{"error":"x-api-key ant-provider-secret invalid"}`,
			sensitive: []string{
				"ant-provider-secret",
				"x-api-key",
			},
		},
		{
			name:     "solar",
			provider: "solar",
			client:   NewSolarClient("upstage-provider-secret", "solar-pro"),
			status:   http.StatusServiceUnavailable,
			body:     `{"error":"api_key=upstage-provider-secret unavailable"}`,
			sensitive: []string{
				"upstage-provider-secret",
				"api_key=",
			},
		},
		{
			name:     "hyperclova",
			provider: "hyperclova",
			client:   NewHyperCLOVAClient("clova-provider-secret", "HCX-DASH-001"),
			status:   http.StatusInternalServerError,
			body:     `{"error":"token=clova-provider-secret failed"}`,
			sensitive: []string{
				"clova-provider-secret",
				"token=",
			},
		},
		{
			name:     "exaone",
			provider: "exaone",
			client:   NewExaoneClient("exaone-provider-secret", "exaone-test"),
			status:   http.StatusTeapot,
			body:     `{"error":"provider echoed key exaone-provider-secret"}`,
			sensitive: []string{
				"exaone-provider-secret",
				"provider echoed key",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			installClientHTTP(t, tc.client, &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tc.status,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(tc.body)),
					Request:    req,
				}, nil
			})})

			_, err := tc.client.Complete(context.Background(), "test prompt")
			if err == nil {
				t.Fatal("Complete() error = nil, want sanitized provider error")
			}
			got := err.Error()
			want := tc.provider + ": API error"
			if !strings.Contains(got, want) {
				t.Fatalf("Complete() error = %q, want prefix containing %q", got, want)
			}
			for _, fragment := range tc.sensitive {
				if strings.Contains(got, fragment) {
					t.Fatalf("Complete() error leaked sensitive fragment %q in %q", fragment, got)
				}
			}
		})
	}
}

func TestProviderClientsDoNotReturnRawRequestErrors(t *testing.T) {
	secret := "AIza-provider-secret"
	client := NewGoogleClient(secret, "gemini-2.5-flash")
	installClientHTTP(t, client, &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("Post \"" + req.URL.String() + "\": dial tcp: i/o timeout")
	})})

	_, err := client.Complete(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("Complete() error = nil, want sanitized request error")
	}
	got := err.Error()
	if got != "google: request timeout" {
		t.Fatalf("Complete() error = %q, want %q", got, "google: request timeout")
	}
	for _, fragment := range []string{secret, "key=" + secret, "generativelanguage.googleapis.com"} {
		if strings.Contains(got, fragment) {
			t.Fatalf("Complete() error leaked sensitive fragment %q in %q", fragment, got)
		}
	}
}

func installClientHTTP(t *testing.T, client Client, httpClient *http.Client) {
	t.Helper()
	switch c := client.(type) {
	case *OpenAIClient:
		c.http = httpClient
	case *GoogleClient:
		c.http = httpClient
	case *LlamaClient:
		c.http = httpClient
	case *GrokClient:
		c.http = httpClient
	case *AnthropicClient:
		c.http = httpClient
	case *SolarClient:
		c.http = httpClient
	case *HyperCLOVAClient:
		c.http = httpClient
	case *ExaoneClient:
		c.http = httpClient
	default:
		t.Fatalf("unsupported client type %T", client)
	}
}
