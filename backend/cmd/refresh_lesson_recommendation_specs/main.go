package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/db"
)

type refreshSummary struct {
	DryRun                bool            `json:"dry_run"`
	Limit                 int             `json:"limit"`
	Scanned               int             `json:"scanned"`
	Refreshable           int             `json:"refreshable"`
	Applied               int             `json:"applied"`
	Unchanged             int             `json:"unchanged"`
	SkippedNoProfile      int             `json:"skipped_no_profile"`
	SkippedManualSpec     int             `json:"skipped_manual_spec"`
	SkippedManualLesson   int             `json:"skipped_manual_lesson"`
	SkippedGenerated      int             `json:"skipped_generated"`
	SkippedAlreadyOverlay int             `json:"skipped_already_overlay"`
	SkippedUnsupported    int             `json:"skipped_unsupported"`
	Failed                int             `json:"failed"`
	Samples               []refreshSample `json:"samples,omitempty"`
}

type refreshSample struct {
	LessonID      uuid.UUID  `json:"lesson_id"`
	DraftID       uuid.UUID  `json:"draft_id"`
	GoalProfileID *uuid.UUID `json:"goal_profile_id,omitempty"`
	Title         string     `json:"title"`
	SourceBefore  string     `json:"source_before,omitempty"`
	SourceAfter   string     `json:"source_after,omitempty"`
	QueryBefore   string     `json:"query_before,omitempty"`
	QueryAfter    string     `json:"query_after,omitempty"`
	MustBefore    []string   `json:"must_before,omitempty"`
	MustAfter     []string   `json:"must_after,omitempty"`
	NiceBefore    []string   `json:"nice_before,omitempty"`
	NiceAfter     []string   `json:"nice_after,omitempty"`
	Applied       bool       `json:"applied"`
	Skipped       string     `json:"skipped,omitempty"`
	Error         string     `json:"error,omitempty"`
}

type refreshLessonRow struct {
	LessonID         uuid.UUID
	DraftID          uuid.UUID
	GoalProfileID    *uuid.UUID
	SourceQuery      string
	LearningGoal     *string
	Language         string
	Title            string
	Objective        *string
	OrderIndex       int
	LessonSourceType string
	RawSpec          []byte
	LearningIntent   curriculum.LearningIntentProfile
}

type refreshDecision struct {
	Spec    curriculum.LessonRecommendationSearchSpec
	Skip    string
	Changed bool
}

