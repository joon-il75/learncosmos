package curriculum

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) logCurriculumPatternShadowAsync(userID uuid.UUID, req CreateCourseDraftRequest, codePatternKey, codeSubpatternKey, refinementMode string) {
	mode := strings.TrimSpace(strings.ToLower(os.Getenv("CURRICULUM_PATTERN_MATCH_MODE")))
	if mode == "" {
		mode = "code_only"
	}
	if mode != "shadow_compare" && mode != "embedding_primary" {
		return
	}

	queryText := BuildCurriculumPatternQueryText(req)
	if strings.TrimSpace(queryText) == "" {
		log.Printf("[curriculum/pattern_shadow] user=%s mode=%s code_pattern=%s code_subpattern=%s refinement_mode=%s fallback_reason=empty_query", userID, mode, codePatternKey, codeSubpatternKey, refinementMode)
		return
	}

	go func() {
		startedAt := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		result, err := h.resolveCurriculumPatternMatches(ctx, req, 3)
		if err != nil {
			log.Printf("[curriculum/pattern_shadow] user=%s mode=%s code_pattern=%s code_subpattern=%s refinement_mode=%s fallback_reason=%s duration_ms=%d", userID, mode, codePatternKey, codeSubpatternKey, refinementMode, truncateShadowError(err), time.Since(startedAt).Milliseconds())
			return
		}

		h.logCurriculumPatternMatchResult(userID, mode, codePatternKey, codeSubpatternKey, refinementMode, result, startedAt)
	}()
}

func truncateShadowError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	runes := []rune(message)
	if len(runes) <= 240 {
		return message
	}
	return string(runes[:240])
}
