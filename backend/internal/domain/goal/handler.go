package goal

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/llmjobs"
	"github.com/learnweaver/backend/internal/domain/safety"
)

type Handler struct {
	repo                           *Repository
	service                        *Service
	singlePromptFallbackLLMEnabled bool
	goalChatMode                   string
	goalChatFastPathEnabled        bool
	goalChatWorkerEnabled          bool
	goalChatWorkerSyncWait         time.Duration
	llmGateway                     *llmjobs.Gateway
	safetyService                  *safety.Service
}

func NewHandler(repo *Repository, service *Service) *Handler {
	return NewHandlerWithOptions(repo, service, HandlerOptions{
		SinglePromptFallbackLLMEnabled: envBool("GOAL_CHAT_SINGLE_PROMPT_FALLBACK_ENABLED", false),
		GoalChatMode:                   envString("GOAL_CHAT_MODE", "two_step"),
		GoalChatFastPathEnabled:        envBool("GOAL_CHAT_FAST_PATH_ENABLED", false),
		GoalChatWorkerEnabled:          envBool("LLM_WORKER_FEATURE_GOAL_INTERVIEW", false),
		GoalChatWorkerSyncWait:         envDurationMS("LLM_WORKER_GOAL_INTERVIEW_SYNC_WAIT_MS", 8*time.Second),
	})
}

type HandlerOptions struct {
	SinglePromptFallbackLLMEnabled bool
	GoalChatMode                   string
	GoalChatFastPathEnabled        bool
	GoalChatWorkerEnabled          bool
	GoalChatWorkerSyncWait         time.Duration
}

func NewHandlerWithOptions(repo *Repository, service *Service, opts HandlerOptions) *Handler {
	goalChatMode := strings.TrimSpace(opts.GoalChatMode)
	if goalChatMode == "" {
		goalChatMode = "two_step"
	}
	workerWait := opts.GoalChatWorkerSyncWait
	if workerWait <= 0 {
		workerWait = 8 * time.Second
	}
	return &Handler{
		repo:                           repo,
		service:                        service,
		singlePromptFallbackLLMEnabled: opts.SinglePromptFallbackLLMEnabled,
		goalChatMode:                   goalChatMode,
		goalChatFastPathEnabled:        opts.GoalChatFastPathEnabled,
		goalChatWorkerEnabled:          opts.GoalChatWorkerEnabled,
		goalChatWorkerSyncWait:         workerWait,
	}
}

func envString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return strings.EqualFold(value, "true") || value == "1" || strings.EqualFold(value, "yes")
}

func envDurationMS(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err == nil {
		return parsed
	}
	if ms, err := strconv.Atoi(value); err == nil && ms > 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return fallback
}

func (h *Handler) SetLLMGateway(gateway *llmjobs.Gateway) {
	if h == nil {
		return
	}
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
