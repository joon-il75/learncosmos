package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const exaoneBaseURL = "https://api.lgresearch.ai/v1"

type ExaoneClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewExaoneClient(apiKey, model string) *ExaoneClient {
	if model == "" {
		model = "exaone-3.5-32b-instruct"
	}
	return &ExaoneClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (c *ExaoneClient) Complete(ctx context.Context, prompt string) (string, error) {
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
		return "", fmt.Errorf("exaone: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, exaoneBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("exaone: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("exaone", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("exaone", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("exaone: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("exaone: empty choices")
	}
	return result.Choices[0].Message.Content, nil
}
