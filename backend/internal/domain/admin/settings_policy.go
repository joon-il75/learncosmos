package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// ─── 정책 문서 ───

func normalizePolicyType(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func isSupportedPolicyType(value string) bool {
	switch normalizePolicyType(value) {
	case "terms", "privacy":
		return true
	default:
		return false
	}
}

func normalizePolicyLocale(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func formatNullableTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

// GET /api/v1/public/policies
func (h *AdminHandler) GetPublicPolicies(c *gin.Context) {
	requestedLocale := normalizePolicyLocale(c.Query("locale"))
	rows, err := h.db.Query(c.Request.Context(), `
		WITH active_sets AS (
			SELECT id, document_type
			FROM policy_document_sets
			WHERE active = true
			  AND required = true
		),
		selected_documents AS (
			SELECT DISTINCT ON (pds.id)
				pd.id,
				pd.set_id,
				pd.type,
				pd.title,
				pd.version,
				pd.content,
				pd.is_required,
				pd.is_active,
				pd.locale,
				pd.translation_status,
				($1::text) AS requested_locale,
				(pd.locale <> $1::text) AS fallback_used,
				pd.effective_at,
				pd.published_at,
				pd.updated_at
			FROM active_sets pds
			JOIN policy_documents pd
			  ON pd.set_id = pds.id
			 AND pd.is_active = true
			 AND pd.locale IN ($1::text, 'ko')
			ORDER BY pds.id, (pd.locale = $1::text) DESC, (pd.locale = 'ko') DESC
		)
		SELECT id, set_id, type, title, version, content, is_required, is_active, locale, translation_status,
		       requested_locale, fallback_used, effective_at, published_at, updated_at
		FROM selected_documents
		ORDER BY type
	`, requestedLocale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 조회 실패"})
		return
	}
	defer rows.Close()

	docs := []gin.H{}
	for rows.Next() {
		var id, setID, policyType, title, content, locale, translationStatus, requestedLocale string
		var version int
		var required, active bool
		var fallbackUsed bool
		var effectiveAt, publishedAt, updatedAt time.Time
		if err := rows.Scan(&id, &setID, &policyType, &title, &version, &content, &required, &active, &locale, &translationStatus, &requestedLocale, &fallbackUsed, &effectiveAt, &publishedAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 파싱 실패"})
			return
		}
		docs = append(docs, gin.H{
			"id":                 id,
			"set_id":             setID,
			"type":               policyType,
			"title":              title,
			"version":            version,
			"content":            content,
			"required":           required,
			"active":             active,
			"locale":             locale,
			"translation_status": translationStatus,
			"requested_locale":   requestedLocale,
			"fallback_used":      fallbackUsed,
			"effective_at":       effectiveAt.UTC().Format("2006-01-02T15:04:05Z"),
			"published_at":       publishedAt.UTC().Format("2006-01-02T15:04:05Z"),
			"updated_at":         updatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

// GET /api/v1/public/policies/:type
func (h *AdminHandler) GetPublicPolicyByType(c *gin.Context) {
	policyType := normalizePolicyType(c.Param("type"))
	if !isSupportedPolicyType(policyType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "지원하지 않는 정책 유형입니다"})
		return
	}
	requestedLocale := normalizePolicyLocale(c.Query("locale"))

	var id, setID, title, content, locale, translationStatus string
	var version int
	var required, active bool
	var fallbackUsed bool
	var effectiveAt, publishedAt, updatedAt time.Time
	err := h.db.QueryRow(c.Request.Context(), `
		SELECT
			pd.id,
			pd.set_id,
			pd.type,
			pd.title,
			pd.version,
			pd.content,
			pd.is_required,
			pd.is_active,
			pd.locale,
			pd.translation_status,
			(pd.locale <> $2::text) AS fallback_used,
			pd.effective_at,
			pd.published_at,
			pd.updated_at
		FROM policy_document_sets pds
		JOIN policy_documents pd
		  ON pd.set_id = pds.id
		 AND pd.is_active = true
		 AND pd.locale IN ($2::text, 'ko')
		WHERE pds.document_type = $1
		  AND pds.active = true
		  AND pds.required = true
		ORDER BY (pd.locale = $2::text) DESC, (pd.locale = 'ko') DESC
		LIMIT 1
	`, policyType, requestedLocale).Scan(&id, &setID, &policyType, &title, &version, &content, &required, &active, &locale, &translationStatus, &fallbackUsed, &effectiveAt, &publishedAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "정책 문서를 찾을 수 없습니다"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"document": gin.H{
			"id":                 id,
			"set_id":             setID,
			"type":               policyType,
			"title":              title,
			"version":            version,
			"content":            content,
			"required":           required,
			"active":             active,
			"locale":             locale,
			"translation_status": translationStatus,
			"requested_locale":   requestedLocale,
			"fallback_used":      fallbackUsed,
			"effective_at":       effectiveAt.UTC().Format("2006-01-02T15:04:05Z"),
			"published_at":       publishedAt.UTC().Format("2006-01-02T15:04:05Z"),
			"updated_at":         updatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		},
	})
}

// GET /api/v1/super-admin/settings/policies
func (h *AdminHandler) GetPolicySettings(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), `
		SELECT id, set_id, type, title, version, content, is_required, is_active, is_draft, locale, translation_status, effective_at, published_at, updated_at
		FROM policy_documents
		WHERE is_active = true OR is_draft = true
		ORDER BY type, locale, is_draft DESC, version DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 조회 실패"})
		return
	}
	defer rows.Close()

	documents := []gin.H{}
	for rows.Next() {
		var id, setID, policyType, title, content, locale, translationStatus string
		var version int
		var required, active, draft bool
		var effectiveAt, updatedAt time.Time
		var publishedAt *time.Time
		if err := rows.Scan(&id, &setID, &policyType, &title, &version, &content, &required, &active, &draft, &locale, &translationStatus, &effectiveAt, &publishedAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 파싱 실패"})
			return
		}
		documents = append(documents, gin.H{
			"id":                 id,
			"set_id":             setID,
			"type":               policyType,
			"title":              title,
			"version":            version,
			"content":            content,
			"required":           required,
			"active":             active,
			"draft":              draft,
			"locale":             locale,
			"translation_status": translationStatus,
			"effective_at":       effectiveAt.UTC().Format("2006-01-02T15:04:05Z"),
			"published_at":       formatNullableTime(publishedAt),
			"updated_at":         updatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// PUT /api/v1/super-admin/settings/policies
func (h *AdminHandler) SavePolicyDrafts(c *gin.Context) {
	var req struct {
		Documents []struct {
			Type        string `json:"type"`
			Locale      string `json:"locale"`
			Title       string `json:"title"`
			Content     string `json:"content"`
			Required    bool   `json:"required"`
			EffectiveAt string `json:"effective_at"`
		} `json:"documents"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.db.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "트랜잭션 시작 실패"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	for _, doc := range req.Documents {
		policyType := normalizePolicyType(doc.Type)
		locale := normalizePolicyLocale(doc.Locale)
		title := strings.TrimSpace(doc.Title)
		content := strings.TrimSpace(doc.Content)
		if !isSupportedPolicyType(policyType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "지원하지 않는 정책 유형입니다"})
			return
		}
		if title == "" || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "정책 문서 제목과 본문은 비워둘 수 없습니다"})
			return
		}

		effectiveAt := time.Now().UTC()
		if strings.TrimSpace(doc.EffectiveAt) != "" {
			parsed, err := time.Parse(time.RFC3339, doc.EffectiveAt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "effective_at 형식이 잘못되었습니다"})
				return
			}
			effectiveAt = parsed.UTC()
		}

		var activeVersion int
		var activeSetID *string
		err := tx.QueryRow(c.Request.Context(), `
			SELECT version, set_id
			FROM policy_documents
			WHERE type = $1 AND locale = $2 AND is_active = true
			ORDER BY version DESC
			LIMIT 1
		`, policyType, locale).Scan(&activeVersion, &activeSetID)
		if err != nil && err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 조회 실패"})
			return
		}
		if err == pgx.ErrNoRows && locale == "en" {
			err = tx.QueryRow(c.Request.Context(), `
				SELECT version, set_id
				FROM policy_documents
				WHERE type = $1 AND locale = 'ko' AND is_active = true
				ORDER BY version DESC
				LIMIT 1
			`, policyType).Scan(&activeVersion, &activeSetID)
			if err != nil && err != pgx.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 기준 문서 조회 실패"})
				return
			}
		}

		var draftID string
		err = tx.QueryRow(c.Request.Context(), `
			SELECT id
			FROM policy_documents
			WHERE type = $1 AND locale = $2 AND is_draft = true
			LIMIT 1
		`, policyType, locale).Scan(&draftID)
		if err != nil && err != pgx.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 초안 조회 실패"})
			return
		}

		if err == pgx.ErrNoRows {
			nextVersion := activeVersion
			setID := activeSetID
			if locale == "ko" {
				nextVersion = 1
				if activeVersion > 0 {
					nextVersion = activeVersion + 1
				}
				setID = nil
			}
			if locale == "en" && (setID == nil || nextVersion == 0) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "영어 문서는 먼저 활성 한국어 기준 문서가 있어야 합니다"})
				return
			}
			if setID == nil {
				var newSetID string
				err = tx.QueryRow(c.Request.Context(), `
					INSERT INTO policy_document_sets (
						id, document_type, version, required, active, effective_at, created_at, updated_at
					)
					VALUES (gen_random_uuid(), $1, $2, $3, false, $4, NOW(), NOW())
					ON CONFLICT (document_type, version) DO UPDATE
					SET required = EXCLUDED.required,
					    effective_at = EXCLUDED.effective_at,
					    updated_at = NOW()
					RETURNING id
				`, policyType, nextVersion, doc.Required, effectiveAt).Scan(&newSetID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 기준 생성 실패"})
					return
				}
				setID = &newSetID
			}
			translationStatus := "source"
			if locale == "en" {
				translationStatus = "translated"
			}
			_, err = tx.Exec(c.Request.Context(), `
				INSERT INTO policy_documents (
					id, set_id, type, version, title, content, is_required, is_active, is_draft, locale, translation_status, effective_at, published_at, updated_at
				)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, false, true, $7, $8, $9, $9, NOW())
			`, *setID, policyType, nextVersion, title, content, doc.Required, locale, translationStatus, effectiveAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 초안 저장 실패"})
				return
			}
			continue
		}

		if _, err := tx.Exec(c.Request.Context(), `
			UPDATE policy_documents
			SET title = $1,
			    content = $2,
			    is_required = $3,
			    effective_at = $4,
			    updated_at = NOW()
			WHERE id = $5
		`, title, content, doc.Required, effectiveAt, draftID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 초안 수정 실패"})
			return
		}
	}

	if err := tx.Commit(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 문서 저장 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "정책 초안이 저장되었습니다"})
}

