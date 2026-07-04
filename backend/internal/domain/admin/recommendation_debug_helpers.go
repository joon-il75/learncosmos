package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func labelCountValue(counts gin.H, label string) int64 {
	value, ok := counts[label]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}

func formatNullTimeRFC3339(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339)
}

func formatNullInt32(value sql.NullInt32) string {
	if !value.Valid {
		return ""
	}
	return strconv.Itoa(int(value.Int32))
}

func nullInt32Value(value sql.NullInt32) any {
	if !value.Valid {
		return nil
	}
	return value.Int32
}

func nullFloat64Value(value sql.NullFloat64) any {
	if !value.Valid {
		return nil
	}
	return value.Float64
}

func rawJSONOrEmpty(value string) json.RawMessage {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || !json.Valid([]byte(trimmed)) {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(trimmed)
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case *string:
		if typed == nil {
			return ""
		}
		return strings.TrimSpace(*typed)
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "<nil>" {
			return ""
		}
		return text
	}
}

func jsonStringValue(value any) string {
	if value == nil {
		return "{}"
	}
	raw, err := json.Marshal(value)
	if err != nil || !json.Valid(raw) {
		return "{}"
	}
	return string(raw)
}

func mustJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil || !json.Valid(raw) {
		return "{}"
	}
	return string(raw)
}

func isSafeRecommendationVersionText(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return true
	}
	for _, r := range trimmed {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			continue
		}
		switch r {
		case '-', '_', '.', '/', ':':
			continue
		default:
			return false
		}
	}
	return true
}

func safeRate(part, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return float64(part) / float64(total)
}

func isRecommendationDebugLabelAllowed(label string) bool {
	switch label {
	case "good_fit",
		"goal_fit_stage_weak",
		"stage_fit_goal_weak",
		"irrelevant",
		"duplicate",
		"broken_link",
		"low_quality",
		"unsafe":
		return true
	default:
		return false
	}
}
