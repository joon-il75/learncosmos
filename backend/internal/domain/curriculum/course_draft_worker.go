package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

type courseDraftCreateProcessor struct {
	h *Handler
}

type courseDraftCreateJobRef struct {
	GoalProfileID      string `json:"goal_profile_id"`
	GoalProfileVersion int    `json:"goal_profile_version"`
	LearningLanguage   string `json:"learning_language"`
}

func NewCourseDraftCreateProcessor(handler *Handler) llmjobs.Processor {
	return &courseDraftCreateProcessor{h: handler}
}

func (p *courseDraftCreateProcessor) Process(ctx context.Context, job *llmjobs.Job) (llmjobs.CompleteJobInput, error) {
	if p == nil || p.h == nil || job == nil || job.UserID == nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("llm_generation_failed", "invalid course draft job", false)
	}
	startedAt := time.Now()
	var ref courseDraftCreateJobRef
	if err := json.Unmarshal(job.RequestRef, &ref); err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("llm_generation_failed", "invalid course draft job payload", false)
	}
	goalID, err := uuid.Parse(strings.TrimSpace(ref.GoalProfileID))
	if err != nil {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("confirmed_goal_required", "invalid goal profile reference", false)
	}
	activeGoal, err := p.h.goalRepo.GetGoalByID(ctx, goalID)
	if err != nil {
		code := "confirmed_goal_required"
		if !errors.Is(err, goal.ErrNotFound) {
			code = "llm_generation_failed"
		}
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError(code, err.Error(), false)
	}
	if activeGoal.UserID != *job.UserID || !activeGoal.IsActive || activeGoal.CourseDraftID != nil || activeGoal.InterviewState != goal.StateConfirmed || activeGoal.ConfirmedGoal == nil || strings.TrimSpace(*activeGoal.ConfirmedGoal) == "" {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("confirmed_goal_required", "active confirmed goal is no longer available", false)
	}
	if ref.GoalProfileVersion > 0 && activeGoal.Version != ref.GoalProfileVersion {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("goal_profile_version_changed", "active goal version changed", false)
	}

	release, acquired := p.h.goalBuildLock.TryAcquire(*job.UserID, activeGoal.ID)
	if !acquired {
		return llmjobs.CompleteJobInput{}, llmjobs.NewProcessorError("course_generation_already_in_progress", "course generation already in progress", true)
	}
	defer release()

	cost, err := p.h.resolveDraftBuildBilling(ctx, *job.UserID)
	if err != nil {
		return llmjobs.CompleteJobInput{}, mapDraftCreateProcessorError(err)
	}
	req := buildCreateDraftRequestFromGoal(activeGoal, ref.LearningLanguage)
	result, err := p.h.executeCourseDraftCreate(ctx, *job.UserID, req, cost, startedAt)
	if err != nil {
		return llmjobs.CompleteJobInput{}, err
	}
	return llmjobs.CompleteJobInput{
		ResultRef:              result,
		LLMGenerationElapsedMS: intPointerFromAny(result["llm_generation_elapsed_ms"]),
		PersistElapsedMS:       intPointerFromAny(result["draft_persist_elapsed_ms"]),
	}, nil
}

func buildCreateDraftRequestFromGoal(activeGoal *goal.GoalProfile, learningLanguage string) CreateCourseDraftRequest {
	confirmedGoal := strings.TrimSpace(*activeGoal.ConfirmedGoal)
	language := normalizeLearningLanguage(learningLanguage)
	if language == "ko" {
		language = normalizeLearningLanguage(activeGoal.Language)
	}
	return CreateCourseDraftRequest{
		SourceQuery:         confirmedGoal,
		LearningGoal:        &confirmedGoal,
		GoalProfileID:       &activeGoal.ID,
		GoalProfileVersion:  &activeGoal.Version,
		GoalUserIntent:      &activeGoal.UserIntent,
		GoalMotivation:      activeGoal.Motivation,
		GoalUsageContext:    activeGoal.UsageContext,
		GoalDifficultyLevel: activeGoal.DifficultyLevel,
		GoalTimeHorizon:     activeGoal.TimeHorizon,
		GoalOutputType:      activeGoal.OutputType,
		GoalType:            activeGoal.GoalType,
		LearningIntent:      learningIntentProfileFromGoal(activeGoal.LearningIntent),
		LearningLanguage:    language,
	}
}

