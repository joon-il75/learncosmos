package curriculum

import (
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

type RecommendExplorerContentRequest struct {
	Query                string   `json:"query"`
	QueryLanguage        string   `json:"query_language"`
	MaxResults           int      `json:"max_results"`
	CourseDraftID        string   `json:"course_draft_id"`
	RegionTitle          string   `json:"region_title"`
	RegionDescription    string   `json:"region_description"`
	SubRegionTitle       string   `json:"subregion_title"`
	SubRegionDescription string   `json:"subregion_description"`
	NodeTitle            string   `json:"node_title"`
	NodeSummary          string   `json:"node_summary"`
	PointID              string   `json:"point_id"`
	LessonID             string   `json:"lesson_id"`
	ExcludedContentIDs   []string `json:"excluded_content_ids"`
	ExcludedURLs         []string `json:"excluded_urls"`
}

type RecommendExplorerContentItem struct {
	ContentID    *string `json:"content_id,omitempty"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	URL          *string `json:"url,omitempty"`
	ContentType  string  `json:"content_type"`
	RankScore    float64 `json:"rank_score"`
}

type explorerSearchStage struct {
	name   string
	bundle normalizer.LessonSearchQuery
	filter explorerFilterStage
}

type ExplorerRecommendationContext struct {
	SourceQuery      string
	LearningGoal     string
	GoalQueryHint    string
	RegionTitle      string
	RegionObjective  string
	SubRegionTitle   string
	SubRegionSummary string
	NodeTitle        string
	NodeSummary      string
	FallbackQuery    string
}
