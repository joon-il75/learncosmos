package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const googleGenerativeLanguageURL = "https://generativelanguage.googleapis.com/v1beta/models"

type GoogleClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewGoogleClient(apiKey, model string) *GoogleClient {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &GoogleClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (c *GoogleClient) Complete(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{
				{"text": "You generate structured learning curriculum drafts in strict JSON. Follow the user's topic exactly. Prefer Korean for all human-readable fields unless the user explicitly asks for another language."},
			},
		},
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("google: marshal request: %w", err)
	}

	target := fmt.Sprintf("%s/%s:generateContent?key=%s", googleGenerativeLanguageURL, c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("google: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("google", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("google", resp.StatusCode)
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("google: decode response: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("google: empty candidates")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}
