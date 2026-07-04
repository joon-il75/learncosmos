package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const openAIBaseURL = "https://api.openai.com/v1"

type OpenAIClient struct {
	apiKey   string
	model    string
	http     *http.Client
	usage    CompletionUsage
	hasUsage bool
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (c *OpenAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	c.hasUsage = false
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
		return "", fmt.Errorf("openai: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("openai: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("openai", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("openai", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("openai: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai: empty choices")
	}
	if result.Usage.TotalTokens > 0 || result.Usage.PromptTokens > 0 || result.Usage.CompletionTokens > 0 {
		c.usage = CompletionUsage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
			TotalTokens:  result.Usage.TotalTokens,
		}
		c.hasUsage = true
	}
	return result.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) LastUsage() (CompletionUsage, bool) {
	return c.usage, c.hasUsage
}

func (c *OpenAIClient) ModelName() string {
	return c.model
}
