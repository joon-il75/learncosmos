package explorer

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateSubRegionLimit(t *testing.T) {
	s := NewService()

	if err := s.ValidateSubRegionLimit(2); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if err := s.ValidateSubRegionLimit(3); err != ErrSubRegionLimitExceeded {
		t.Fatalf("expected ErrSubRegionLimitExceeded, got %v", err)
	}
}

func TestValidateExplorationNodeSource_YoutubeOK(t *testing.T) {
	s := NewService()
	url := "https://youtube.com/watch?v=abc"

	req := CreateExplorationNodeRequest{
		ParentKind: ParentKindRegion,
		ParentID:   uuid.New(),
		Title:      "극한 정의",
		SourceType: SourceTypeYoutube,
		SourceURL:  &url,
	}

	if err := s.ValidateExplorationNodeSource(req); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateExplorationNodeSource_YoutubeMissingURL(t *testing.T) {
	s := NewService()

	req := CreateExplorationNodeRequest{
		ParentKind: ParentKindRegion,
		ParentID:   uuid.New(),
		Title:      "극한 정의",
		SourceType: SourceTypeYoutube,
	}

	if err := s.ValidateExplorationNodeSource(req); err != ErrExplorationNodeBadSource {
		t.Fatalf("expected ErrExplorationNodeBadSource, got %v", err)
	}
}

func TestValidateExplorationNodeSource_InternalOK(t *testing.T) {
	s := NewService()
	contentID := uuid.New()

	req := CreateExplorationNodeRequest{
		ParentKind: ParentKindRegion,
		ParentID:   uuid.New(),
		Title:      "내부 문서",
		SourceType: SourceTypeInternal,
		ContentID:  &contentID,
	}

	if err := s.ValidateExplorationNodeSource(req); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateExplorationNodeSource_InternalMissingContentID(t *testing.T) {
	s := NewService()

	req := CreateExplorationNodeRequest{
		ParentKind: ParentKindRegion,
		ParentID:   uuid.New(),
		Title:      "내부 문서",
		SourceType: SourceTypeInternal,
	}

	if err := s.ValidateExplorationNodeSource(req); err != ErrExplorationNodeBadSource {
		t.Fatalf("expected ErrExplorationNodeBadSource, got %v", err)
	}
}
