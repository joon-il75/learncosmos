package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const hyperclovaBaseURL = "https://clovastudio.stream.naver.com/testapp/v1/chat-completions"

type HyperCLOVAClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewHyperCLOVAClient(apiKey, model string) *HyperCLOVAClient {
	if model == "" {
		model = "HCX-DASH-001"
	}
	return &HyperCLOVAClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

func (c *HyperCLOVAClient) Complete(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]any{
		"messages": []map[string]string{
			{"role": "system", "content": "You generate structured learning curriculum drafts in strict JSON. Follow the user's topic exactly. Prefer Korean for all human-readable fields unless the user explicitly asks for another language."},
			{"role": "user", "content": prompt},
		},
		"maxTokens":   4096,
		"temperature": 0.5,
		"topP":        0.8,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("hyperclova: marshal request: %w", err)
	}

	target := fmt.Sprintf("%s/%s", hyperclovaBaseURL, c.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("hyperclova: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-NCP-CLOVASTUDIO-REQUEST-ID", "learnweaver")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", safeLLMRequestError("hyperclova", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", safeLLMAPIError("hyperclova", resp.StatusCode)
	}

	var result struct {
		Result struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("hyperclova: decode response: %w", err)
	}
	if result.Result.Message.Content == "" {
		return "", fmt.Errorf("hyperclova: empty content")
	}
	return result.Result.Message.Content, nil
}
