package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/learnweaver/backend/internal/pkg/safehttp"
)

// llama 기본 엔드포인트 (Groq 호환). BYOK 설정의 endpoint_url로 재지정 가능.
const llamaDefaultBaseURL = "https://api.groq.com/openai/v1"

type LlamaClient struct {
	apiKey  string
	model   string
	baseURL string
	http    *http.Client
}

// NewLlamaClient creates a Llama client. endpointURL이 비어 있으면 Groq 기본 URL을 사용한다.
func NewLlamaClient(apiKey, model, endpointURL string) *LlamaClient {
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}
	base := llamaDefaultBaseURL
	if endpointURL != "" {
		base = strings.TrimRight(strings.TrimSpace(endpointURL), "/")
	}
	return &LlamaClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: base,
		http: safehttp.NewClient(safehttp.ClientConfig{
			Timeout:      30 * time.Second,
			MaxRedirects: 3,
		}),
	}
}

func (c *LlamaClient) Complete(ctx context.Context, prompt string) (string, error) {
	if err := validateLlamaBaseURL(c.baseURL); err != nil {
		return "", err
	}
	reqBody := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": "You generate structured learning curriculum drafts in strict JSON. Follow the user's topic exactly. Prefer Korean for all human-readable fields unless the user explicitly asks for another language."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("llama: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llama: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("llama", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("llama", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("llama: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("llama: empty choices")
	}
	return result.Choices[0].Message.Content, nil
}

func validateLlamaBaseURL(baseURL string) error {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsed == nil || parsed.Host == "" || parsed.Hostname() == "" {
		return fmt.Errorf("llama: invalid endpoint")
	}
	if parsed.User != nil {
		return fmt.Errorf("llama: invalid endpoint")
	}
	if strings.ToLower(parsed.Scheme) != "https" {
		return fmt.Errorf("llama: endpoint must use https")
	}
	return nil
}
