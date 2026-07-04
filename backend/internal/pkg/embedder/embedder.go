package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	Model     = "text-embedding-3-small"
	Dimension = 1536
	apiURL    = "https://api.openai.com/v1/embeddings"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

func NewProviderClient(provider, apiKey string) (Embedder, error) {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "openai":
		return NewClient(apiKey), nil
	case "google":
		return NewGoogleClient(apiKey), nil
	case "embedding_gemma", "local_http":
		return NewLocalHTTPClient(apiKey), nil
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", provider)
	}
}

type LocalHTTPClient struct {
	endpoint   string
	model      string
	httpClient *http.Client
}

func NewLocalHTTPClient(endpoint string) *LocalHTTPClient {
	return NewLocalHTTPClientWithModel(endpoint, defaultStringFromEnv("EMBEDDING_GEMMA_MODEL", "embedding-gemma"))
}

func NewLocalHTTPClientWithModel(endpoint, model string) *LocalHTTPClient {
	model = strings.TrimSpace(model)
	if model == "" {
		model = "embedding-gemma"
	}
	return &LocalHTTPClient{
		endpoint:   strings.TrimSpace(endpoint),
		model:      model,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type localEmbedRequest struct {
	Input string `json:"input"`
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

type localEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
	Vector    []float32 `json:"vector"`
	Data      []struct {
		Embedding []float32 `json:"embedding"`
		Vector    []float32 `json:"vector"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *LocalHTTPClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if strings.TrimSpace(c.endpoint) == "" {
		return nil, fmt.Errorf("EMBEDDING_GEMMA_ENDPOINT is not configured")
	}
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("임베딩할 텍스트가 비어있습니다")
	}
	runes := []rune(text)
	if len(runes) > 8000 {
		text = string(runes[:8000])
	}

	body, _ := json.Marshal(localEmbedRequest{Input: text, Text: text, Model: c.model})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = resp.Status
		}
		return nil, fmt.Errorf("EmbeddingGemma API HTTP %d: %s", resp.StatusCode, message)
	}

	var result localEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("EmbeddingGemma API 오류: %s", result.Error.Message)
	}
	vector := result.Embedding
	if len(vector) == 0 {
		vector = result.Vector
	}
	if len(vector) == 0 && len(result.Data) > 0 {
		vector = result.Data[0].Embedding
		if len(vector) == 0 {
			vector = result.Data[0].Vector
		}
	}
	if len(vector) == 0 {
		return nil, fmt.Errorf("임베딩 결과가 비어있습니다")
	}
	return vector, nil
}

func defaultStringFromEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

type embedRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type embedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Embed — 텍스트를 1536차원 벡터로 변환
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY가 설정되지 않았습니다")
	}
	if text == "" {
		return nil, fmt.Errorf("임베딩할 텍스트가 비어있습니다")
	}

	// 텍스트 최대 8000자 트림 (토큰 한도 초과 방지)
	runes := []rune(text)
	if len(runes) > 8000 {
		text = string(runes[:8000])
	}

	body, _ := json.Marshal(embedRequest{Input: text, Model: Model})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("OpenAI API 오류: %s", result.Error.Message)
	}
	if len(result.Data) == 0 || len(result.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("임베딩 결과가 비어있습니다")
	}

	return result.Data[0].Embedding, nil
}

// BuildText — title + description 조합 (임베딩용 입력 텍스트 생성)
func BuildText(title, description string) string {
	if description == "" {
		return title
	}
	return title + "\n\n" + description
}
