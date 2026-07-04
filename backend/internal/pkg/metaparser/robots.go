package metaparser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/learnweaver/backend/internal/pkg/safehttp"
)

// IsAllowed — robots.txt 확인 후 크롤링 허용 여부 반환
// 확인 실패 시 허용(true)으로 처리 (과도한 차단 방지)
func IsAllowed(ctx context.Context, rawURL string) bool {
	if _, err := safehttp.ValidateHTTPURL(ctx, rawURL); err != nil {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return true
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
	if err != nil {
		return true
	}

	client := safehttp.NewClient(safehttp.ClientConfig{Timeout: 5 * time.Second})
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return true // robots.txt 없으면 허용
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32*1024)) // 최대 32KB
	if err != nil {
		return true
	}

	path := u.Path
	if path == "" {
		path = "/"
	}

	// 단순 파싱: User-agent: * 섹션의 Disallow 확인
	inOurSection := false
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			agent := strings.TrimSpace(line[len("user-agent:"):])
			inOurSection = agent == "*"
		}
		if inOurSection && strings.HasPrefix(strings.ToLower(line), "disallow:") {
			disallowPath := strings.TrimSpace(line[len("disallow:"):])
			if disallowPath != "" && strings.HasPrefix(path, disallowPath) {
				return false
			}
		}
	}
	return true
}
