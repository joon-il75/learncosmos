package curriculum

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/safety"
	"github.com/learnweaver/backend/internal/pkg/objectstorage"
	"github.com/learnweaver/backend/internal/pkg/systemsettings"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo                             *Repository
	goalRepo                         *goal.Repository
	service                          *Service
	settingsStore                    *systemsettings.Store
	openAIAPIKey                     string
	youtubeAPIKey                    string
	naverClientID                    string
	naverSecret                      string
	objectStorage                    *objectstorage.Client
	redisClient                      *redis.Client
	llmGateway                       *llmjobs.Gateway
	safetyService                    *safety.Service
	buildLimiter                     *draftBuildLimiter
	goalBuildLock                    *courseGenerationGoalLock
	queueTimeout                     time.Duration
	genTimeout                       time.Duration
	reviewTimeout                    time.Duration
	pointAISummaryWorkerEnabled      bool
	pointAISummaryWorkerSyncWait     time.Duration
	pointQuestionWorkerEnabled       bool
	pointQuestionWorkerSyncWait      time.Duration
	pointFeedbackWorkerEnabled       bool
	pointFeedbackWorkerSyncWait      time.Duration
	pointSelfEvalDraftWorkerEnabled  bool
	pointSelfEvalDraftWorkerSyncWait time.Duration
}

type draftBuildDecision struct {
	Draft             *DraftAggregate
	BuildMode         string
	BillingStatus     string
	EffectiveCost     int
	ShouldCharge      bool
	ReviewStatus      string
	QueueWaitMS       int64
	PatternMatchMS    int64
	LLMElapsedMS      int64
	TotalBuildMS      int64
	LimiterInFlight   int32
	LimiterQueueDepth int32
	PatternKey        string
	SubpatternKey     string
	RefinementMode    string
}

type lessonSearchExecution struct {
	Provider      string
	APIKey        string
	BillingStatus string
	EffectiveCost int
	ShouldCharge  bool
}

type evaluationExecution struct {
	Provider      string
	APIKey        string
	BillingStatus string
	EffectiveCost int
	ShouldCharge  bool
}

type systemLLMSelection struct {
	Feature string
	Setting *LLMSetting
}

type pointAIQuestionSuggestion struct {
	Question     string `json:"question"`
	QuestionType string `json:"question_type"`
}

func NewHandler(repo *Repository, goalRepo *goal.Repository, service *Service) *Handler {
	objectStorageConfig := objectstorage.LoadConfigFromEnv()
	objectStorageClient, err := objectstorage.NewClient(objectStorageConfig)
	if err != nil {
		log.Printf("[object-storage] disabled: %v", err)
	} else {
		log.Printf("[object-storage] configured bucket=%s endpoint=%s", objectStorageConfig.Bucket, objectStorageConfig.Endpoint)
	}
	return &Handler{
		repo:                             repo,
		goalRepo:                         goalRepo,
		service:                          service,
		settingsStore:                    systemsettings.NewStore(repo.pool),
		openAIAPIKey:                     os.Getenv("OPENAI_API_KEY"),
		youtubeAPIKey:                    os.Getenv("YOUTUBE_API_KEY"),
		naverClientID:                    os.Getenv("NAVER_CLIENT_ID"),
		naverSecret:                      os.Getenv("NAVER_CLIENT_SECRET"),
		objectStorage:                    objectStorageClient,
		buildLimiter:                     newDraftBuildLimiter(getEnvInt("CURRICULUM_BUILD_MAX_CONCURRENCY", 4), 0.75, getEnvInt("CURRICULUM_BUILD_MAX_QUEUE_DEPTH", 8)),
		goalBuildLock:                    newCourseGenerationGoalLock(),
		queueTimeout:                     getEnvDuration("CURRICULUM_BUILD_QUEUE_TIMEOUT", 12*time.Second),
		genTimeout:                       getEnvDuration("CURRICULUM_GENERATION_TIMEOUT", 18*time.Second),
		reviewTimeout:                    getEnvDuration("CURRICULUM_REVIEW_TIMEOUT", 8*time.Second),
		pointAISummaryWorkerEnabled:      getEnvBool("LLM_WORKER_FEATURE_POINT_AI_SUMMARY", false),
		pointAISummaryWorkerSyncWait:     getEnvDuration("LLM_WORKER_POINT_AI_SUMMARY_SYNC_WAIT", 6*time.Second),
		pointQuestionWorkerEnabled:       getEnvBool("LLM_WORKER_FEATURE_POINT_QUESTION_GENERATE", false),
		pointQuestionWorkerSyncWait:      getEnvDuration("LLM_WORKER_POINT_QUESTION_GENERATE_SYNC_WAIT", 6*time.Second),
		pointFeedbackWorkerEnabled:       getEnvBool("LLM_WORKER_FEATURE_POINT_FEEDBACK_GENERATE", false),
		pointFeedbackWorkerSyncWait:      getEnvDuration("LLM_WORKER_POINT_FEEDBACK_GENERATE_SYNC_WAIT", 6*time.Second),
		pointSelfEvalDraftWorkerEnabled:  getEnvBool("LLM_WORKER_FEATURE_POINT_SELF_EVALUATION_DRAFT", false),
		pointSelfEvalDraftWorkerSyncWait: getEnvDuration("LLM_WORKER_POINT_SELF_EVALUATION_DRAFT_SYNC_WAIT", 6*time.Second),
	}
}

