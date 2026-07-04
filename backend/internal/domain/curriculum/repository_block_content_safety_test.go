package curriculum

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSanitizeResearchBlockContentSanitizesTextHTML(t *testing.T) {
	raw := json.RawMessage(`{
		"title": "Research note",
		"text": "<p onclick=\"bad()\">hello<script>bad()</script><a href=\"javascript:bad()\">bad</a><img src=\"data:image/png;base64,aaaa\" alt=\"bad\"><a href=\"https://example.com\">ok</a></p>"
	}`)

	cleanedRaw, err := sanitizeResearchBlockContent(raw)
	if err != nil {
		t.Fatalf("sanitizeResearchBlockContent() error = %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(cleanedRaw, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload["title"] != "Research note" {
		t.Fatalf("title = %q, want unchanged", payload["title"])
	}
	text := payload["text"]
	for _, forbidden := range []string{"onclick", "script", "javascript:", "data:image", "<img"} {
		if strings.Contains(strings.ToLower(text), forbidden) {
			t.Fatalf("sanitized text kept forbidden fragment %q in %q", forbidden, text)
		}
	}
	if !strings.Contains(text, `<a href="https://example.com" rel="noopener noreferrer" target="_blank">ok</a>`) {
		t.Fatalf("sanitized text did not keep safe link: %q", text)
	}
	if !strings.Contains(text, "hello") {
		t.Fatalf("sanitized text removed safe text: %q", text)
	}
}

func TestSanitizeResearchBlockContentKeepsNonTextFieldsOnExistingPath(t *testing.T) {
	raw := json.RawMessage(`{
		"title": "<b>Title</b><script>bad()</script>",
		"caption": "<i>Caption</i>",
		"url": "https://example.com/?q=<tag>",
		"text": "<p>Body</p>"
	}`)

	cleanedRaw, err := sanitizeResearchBlockContent(raw)
	if err != nil {
		t.Fatalf("sanitizeResearchBlockContent() error = %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(cleanedRaw, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload["title"] != "<b>Title</b>" {
		t.Fatalf("title = %q, want existing non-text behavior", payload["title"])
	}
	if payload["caption"] != "<i>Caption</i>" {
		t.Fatalf("caption = %q, want existing non-text behavior", payload["caption"])
	}
	if payload["url"] != "https://example.com/?q=<tag>" {
		t.Fatalf("url = %q, want existing non-text behavior", payload["url"])
	}
	if payload["text"] != "<p>Body</p>" {
		t.Fatalf("text = %q, want sanitized equivalent", payload["text"])
	}
}
