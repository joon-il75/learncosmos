package explorer

import (
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ValidateSubRegionLimit: 현재 active 서브지역 수가 3 이상이면 에러
func (s *Service) ValidateSubRegionLimit(current int) error {
	if current >= 3 {
		return ErrSubRegionLimitExceeded
	}
	return nil
}

// ValidateExplorationNodeSource:
// 탐험지점 생성/수정 시 sourceType에 따라 최소 입력 조건 검증
func (s *Service) ValidateExplorationNodeSource(req CreateExplorationNodeRequest) error {
	switch req.SourceType {
	case SourceTypeYoutube, SourceTypeWeb:
		if req.SourceURL == nil || strings.TrimSpace(*req.SourceURL) == "" {
			return ErrExplorationNodeBadSource
		}
	case SourceTypeInternal, SourceTypeCreator:
		if req.ContentID == nil {
			return ErrExplorationNodeBadSource
		}
	default:
		return ErrExplorationNodeBadSource
	}
	return nil
}

// BuildCourseAggregate:
// flat region/subregion/node 슬라이스를 API 응답용 중첩 구조로 조립
func (s *Service) BuildCourseAggregate(
	courseDraftID uuid.UUID,
	title string,
	updatedAt time.Time,
	regions []Region,
	subregions []SubRegion,
	nodes []Node,
) CourseAggregate {
	subByRegion := make(map[uuid.UUID][]SubRegion, len(subregions))
	for _, sr := range subregions {
		subByRegion[sr.RegionID] = append(subByRegion[sr.RegionID], sr)
	}

	nodesByParent := make(map[uuid.UUID][]Node, len(nodes))
	for _, n := range nodes {
		nodesByParent[n.ParentID] = append(nodesByParent[n.ParentID], n)
	}

	sortNodes := func(ns []Node) {
		sort.Slice(ns, func(i, j int) bool {
			if ns[i].OrderIndex == ns[j].OrderIndex {
				return ns[i].CreatedAt.Before(ns[j].CreatedAt)
			}
			return ns[i].OrderIndex < ns[j].OrderIndex
		})
	}

	regionAggregates := make([]RegionAggregate, 0, len(regions))
	for _, reg := range regions {
		subs := subByRegion[reg.ID]
		sort.Slice(subs, func(i, j int) bool {
			if subs[i].OrderIndex == subs[j].OrderIndex {
				return subs[i].CreatedAt.Before(subs[j].CreatedAt)
			}
			return subs[i].OrderIndex < subs[j].OrderIndex
		})

		subAggregates := make([]SubRegionAggregate, 0, len(subs))
		for _, sr := range subs {
			ns := nodesByParent[sr.ID]
			sortNodes(ns)
			subAggregates = append(subAggregates, SubRegionAggregate{
				SubRegion: sr,
				Nodes:     ns,
			})
		}

		regionNodes := nodesByParent[reg.ID]
		sortNodes(regionNodes)

		regionAggregates = append(regionAggregates, RegionAggregate{
			Region:     reg,
			SubRegions: subAggregates,
			Nodes:      regionNodes,
		})
	}

	return CourseAggregate{
		CourseDraftID: courseDraftID,
		Title:         title,
		UpdatedAt:     updatedAt,
		Regions:       regionAggregates,
	}
}
