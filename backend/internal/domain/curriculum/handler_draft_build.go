package curriculum

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

func (h *Handler) buildEmptyDraft(userID uuid.UUID, req CreateCourseDraftRequest) *draftBuildDecision {
	draft := &DraftAggregate{
		Draft: CourseDraft{
			ID:                 uuid.New(),
			UserID:             userID,
			SourceQuery:        req.SourceQuery,
			LearningGoal:       req.LearningGoal,
			GoalProfileID:      req.GoalProfileID,
			GoalProfileVersion: req.GoalProfileVersion,
			CurrentLevel:       req.CurrentLevel,
			DurationWeeks:      req.DurationWeeks,
			StudyHoursPerWeek:  req.StudyHoursPerWeek,
			PreferredFormat:    req.PreferredFormat,
			GenerationLanguage: normalizeLearningLanguage(req.LearningLanguage),
			Title:              req.SourceQuery,
			Status:             DraftStatusDraft,
		},
		Lessons: []DraftLessonTree{},
	}
	return &draftBuildDecision{
		Draft:         draft,
		BuildMode:     "skip_generation",
		BillingStatus: "free",
		ShouldCharge:  false,
	}
}

func (h *Handler) buildDraft(ctx context.Context, userID uuid.UUID, req CreateCourseDraftRequest, policyCost int) (*draftBuildDecision, error) {
	totalBuildStartedAt := time.Now()
	patternMatchMS := int64(0)
	routing := analyzeGoalRouting(req)
	patternKey := formatGoalPatternKeyForLog(routing.PatternKey)
	subpatternKey := strings.TrimSpace(routing.SubpatternKey)
	if subpatternKey == "" {
		subpatternKey = "none"
	}
	refinementMode := strings.TrimSpace(routing.RefinementMode)
	if refinementMode == "" {
		refinementMode = "unclassified"
	}
	patternMatchMode := strings.TrimSpace(strings.ToLower(os.Getenv("CURRICULUM_PATTERN_MATCH_MODE")))
	if patternMatchMode == "embedding_primary" {
		matchStartedAt := time.Now()
		matchCtx, cancelMatch := context.WithTimeout(ctx, 8*time.Second)
		result, matchErr := h.resolveCurriculumPatternMatches(matchCtx, req, 3)
		cancelMatch()
		patternMatchMS = time.Since(matchStartedAt).Milliseconds()
		if matchErr != nil {
			log.Printf("[curriculum/pattern_shadow] user=%s mode=%s code_pattern=%s code_subpattern=%s refinement_mode=%s fallback_reason=%s duration_ms=%d", userID, patternMatchMode, patternKey, subpatternKey, refinementMode, truncateShadowError(matchErr), patternMatchMS)
		} else {
			h.logCurriculumPatternMatchResult(userID, patternMatchMode, patternKey, subpatternKey, refinementMode, result, matchStartedAt)
			guidanceMatches := FilterCurriculumPatternGuidanceMatches(result.Matches, routing.PatternKey.DomainAxis, routing.PatternKey.GoalModePrimary, routing.SubpatternKey, 2, 0.30)
			req.PatternMatchGuidance = BuildCurriculumPatternMatchGuidanceForLanguage(guidanceMatches, 2, req.LearningLanguage)
			if len(guidanceMatches) > 0 {
				req.PatternSearchSpecTemplate = guidanceMatches[0].RecommendationSearchSpecTemplate
			}
		}
	} else {
		h.logCurriculumPatternShadowAsync(userID, req, patternKey, subpatternKey, refinementMode)
	}

	queueStartedAt := time.Now()
	acquireCtx, cancelAcquire := context.WithTimeout(ctx, h.queueTimeout)
	defer cancelAcquire()
	if err := h.buildLimiter.Acquire(acquireCtx); err != nil {
		snapshot := h.buildLimiter.Snapshot()
		log.Printf("[curriculum/create] user=%s mode=build_limiter_busy queue_wait_ms=%d limiter_inflight_count=%d limiter_queue_depth=%d limiter_capacity=%d source_query_summary=%s err=%s", userID, time.Since(queueStartedAt).Milliseconds(), snapshot.InFlight, snapshot.QueueDepth, snapshot.Capacity, logsafe.Summary(req.SourceQuery), logsafe.Error(err))
		return nil, err
	}
	defer h.buildLimiter.Release()
	limiterSnapshot := h.buildLimiter.Snapshot()

	skipReview := true
	reviewStatus := "disabled_by_policy"
	buildOpts := DraftBuildOptions{
		SkipReview:        skipReview,
		GenerationTimeout: h.genTimeout,
		ReviewTimeout:     h.reviewTimeout,
	}
	queueWaitMS := time.Since(queueStartedAt).Milliseconds()

	userConfig, err := h.repo.GetUserRuntimeAIConfig(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userConfig != nil && userConfig.Mode == "byok" && strings.TrimSpace(userConfig.APIKey) != "" {
		provider := strings.TrimSpace(strings.ToLower(userConfig.Provider))
		if isFreeBYOKProvider(provider) {
			client, err := llm.NewClient(provider, userConfig.APIKey, "")
			if err != nil {
				return nil, fmt.Errorf("byok llm client init failed: %w", err)
			}
			client = h.trackBYOKLLMUsage(client, userID, provider, "course_draft_create", "byok_no_charge")
			timedClient := newTimedLLMClient(client)

			llmCtx, cancel := context.WithTimeout(ctx, h.genTimeout+h.reviewTimeout+2*time.Second)
			defer cancel()

			draft, llmErr := h.service.BuildDraftWithLLM(llmCtx, timedClient, userID, req, buildOpts)
			if llmErr != nil {
				go h.repo.RecordBYOKValidation(context.Background(), userID, provider, false, llmErr.Error())
				return nil, fmt.Errorf("byok llm generation failed: %w", llmErr)
			}
			go h.repo.RecordBYOKValidation(context.Background(), userID, provider, true, "")
			return &draftBuildDecision{
				Draft:             draft,
				BuildMode:         "user_byok",
				BillingStatus:     "byok_no_charge",
				EffectiveCost:     0,
				ShouldCharge:      false,
				ReviewStatus:      reviewStatus,
				QueueWaitMS:       queueWaitMS,
				PatternMatchMS:    patternMatchMS,
				LLMElapsedMS:      timedClient.ElapsedMS(),
				TotalBuildMS:      time.Since(totalBuildStartedAt).Milliseconds(),
				LimiterInFlight:   limiterSnapshot.InFlight,
				LimiterQueueDepth: limiterSnapshot.QueueDepth,
				PatternKey:        patternKey,
				SubpatternKey:     subpatternKey,
				RefinementMode:    refinementMode,
			}, nil
		}
	}

	selection, err := h.resolveDraftSystemLLM(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("system llm not configured: %w", err)
	}

	apiKey, err := h.settingsStore.ResolveAPIKey(ctx, selection.Setting.Provider, h.openAIAPIKey)
	if err != nil || strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("system llm api key unavailable")
	}

	client, err := llm.NewClient(strings.TrimSpace(selection.Setting.Provider), apiKey, strings.TrimSpace(selection.Setting.Model))
	if err != nil {
		return nil, fmt.Errorf("system llm client init failed: %w", err)
	}

	timedClient := newTimedLLMClient(client)

	llmCtx, cancel := context.WithTimeout(ctx, h.genTimeout+h.reviewTimeout+2*time.Second)
	defer cancel()

	draft, err := h.service.BuildDraftWithLLM(llmCtx, timedClient, userID, req, buildOpts)
	if err != nil {
		return nil, fmt.Errorf("system llm generation failed: %w", err)
	}
	return &draftBuildDecision{
		Draft:             draft,
		BuildMode:         buildModeForSystemFeature(selection.Feature),
		BillingStatus:     "charged",
		EffectiveCost:     policyCost,
		ShouldCharge:      true,
		ReviewStatus:      reviewStatus,
		QueueWaitMS:       queueWaitMS,
		PatternMatchMS:    patternMatchMS,
		LLMElapsedMS:      timedClient.ElapsedMS(),
		TotalBuildMS:      time.Since(totalBuildStartedAt).Milliseconds(),
		LimiterInFlight:   limiterSnapshot.InFlight,
		LimiterQueueDepth: limiterSnapshot.QueueDepth,
		PatternKey:        patternKey,
		SubpatternKey:     subpatternKey,
		RefinementMode:    refinementMode,
	}, nil
}

func (h *Handler) resolveDraftSystemLLM(ctx context.Context, userID uuid.UUID) (*systemLLMSelection, error) {
	feature := "default"

	premiumAccess, err := h.repo.UserHasPremiumAccess(ctx, userID)
	if err == nil && premiumAccess {
		if setting, settingErr := h.repo.GetLLMSetting(ctx, "pro_curriculum"); settingErr == nil {
			return &systemLLMSelection{
				Feature: "pro_curriculum",
				Setting: setting,
			}, nil
		}
	}

	setting, err := h.repo.GetLLMSetting(ctx, feature)
	if err != nil {
		return nil, err
	}
	return &systemLLMSelection{
		Feature: feature,
		Setting: setting,
	}, nil
}

func buildModeForSystemFeature(feature string) string {
	switch strings.TrimSpace(strings.ToLower(feature)) {
	case "pro_curriculum":
		return "system_llm_pro_curriculum"
	default:
		return "system_llm_default"
	}
}
