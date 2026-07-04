package curriculum

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

type courseGenerationWaitEstimate struct {
	ActiveJobs    int    `json:"active_jobs"`
	QueuePosition *int   `json:"queue_position,omitempty"`
	MaxActiveJobs int    `json:"max_active_jobs"`
	MinSeconds    int    `json:"min_seconds"`
	MaxSeconds    int    `json:"max_seconds"`
	Label         string `json:"label"`
}

func (h *Handler) courseGenerationWaitEstimate(ctx context.Context, job *llmjobs.Job) *courseGenerationWaitEstimate {
	if job == nil || job.Feature != llmjobs.FeatureCourseDraftCreate || terminalCourseGenerationJobStatus(job.Status) {
		return nil
	}

	estimate, err := llmjobs.NewRepository(h.repo.pool).GetActiveJobEstimate(ctx, llmjobs.FeatureCourseDraftCreate, job.ID)
	if err != nil {
		return nil
	}
	return buildCourseGenerationWaitEstimate(estimate, courseGenerationAsyncMaxActiveJobs())
}

func terminalCourseGenerationJobStatus(status llmjobs.Status) bool {
	switch status {
	case llmjobs.StatusSucceeded, llmjobs.StatusFailed, llmjobs.StatusCanceled, llmjobs.StatusExpired:
		return true
	default:
		return false
	}
}

func buildCourseGenerationWaitEstimate(active llmjobs.ActiveJobEstimate, maxActiveJobs int) *courseGenerationWaitEstimate {
	minSeconds, maxSeconds := courseGenerationWaitSeconds(active.QueuePosition, active.ActiveCount)
	return &courseGenerationWaitEstimate{
		ActiveJobs:    active.ActiveCount,
		QueuePosition: active.QueuePosition,
		MaxActiveJobs: maxActiveJobs,
		MinSeconds:    minSeconds,
		MaxSeconds:    maxSeconds,
		Label:         formatCourseGenerationWaitLabel(minSeconds, maxSeconds),
	}
}

func courseGenerationWaitSeconds(queuePosition *int, activeJobs int) (int, int) {
	if queuePosition != nil {
		position := *queuePosition
		if position <= 0 {
			return 5, 20
		}
		if position <= 8 {
			return 10, 20
		}
		if position <= 16 {
			return 15, 30
		}
		if position <= 24 {
			return 20, 40
		}
		if position <= 32 {
			return 30, 55
		}
		return 40, 75
	}

	if activeJobs <= 8 {
		return 10, 20
	}
	if activeJobs <= 16 {
		return 15, 30
	}
	if activeJobs <= 24 {
		return 20, 40
	}
	if activeJobs <= 32 {
		return 30, 55
	}
	return 40, 75
}

func formatCourseGenerationWaitLabel(minSeconds, maxSeconds int) string {
	if minSeconds <= 0 && maxSeconds <= 0 {
		return "곧 시작됩니다"
	}
	if minSeconds == maxSeconds {
		return fmt.Sprintf("약 %d초", maxSeconds)
	}
	return fmt.Sprintf("약 %d~%d초", minSeconds, maxSeconds)
}

func courseGenerationWaitEstimateJSON(estimate *courseGenerationWaitEstimate) gin.H {
	if estimate == nil {
		return nil
	}
	payload := gin.H{
		"active_jobs":     estimate.ActiveJobs,
		"max_active_jobs": estimate.MaxActiveJobs,
		"min_seconds":     estimate.MinSeconds,
		"max_seconds":     estimate.MaxSeconds,
		"label":           estimate.Label,
	}
	if estimate.QueuePosition != nil {
		payload["queue_position"] = *estimate.QueuePosition
	}
	return payload
}
