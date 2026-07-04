package curriculum

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"
)

const maxResearchBlockStringRunes = 50000

var (
	researchBlockScriptPattern  = regexp.MustCompile(`(?is)<\s*script\b[^>]*>.*?<\s*/\s*script\s*>`)
	researchBlockStylePattern   = regexp.MustCompile(`(?is)<\s*style\b[^>]*>.*?<\s*/\s*style\s*>`)
	researchBlockHTMLTagPattern = regexp.MustCompile(`(?is)<[^>]+>`)
)

func sanitizeResearchBlockContent(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return raw, nil
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, errLearningPointBlockInvalidInput
	}

	for key, value := range payload {
		text, ok := value.(string)
		if !ok {
			continue
		}
		sanitized, err := sanitizeResearchBlockField(key, text)
		if err != nil {
			return nil, err
		}
		payload[key] = sanitized
	}

	cleaned, err := json.Marshal(payload)
	if err != nil {
		return nil, errLearningPointBlockInvalidInput
	}
	return cleaned, nil
}

func sanitizeResearchBlockField(key, value string) (string, error) {
	if key == "text" {
		value = sanitizeTiptapHTML(value)
	}
	return sanitizeResearchBlockString(value)
}

func sanitizeResearchBlockString(value string) (string, error) {
	cleaned := researchBlockStylePattern.ReplaceAllString(value, "")
	cleaned = researchBlockScriptPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.TrimSpace(cleaned)
	if utf8.RuneCountInString(cleaned) > maxResearchBlockStringRunes {
		return "", errLearningPointBlockInvalidInput
	}
	return cleaned, nil
}

func researchBlockPlainText(value string) string {
	plain := researchBlockHTMLTagPattern.ReplaceAllString(value, " ")
	plain = strings.ReplaceAll(plain, "&nbsp;", " ")
	plain = strings.ReplaceAll(plain, "&#160;", " ")
	return strings.TrimSpace(strings.Join(strings.Fields(plain), " "))
}

func validateLearningResearchTextBlockContent(raw json.RawMessage) error {
	if len(raw) == 0 {
		return errLearningPointBlockInvalidInput
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return errLearningPointBlockInvalidInput
	}
	if _, ok := payload["text"]; !ok {
		return nil
	}
	title, _ := payload["title"].(string)
	if strings.TrimSpace(title) == "" {
		return errLearningPointBlockInvalidInput
	}
	text, _ := payload["text"].(string)
	if researchBlockPlainText(text) == "" {
		return errLearningPointBlockInvalidInput
	}
	return nil
}
