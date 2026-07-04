package goal

import (
	"context"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

type usageTrackingLLMClient struct {
	base          llm.Client
	repo          *Repository
	userID        uuid.UUID
	provider      string
	model         string
	feature       string
	billingStatus string
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
		})
	}()
}

func (h *Handler) trackBYOKLLMUsage(client llm.Client, userID uuid.UUID, provider, model, feature, billingStatus string) llm.Client {
	if client == nil || billingStatus != "byok_no_charge" {
		return client
	}
	return &usageTrackingLLMClient{
		base:          client,
		repo:          h.repo,
		userID:        userID,
		provider:      strings.TrimSpace(strings.ToLower(provider)),
		model:         strings.TrimSpace(model),
		feature:       feature,
		billingStatus: billingStatus,
	}
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

func estimateBYOKCostUSD(provider, model string, inputTokens, outputTokens int) *float64 {
	if strings.TrimSpace(strings.ToLower(provider)) != "openai" || !strings.HasPrefix(strings.TrimSpace(model), "gpt-4o-mini") {
		return nil
	}
	cost := float64(inputTokens)*0.15/1_000_000 + float64(outputTokens)*0.60/1_000_000
	return &cost
}

func trimUsageError(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) > 180 {
		return string(runes[:180])
	}
	return value
}