func (h *Handler) executeCourseDraftCreate(ctx context.Context, userID uuid.UUID, req CreateCourseDraftRequest, cost int, startedAt time.Time) (map[string]any, error) {
	buildStartedAt := time.Now()
	sourceQuerySummary := logsafe.Summary(req.SourceQuery)
	decision, buildErr := h.buildDraft(ctx, userID, req, cost)
	if buildErr != nil {
		log.Printf("[curriculum/create-worker] user=%s mode=build_failed source_query_summary=%s elapsed_ms=%d err=%s", userID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(buildErr))
		if errors.Is(buildErr, errDraftBuildBusy) {
			return nil, llmjobs.NewProcessorError("course_generation_busy", "현재 코스 생성 요청이 많습니다. 잠시 후 다시 시도해 주세요.", true)
		}
		if errors.Is(buildErr, errDraftGenerationRetryable) {
			return nil, llmjobs.NewProcessorError("draft_generation_retry", "draft generation retry required", true)
		}
		return nil, llmjobs.NewProcessorError("llm_generation_failed", buildErr.Error(), false)
	}
	log.Printf("[curriculum/create-worker] user=%s build_mode=%s review_status=%s queue_wait_ms=%d llm_generation_elapsed_ms=%d pattern_match_elapsed_ms=%d total_build_elapsed_ms=%d limiter_inflight_count=%d limiter_queue_depth=%d pattern_key=%s subpattern_key=%s refinement_mode=%s source_query_summary=%s build_elapsed_ms=%d", userID, decision.BuildMode, decision.ReviewStatus, decision.QueueWaitMS, decision.LLMElapsedMS, decision.PatternMatchMS, decision.TotalBuildMS, decision.LimiterInFlight, decision.LimiterQueueDepth, decision.PatternKey, decision.SubpatternKey, decision.RefinementMode, sourceQuerySummary, time.Since(buildStartedAt).Milliseconds())

	chargeAmount := 0
	if decision.ShouldCharge {
		chargeAmount = decision.EffectiveCost
	}
	confirmed, err := h.service.BuildConfirmedCourseFromDraft(*decision.Draft)
	if err != nil {
		return nil, llmjobs.NewProcessorError("llm_generation_failed", err.Error(), false)
	}
	persistStartedAt := time.Now()
	if err := h.repo.CreateDraftAndStartLearning(ctx, decision.Draft, confirmed, chargeAmount); err != nil {
		if err.Error() == "insufficient_points" {
			return nil, llmjobs.NewProcessorError("insufficient_points", "insufficient points", false)
		}
		if errors.Is(err, ErrActiveGoalAlreadyAttached) {
			return nil, llmjobs.NewProcessorError("active_goal_already_attached", "active goal already attached", false)
		}
		return nil, llmjobs.NewProcessorError("llm_generation_failed", err.Error(), false)
	}
	draftPersistElapsedMS := time.Since(persistStartedAt).Milliseconds()
	fetchedDraft, err := h.repo.GetDraftByID(ctx, decision.Draft.Draft.ID, userID)
	if err != nil {
		return nil, llmjobs.NewProcessorError("llm_generation_failed", "draft created but failed to reload", false)
	}
	freeBalance, paidBalance, err := h.repo.GetPointBalances(ctx, userID)
	if err != nil {
		return nil, llmjobs.NewProcessorError("llm_generation_failed", "draft created but failed to load point balance", false)
	}
	h.observeCourseDraftAIOutput(ctx, userID, fetchedDraft, confirmed.Course.ID, "course_draft_create_worker", map[string]any{
		"build_mode":      decision.BuildMode,
		"billing_status":  decision.BillingStatus,
		"pattern_key":     decision.PatternKey,
		"subpattern_key":  decision.SubpatternKey,
		"refinement_mode": decision.RefinementMode,
	})
	log.Printf("[curriculum/create-worker] user=%s draft_id=%s build_mode=%s billing_status=%s should_charge=%t queue_wait_ms=%d llm_generation_elapsed_ms=%d pattern_match_elapsed_ms=%d draft_persist_elapsed_ms=%d total_build_elapsed_ms=%d limiter_inflight_count=%d limiter_queue_depth=%d pattern_key=%s subpattern_key=%s refinement_mode=%s source_query_summary=%s total_elapsed_ms=%d lesson_count=%d", userID, decision.Draft.Draft.ID, decision.BuildMode, decision.BillingStatus, decision.ShouldCharge, decision.QueueWaitMS, decision.LLMElapsedMS, decision.PatternMatchMS, draftPersistElapsedMS, decision.TotalBuildMS, decision.LimiterInFlight, decision.LimiterQueueDepth, decision.PatternKey, decision.SubpatternKey, decision.RefinementMode, sourceQuerySummary, time.Since(startedAt).Milliseconds(), len(fetchedDraft.Lessons))
	return map[string]any{
		"draft_id":                  fetchedDraft.Draft.ID,
		"course_id":                 confirmed.Course.ID,
		"billing_status":            decision.BillingStatus,
		"build_mode":                decision.BuildMode,
		"policy_cost":               cost,
		"point_preview":             CourseDraftPointPreview{Cost: decision.EffectiveCost, FreeBalance: freeBalance, PaidBalance: paidBalance, TotalBalance: freeBalance + paidBalance},
		"build_metrics":             map[string]any{"queue_wait_ms": decision.QueueWaitMS, "llm_generation_elapsed_ms": decision.LLMElapsedMS, "pattern_match_elapsed_ms": decision.PatternMatchMS, "draft_persist_elapsed_ms": draftPersistElapsedMS, "total_build_elapsed_ms": decision.TotalBuildMS},
		"queue_wait_ms":             decision.QueueWaitMS,
		"llm_generation_elapsed_ms": decision.LLMElapsedMS,
		"pattern_match_elapsed_ms":  decision.PatternMatchMS,
		"draft_persist_elapsed_ms":  draftPersistElapsedMS,
		"total_build_elapsed_ms":    decision.TotalBuildMS,
		"limiter_inflight_count":    decision.LimiterInFlight,
		"limiter_queue_depth":       decision.LimiterQueueDepth,
		"pattern_key":               decision.PatternKey,
		"subpattern_key":            decision.SubpatternKey,
		"refinement_mode":           decision.RefinementMode,
		"lesson_count":              len(fetchedDraft.Lessons),
	}, nil
}

func mapDraftCreateProcessorError(err error) error {
	if err == nil {
		return nil
	}
	if err.Error() == "insufficient_points" {
		return llmjobs.NewProcessorError("insufficient_points", "insufficient points", false)
	}
	switch {
	case strings.Contains(err.Error(), "point policy"):
		return llmjobs.NewProcessorError("llm_generation_failed", "failed to load point policy", false)
	case strings.Contains(err.Error(), "point balance"):
		return llmjobs.NewProcessorError("llm_generation_failed", "failed to load point balance", false)
	default:
		return llmjobs.NewProcessorError("llm_generation_failed", fmt.Sprintf("failed to load ai config: %v", err), false)
	}
}

func intPointerFromAny(value any) *int {
	switch v := value.(type) {
	case int:
		return &v
	case int64:
		vv := int(v)
		return &vv
	case float64:
		vv := int(v)
		return &vv
	default:
		return nil
	}
}
