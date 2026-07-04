package keyvalidator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/learnweaver/backend/internal/pkg/safehttp"
)

type Result struct {
	Valid   bool
	Message string
}

func Validate(ctx context.Context, provider, apiKey string, endpointURL *string) Result {
	provider = strings.TrimSpace(strings.ToLower(provider))
	apiKey = strings.TrimSpace(apiKey)
	if provider == "" {
		return Result{Valid: false, Message: "AI 제공자를 선택해주세요"}
	}
	if apiKey == "" {
		return Result{Valid: false, Message: "API 키를 입력해주세요"}
	}

	switch provider {
	case "openai":
		return validateOpenAI(ctx, apiKey)
	case "anthropic":
		return validateAnthropic(ctx, apiKey)
	case "google":
		return validateGoogle(ctx, apiKey)
	case "grok":
		return validateGrok(ctx, apiKey)
	case "llama", "exaone":
		return validateOpenAICompatible(ctx, apiKey, endpointURL)
	case "solar":
		return Result{Valid: false, Message: "Solar 키 검증은 아직 지원하지 않습니다"}
	case "hyperclova":
		return Result{Valid: false, Message: "HyperCLOVA X는 보조 키가 필요해 학습자 검증 UI를 추가해야 합니다"}
	default:
		return Result{Valid: false, Message: "지원하지 않는 AI 제공자입니다"}
	}
}

func validateOpenAI(ctx context.Context, apiKey string) Result {
	return doRequest(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil, map[string]string{
		"Authorization": "Bearer " + apiKey,
	})
}

func validateAnthropic(ctx context.Context, apiKey string) Result {
	return doRequest(ctx, http.MethodGet, "https://api.anthropic.com/v1/models", nil, map[string]string{
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
	})
}

func validateGoogle(ctx context.Context, apiKey string) Result {
	u, _ := url.Parse("https://generativelanguage.googleapis.com/v1beta/models")
	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()
	return doRequest(ctx, http.MethodGet, u.String(), nil, nil)
}

func validateGrok(ctx context.Context, apiKey string) Result {
	return doRequest(ctx, http.MethodGet, "https://api.x.ai/v1/models", nil, map[string]string{
		"Authorization": "Bearer " + apiKey,
	})
}

func validateOpenAICompatible(ctx context.Context, apiKey string, endpointURL *string) Result {
	if endpointURL == nil || strings.TrimSpace(*endpointURL) == "" {
		return Result{Valid: false, Message: "엔드포인트 URL을 입력해주세요"}
	}

	base := strings.TrimRight(strings.TrimSpace(*endpointURL), "/")
	target := base
	if !strings.HasSuffix(base, "/v1/models") {
		if strings.HasSuffix(base, "/v1") {
			target = base + "/models"
		} else if strings.HasSuffix(base, "/models") {
			target = base
		} else {
			target = base + "/v1/models"
		}
	}

	return doRequest(ctx, http.MethodGet, target, nil, map[string]string{
		"Authorization": "Bearer " + apiKey,
	})
}

func doRequest(ctx context.Context, method, target string, body io.Reader, headers map[string]string) Result {
	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, method, target, body)
	if err != nil {
		return Result{Valid: false, Message: "검증 요청 생성 실패"}
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := safehttp.DefaultClient().Do(req)
	if err != nil {
		return Result{Valid: false, Message: "검증 요청 실패"}
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return Result{Valid: true, Message: "유효한 API 키입니다"}
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return Result{Valid: false, Message: "유효하지 않거나 권한이 없는 API 키입니다"}
	default:
		return Result{Valid: false, Message: fmt.Sprintf("검증 실패 (%d)", resp.StatusCode)}
	}
}