// POST /api/v1/super-admin/settings/policies/confirm
func (h *AdminHandler) ConfirmPolicyDrafts(c *gin.Context) {
	var req struct {
		Types  []string `json:"types"`
		Locale string   `json:"locale"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Types) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "확정할 정책 유형이 필요합니다"})
		return
	}

	tx, err := h.db.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "트랜잭션 시작 실패"})
		return
	}
	defer tx.Rollback(c.Request.Context())

	for _, rawType := range req.Types {
		policyType := normalizePolicyType(rawType)
		locale := normalizePolicyLocale(req.Locale)
		if !isSupportedPolicyType(policyType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "지원하지 않는 정책 유형입니다"})
			return
		}

		var draftID string
		var draftSetID string
		err := tx.QueryRow(c.Request.Context(), `
			SELECT id, set_id
			FROM policy_documents
			WHERE type = $1 AND locale = $2 AND is_draft = true
			LIMIT 1
		`, policyType, locale).Scan(&draftID, &draftSetID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"error": "확정할 정책 초안이 없습니다"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 초안 조회 실패"})
			return
		}

		if _, err := tx.Exec(c.Request.Context(), `
			UPDATE policy_documents
			SET is_active = false, updated_at = NOW()
			WHERE type = $1 AND locale = $2 AND is_active = true
		`, policyType, locale); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "기존 정책 비활성화 실패"})
			return
		}

		if _, err := tx.Exec(c.Request.Context(), `
			UPDATE policy_documents
			SET is_active = true,
			    is_draft = false,
			    published_at = NOW(),
			    updated_at = NOW()
			WHERE id = $1
		`, draftID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 확정 실패"})
			return
		}

		if locale == "ko" {
			if _, err := tx.Exec(c.Request.Context(), `
				UPDATE policy_document_sets
				SET active = false,
				    updated_at = NOW()
				WHERE document_type = $1 AND active = true
			`, policyType); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "기존 정책 기준 비활성화 실패"})
				return
			}

			if _, err := tx.Exec(c.Request.Context(), `
				UPDATE policy_document_sets
				SET active = true,
				    required = (
				      SELECT is_required
				      FROM policy_documents
				      WHERE id = $1
				    ),
				    effective_at = (
				      SELECT effective_at
				      FROM policy_documents
				      WHERE id = $1
				    ),
				    updated_at = NOW()
				WHERE id = $2
			`, draftID, draftSetID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 기준 확정 실패"})
				return
			}
		}
	}

	if err := tx.Commit(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "정책 확정 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "정책 문서가 확정되었습니다"})
}
