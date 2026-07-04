package admin

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/auth"
	"github.com/learnweaver/backend/internal/domain/demoaccount"
	"github.com/learnweaver/backend/internal/pkg/objectstorage"
)

type createDemoAccountsRequest struct {
	Count            int    `json:"count"`
	StartNumber      int    `json:"start_number"`
	EmailPrefix      string `json:"email_prefix"`
	EmailDomain      string `json:"email_domain"`
	LabelPrefix      string `json:"label_prefix"`
	UILocale         string `json:"ui_locale"`
	LearningLanguage string `json:"learning_language"`
	InitialPoints    int    `json:"initial_points"`
	ExpiresAt        string `json:"expires_at"`
	AssignedTo       string `json:"assigned_to"`
	AssignmentNote   string `json:"assignment_note"`
}

type updateDemoAccountRequest struct {
	AssignedTo     *string `json:"assigned_to"`
	AssignmentNote *string `json:"assignment_note"`
	ExpiresAt      *string `json:"expires_at"`
	Disabled       *bool   `json:"disabled"`
	DisabledReason *string `json:"disabled_reason"`
}

type purgeDemoAccountRequest struct {
	Confirm        string `json:"confirm"`
	IncludeAccount bool   `json:"include_account"`
}

func (h *AdminHandler) listDemoAccountsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	repo := demoaccount.NewRepository(h.db)
	result, err := repo.List(c.Request.Context(), demoaccount.ListOptions{
		Query:  c.Query("q"),
		Status: c.DefaultQuery("status", "all"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 계정 목록을 불러오지 못했습니다"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) createDemoAccountsHandler(c *gin.Context) {
	if strings.TrimSpace(h.cfg.DemoLoginSecret) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DEMO_LOGIN_SECRET 설정이 필요합니다"})
		return
	}

	var req createDemoAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}
	if req.Count < 1 || req.Count > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "생성 개수는 1~100 사이여야 합니다"})
		return
	}

	expiresAt, err := parseOptionalDemoExpiresAt(req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "만료일을 확인해 주세요"})
		return
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "만료일은 현재 시각 이후여야 합니다"})
		return
	}

	var createdBy *string
	if userID, ok := auth.GetCurrentUserID(c); ok {
		if _, err := uuid.Parse(userID); err == nil {
			createdBy = &userID
		}
	}

	repo := demoaccount.NewRepository(h.db)
	items, err := repo.CreateBatch(c.Request.Context(), demoaccount.CreateBatchRequest{
		Count:            req.Count,
		StartNumber:      req.StartNumber,
		EmailPrefix:      req.EmailPrefix,
		EmailDomain:      req.EmailDomain,
		LabelPrefix:      req.LabelPrefix,
		UILocale:         req.UILocale,
		LearningLanguage: req.LearningLanguage,
		InitialPoints:    req.InitialPoints,
		ExpiresAt:        expiresAt,
		AssignedTo:       req.AssignedTo,
		AssignmentNote:   req.AssignmentNote,
		CreatedBy:        createdBy,
		CreatedByActor:   "super_admin",
	}, h.cfg.DemoLoginSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 계정 생성에 실패했습니다"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"items": items})
}

func (h *AdminHandler) rotateDemoAccountCodeHandler(c *gin.Context) {
	if strings.TrimSpace(h.cfg.DemoLoginSecret) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DEMO_LOGIN_SECRET 설정이 필요합니다"})
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데모 계정 ID입니다"})
		return
	}

	repo := demoaccount.NewRepository(h.db)
	code, err := repo.RotateCode(c.Request.Context(), id, h.cfg.DemoLoginSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 로그인 코드 재발급에 실패했습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"login_code": code})
}

func (h *AdminHandler) updateDemoAccountHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데모 계정 ID입니다"})
		return
	}

	var req updateDemoAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}

	var expiresValue *time.Time
	var expiresPtr **time.Time
	if req.ExpiresAt != nil {
		parsed, err := parseOptionalDemoExpiresAt(*req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "만료일을 확인해 주세요"})
			return
		}
		expiresValue = parsed
		expiresPtr = &expiresValue
	}

	repo := demoaccount.NewRepository(h.db)
	if err := repo.Update(c.Request.Context(), id, demoaccount.UpdateRequest{
		AssignedTo:     req.AssignedTo,
		AssignmentNote: req.AssignmentNote,
		ExpiresAt:      expiresPtr,
		Disabled:       req.Disabled,
		DisabledReason: req.DisabledReason,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 계정 수정에 실패했습니다"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AdminHandler) previewDemoAccountPurgeHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데모 계정 ID입니다"})
		return
	}

	repo := demoaccount.NewRepository(h.db)
	preview, err := repo.PurgePreview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 계정 삭제 미리보기에 실패했습니다"})
		return
	}
	c.JSON(http.StatusOK, preview)
}

func (h *AdminHandler) purgeDemoAccountHandler(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데모 계정 ID입니다"})
		return
	}

	var req purgeDemoAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
		return
	}
	if req.Confirm != "PURGE_DEMO_ACCOUNT_DATA" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "확인 문자열이 일치하지 않습니다"})
		return
	}
	if req.IncludeAccount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "계정 row 삭제는 아직 지원하지 않습니다"})
		return
	}

	repo := demoaccount.NewRepository(h.db)
	preview, err := repo.PurgeLearningData(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "데모 계정 데이터 삭제에 실패했습니다"})
		return
	}

	result := demoaccount.PurgeResult{
		PurgePreview:   preview,
		AccountDeleted: false,
	}
	objectStorageClient, objectStorageErr := objectstorage.NewClient(objectstorage.LoadConfigFromEnv())
	for _, key := range preview.ObjectKeys {
		if objectStorageErr != nil {
			result.ObjectDeleteFailedKeys = append(result.ObjectDeleteFailedKeys, key)
			result.ObjectDeleteFailureText = append(result.ObjectDeleteFailureText, objectStorageErr.Error())
			continue
		}
		if err := objectStorageClient.DeleteObject(c.Request.Context(), key); err != nil {
			result.ObjectDeleteFailedKeys = append(result.ObjectDeleteFailedKeys, key)
			result.ObjectDeleteFailureText = append(result.ObjectDeleteFailureText, err.Error())
			continue
		}
		result.DeletedObjectKeys = append(result.DeletedObjectKeys, key)
	}

	c.JSON(http.StatusOK, result)
}

func parseOptionalDemoExpiresAt(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		defaultExpiresAt := time.Now().AddDate(0, 0, 30)
		return &defaultExpiresAt, nil
	}
	if strings.EqualFold(value, "none") {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return &t, nil
	}
	date, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return nil, err
	}
	endOfDay := date.Add(24*time.Hour - time.Nanosecond)
	return &endOfDay, nil
}
