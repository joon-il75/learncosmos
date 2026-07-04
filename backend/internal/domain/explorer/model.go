package explorer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSubRegionLimitExceeded   = errors.New("subregion limit exceeded: max 3 per region")
	ErrExplorationNodeBadSource = errors.New("exploration node requires source_url or content_id")
	ErrParentInactive           = errors.New("parent is inactive: activate parent first")
	ErrNotFound                 = errors.New("not found")
	ErrForbidden                = errors.New("forbidden")
)

type NodeType string
type SourceType string
type ResearchType string
type LayoutType string
type ParentKind string
type ItemStatus string

const (
	NodeTypeExploration NodeType = "exploration"
	NodeTypeResearch    NodeType = "research"

	SourceTypeYoutube  SourceType = "youtube"
	SourceTypeWeb      SourceType = "web"
	SourceTypeCreator  SourceType = "creator"
	SourceTypeInternal SourceType = "internal"

	ResearchTypeConcept  ResearchType = "concept"
	ResearchTypePractice ResearchType = "practice"
	ResearchTypeProblem  ResearchType = "problem"
	ResearchTypeFree     ResearchType = "free"

	LayoutTypeBasic     LayoutType = "basic"
	LayoutTypeTwoColumn LayoutType = "two-column"
	LayoutTypeNoteCard  LayoutType = "note-card"

	ParentKindRegion    ParentKind = "region"
	ParentKindSubRegion ParentKind = "subregion"

	ItemStatusActive   ItemStatus = "active"
	ItemStatusInactive ItemStatus = "inactive"
)

type Region struct {
	ID                  uuid.UUID  `json:"id"`
	CourseDraftID       uuid.UUID  `json:"course_draft_id"`
	CourseDraftLessonID *uuid.UUID `json:"course_draft_lesson_id,omitempty"`
	Name                string     `json:"name"`
	Description         *string    `json:"description"`
	OrderIndex          int        `json:"order_index"`
	Status              ItemStatus `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type SubRegion struct {
	ID                  uuid.UUID  `json:"id"`
	RegionID            uuid.UUID  `json:"region_id"`
	CourseDraftLessonID *uuid.UUID `json:"course_draft_lesson_id,omitempty"`
	Name                string     `json:"name"`
	Description         *string    `json:"description"`
	OrderIndex          int        `json:"order_index"`
	Status              ItemStatus `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type Node struct {
	ID             uuid.UUID     `json:"id"`
	ParentKind     ParentKind    `json:"parent_kind"`
	ParentID       uuid.UUID     `json:"parent_id"`
	NodeType       NodeType      `json:"node_type"`
	DraftPointID   *uuid.UUID    `json:"draft_point_id,omitempty"`
	LearningStatus *string       `json:"learning_status,omitempty"`
	Title          string        `json:"title"`
	OrderIndex     int           `json:"order_index"`
	Status         ItemStatus    `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	SourceType     *SourceType   `json:"source_type,omitempty"`
	SourceURL      *string       `json:"source_url,omitempty"`
	ContentID      *uuid.UUID    `json:"content_id,omitempty"`
	Summary        *string       `json:"summary,omitempty"`
	ResearchType   *ResearchType `json:"research_type,omitempty"`
	LayoutType     *LayoutType   `json:"layout_type,omitempty"`
	BlockCount     int           `json:"block_count"`
}

type SubRegionAggregate struct {
	SubRegion SubRegion `json:"subregion"`
	Nodes     []Node    `json:"nodes"`
}

type RegionAggregate struct {
	Region     Region               `json:"region"`
	SubRegions []SubRegionAggregate `json:"subregions"`
	Nodes      []Node               `json:"nodes"`
}

type CourseAggregate struct {
	CourseDraftID                           uuid.UUID         `json:"course_draft_id"`
	Title                                   string            `json:"title"`
	Status                                  string            `json:"status"`
	IsInactive                              bool              `json:"is_inactive"`
	Progress                                *float64          `json:"progress,omitempty"`
	PlanetTextureMapID                      *uuid.UUID        `json:"planet_texture_map_id,omitempty"`
	PlanetTextureMapName                    *string           `json:"planet_texture_map_name,omitempty"`
	PlanetTextureMapAsset                   *string           `json:"planet_texture_map_asset,omitempty"`
	PlanetTextureMapRotationDurationSeconds *int              `json:"planet_texture_map_rotation_duration_seconds,omitempty"`
	PlanetTextureMapRotationDirection       *string           `json:"planet_texture_map_rotation_direction,omitempty"`
	UpdatedAt                               time.Time         `json:"updated_at"`
	Regions                                 []RegionAggregate `json:"regions"`
}

type CreateRegionRequest struct {
	CourseDraftID uuid.UUID `json:"course_draft_id" binding:"required"`
	Name          string    `json:"name" binding:"required,max=100"`
	Description   *string   `json:"description"`
	OrderIndex    int       `json:"order_index"`
}

type UpdateRegionRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	OrderIndex  *int    `json:"order_index"`
}

type UpdateStatusRequest struct {
	Status       ItemStatus `json:"status" binding:"required"`
	WithChildren bool       `json:"with_children"`
}

type CreateSubRegionRequest struct {
	RegionID    uuid.UUID `json:"region_id" binding:"required"`
	Name        string    `json:"name" binding:"required,max=100"`
	Description *string   `json:"description"`
	OrderIndex  int       `json:"order_index"`
}

type UpdateSubRegionRequest struct {
	RegionID    *uuid.UUID `json:"region_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	OrderIndex  *int       `json:"order_index"`
}

type CreateExplorationNodeRequest struct {
	ParentKind ParentKind `json:"parent_kind" binding:"required"`
	ParentID   uuid.UUID  `json:"parent_id" binding:"required"`
	Title      string     `json:"title" binding:"required,max=200"`
	SourceType SourceType `json:"source_type" binding:"required"`
	SourceURL  *string    `json:"source_url"`
	ContentID  *uuid.UUID `json:"content_id"`
	OrderIndex int        `json:"order_index"`
}

type CreateResearchNodeRequest struct {
	ParentKind ParentKind `json:"parent_kind" binding:"required"`
	ParentID   uuid.UUID  `json:"parent_id" binding:"required"`
	Title      string     `json:"title" binding:"required,max=200"`
	OrderIndex int        `json:"order_index"`
}

type SetResearchNodeTypeRequest struct {
	ResearchType ResearchType `json:"research_type" binding:"required"`
	LayoutType   LayoutType   `json:"layout_type" binding:"required"`
}

type UpdateNodeRequest struct {
	ParentKind   *ParentKind   `json:"parent_kind"`
	ParentID     *uuid.UUID    `json:"parent_id"`
	Title        *string       `json:"title"`
	SourceType   *SourceType   `json:"source_type"`
	SourceURL    *string       `json:"source_url"`
	ContentID    *uuid.UUID    `json:"content_id"`
	Summary      *string       `json:"summary"`
	ResearchType *ResearchType `json:"research_type"`
	LayoutType   *LayoutType   `json:"layout_type"`
	OrderIndex   *int          `json:"order_index"`
}

type SavePlanRequest struct {
	CourseDraftID uuid.UUID   `json:"course_draft_id" binding:"required"`
	RegionOrder   []uuid.UUID `json:"region_order"`
}
