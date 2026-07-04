package curriculum

import "github.com/google/uuid"

type TodayTaskKind string

const (
	TodayTaskKindCreateCourse   TodayTaskKind = "create_course"
	TodayTaskKindStartPlanet    TodayTaskKind = "start_planet"
	TodayTaskKindContinuePoint  TodayTaskKind = "continue_point"
	TodayTaskKindCompletePlanet TodayTaskKind = "complete_planet"
	TodayTaskKindReviewRecords  TodayTaskKind = "review_records"
)

type TodayTaskResponse struct {
	Task *TodayTask `json:"task"`
}

type TodayTask struct {
	Kind        TodayTaskKind `json:"kind"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CTALabel    string        `json:"cta_label"`
	Href        string        `json:"href"`
	PlanetID    *uuid.UUID    `json:"planet_id,omitempty"`
	PointID     *uuid.UUID    `json:"point_id,omitempty"`
}
