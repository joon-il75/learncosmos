package llmjobs

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo    *Repository
	gateway *Gateway
}

func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, redisClient *redis.Client, authMiddleware gin.HandlerFunc, headerOnlyAuthMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc) {
	repo := NewRepository(pool)
	queue := NewRedisStreamQueue(redisClient, DefaultStreamKey, DefaultConsumerGroup, DefaultDLQKey)
	handler := &Handler{
		repo:    repo,
		gateway: NewGateway(repo, queue),
	}

	userGroup := rg.Group("/llm-jobs", authMiddleware, consentMiddleware)
	{
		userGroup.GET("/:id", handler.GetJob)
	}

	superAdminGroup := rg.Group("/llm-jobs", headerOnlyAuthMiddleware)
	{
		superAdminGroup.GET("/observations", handler.GetJobObservations)
		superAdminGroup.POST("/dummy", handler.CreateDummyJob)
	}
}

func (h *Handler) GetJobObservations(c *gin.Context) {
	role, roleOK := auth.GetCurrentRole(c)
	if !roleOK || role != auth.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "error_code": "forbidden"})
		return
	}
	policy := DefaultObservationPolicy()
	policy.QueuedTimeout = parseObservationDuration(c.Query("queued_timeout"), policy.QueuedTimeout)
	policy.RunningTimeout = parseObservationDuration(c.Query("running_timeout"), policy.RunningTimeout)
	now := time.Now()
	counts, err := h.repo.CountJobObservations(c.Request.Context(), now, policy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "llm_job_observations_failed", "error_code": "llm_job_observations_failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"generated_at": now,
		"policy": gin.H{
			"queued_timeout_ms":  policy.QueuedTimeout.Milliseconds(),
			"running_timeout_ms": policy.RunningTimeout.Milliseconds(),
		},
		"counts": counts,
	})
}

func (h *Handler) GetJob(c *gin.Context) {
	userID, ok := currentUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_code": "unauthorized"})
		return
	}
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_job_id", "error_code": "invalid_job_id"})
		return
	}
	job, err := h.repo.GetJobForUser(c.Request.Context(), jobID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "llm_job_load_failed"
		if errors.Is(err, ErrJobNotFound) || errors.Is(err, ErrJobAccessDenied) {
			status = http.StatusNotFound
			code = "llm_job_not_found"
		}
		c.JSON(status, gin.H{"error": code, "error_code": code})
		return
	}
	c.JSON(http.StatusOK, jobResponse(job))
}

func (h *Handler) CreateDummyJob(c *gin.Context) {
	role, roleOK := auth.GetCurrentRole(c)
	if !roleOK || role != auth.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "error_code": "forbidden"})
		return
	}
	userID, ok := currentUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_code": "unauthorized"})
		return
	}
	job, err := h.gateway.SubmitAndWait(c.Request.Context(), CreateJobInput{
		UserID:  &userID,
		Feature: FeatureDummy,
		RequestRef: map[string]any{
			"source": "manual_dummy",
		},
		PromptInputRef: map[string]any{},
	}, 2*time.Second)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "llm_job_queue_failed", "error_code": "llm_job_queue_failed"})
		return
	}
	status := http.StatusAccepted
	if terminalStatus(job.Status) {
		status = http.StatusOK
	}
	c.JSON(status, jobResponse(job))
}

func currentUserUUID(c *gin.Context) (uuid.UUID, bool) {
	raw, ok := auth.GetCurrentUserID(c)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func jobResponse(job *Job) gin.H {
	resp := gin.H{
		"job_id":     job.ID,
		"feature":    job.Feature,
		"status":     job.Status,
		"phase":      job.Phase,
		"result_ref": job.ResultRef,
		"error_code": job.ErrorCode,
		"metrics": gin.H{
			"queue_wait_ms":             job.QueueWaitMS,
			"provider_wait_ms":          job.ProviderWaitMS,
			"llm_generation_elapsed_ms": job.LLMGenerationElapsedMS,
			"parse_validate_elapsed_ms": job.ParseValidateElapsedMS,
			"persist_elapsed_ms":        job.PersistElapsedMS,
			"total_job_elapsed_ms":      job.TotalJobElapsedMS,
			"provider_inflight_count":   job.ProviderInFlightCount,
			"provider_queue_depth":      job.ProviderQueueDepth,
			"worker_id":                 job.WorkerID,
			"attempts":                  job.Attempts,
		},
		"created_at":  job.CreatedAt,
		"started_at":  job.StartedAt,
		"finished_at": job.FinishedAt,
	}
	return resp
}