func main() {
	var limit int
	var apply bool
	flag.IntVar(&limit, "limit", 100, "maximum draft lessons to inspect")
	flag.BoolVar(&apply, "apply", false, "write refreshed recommendation search specs")
	flag.Parse()

	if limit <= 0 {
		log.Fatal("-limit must be greater than 0")
	}

	loadEnv()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	rows, err := pool.Query(ctx, `
		SELECT cdl.id, cdl.course_draft_id, cgp.id,
		       cd.source_query, cd.learning_goal, COALESCE(cd.generation_language, 'ko'),
		       cdl.title, cdl.objective, cdl.order_index, COALESCE(cdl.source_type, ''),
		       COALESCE(cdl.recommendation_search_spec, '{}'::jsonb),
		       COALESCE(cgp.learning_intent_profile, '{}'::jsonb)
		FROM course_draft_lessons cdl
		JOIN course_drafts cd ON cd.id = cdl.course_draft_id
		LEFT JOIN LATERAL (
			SELECT gp.id, gp.learning_intent_profile
			FROM course_goal_profiles gp
			WHERE gp.user_id = cd.user_id
			  AND gp.is_active = true
			  AND (gp.id = cd.goal_profile_id OR gp.course_draft_id = cd.id)
			ORDER BY CASE WHEN gp.id = cd.goal_profile_id THEN 0 ELSE 1 END, gp.updated_at DESC
			LIMIT 1
		) cgp ON true
		WHERE cd.status IN ('draft', 'confirmed', 'learning')
		ORDER BY cdl.updated_at DESC, cdl.id ASC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Fatalf("failed to query draft lessons: %v", err)
	}
	defer rows.Close()

	result := refreshSummary{DryRun: !apply, Limit: limit}
	for rows.Next() {
		var row refreshLessonRow
		var rawProfile []byte
		if err := rows.Scan(
			&row.LessonID, &row.DraftID, &row.GoalProfileID,
			&row.SourceQuery, &row.LearningGoal, &row.Language,
			&row.Title, &row.Objective, &row.OrderIndex, &row.LessonSourceType,
			&row.RawSpec, &rawProfile,
		); err != nil {
			result.Failed++
			continue
		}
		result.Scanned++
		if err := json.Unmarshal(rawProfile, &row.LearningIntent); err != nil {
			result.Failed++
			addSample(&result, sampleFromRow(row, refreshDecision{Skip: "profile_unmarshal_failed"}, false, fmt.Sprintf("profile unmarshal: %v", err)))
			continue
		}

		decision := buildRefreshDecision(row)
		countDecision(&result, decision)
		if decision.Skip != "" || !decision.Changed {
			addSample(&result, sampleFromRow(row, decision, false, ""))
			continue
		}

		result.Refreshable++
		applied := false
		if apply {
			payload, err := curriculum.MarshalLessonRecommendationSearchSpec(decision.Spec)
			if err != nil {
				result.Failed++
				addSample(&result, sampleFromRow(row, decision, false, fmt.Sprintf("marshal: %v", err)))
				continue
			}
			current := strings.TrimSpace(string(row.RawSpec))
			if current == "" {
				current = "{}"
			}
			tag, err := pool.Exec(ctx, `
				UPDATE course_draft_lessons
				SET recommendation_search_spec = $2::jsonb,
				    updated_at = NOW()
				WHERE id = $1
				  AND COALESCE(recommendation_search_spec, '{}'::jsonb) = $3::jsonb
			`, row.LessonID, payload, current)
			if err != nil {
				result.Failed++
				addSample(&result, sampleFromRow(row, decision, false, fmt.Sprintf("update: %v", err)))
				continue
			}
			if tag.RowsAffected() > 0 {
				result.Applied++
				applied = true
			}
		}
		addSample(&result, sampleFromRow(row, decision, applied, ""))
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("failed to iterate draft lessons: %v", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("failed to print summary: %v", err)
	}
}

func buildRefreshDecision(row refreshLessonRow) refreshDecision {
	profile := curriculum.NormalizeLearningIntentProfile(row.LearningIntent)
	if isEmptyLearningIntentProfile(profile) {
		return refreshDecision{Skip: "no_profile"}
	}
	current := curriculum.ParseLessonRecommendationSearchSpec(row.RawSpec)
	if strings.EqualFold(row.LessonSourceType, string(curriculum.LessonSourceManual)) {
		return refreshDecision{Spec: current, Skip: "manual_lesson"}
	}
	switch current.Source {
	case curriculum.SearchSpecSourceManual:
		return refreshDecision{Spec: current, Skip: "manual_spec"}
	case curriculum.SearchSpecSourceGenerated:
		return refreshDecision{Spec: current, Skip: "generated"}
	case curriculum.SearchSpecSourcePatternTemplateOverlay, curriculum.SearchSpecSourceFallbackOverlay:
		return refreshDecision{Spec: current, Skip: "already_overlay"}
	case "", curriculum.SearchSpecSourceFallback, curriculum.SearchSpecSourcePatternTemplate:
		// refreshable
	default:
		return refreshDecision{Spec: current, Skip: "unsupported_source"}
	}

	next := current
	if strings.TrimSpace(current.PrimaryQuery) == "" {
		next = curriculum.BuildFallbackLessonRecommendationSearchSpecForLessonWithProfile(
			row.SourceQuery,
			row.LearningGoal,
			row.Language,
			row.Title,
			row.Objective,
			row.OrderIndex,
			profile,
		)
	} else {
		next = curriculum.ApplyLearningIntentProfileToLessonSearchSpec(current, profile)
	}
	if strings.TrimSpace(next.PrimaryQuery) == "" || lessonSearchSpecsEqual(current, next) {
		return refreshDecision{Spec: current, Changed: false}
	}
	return refreshDecision{Spec: next, Changed: true}
}

func lessonSearchSpecsEqual(left, right curriculum.LessonRecommendationSearchSpec) bool {
	left = curriculum.NormalizeLessonRecommendationSearchSpec(left, curriculum.LessonRecommendationSearchSpec{})
	right = curriculum.NormalizeLessonRecommendationSearchSpec(right, curriculum.LessonRecommendationSearchSpec{})
	leftPayload, leftErr := json.Marshal(left)
	rightPayload, rightErr := json.Marshal(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	return string(leftPayload) == string(rightPayload)
}

func countDecision(result *refreshSummary, decision refreshDecision) {
	switch decision.Skip {
	case "no_profile":
		result.SkippedNoProfile++
	case "manual_spec":
		result.SkippedManualSpec++
	case "manual_lesson":
		result.SkippedManualLesson++
	case "generated":
		result.SkippedGenerated++
	case "already_overlay":
		result.SkippedAlreadyOverlay++
	case "unsupported_source":
		result.SkippedUnsupported++
	case "":
		if !decision.Changed {
			result.Unchanged++
		}
	}
}

func sampleFromRow(row refreshLessonRow, decision refreshDecision, applied bool, errText string) refreshSample {
	current := curriculum.ParseLessonRecommendationSearchSpec(row.RawSpec)
	return refreshSample{
		LessonID:      row.LessonID,
		DraftID:       row.DraftID,
		GoalProfileID: row.GoalProfileID,
		Title:         row.Title,
		SourceBefore:  current.Source,
		SourceAfter:   decision.Spec.Source,
		QueryBefore:   current.PrimaryQuery,
		QueryAfter:    decision.Spec.PrimaryQuery,
		MustBefore:    current.MustInclude,
		MustAfter:     decision.Spec.MustInclude,
		NiceBefore:    current.NiceToHave,
		NiceAfter:     decision.Spec.NiceToHave,
		Applied:       applied,
		Skipped:       decision.Skip,
		Error:         errText,
	}
}

func addSample(result *refreshSummary, sample refreshSample) {
	if len(result.Samples) >= 12 {
		return
	}
	if sampleLooksUnchanged(sample) && len(result.Samples) >= 2 {
		return
	}
	if sample.Skipped == "no_profile" && len(result.Samples) >= 4 {
		return
	}
	result.Samples = append(result.Samples, sample)
}

func sampleLooksUnchanged(sample refreshSample) bool {
	return sample.Skipped == "" && !sample.Applied && sample.SourceBefore == sample.SourceAfter && sample.QueryBefore == sample.QueryAfter &&
		stringSlicesEqual(sample.MustBefore, sample.MustAfter) && stringSlicesEqual(sample.NiceBefore, sample.NiceAfter)
}

func stringSlicesEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func isEmptyLearningIntentProfile(profile curriculum.LearningIntentProfile) bool {
	profile = curriculum.NormalizeLearningIntentProfile(profile)
	return profile.ConfirmedGoal == "" && profile.LearnerLevel == "" && profile.Purpose == "" && profile.DesiredOutput == "" &&
		len(profile.CurrentBlockers) == 0 && len(profile.PreferredActivities) == 0 && len(profile.SuccessCriteria) == 0
}

func loadEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}
