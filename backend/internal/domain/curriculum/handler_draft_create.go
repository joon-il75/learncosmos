package curriculum

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
)

func (h *Handler) resolveDraftBuildBilling(ctx context.Context, userID uuid.UUID) (int, error) {
	cost, err := h.repo.GetPointSetting(ctx, "course_gen_cost")
	if err != nil {
		return 0, fmt.Errorf("failed to load point policy: %w", err)
	}

	freeBalance, paidBalance, err := h.repo.GetPointBalances(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to load point balance: %w", err)
	}

	userConfig, err := h.repo.GetUserRuntimeAIConfig(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to load ai config: %w", err)
	}

	usesFreeBYOK := userConfig != nil &&
		userConfig.Mode == "byok" &&
		strings.TrimSpace(userConfig.APIKey) != "" &&
		isFreeBYOKProvider(strings.TrimSpace(strings.ToLower(userConfig.Provider)))
	if !usesFreeBYOK && cost > 0 && freeBalance+paidBalance < cost {
		return 0, fmt.Errorf("insufficient_points")
	}

	return cost, nil
}

func (h *Handler) CreateDraft(c *gin.Context) {
	startedAt := time.Now()
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateCourseDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.SourceQuery = strings.TrimSpace(req.SourceQuery)
	req.LearningLanguage = h.repo.GetUserLearningLanguage(c.Request.Context(), userID)
	if req.SourceQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_query is required"})
		return
	}
	if h.enforceSafety(c, userID, "course_draft_source_query", nil, req.SourceQuery, req.LearningLanguage, map[string]any{"field": "source_query"}) {
		return
	}
	sourceQuerySummary := logsafe.Summary(req.SourceQuery)
	var releaseGoalBuild func()
	var activeGoal *goal.GoalProfile
	if !req.SkipGeneration {
		var err error
		activeGoal, err = h.goalRepo.GetActiveGoalByUser(c.Request.Context(), userID)
		if err == goal.ErrNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed_goal_required"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load active goal"})
			return
		}
		if activeGoal.InterviewState != goal.StateConfirmed || activeGoal.ConfirmedGoal == nil || strings.TrimSpace(*activeGoal.ConfirmedGoal) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed_goal_required"})
			return
		}
		var acquired bool
		releaseGoalBuild, acquired = h.goalBuildLock.TryAcquire(userID, activeGoal.ID)
		if !acquired {
			log.Printf("[curriculum/create] user=%s mode=goal_build_conflict goal_profile_id=%s source_query_summary=%s elapsed_ms=%d", userID, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds())
			c.JSON(http.StatusConflict, gin.H{"error": "course_generation_already_in_progress"})
			return
		}
		defer releaseGoalBuild()

		confirmedGoal := strings.TrimSpace(*activeGoal.ConfirmedGoal)
		req.LearningGoal = &confirmedGoal
		req.GoalProfileID = &activeGoal.ID
		req.GoalProfileVersion = &activeGoal.Version
		req.GoalUserIntent = &activeGoal.UserIntent
		req.GoalMotivation = activeGoal.Motivation
		req.GoalUsageContext = activeGoal.UsageContext
		req.GoalDifficultyLevel = activeGoal.DifficultyLevel
		req.GoalTimeHorizon = activeGoal.TimeHorizon
		req.GoalOutputType = activeGoal.OutputType
		req.GoalType = activeGoal.GoalType
		req.LearningIntent = learningIntentProfileFromGoal(activeGoal.LearningIntent)
	}

	cost, err := h.resolveDraftBuildBilling(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "insufficient_points" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
			return
		}
		switch {
		case strings.Contains(err.Error(), "point policy"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point policy"})
		case strings.Contains(err.Error(), "point balance"):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point balance"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ai config"})
		}
		return
	}

	if !req.SkipGeneration && courseGenerationAsyncEnabled() {
		if h.llmGateway == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
			return
		}
		jobInput := CreateCourseDraftJobInput(userID, activeGoal, req.LearningLanguage, cost)
		if jobInput.IdempotencyKey != nil {
			existing, existingErr := h.llmGateway.GetActiveJobByIdempotencyKey(c.Request.Context(), jobInput.Feature, *jobInput.IdempotencyKey)
			if existingErr == nil && existing != nil {
				log.Printf("[curriculum/create] user=%s mode=async_existing_job job_id=%s goal_profile_id=%s source_query_summary=%s elapsed_ms=%d", userID, existing.ID, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds())
				h.courseGenerationJobAccepted(c, existing, cost, "코스 초안 생성이 이미 진행 중이에요.")
				return
			}
		}
		jobInput.MaxActiveJobs = courseGenerationAsyncMaxActiveJobs()
		job, submitErr := h.llmGateway.Submit(c.Request.Context(), jobInput)
		if submitErr != nil {
			if errors.Is(submitErr, llmjobs.ErrActiveJobLimitExceeded) {
				log.Printf("[curriculum/create] user=%s mode=async_active_limit_busy max_active_jobs=%d goal_profile_id=%s source_query_summary=%s elapsed_ms=%d", userID, jobInput.MaxActiveJobs, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds())
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":           "course_generation_busy",
					"error_code":      "course_generation_busy",
					"message":         "현재 코스 생성 요청이 많습니다. 잠시 후 다시 시도해 주세요.",
					"max_active_jobs": jobInput.MaxActiveJobs,
				})
				return
			}
			if errors.Is(submitErr, llmjobs.ErrDuplicateJob) && jobInput.IdempotencyKey != nil {
				existing, existingErr := h.llmGateway.GetActiveJobByIdempotencyKey(c.Request.Context(), jobInput.Feature, *jobInput.IdempotencyKey)
				if existingErr == nil && existing != nil {
					log.Printf("[curriculum/create] user=%s mode=async_existing_job job_id=%s goal_profile_id=%s source_query_summary=%s elapsed_ms=%d", userID, existing.ID, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds())
					h.courseGenerationJobAccepted(c, existing, cost, "코스 초안 생성이 이미 진행 중이에요.")
					return
				}
			}
			log.Printf("[curriculum/create] user=%s mode=async_enqueue_failed goal_profile_id=%s source_query_summary=%s elapsed_ms=%d err=%s", userID, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(submitErr))
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
			return
		}
		log.Printf("[curriculum/create] user=%s mode=async_accepted job_id=%s goal_profile_id=%s source_query_summary=%s elapsed_ms=%d", userID, job.ID, activeGoal.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds())
		h.courseGenerationJobAccepted(c, job, cost, "코스 초안 생성을 시작했어요.")
		return
	}

	var decision *draftBuildDecision
	if req.SkipGeneration {
		decision = h.buildEmptyDraft(userID, req)
	} else {
		buildStartedAt := time.Now()
		var buildErr error
		decision, buildErr = h.buildDraft(c.Request.Context(), userID, req, cost)
		if buildErr != nil {
			log.Printf("[curriculum/create] user=%s mode=build_failed source_query_summary=%s elapsed_ms=%d err=%s", userID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(buildErr))
			if errors.Is(buildErr, errDraftBuildBusy) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":   "course_generation_busy",
					"message": "현재 코스 생성 요청이 많습니다. 잠시 후 다시 시도해 주세요.",
				})
				return
			}
			if errors.Is(buildErr, errDraftGenerationRetryable) {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "draft_generation_retry"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": buildErr.Error()})
			return
		}
		log.Printf("[curriculum/create] user=%s build_mode=%s review_status=%s queue_wait_ms=%d llm_generation_elapsed_ms=%d pattern_match_elapsed_ms=%d total_build_elapsed_ms=%d limiter_inflight_count=%d limiter_queue_depth=%d pattern_key=%s subpattern_key=%s refinement_mode=%s source_query_summary=%s build_elapsed_ms=%d", userID, decision.BuildMode, decision.ReviewStatus, decision.QueueWaitMS, decision.LLMElapsedMS, decision.PatternMatchMS, decision.TotalBuildMS, decision.LimiterInFlight, decision.LimiterQueueDepth, decision.PatternKey, decision.SubpatternKey, decision.RefinementMode, sourceQuerySummary, time.Since(buildStartedAt).Milliseconds())
	}

	chargeAmount := 0
	if decision.ShouldCharge {
		chargeAmount = decision.EffectiveCost
	}

	confirmed, err := h.service.BuildConfirmedCourseFromDraft(*decision.Draft)
	if err != nil {
		log.Printf("[curriculum/create] user=%s mode=confirm_build_failed source_query_summary=%s elapsed_ms=%d err=%s", userID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persistStartedAt := time.Now()
	if err := h.repo.CreateDraftAndStartLearning(c.Request.Context(), decision.Draft, confirmed, chargeAmount); err != nil {
		if err.Error() == "insufficient_points" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
			return
		}
		if errors.Is(err, ErrActiveGoalAlreadyAttached) {
			log.Printf("[curriculum/create] user=%s mode=active_goal_already_attached goal_profile_id=%s source_query_summary=%s elapsed_ms=%d err=%s", userID, formatOptionalUUID(decision.Draft.Draft.GoalProfileID), sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(err))
			c.JSON(http.StatusConflict, gin.H{"error": "active_goal_already_attached"})
			return
		}
		log.Printf("[curriculum/create] user=%s mode=create_failed source_query_summary=%s elapsed_ms=%d err=%s", userID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create draft"})
		return
	}

	draftPersistElapsedMS := time.Since(persistStartedAt).Milliseconds()

	fetchedDraft, err := h.repo.GetDraftByID(c.Request.Context(), decision.Draft.Draft.ID, userID)
	if err != nil {
		log.Printf("[curriculum/create] user=%s mode=reload_failed draft_id=%s source_query_summary=%s elapsed_ms=%d err=%s", userID, decision.Draft.Draft.ID, sourceQuerySummary, time.Since(startedAt).Milliseconds(), logsafe.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "draft created but failed to reload"})
		return
	}

	freeBalance, paidBalance, err := h.repo.GetPointBalances(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "draft created but failed to load point balance"})
		return
	}
	if !req.SkipGeneration {
		h.observeCourseDraftAIOutput(c.Request.Context(), userID, fetchedDraft, confirmed.Course.ID, "course_draft_create_sync", map[string]any{
			"build_mode":      decision.BuildMode,
			"billing_status":  decision.BillingStatus,
			"pattern_key":     decision.PatternKey,
			"subpattern_key":  decision.SubpatternKey,
			"refinement_mode": decision.RefinementMode,
		})
	}

	log.Printf(
		"[curriculum/create] user=%s draft_id=%s build_mode=%s billing_status=%s should_charge=%t queue_wait_ms=%d llm_generation_elapsed_ms=%d pattern_match_elapsed_ms=%d draft_persist_elapsed_ms=%d total_build_elapsed_ms=%d limiter_inflight_count=%d limiter_queue_depth=%d pattern_key=%s subpattern_key=%s refinement_mode=%s source_query_summary=%s total_elapsed_ms=%d lesson_count=%d",
		userID,
		decision.Draft.Draft.ID,
		decision.BuildMode,
		decision.BillingStatus,
		decision.ShouldCharge,
		decision.QueueWaitMS,
		decision.LLMElapsedMS,
		decision.PatternMatchMS,
		draftPersistElapsedMS,
		decision.TotalBuildMS,
		decision.LimiterInFlight,
		decision.LimiterQueueDepth,
		decision.PatternKey,
		decision.SubpatternKey,
		decision.RefinementMode,
		sourceQuerySummary,
		time.Since(startedAt).Milliseconds(),
		len(fetchedDraft.Lessons),
	)

	c.JSON(http.StatusCreated, gin.H{
		"draft":     fetchedDraft,
		"course_id": confirmed.Course.ID,
		"point_preview": CourseDraftPointPreview{
			Cost:         decision.EffectiveCost,
			FreeBalance:  freeBalance,
			PaidBalance:  paidBalance,
			TotalBalance: freeBalance + paidBalance,
		},
		"billing_status": decision.BillingStatus,
		"build_mode":     decision.BuildMode,
		"policy_cost":    cost,
		"build_metrics": gin.H{
			"queue_wait_ms":             decision.QueueWaitMS,
			"llm_generation_elapsed_ms": decision.LLMElapsedMS,
			"pattern_match_elapsed_ms":  decision.PatternMatchMS,
			"draft_persist_elapsed_ms":  draftPersistElapsedMS,
			"total_build_elapsed_ms":    decision.TotalBuildMS,
		},
	})
}
