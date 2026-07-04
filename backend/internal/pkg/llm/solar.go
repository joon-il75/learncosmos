package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const solarBaseURL = "https://api.upstage.ai/v1/solar"

type SolarClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewSolarClient(apiKey, model string) *SolarClient {
	if model == "" {
		model = "solar-pro"
	}
	return &SolarClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (c *SolarClient) Complete(ctx context.Context, prompt string) (string, error) {
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
		return "", fmt.Errorf("solar: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, solarBaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("solar: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("solar", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("solar", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("solar: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("solar: empty choices")
	}
	return result.Choices[0].Message.Content, nil
}
