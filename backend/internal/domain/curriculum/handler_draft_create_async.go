package curriculum

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

const defaultCourseGenerationAsyncMaxActiveJobs = 40

func courseGenerationAsyncEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("COURSE_GENERATION_ASYNC_ENABLED")), "true")
}

func courseGenerationAsyncMaxActiveJobs() int {
	return getEnvInt("COURSE_GENERATION_ASYNC_MAX_ACTIVE_JOBS", defaultCourseGenerationAsyncMaxActiveJobs)
}

func (h *Handler) courseGenerationJobAccepted(c *gin.Context, job *llmjobs.Job, policyCost int, message string) {
	c.JSON(http.StatusAccepted, gin.H{
		"job_id":         job.ID,
		"status":         job.Status,
		"message":        message,
		"poll_url":       fmt.Sprintf("/api/v1/course-generation-jobs/%s", job.ID),
		"policy_cost":    policyCost,
		"estimated_wait": courseGenerationWaitEstimateJSON(h.courseGenerationWaitEstimate(c.Request.Context(), job)),
	})
}

func CreateCourseDraftJobInput(userID uuid.UUID, activeGoal *goal.GoalProfile, learningLanguage string, policyCost int) llmjobs.CreateJobInput {
	idempotencyKey := fmt.Sprintf("course_draft_create:%s:%s", userID, activeGoal.ID)
	providerMode := "pending"
	feature := llmjobs.FeatureCourseDraftCreate
	return llmjobs.CreateJobInput{
		UserID:         &userID,
		Feature:        feature,
		IdempotencyKey: &idempotencyKey,
		RequestRef: map[string]any{
			"goal_profile_id":      activeGoal.ID.String(),
			"goal_profile_version": activeGoal.Version,
			"learning_language":    normalizeLearningLanguage(learningLanguage),
		},
		PromptInputRef: map[string]any{},
		ProviderMode:   &providerMode,
		PointCost:      &policyCost,
	}
}
