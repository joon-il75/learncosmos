package curriculum

import (
	"bytes"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var tiptapAllowedElements = map[string]struct{}{
	"a": block{}, "blockquote": block{}, "br": block{}, "em": block{}, "h1": block{}, "h2": block{}, "h3": block{},
	"li": block{}, "mark": block{}, "ol": block{}, "p": block{}, "s": block{}, "strong": block{}, "sub": block{},
	"sup": block{}, "table": block{}, "tbody": block{}, "td": block{}, "th": block{}, "thead": block{}, "tr": block{},
	"u": block{}, "ul": block{}, "img": block{},
}

var tiptapDropElements = map[string]struct{}{
	"script": block{}, "style": block{}, "iframe": block{}, "object": block{}, "embed": block{}, "svg": block{},
	"math": block{}, "meta": block{}, "link": block{}, "base": block{}, "form": block{}, "input": block{},
	"button": block{}, "textarea": block{}, "select": block{},
}

type block struct{}

func sanitizeTiptapHTML(value string) string {
	if strings.TrimSpace(value) == "" {
		return emptyResearchBlockHTML
	}

	doc, err := html.Parse(strings.NewReader("<!doctype html><html><body>" + value + "</body></html>"))
	if err != nil {
		return html.EscapeString(strings.TrimSpace(value))
	}

	var out bytes.Buffer
	body := findHTMLBodyNode(doc)
	if body == nil {
		return emptyResearchBlockHTML
	}
	for child := body.FirstChild; child != nil; child = child.NextSibling {
		writeSanitizedHTMLNode(&out, child)
	}
	cleaned := strings.TrimSpace(out.String())
	if cleaned == "" {
		return emptyResearchBlockHTML
	}
	return cleaned
}

const emptyResearchBlockHTML = "<p></p>"

func SanitizeTiptapHTMLForStorage(value string) string {
	return sanitizeTiptapHTML(value)
}

func findHTMLBodyNode(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}
	if node.Type == html.ElementNode && strings.EqualFold(node.Data, "body") {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findHTMLBodyNode(child); found != nil {
			return found
		}
	}
	return nil
}

func writeSanitizedHTMLNode(out *bytes.Buffer, node *html.Node) {
	switch node.Type {
	case html.TextNode:
		out.WriteString(html.EscapeString(node.Data))
	case html.ElementNode:
		name := strings.ToLower(node.Data)
		if _, drop := tiptapDropElements[name]; drop {
			return
		}
		if _, ok := tiptapAllowedElements[name]; !ok {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				writeSanitizedHTMLNode(out, child)
			}
			return
		}

		attrs, keep := sanitizedTiptapAttrs(name, node.Attr)
		if !keep {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				writeSanitizedHTMLNode(out, child)
			}
			return
		}

		out.WriteByte('<')
		out.WriteString(name)
		for _, attr := range attrs {
			out.WriteByte(' ')
			out.WriteString(attr.Key)
			out.WriteString(`="`)
			out.WriteString(html.EscapeString(attr.Val))
			out.WriteByte('"')
		}
		out.WriteByte('>')

		if name != "br" && name != "img" {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				writeSanitizedHTMLNode(out, child)
			}
			out.WriteString("</")
			out.WriteString(name)
			out.WriteByte('>')
		}
	case html.DocumentNode:
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			writeSanitizedHTMLNode(out, child)
		}
	}
}

func sanitizedTiptapAttrs(element string, attrs []html.Attribute) ([]html.Attribute, bool) {
	clean := make([]html.Attribute, 0, len(attrs)+2)
	for _, attr := range attrs {
		key := strings.ToLower(strings.TrimSpace(attr.Key))
		value := strings.TrimSpace(attr.Val)
		if key == "" || strings.HasPrefix(key, "on") {
			continue
		}

		switch element {
		case "a":
			if key == "href" {
				if !isSafeTiptapHTTPSURL(value) {
					return nil, false
				}
				clean = append(clean, html.Attribute{Key: "href", Val: value})
			}
		case "img":
			switch key {
			case "src":
				if !isSafeTiptapHTTPSURL(value) {
					return nil, false
				}
				clean = append(clean, html.Attribute{Key: "src", Val: value})
			case "alt", "title":
				clean = append(clean, html.Attribute{Key: key, Val: value})
			case "loading":
				if strings.EqualFold(value, "lazy") {
					clean = append(clean, html.Attribute{Key: "loading", Val: "lazy"})
				}
			}
		case "table":
			if key == "class" && containsAllowedClass(value, "lw-tiptap-table") {
				clean = append(clean, html.Attribute{Key: "class", Val: "lw-tiptap-table"})
			}
		case "td", "th":
			if (key == "colspan" || key == "rowspan") && isSafePositiveIntAttr(value, 1, 24) {
				clean = append(clean, html.Attribute{Key: key, Val: value})
			}
		}

		if key == "style" {
			if style := sanitizeTiptapStyle(value); style != "" {
				clean = append(clean, html.Attribute{Key: "style", Val: style})
			}
		}
	}

	if element == "a" {
		clean = upsertAttr(clean, "rel", "noopener noreferrer")
		clean = upsertAttr(clean, "target", "_blank")
	}
	if element == "img" && !hasAttr(clean, "src") {
		return nil, false
	}

	sort.SliceStable(clean, func(i, j int) bool {
		return clean[i].Key < clean[j].Key
	})
	return clean, true
}

func sanitizeTiptapStyle(value string) string {
	for _, part := range strings.Split(value, ";") {
		key, rawValue, ok := strings.Cut(part, ":")
		if !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(key), "text-align") {
			align := strings.ToLower(strings.TrimSpace(rawValue))
			switch align {
			case "left", "right", "center", "justify":
				return "text-align: " + align
			}
		}
	}
	return ""
}

func isSafeTiptapHTTPSURL(value string) bool {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil {
		return false
	}
	return parsed.Scheme == "https" && parsed.Host != "" && parsed.User == nil
}

func isSafePositiveIntAttr(value string, minValue, maxValue int) bool {
	n, err := strconv.Atoi(value)
	return err == nil && n >= minValue && n <= maxValue
}

func containsAllowedClass(value, allowed string) bool {
	for _, token := range strings.Fields(value) {
		if token == allowed {
			return true
		}
	}
	return false
}

func upsertAttr(attrs []html.Attribute, key, value string) []html.Attribute {
	for i := range attrs {
		if attrs[i].Key == key {
			attrs[i].Val = value
			return attrs
		}
	}
	return append(attrs, html.Attribute{Key: key, Val: value})
}

func hasAttr(attrs []html.Attribute, key string) bool {
	for _, attr := range attrs {
		if attr.Key == key {
			return true
		}
	}
	return false
}
