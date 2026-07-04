package curriculum

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func buildLessonSearchQuery(draft CourseDraft, mainLesson CourseDraftLesson, subLesson CourseDraftLesson) normalizer.LessonSearchQuery {
	mainObjective := ""
	if mainLesson.Objective != nil {
		mainObjective = *mainLesson.Objective
	}
	lessonObjective := ""
	if subLesson.Objective != nil {
		lessonObjective = *subLesson.Objective
	}
	lessonSummary := ""
	if subLesson.Summary != nil {
		lessonSummary = *subLesson.Summary
	}
	return normalizer.BuildLessonSearchQuery(
		draft.SourceQuery,
		mainLesson.Title,
		mainObjective,
		subLesson.Title,
		lessonObjective,
		lessonSummary,
	)
}

func (h *Handler) resolveLessonSearchExecution(ctx context.Context, userID uuid.UUID, policyCost int) (*lessonSearchExecution, error) {
	userConfig, err := h.repo.GetUserRuntimeAIConfig(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userConfig != nil && userConfig.Mode == "byok" && isFreeBYOKProvider(userConfig.Provider) && strings.TrimSpace(userConfig.APIKey) != "" {
		return &lessonSearchExecution{
			Provider:      strings.TrimSpace(strings.ToLower(userConfig.Provider)),
			APIKey:        "",
			BillingStatus: "byok_no_charge",
			EffectiveCost: 0,
			ShouldCharge:  false,
		}, nil
	}
	return decideLessonSearchExecution(policyCost), nil
}

func decideLessonSearchExecution(policyCost int) *lessonSearchExecution {
	return &lessonSearchExecution{
		Provider:      primaryEmbeddingProvider(),
		APIKey:        "",
		BillingStatus: "charged",
		EffectiveCost: policyCost,
		ShouldCharge:  policyCost > 0,
	}
}

func float32VectorToPGVector(vector []float32) string {
	parts := make([]string, len(vector))
	for i, v := range vector {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
