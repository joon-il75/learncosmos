package safety

import (
	"fmt"
	"strings"
)

var safeMetadataKeys = map[string]struct{}{
	"billing_status":  {},
	"build_mode":      {},
	"content_length":  {},
	"course_draft_id": {},
	"direction":       {},
	"field":           {},
	"goal_profile_id": {},
	"pattern_key":     {},
	"refinement_mode": {},
	"source":          {},
	"source_feature":  {},
	"subpattern_key":  {},
}

func sanitizeSafetyMetadata(metadata map[string]any) map[string]any {
	if len(metadata) == 0 {
		return map[string]any{}
	}
	cleaned := make(map[string]any, len(metadata))
	for key, value := range metadata {
		key = strings.TrimSpace(key)
		if _, ok := safeMetadataKeys[key]; !ok {
			continue
		}
		if cleanedValue, ok := sanitizeSafetyMetadataValue(value); ok {
			cleaned[key] = cleanedValue
		}
	}
	return cleaned
}

func sanitizeSafetyMetadataValue(value any) (any, bool) {
	switch v := value.(type) {
	case nil:
		return nil, false
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil, false
		}
		if len([]rune(trimmed)) > 120 {
			trimmed = string([]rune(trimmed)[:120])
		}
		return trimmed, true
	case bool:
		return v, true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return v, true
	default:
		text := strings.TrimSpace(fmt.Sprint(v))
		if text == "" || strings.Contains(text, "map[") || strings.Contains(text, "[") {
			return nil, false
		}
		if len([]rune(text)) > 120 {
			text = string([]rune(text)[:120])
		}
		return text, true
	}
}