func NewHandlerWithRedis(repo *Repository, goalRepo *goal.Repository, service *Service, redisClient *redis.Client) *Handler {
	handler := NewHandler(repo, goalRepo, service)
	handler.redisClient = redisClient
	return handler
}

func (h *Handler) SetLLMGateway(gateway *llmjobs.Gateway) {
	h.llmGateway = gateway
}

func (h *Handler) SetSafetyService(service *safety.Service) {
	if h == nil {
		return
	}
	h.safetyService = service
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	idStr, ok := raw.(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func learningPointStatuses() []DraftStatus {
	return []DraftStatus{
		DraftStatusConfirmed,
		DraftStatusLearning,
	}
}

func sharedPointStatuses() []DraftStatus {
	return []DraftStatus{
		DraftStatusArchived,
	}
}

func (h *Handler) respondWithLearningPointDetail(
	c *gin.Context,
	userID uuid.UUID,
	courseID uuid.UUID,
	pointID uuid.UUID,
	statusCode int,
	extra gin.H,
	reloadError string,
) {
	point, err := h.repo.GetPlanetPointDetailByIDAndDraftStatuses(c.Request.Context(), userID, courseID, pointID, learningPointStatuses())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": reloadError})
		return
	}

	response := gin.H{
		"planet": point.Planet,
		"point":  point,
	}
	for key, value := range extra {
		response[key] = value
	}

	c.JSON(statusCode, response)
}

// resolveEvaluationExecution resolves the LLM key and billing info for on-demand AI features.
// Priority: BYOK(OpenAI) free → system key paid → empty key fallback.
func (h *Handler) resolveEvaluationExecution(ctx context.Context, userID uuid.UUID, policyCost int) (*evaluationExecution, error) {
	userConfig, err := h.repo.GetUserRuntimeAIConfig(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userConfig != nil && userConfig.Mode == "byok" && isFreeBYOKProvider(userConfig.Provider) && strings.TrimSpace(userConfig.APIKey) != "" {
		return &evaluationExecution{
			Provider:      strings.TrimSpace(strings.ToLower(userConfig.Provider)),
			APIKey:        strings.TrimSpace(userConfig.APIKey),
			BillingStatus: "byok_no_charge",
			EffectiveCost: 0,
			ShouldCharge:  false,
		}, nil
	}
	selection, err := h.resolveDraftSystemLLM(ctx, userID)
	if err != nil || selection == nil || selection.Setting == nil {
		return &evaluationExecution{}, nil
	}
	apiKey, err := h.settingsStore.ResolveAPIKey(ctx, selection.Setting.Provider, h.openAIAPIKey)
	if err != nil || strings.TrimSpace(apiKey) == "" {
		return &evaluationExecution{}, nil
	}
	return &evaluationExecution{
		Provider:      strings.TrimSpace(selection.Setting.Provider),
		APIKey:        strings.TrimSpace(apiKey),
		BillingStatus: "charged",
		EffectiveCost: policyCost,
		ShouldCharge:  policyCost > 0,
	}, nil
}

func isFreeBYOKProvider(provider string) bool {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "openai":
		return true
	default:
		return false
	}
}

func getEnvBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	return strings.EqualFold(raw, "true") || raw == "1" || strings.EqualFold(raw, "yes")
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
