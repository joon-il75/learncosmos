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
	"github.com/learnweaver/backend/internal/domain/goal"
	"github.com/learnweaver/backend/internal/pkg/db"
)

type summary struct {
	DryRun       bool           `json:"dry_run"`
	Limit        int            `json:"limit"`
	Scanned      int            `json:"scanned"`
	Generated    int            `json:"generated"`
	EmptyProfile int            `json:"empty_profile"`
	Applied      int            `json:"applied"`
	Failed       int            `json:"failed"`
	Samples      []sampleRecord `json:"samples,omitempty"`
}

type sampleRecord struct {
	GoalProfileID uuid.UUID                  `json:"goal_profile_id"`
	ConfirmedGoal string                     `json:"confirmed_goal,omitempty"`
	Profile       goal.LearningIntentProfile `json:"profile"`
	Applied       bool                       `json:"applied"`
	Error         string                     `json:"error,omitempty"`
}

func main() {
	var limit int
	var apply bool
	flag.IntVar(&limit, "limit", 100, "maximum goal profiles to inspect")
	flag.BoolVar(&apply, "apply", false, "write generated learning intent profiles")
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
		SELECT id, course_draft_id, user_id,
		       user_intent, motivation, usage_context,
		       confirmed_goal, goal_type, output_type, difficulty_level, time_horizon,
		       COALESCE(language, 'ko') AS language,
		       summarized_context
		FROM course_goal_profiles
		WHERE COALESCE(learning_intent_profile, '{}'::jsonb) = '{}'::jsonb
		  AND (
		    btrim(COALESCE(confirmed_goal, '')) <> ''
		    OR btrim(COALESCE(motivation, '')) <> ''
		    OR btrim(COALESCE(usage_context, '')) <> ''
		    OR btrim(COALESCE(output_type, '')) <> ''
		    OR btrim(COALESCE(goal_type, '')) <> ''
		  )
		ORDER BY updated_at DESC, id ASC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Fatalf("failed to query target goal profiles: %v", err)
	}
	defer rows.Close()

	result := summary{DryRun: !apply, Limit: limit}
	for rows.Next() {
		var profile goal.GoalProfile
		if err := rows.Scan(
			&profile.ID, &profile.CourseDraftID, &profile.UserID,
			&profile.UserIntent, &profile.Motivation, &profile.UsageContext,
			&profile.ConfirmedGoal, &profile.GoalType, &profile.OutputType, &profile.DifficultyLevel, &profile.TimeHorizon,
			&profile.Language, &profile.SummarizedContext,
		); err != nil {
			result.Failed++
			continue
		}
		result.Scanned++

		learningIntent := goal.BuildLearningIntentProfileWithSource(&profile, goal.LearningIntentSourceBackfill)
		sample := sampleRecord{
			GoalProfileID: profile.ID,
			ConfirmedGoal: strings.TrimSpace(deref(profile.ConfirmedGoal)),
			Profile:       learningIntent,
		}
		if goal.IsEmptyLearningIntentProfile(learningIntent) || !goal.IsActionableLearningIntentProfile(learningIntent) {
			result.EmptyProfile++
			addSample(&result, sample)
			continue
		}
		result.Generated++

		if apply {
			payload, err := json.Marshal(learningIntent)
			if err != nil {
				result.Failed++
				sample.Error = fmt.Sprintf("marshal: %v", err)
				addSample(&result, sample)
				continue
			}
			tag, err := pool.Exec(ctx, `
				UPDATE course_goal_profiles
				SET learning_intent_profile = $2,
				    updated_at = now()
				WHERE id = $1
				  AND COALESCE(learning_intent_profile, '{}'::jsonb) = '{}'::jsonb
			`, profile.ID, payload)
			if err != nil {
				result.Failed++
				sample.Error = fmt.Sprintf("update: %v", err)
				addSample(&result, sample)
				continue
			}
			if tag.RowsAffected() > 0 {
				result.Applied++
				sample.Applied = true
			}
		}
		addSample(&result, sample)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("failed to iterate target goal profiles: %v", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("failed to print summary: %v", err)
	}
}

func loadEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}

func addSample(result *summary, sample sampleRecord) {
	if len(result.Samples) >= 5 {
		return
	}
	result.Samples = append(result.Samples, sample)
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
