package curriculum

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(rg *gin.RouterGroup, pool *pgxpool.Pool, redisClient *redis.Client, authMiddleware gin.HandlerFunc, consentMiddleware gin.HandlerFunc, alphaAccessMiddleware gin.HandlerFunc) {
	repo := NewRepository(pool)
	goalRepo := goal.NewRepository(pool)
	service := NewService()
	handler := NewHandlerWithRedis(repo, goalRepo, service, redisClient)
	handler.SetSafetyService(safety.NewService(safety.NewRepository(pool)))
	llmJobRepo := llmjobs.NewRepository(pool)
	llmQueue := llmjobs.NewRedisStreamQueue(redisClient, llmjobs.DefaultStreamKey, llmjobs.DefaultConsumerGroup, llmjobs.DefaultDLQKey)
	handler.SetLLMGateway(llmjobs.NewGateway(llmJobRepo, llmQueue))

	generationJobs := rg.Group("/course-generation-jobs", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		generationJobs.GET("/:id", handler.GetCourseGenerationJob)
	}

	drafts := rg.Group("/course-drafts", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		drafts.GET("", handler.ListDrafts)
		drafts.POST("", handler.CreateDraft)
		drafts.GET("/:id", handler.GetDraft)
		drafts.PATCH("/:id", handler.UpdateDraft)
		drafts.POST("/:id/goal/apply", handler.ApplyDraftGoalContext)
		drafts.POST("/:id/activate", handler.ActivateDraft)
		drafts.DELETE("/:id", handler.DeleteDraft)
		drafts.POST("/:id/search", handler.SearchDraftLessons)
		drafts.PATCH("/:id/structure", handler.UpdateDraftStructure)
		drafts.PATCH("/:id/lessons/:lesson_id", handler.UpdateDraftLesson)
		drafts.PATCH("/:id/lessons/:lesson_id/memo", handler.UpdateDraftLessonMemo)
		drafts.PATCH("/:id/lessons/:lesson_id/journal", handler.UpdateDraftLessonJournal)
		drafts.PATCH("/:id/lessons/:lesson_id/record", handler.UpdateDraftLessonRecord)
		drafts.PATCH("/:id/lessons/:lesson_id/artifact", handler.UpdateDraftLessonArtifact)
		drafts.POST("/:id/lessons/:lesson_id/resources", handler.AttachDraftLessonResource)
		drafts.POST("/:id/resources/:resource_id/select", handler.SelectDraftResource)
		// Phase 5 public routes expose only lesson/point namespaces.
		drafts.PATCH("/:id/lessons/main/:main_lesson_id", handler.UpdateDraftLevel)
		drafts.PATCH("/:id/lessons/main/:main_lesson_id/memo", handler.UpdateDraftLevelMemo)
		drafts.POST("/:id/lessons/main/:main_lesson_id/sub", handler.CreateSubLesson)
		drafts.POST("/:id/lessons/main/:main_lesson_id/points", handler.CreateResearchNode)
		drafts.PATCH("/:id/points/:point_id", handler.UpdateResearchPoint)
		drafts.DELETE("/:id/points/:point_id", handler.DeleteResearchPoint)
		drafts.GET("/:id/points/:point_id/blocks", handler.GetResearchPointBlocks)
		drafts.POST("/:id/points/:point_id/blocks", handler.CreateResearchPointBlock)
		drafts.PATCH("/:id/points/:point_id/blocks/:block_id", handler.UpdateResearchPointBlock)
		drafts.DELETE("/:id/points/:point_id/blocks/:block_id", handler.DeleteResearchPointBlock)
		drafts.POST("/:id/points/:point_id/complete", handler.CompleteResearchPoint)
		drafts.DELETE("/:id/points/:point_id/complete", handler.UncompleteResearchPoint)
	}

	explorerRec := rg.Group("/explorer", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		explorerRec.POST("/recommend", handler.RecommendExplorerContent)
	}

	planets := rg.Group("/planets", authMiddleware, consentMiddleware, alphaAccessMiddleware)
	{
		planets.GET("/today-task", handler.GetTodayTask)
		planets.GET("/learning", handler.ListLearningPlanets)
		planets.GET("/learning/:id", handler.GetLearningPlanet)
		planets.GET("/learning/:id/records/feed", handler.GetLearningPlanetRecordFeed)
		planets.GET("/learning/:id/records", handler.GetLearningPlanetRecords)
		planets.GET("/learning/:id/results", handler.GetLearningPlanetResults)
		planets.GET("/learning/:id/points/:point_id", handler.GetLearningPlanetPoint)
		planets.POST("/learning/:id/start", handler.StartLearningPlanet)
		planets.PATCH("/learning/:id/points/:point_id/runtime", handler.UpdateLearningPointRuntime)
		planets.POST("/learning/:id/points/:point_id/sessions", handler.CreateLearningPointSession)
		planets.PATCH("/learning/:id/points/:point_id/sessions/:session_id/heartbeat", handler.HeartbeatLearningPointSession)
		planets.POST("/learning/:id/points/:point_id/sessions/:session_id/end", handler.EndLearningPointSession)
		planets.POST("/learning/:id/points/:point_id/research-material/confirm", handler.ConfirmLearningPointResearchMaterial)
		planets.DELETE("/learning/:id/points/:point_id/research-material/confirm", handler.UnconfirmLearningPointResearchMaterial)
		planets.PATCH("/learning/:id/points/:point_id/goal", handler.UpdateLearningPointGoal)
		planets.POST("/learning/:id/points/:point_id/questions", handler.CreateLearningPointQuestion)
		planets.POST("/learning/:id/points/:point_id/questions/ai-generate", handler.GenerateLearningPointQuestion)
		planets.POST("/learning/:id/points/:point_id/questions/:question_id/ai-feedback", handler.GenerateLearningPointQuestionFeedback)
		planets.PATCH("/learning/:id/points/:point_id/questions/:question_id", handler.UpdateLearningPointQuestion)
		planets.DELETE("/learning/:id/points/:point_id/questions/:question_id", handler.DeleteLearningPointQuestion)
		planets.POST("/learning/:id/points/:point_id/self-evaluation/ai-draft", handler.GenerateLearningPointSelfEvaluationDraft)
		planets.PATCH("/learning/:id/points/:point_id/self-evaluation", handler.UpsertLearningPointSelfEvaluation)
		planets.PATCH("/learning/:id/points/:point_id/journal", handler.UpdateLearningPointJournal)
		planets.POST("/learning/:id/points/:point_id/observation-notes", handler.CreateLearningPointObservationNote)
		planets.PATCH("/learning/:id/points/:point_id/observation-notes/:note_id", handler.UpdateLearningPointObservationNote)
		planets.DELETE("/learning/:id/points/:point_id/observation-notes/:note_id", handler.DeleteLearningPointObservationNote)
		planets.PATCH("/learning/:id/points/:point_id/record", handler.UpdateLearningPointRecord)
		planets.PATCH("/learning/:id/points/:point_id/artifact", handler.UpdateLearningPointArtifact)
		planets.GET("/learning/:id/points/:point_id/artifacts", handler.ListLearningPointArtifacts)
		planets.POST("/learning/:id/points/:point_id/artifacts", handler.CreateLearningPointArtifact)
		planets.PATCH("/learning/:id/points/:point_id/artifacts/:artifact_id", handler.UpdateLearningPointArtifactItem)
		planets.DELETE("/learning/:id/points/:point_id/artifacts/:artifact_id", handler.DeleteLearningPointArtifact)
		planets.GET("/learning/:id/points/:point_id/attachments", handler.ListLearningPointAttachments)
		planets.POST("/learning/:id/points/:point_id/attachments", handler.CreateLearningPointAttachment)
		planets.POST("/learning/:id/points/:point_id/attachments/upload", handler.UploadLearningPointAttachment)
		planets.GET("/learning/:id/points/:point_id/attachments/:attachment_id/open", handler.OpenLearningPointAttachment)
		planets.GET("/learning/:id/points/:point_id/attachments/:attachment_id/inline", handler.InlineLearningPointAttachment)
		planets.PATCH("/learning/:id/points/:point_id/attachments/:attachment_id", handler.UpdateLearningPointAttachment)
		planets.DELETE("/learning/:id/points/:point_id/attachments/:attachment_id", handler.DeleteLearningPointAttachment)
		planets.POST("/learning/:id/points/:point_id/report-material", handler.ReportLearningPointMaterial)
		planets.POST("/learning/:id/points/:point_id/report-material/:report_id/cancel", handler.CancelLearningPointMaterialReport)
		planets.POST("/learning/:id/points/:point_id/replace-material", handler.ReplaceLearningPointMaterial)
		planets.GET("/learning/:id/points/:point_id/events", handler.ListLearningPointEvents)
		planets.POST("/learning/:id/points/:point_id/practice-logs", handler.CreateLearningPointPracticeLog)
		planets.PATCH("/learning/:id/points/:point_id/practice-logs/:log_id", handler.UpdateLearningPointPracticeLog)
		planets.DELETE("/learning/:id/points/:point_id/practice-logs/:log_id", handler.DeleteLearningPointPracticeLog)
		planets.POST("/learning/:id/points/:point_id/ai-summary", handler.GenerateLearningPointAISummary)
		planets.POST("/learning/:id/points/:point_id/blocks", handler.CreateLearningPointBlock)
		planets.PATCH("/learning/:id/points/:point_id/blocks/:block_id", handler.UpdateLearningPointBlock)
		planets.DELETE("/learning/:id/points/:point_id/blocks/:block_id", handler.DeleteLearningPointBlock)
		planets.POST("/learning/:id/complete", handler.CompleteLearningPlanet)
		planets.GET("/shared", handler.ListSharedPlanets)
		planets.GET("/shared/:id", handler.GetSharedPlanet)
		planets.GET("/shared/:id/records/feed", handler.GetSharedPlanetRecordFeed)
		planets.GET("/shared/:id/records", handler.GetSharedPlanetRecords)
		planets.GET("/shared/:id/results", handler.GetSharedPlanetResults)
		planets.GET("/shared/:id/points/:point_id", handler.GetSharedPlanetPoint)
	}
}
