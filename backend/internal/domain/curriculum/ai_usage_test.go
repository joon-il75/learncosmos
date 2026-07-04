package curriculum

import (
	"strings"
	"testing"
)

func TestTrimUsageErrorRedactsProviderSecretMaterial(t *testing.T) {
	cases := []struct {
		name      string
		in        string
		want      string
		sensitive []string
	}{
		{
			name: "google url key in network error",
			in:   "google: do request: Post \"https://generativelanguage.googleapis.com/v1beta/models/gemini:generateContent?key=AIza-secret-value\": dial tcp: i/o timeout",
			want: "provider_timeout",
			sensitive: []string{
				"AIza-secret-value",
				"key=AIza-secret-value",
				"generativelanguage.googleapis.com",
			},
		},
		{
			name: "openai invalid key response",
			in:   "openai: API error 401: invalid api key sk-secret-value",
			want: "provider_unauthorized",
			sensitive: []string{
				"sk-secret-value",
				"invalid api key sk-secret-value",
			},
		},
		{
			name: "google quota response",
			in:   "google: API error 429: quota exceeded for key AIza-secret-value",
			want: "provider_rate_limited",
			sensitive: []string{
				"AIza-secret-value",
				"quota exceeded for key",
			},
		},
		{
			name: "llama openai compatible bearer response",
			in:   "llama: API error 401: Authorization Bearer gsk-provider-secret rejected by upstream",
			want: "provider_unauthorized",
			sensitive: []string{
				"gsk-provider-secret",
				"Authorization Bearer",
			},
		},
		{
			name: "grok openai compatible upstream body",
			in:   "grok: API error 403: forbidden for bearer xai-provider-secret on /v1/chat/completions",
			want: "provider_forbidden",
			sensitive: []string{
				"xai-provider-secret",
				"/v1/chat/completions",
			},
		},
		{
			name: "anthropic x api key response",
			in:   "anthropic: API error 401: x-api-key ant-provider-secret is invalid",
			want: "provider_unauthorized",
			sensitive: []string{
				"ant-provider-secret",
				"x-api-key",
			},
		},
		{
			name: "solar upstream server body",
			in:   "solar: API error 503: upstream unavailable request_id=req-secret-123 api_key=upstage-provider-secret",
			want: "provider_unavailable",
			sensitive: []string{
				"upstage-provider-secret",
				"req-secret-123",
			},
		},
		{
			name: "hyperclova raw api error",
			in:   "hyperclova: API error 500: internal error token=clova-provider-secret",
			want: "provider_unavailable",
			sensitive: []string{
				"clova-provider-secret",
				"token=",
			},
		},
		{
			name: "exaone generic api error",
			in:   "exaone: API error 418: provider echoed key exaone-provider-secret",
			want: "provider_api_error",
			sensitive: []string{
				"exaone-provider-secret",
				"provider echoed key",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := trimUsageError(tc.in)
			if got != tc.want {
				t.Fatalf("trimUsageError() = %q, want %q", got, tc.want)
			}
			for _, fragment := range tc.sensitive {
				if strings.Contains(got, fragment) {
					t.Fatalf("trimUsageError returned sensitive fragment %q in %q", fragment, got)
				}
			}
			if strings.Contains(got, tc.in) {
				t.Fatalf("trimUsageError returned raw provider error")
			}
		})
	}
}
