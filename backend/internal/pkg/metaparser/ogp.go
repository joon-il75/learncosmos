package metaparser

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/learnweaver/backend/internal/pkg/logsafe"
	"github.com/learnweaver/backend/internal/pkg/safehttp"

	"golang.org/x/net/html"
)

// ParseOGP — URL에서 HTML을 가져와 OGP 메타 태그 파싱
func ParseOGP(ctx context.Context, rawURL string) (*MetaInfo, error) {
	if _, err := safehttp.ValidateHTTPURL(ctx, rawURL); err != nil {
		return nil, err
	}
	// robots.txt 확인
	if !IsAllowed(ctx, rawURL) {
		return nil, fmt.Errorf("robots.txt에 의해 크롤링이 제한된 URL입니다")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LearnWeaver/1.0; +https://learnweavr.com)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	client := safehttp.NewClient(safehttp.ClientConfig{Timeout: 10 * time.Second})
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, logsafe.URL(rawURL))
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		return nil, fmt.Errorf("HTML 콘텐츠가 아닙니다: %s", ct)
	}

	info := &MetaInfo{
		SourceURL:   rawURL,
		ContentType: DetectContentType(rawURL),
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "meta":
				prop := attrVal(n, "property")
				name := attrVal(n, "name")
				content := attrVal(n, "content")

				switch prop {
				case "og:title":
					if info.Title == "" {
						info.Title = content
					}
				case "og:description":
					if info.Description == "" {
						info.Description = content
					}
				case "og:image":
					if info.ThumbnailURL == "" {
						info.ThumbnailURL = content
					}
				case "og:site_name":
					if info.Author == "" {
						info.Author = content
					}
				}
				if prop == "" {
					switch name {
					case "twitter:title":
						if info.Title == "" {
							info.Title = content
						}
					case "twitter:description":
						if info.Description == "" {
							info.Description = content
						}
					case "twitter:image":
						if info.ThumbnailURL == "" {
							info.ThumbnailURL = content
						}
					case "description":
						if info.Description == "" {
							info.Description = content
						}
					}
				}

			case "title":
				if info.Title == "" && n.FirstChild != nil {
					info.Title = strings.TrimSpace(n.FirstChild.Data)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if len([]rune(info.Description)) > 300 {
		runes := []rune(info.Description)
		info.Description = string(runes[:300]) + "..."
	}

	return info, nil
}

func attrVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
