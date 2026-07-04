package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

const googleEmbeddingURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent"

type GoogleClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGoogleClient(apiKey string) *GoogleClient {
	return &GoogleClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *GoogleClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY가 설정되지 않았습니다")
	}
	if text == "" {
		return nil, fmt.Errorf("임베딩할 텍스트가 비어있습니다")
	}

	runes := []rune(text)
	if len(runes) > 8000 {
		text = string(runes[:8000])
	}

	body, _ := json.Marshal(map[string]any{
		"content": map[string]any{
			"parts": []map[string]string{
				{"text": text},
			},
		},
		"output_dimensionality": Dimension,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleEmbeddingURL+"?key="+c.apiKey, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Embedding *struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("Google API 오류: %s", result.Error.Message)
	}
	if result.Embedding == nil || len(result.Embedding.Values) == 0 {
		return nil, fmt.Errorf("임베딩 결과가 비어있습니다")
	}

	return normalizeEmbedding(result.Embedding.Values), nil
}

func normalizeEmbedding(values []float32) []float32 {
	if len(values) == 0 {
		return values
	}

	var sum float64
	for _, v := range values {
		sum += float64(v * v)
	}
	if sum == 0 {
		return values
	}

	norm := float32(math.Sqrt(sum))
	out := make([]float32, len(values))
	for i, v := range values {
		out[i] = v / norm
	}
	return out
}
