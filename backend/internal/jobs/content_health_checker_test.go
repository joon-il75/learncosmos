package jobs

import (
	"net/http"
	"testing"

	"github.com/learnweaver/backend/internal/domain/content"
)

func TestInterpretHealthHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		wantStatus content.ContentStatus
		wantScore  float64
	}{
		{name: "active", code: 200, wantStatus: content.ContentStatusActive, wantScore: 1.0},
		{name: "redirect", code: http.StatusMovedPermanently, wantStatus: content.ContentStatusRedirect, wantScore: 0.8},
		{name: "restricted", code: http.StatusForbidden, wantStatus: content.ContentStatusRestricted, wantScore: 0.4},
		{name: "broken", code: http.StatusNotFound, wantStatus: content.ContentStatusBroken, wantScore: 0.0},
		{name: "unknown", code: 429, wantStatus: content.ContentStatusUnknown, wantScore: 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotScore, gotCode := interpretHealthHTTPStatus(tt.code)
			if gotStatus != tt.wantStatus {
				t.Fatalf("expected status %s, got %s", tt.wantStatus, gotStatus)
			}
			if gotScore != tt.wantScore {
				t.Fatalf("expected score %f, got %f", tt.wantScore, gotScore)
			}
			if gotCode == nil || *gotCode != tt.code {
				t.Fatalf("expected http status pointer %d, got %+v", tt.code, gotCode)
			}
		})
	}
}
