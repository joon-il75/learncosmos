package curriculum

import (
	"strings"
	"testing"
)

func TestNormalizePointArtifactInputSanitizesDescriptionHTML(t *testing.T) {
	artifactType, title, artifactURL, description, err := normalizePointArtifactInput(
		" document ",
		" Portfolio ",
		" https://example.com/result ",
		`<p onclick="bad()">hello<script>bad()</script><a href="javascript:bad()">bad</a><img src="data:image/png;base64,aaaa" alt="bad"><a href="https://example.com">ok</a></p>`,
	)
	if err != nil {
		t.Fatalf("normalizePointArtifactInput() error = %v", err)
	}

	if artifactType != "document" {
		t.Fatalf("artifactType = %q, want document", artifactType)
	}
	if title != "Portfolio" {
		t.Fatalf("title = %q, want Portfolio", title)
	}
	if artifactURL != "https://example.com/result" {
		t.Fatalf("url = %q, want trimmed URL", artifactURL)
	}
	for _, forbidden := range []string{"onclick", "script", "javascript:", "data:image", "<img"} {
		if strings.Contains(strings.ToLower(description), forbidden) {
			t.Fatalf("description kept forbidden fragment %q in %q", forbidden, description)
		}
	}
	if !strings.Contains(description, `<a href="https://example.com" rel="noopener noreferrer" target="_blank">ok</a>`) {
		t.Fatalf("description did not keep safe link: %q", description)
	}
	if !strings.Contains(description, "hello") {
		t.Fatalf("description removed safe text: %q", description)
	}
}

func TestNormalizePointArtifactInputPreservesSupportedDescriptionFormatting(t *testing.T) {
	_, _, _, description, err := normalizePointArtifactInput(
		"",
		"Result",
		"",
		`<h3 style="text-align: right; color: red">Title</h3><table class="x lw-tiptap-table"><tbody><tr><td colspan="2">cell</td></tr></tbody></table>`,
	)
	if err != nil {
		t.Fatalf("normalizePointArtifactInput() error = %v", err)
	}

	for _, want := range []string{
		`<h3 style="text-align: right">Title</h3>`,
		`<table class="lw-tiptap-table">`,
		`<td colspan="2">cell</td>`,
	} {
		if !strings.Contains(description, want) {
			t.Fatalf("description missing %q in %q", want, description)
		}
	}
	if strings.Contains(description, "color:") || strings.Contains(description, `class="x`) {
		t.Fatalf("description kept unsupported style/class: %q", description)
	}
}

func TestNormalizePointArtifactInputRejectsUnsafeURL(t *testing.T) {
	for _, input := range []string{
		"javascript:alert(1)",
		"//example.com/result",
		"https://user:pass@example.com/result",
		"https://example.com/a b",
	} {
		if _, _, _, _, err := normalizePointArtifactInput("document", "Result", input, ""); err == nil {
			t.Fatalf("normalizePointArtifactInput(%q) error = nil, want error", input)
		}
	}
}

func TestNormalizePointAttachmentInputRejectsUnsafeURL(t *testing.T) {
	title := "Resource"
	for _, input := range []string{
		"javascript:alert(1)",
		"//example.com/resource",
		"https://user:pass@example.com/resource",
		"https://example.com/a b",
	} {
		urlValue := input
		if _, _, _, _, _, _, _, _, err := normalizePointAttachmentInput(nil, nil, nil, &title, &urlValue, nil, nil, nil); err == nil {
			t.Fatalf("normalizePointAttachmentInput(%q) error = nil, want error", input)
		}
	}
}
