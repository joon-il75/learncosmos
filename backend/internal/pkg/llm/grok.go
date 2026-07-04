package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const grokBaseURL = "https://api.x.ai/v1"

type GrokClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewGrokClient(apiKey, model string) *GrokClient {
	if model == "" {
		model = "grok-beta"
	}
	return &GrokClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

// Complete implements the LLMClient interface (OpenAI-compatible endpoint)
func (c *GrokClient) Complete(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("grok: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		grokBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("grok: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("grok", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("grok", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("grok: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("grok: empty choices")
	}
	return result.Choices[0].Message.Content, nil
}
