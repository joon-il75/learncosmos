package curriculum

import (
	"strings"
	"testing"
)

func TestSanitizeTiptapHTMLPreservesSupportedFormatting(t *testing.T) {
	input := `<h2 style="text-align: center; color: red">Title</h2><p><strong>bold</strong> <em>em</em> <u>u</u> <mark>mark</mark> <sub>2</sub><sup>3</sup></p><blockquote>quote</blockquote><ul><li>one</li></ul>`
	got := sanitizeTiptapHTML(input)

	for _, want := range []string{
		`<h2 style="text-align: center">Title</h2>`,
		`<strong>bold</strong>`,
		`<em>em</em>`,
		`<u>u</u>`,
		`<mark>mark</mark>`,
		`<sub>2</sub>`,
		`<sup>3</sup>`,
		`<blockquote>quote</blockquote>`,
		`<ul><li>one</li></ul>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("sanitizeTiptapHTML() missing %q in %q", want, got)
		}
	}
	if strings.Contains(got, "color:") {
		t.Fatalf("sanitizeTiptapHTML() kept unsafe style: %q", got)
	}
}

func TestSanitizeTiptapHTMLPreservesSafeLinksImagesAndTables(t *testing.T) {
	input := `<p><a href="https://example.com/path?q=1" onclick="bad()">link</a><img src="https://cdn.example.com/a.png" alt="A" loading="lazy" onerror="bad()"></p><table class="x lw-tiptap-table"><tbody><tr><td colspan="2" rowspan="1" onclick="bad()">cell</td></tr></tbody></table>`
	got := sanitizeTiptapHTML(input)

	for _, want := range []string{
		`<a href="https://example.com/path?q=1" rel="noopener noreferrer" target="_blank">link</a>`,
		`<img alt="A" loading="lazy" src="https://cdn.example.com/a.png">`,
		`<table class="lw-tiptap-table">`,
		`<td colspan="2" rowspan="1">cell</td>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("sanitizeTiptapHTML() missing %q in %q", want, got)
		}
	}
	for _, forbidden := range []string{"onclick", "onerror", "class=\"x"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("sanitizeTiptapHTML() kept forbidden fragment %q in %q", forbidden, got)
		}
	}
}

func TestSanitizeTiptapHTMLDropsExecutableContent(t *testing.T) {
	input := `<p onclick="bad()">hello<script>bad()</script><style>body{}</style><iframe src="https://example.com"></iframe><svg onload="bad()"><text>x</text></svg><a href="javascript:bad()">bad link</a><img src="data:image/png;base64,aaaa" alt="bad"></p>`
	got := sanitizeTiptapHTML(input)

	for _, forbidden := range []string{"script", "style", "iframe", "svg", "onload", "onclick", "javascript:", "data:image"} {
		if strings.Contains(strings.ToLower(got), forbidden) {
			t.Fatalf("sanitizeTiptapHTML() kept forbidden fragment %q in %q", forbidden, got)
		}
	}
	if strings.Contains(got, "<a") || strings.Contains(got, "<img") {
		t.Fatalf("sanitizeTiptapHTML() kept unsafe link/image tag: %q", got)
	}
	if !strings.Contains(got, "hello") {
		t.Fatalf("sanitizeTiptapHTML() removed safe text: %q", got)
	}
}

func TestSanitizeTiptapHTMLEmptyFallback(t *testing.T) {
	if got := sanitizeTiptapHTML("   "); got != emptyResearchBlockHTML {
		t.Fatalf("sanitizeTiptapHTML(empty) = %q, want %q", got, emptyResearchBlockHTML)
	}
	if got := sanitizeTiptapHTML(`<script>bad()</script>`); got != emptyResearchBlockHTML {
		t.Fatalf("sanitizeTiptapHTML(script only) = %q, want %q", got, emptyResearchBlockHTML)
	}
}
