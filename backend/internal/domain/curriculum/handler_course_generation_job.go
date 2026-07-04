package curriculum

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
)

func (h *Handler) GetCourseGenerationJob(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_code": "unauthorized"})
		return
	}
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_job_id", "error_code": "invalid_job_id"})
		return
	}
	repo := llmjobs.NewRepository(h.repo.pool)
	job, err := repo.GetJobForUser(c.Request.Context(), jobID, userID)
	if err != nil || job.Feature != llmjobs.FeatureCourseDraftCreate {
		status := http.StatusInternalServerError
		code := "llm_job_load_failed"
		if errors.Is(err, llmjobs.ErrJobNotFound) || errors.Is(err, llmjobs.ErrJobAccessDenied) || job == nil || job.Feature != llmjobs.FeatureCourseDraftCreate {
			status = http.StatusNotFound
			code = "llm_job_not_found"
		}
		c.JSON(status, gin.H{"error": code, "error_code": code})
		return
	}

	result := gin.H{}
	if len(job.ResultRef) > 0 {
		_ = json.Unmarshal(job.ResultRef, &result)
	}
	resp := gin.H{
		"job_id":         job.ID,
		"status":         job.Status,
		"phase":          job.Phase,
		"error_code":     job.ErrorCode,
		"message":        courseGenerationJobMessage(job),
		"estimated_wait": courseGenerationWaitEstimateJSON(h.courseGenerationWaitEstimate(c.Request.Context(), job)),
		"metrics": gin.H{
			"queue_wait_ms":             job.QueueWaitMS,
			"provider_wait_ms":          job.ProviderWaitMS,
			"llm_generation_elapsed_ms": job.LLMGenerationElapsedMS,
			"total_job_elapsed_ms":      job.TotalJobElapsedMS,
			"worker_id":                 job.WorkerID,
			"attempts":                  job.Attempts,
		},
		"created_at":  job.CreatedAt,
		"started_at":  job.StartedAt,
		"finished_at": job.FinishedAt,
	}
	for key, value := range result {
		resp[key] = value
	}
	c.JSON(http.StatusOK, resp)
}

func courseGenerationJobMessage(job *llmjobs.Job) string {
	if job == nil {
		return "코스 초안 생성 상태를 확인할 수 없습니다."
	}
	switch job.Status {
	case llmjobs.StatusQueued:
		return "코스 초안 생성을 준비하고 있습니다."
	case llmjobs.StatusRunning:
		return "목표에 맞춰 코스 초안을 생성하고 있습니다."
	case llmjobs.StatusSucceeded:
		return "코스 초안 생성이 완료되었습니다."
	case llmjobs.StatusFailed:
		if job.ErrorCode != nil && *job.ErrorCode == "course_generation_busy" {
			return "현재 코스 생성 요청이 많습니다. 잠시 후 다시 시도해 주세요."
		}
		return "코스 초안 생성에 실패했습니다."
	case llmjobs.StatusCanceled:
		return "코스 초안 생성이 취소되었습니다."
	case llmjobs.StatusExpired:
		return "코스 초안 생성 요청이 만료되었습니다."
	default:
		return "코스 초안 생성 상태를 확인하고 있습니다."
	}
}
