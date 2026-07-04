package curriculum

import (
	"context"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

type usageTrackingLLMClient struct {
	base          llm.Client
	repo          *Repository
	userID        uuid.UUID
	provider      string
	model         string
	feature       string
	billingStatus string
	metadata      map[string]any
}

func (c *usageTrackingLLMClient) Complete(ctx context.Context, prompt string) (string, error) {
	raw, err := c.base.Complete(ctx, prompt)
	inputTokens := estimateTokenCount(prompt)
	outputTokens := 0
	if err == nil {
		outputTokens = estimateTokenCount(raw)
	}
	model := c.model
	if reporter, ok := c.base.(llm.UsageReporter); ok {
		if reportedModel := strings.TrimSpace(reporter.ModelName()); reportedModel != "" {
			model = reportedModel
		}
		if usage, ok := reporter.LastUsage(); ok {
			inputTokens = usage.InputTokens
			outputTokens = usage.OutputTokens
		}
	}
	c.record(model, inputTokens, outputTokens, err)
	return raw, err
}

func (c *usageTrackingLLMClient) record(model string, inputTokens, outputTokens int, callErr error) {
	if c == nil || c.repo == nil || c.userID == uuid.Nil {
		return
	}
	inputPtr := &inputTokens
	outputPtr := &outputTokens
	cost := estimateBYOKCostUSD(c.provider, model, inputTokens, outputTokens)
	errorCode := ""
	if callErr != nil {
		errorCode = trimUsageError(callErr.Error())
	}
	go func() {
		_ = c.repo.RecordAIUsageEvent(context.Background(), AIUsageEventInput{
			UserID:           c.userID,
			Source:           "byok",
			Provider:         c.provider,
			Model:            model,
			Feature:          c.feature,
			BillingStatus:    c.billingStatus,
			InputTokens:      inputPtr,
			OutputTokens:     outputPtr,
			EstimatedCostUSD: cost,
			Success:          callErr == nil,
			ErrorCode:        errorCode,
			Metadata:         c.metadata,
		})
	}()
}

func (h *Handler) trackBYOKLLMUsage(client llm.Client, userID uuid.UUID, provider, feature, billingStatus string) llm.Client {
	return h.trackBYOKLLMUsageWithMetadata(client, userID, provider, feature, billingStatus, nil)
}

func (h *Handler) trackBYOKLLMUsageWithMetadata(client llm.Client, userID uuid.UUID, provider, feature, billingStatus string, metadata map[string]any) llm.Client {
	if client == nil || billingStatus != "byok_no_charge" {
		return client
	}
	return &usageTrackingLLMClient{
		base:          client,
		repo:          h.repo,
		userID:        userID,
		provider:      strings.TrimSpace(strings.ToLower(provider)),
		model:         defaultBYOKModel(provider),
		feature:       feature,
		billingStatus: billingStatus,
		metadata:      metadata,
	}
}

func (h *Handler) recordBYOKEmbeddingUsage(userID uuid.UUID, provider, feature, billingStatus, input string, callErr error) {
	if billingStatus != "byok_no_charge" {
		return
	}
	inputTokens := estimateTokenCount(input)
	outputTokens := 0
	inputPtr := &inputTokens
	outputPtr := &outputTokens
	errorCode := ""
	if callErr != nil {
		errorCode = trimUsageError(callErr.Error())
	}
	go func() {
		_ = h.repo.RecordAIUsageEvent(context.Background(), AIUsageEventInput{
			UserID:        userID,
			Source:        "byok",
			Provider:      strings.TrimSpace(strings.ToLower(provider)),
			Model:         "text-embedding-3-small",
			Feature:       feature,
			BillingStatus: billingStatus,
			InputTokens:   inputPtr,
			OutputTokens:  outputPtr,
			Success:       callErr == nil,
			ErrorCode:     errorCode,
		})
	}()
}

func estimateTokenCount(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	runes := utf8.RuneCountInString(value)
	if runes == 0 {
		return 0
	}
	return int(math.Ceil(float64(runes) / 3.0))
}

func defaultBYOKModel(provider string) string {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "openai":
		return "gpt-4o-mini"
	case "google":
		return "gemini-2.5-flash"
	case "anthropic":
		return "claude-3-5-haiku-20241022"
	case "grok":
		return "grok-2-latest"
	case "solar":
		return "solar-pro"
	case "hyperclova":
		return "HCX-003"
	case "llama":
		return "llama-3.1-8b-instant"
	case "exaone":
		return "exaone-3.5-32b-instruct"
	default:
		return ""
	}
}

func estimateBYOKCostUSD(provider, model string, inputTokens, outputTokens int) *float64 {
	if strings.TrimSpace(strings.ToLower(provider)) != "openai" || !strings.HasPrefix(strings.TrimSpace(model), "gpt-4o-mini") {
		return nil
	}
	cost := float64(inputTokens)*0.15/1_000_000 + float64(outputTokens)*0.60/1_000_000
	return &cost
}

func trimUsageError(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return logsafe.ProviderErrorCode(value)
}
