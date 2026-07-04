package goal

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, redisClient *redis.Client, authMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc, alphaAccessMiddleware gin.HandlerFunc) {
	repo := NewRepository(pool)
	service := NewService()
	handler := NewHandler(repo, service)
	handler.SetSafetyService(safety.NewService(safety.NewRepository(pool)))
	if redisClient != nil {
		llmJobRepo := llmjobs.NewRepository(pool)
		llmQueue := llmjobs.NewRedisStreamQueue(redisClient, llmjobs.DefaultStreamKey, llmjobs.DefaultConsumerGroup, llmjobs.DefaultDLQKey)
		handler.SetLLMGateway(llmjobs.NewGateway(llmJobRepo, llmQueue))
	}

	goals := rg.Group("/goals", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		goals.GET("/active", handler.GetActiveGoal)
		goals.DELETE("/active", handler.ResetActiveGoal)
		goals.POST("/start", handler.StartInterviewForUser)
		goals.POST("/interview", handler.SendMessageForUser)
		goals.POST("/confirm", handler.ConfirmGoalForUser)
		goals.POST("/revise", handler.ReviseGoalForUser)
		goals.POST("/cancel-revision", handler.CancelGoalRevisionForUser)
		goals.POST("/rebuild-decision", handler.RebuildDecisionForUser)
		goals.POST("/attach-draft", handler.AttachDraft)
	}

	drafts := rg.Group("/course-drafts", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		drafts.GET("/:id/goal", handler.GetGoal)
		drafts.POST("/:id/goal/start", handler.StartInterview)
		drafts.POST("/:id/goal/interview", handler.SendMessage)
		drafts.POST("/:id/goal/confirm", handler.ConfirmGoal)
		drafts.POST("/:id/goal/revise", handler.ReviseGoal)
		drafts.POST("/:id/goal/cancel-revision", handler.CancelGoalRevision)
		drafts.POST("/:id/goal/rebuild-decision", handler.RebuildDecision)
	}
}
