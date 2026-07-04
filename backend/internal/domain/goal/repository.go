package goal

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	secretmanager "github.com/learnweaver/backend/internal/pkg/secretmanager"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

type UserRuntimeAIConfig struct {
	Mode        string
	Provider    string
	APIKey      string
	EndpointURL *string
}

type AIUsageEventInput struct {
	UserID           uuid.UUID
	Source           string
	Provider         string
	Model            string
	Feature          string
	BillingStatus    string
	InputTokens      *int
	OutputTokens     *int
	EstimatedCostUSD *float64
	Success          bool
	ErrorCode        string
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func normalizeLearningLanguage(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func (r *Repository) GetUserLearningLanguage(ctx context.Context, userID uuid.UUID) string {
	var language string
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(learning_language, 'ko')
		FROM users
		WHERE id = $1
	`, userID).Scan(&language)
	if err != nil {
		return "ko"
	}
	return normalizeLearningLanguage(language)
}

func (r *Repository) GetActiveGoalByUser(ctx context.Context, userID uuid.UUID) (*GoalProfile, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, course_draft_id, user_id,
		       user_intent, motivation, usage_context,
		       confirmed_goal, goal_type, output_type, difficulty_level, time_horizon,
		       COALESCE(language, 'ko') AS language,
		       summarized_context, COALESCE(learning_intent_profile, '{}'::jsonb), rebuild_decision, revision_snapshot,
		       interview_state, interview_messages,
		       version, is_active, created_at, updated_at
		FROM course_goal_profiles
		WHERE user_id = $1
		  AND course_draft_id IS NULL
		  AND is_active = true
		ORDER BY version DESC, updated_at DESC
		LIMIT 1
	`, userID)
	return scanGoalProfile(row)
}

func (r *Repository) GetActiveGoal(ctx context.Context, courseDraftID uuid.UUID) (*GoalProfile, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, course_draft_id, user_id,
		       user_intent, motivation, usage_context,
		       confirmed_goal, goal_type, output_type, difficulty_level, time_horizon,
		       COALESCE(language, 'ko') AS language,
		       summarized_context, COALESCE(learning_intent_profile, '{}'::jsonb), rebuild_decision, revision_snapshot,
		       interview_state, interview_messages,
		       version, is_active, created_at, updated_at
		FROM course_goal_profiles
		WHERE course_draft_id = $1 AND is_active = true
		ORDER BY version DESC
		LIMIT 1
	`, courseDraftID)
	return scanGoalProfile(row)
}

func (r *Repository) GetGoalByID(ctx context.Context, id uuid.UUID) (*GoalProfile, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, course_draft_id, user_id,
		       user_intent, motivation, usage_context,
		       confirmed_goal, goal_type, output_type, difficulty_level, time_horizon,
		       COALESCE(language, 'ko') AS language,
		       summarized_context, COALESCE(learning_intent_profile, '{}'::jsonb), rebuild_decision, revision_snapshot,
		       interview_state, interview_messages,
		       version, is_active, created_at, updated_at
		FROM course_goal_profiles
		WHERE id = $1
	`, id)
	return scanGoalProfile(row)
}

func (r *Repository) CreateGoalProfile(ctx context.Context, p *GoalProfile) error {
	msgs, err := json.Marshal(p.Messages)
	if err != nil {
		return err
	}
	if isEmptyLearningIntentProfile(p.LearningIntent) {
		p.LearningIntent = BuildLearningIntentProfile(p)
	} else {
		p.LearningIntent = NormalizeLearningIntentProfile(p.LearningIntent)
	}
	learningIntent, err := json.Marshal(p.LearningIntent)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx, `
		INSERT INTO course_goal_profiles
		  (course_draft_id, user_id, user_intent, language, learning_intent_profile, interview_state, interview_messages, version, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
		RETURNING id, created_at, updated_at
	`, p.CourseDraftID, p.UserID, p.UserIntent, normalizeLearningLanguage(p.Language), learningIntent, p.InterviewState, msgs, p.Version,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *Repository) UpdateGoalProfile(ctx context.Context, p *GoalProfile) error {
	msgs, err := json.Marshal(p.Messages)
	if err != nil {
		return err
	}
	var rebuildDecision *string
	if p.RebuildDecision != nil {
		value := string(*p.RebuildDecision)
		rebuildDecision = &value
	}
	if isEmptyLearningIntentProfile(p.LearningIntent) {
		p.LearningIntent = BuildLearningIntentProfile(p)
	} else {
		p.LearningIntent = NormalizeLearningIntentProfile(p.LearningIntent)
	}
	learningIntent, err := json.Marshal(p.LearningIntent)
	if err != nil {
		return err
	}
	var snapshotJSON []byte
	if p.RevisionSnapshot != nil {
		snapshotJSON, err = json.Marshal(p.RevisionSnapshot)
		if err != nil {
			return err
		}
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_goal_profiles SET
			user_intent        = $2,
			motivation         = $3,
			usage_context      = $4,
			confirmed_goal     = $5,
			goal_type          = $6,
			output_type        = $7,
			difficulty_level   = $8,
			time_horizon       = $9,
			language           = $10,
			summarized_context = $11,
			learning_intent_profile = $12,
			rebuild_decision   = $13,
			revision_snapshot  = $14,
			interview_state    = $15,
			interview_messages = $16,
			version            = $17,
			updated_at         = now()
		WHERE id = $1
	`, p.ID, p.UserIntent, p.Motivation, p.UsageContext,
		p.ConfirmedGoal, p.GoalType, p.OutputType, p.DifficultyLevel, p.TimeHorizon,
		normalizeLearningLanguage(p.Language), p.SummarizedContext, learningIntent, rebuildDecision, snapshotJSON, p.InterviewState, msgs, p.Version,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if p.CourseDraftID != nil && p.ConfirmedGoal != nil && strings.TrimSpace(*p.ConfirmedGoal) != "" {
		if _, err := r.pool.Exec(ctx, `
			UPDATE course_drafts
			SET learning_goal = $1,
			    goal_profile_id = $2,
			    goal_profile_version = $3,
			    updated_at = now()
			WHERE id = $4
			  AND user_id = $5
		`, strings.TrimSpace(*p.ConfirmedGoal), p.ID, p.Version, *p.CourseDraftID, p.UserID); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) CancelGoalRevision(ctx context.Context, goalID uuid.UUID) (*GoalProfile, error) {
	profile, err := r.GetGoalByID(ctx, goalID)
	if err != nil {
		return nil, err
	}
	if profile.RevisionSnapshot == nil {
		profile.InterviewState = StateConfirmed
		profile.RebuildDecision = nil
		if err := r.UpdateGoalProfile(ctx, profile); err != nil {
			return nil, err
		}
		return profile, nil
	}
	snapshot := profile.RevisionSnapshot
	profile.UserIntent = snapshot.UserIntent
	profile.Motivation = snapshot.Motivation
	profile.UsageContext = snapshot.UsageContext
	profile.ConfirmedGoal = snapshot.ConfirmedGoal
	profile.GoalType = snapshot.GoalType
	profile.OutputType = snapshot.OutputType
	profile.DifficultyLevel = snapshot.DifficultyLevel
	profile.TimeHorizon = snapshot.TimeHorizon
	profile.SummarizedContext = snapshot.SummarizedContext
	profile.LearningIntent = snapshot.LearningIntent
	profile.InterviewState = snapshot.InterviewState
	if profile.InterviewState != StateConfirmed {
		profile.InterviewState = StateConfirmed
	}
	profile.Messages = append([]InterviewMessage(nil), snapshot.Messages...)
	profile.Version = snapshot.Version
	profile.RebuildDecision = nil
	profile.RevisionSnapshot = nil
	if err := r.UpdateGoalProfile(ctx, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *Repository) UpdateRebuildDecision(ctx context.Context, goalID uuid.UUID, decision RebuildDecision) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_goal_profiles
		SET rebuild_decision = $2,
		    interview_state  = 'confirmed',
		    updated_at       = now()
		WHERE id = $1
	`, goalID, decision)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeactivateGoals(ctx context.Context, courseDraftID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE course_goal_profiles SET is_active = false WHERE course_draft_id = $1`,
		courseDraftID,
	)
	return err
}

func (r *Repository) DeactivatePredraftGoalsByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE course_goal_profiles
		SET is_active = false,
		    updated_at = now()
		WHERE user_id = $1
		  AND course_draft_id IS NULL
		  AND is_active = true
	`, userID)
	return err
}

func (r *Repository) AttachDraftToGoal(ctx context.Context, goalID, userID, courseDraftID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE course_goal_profiles
		SET course_draft_id = $3,
		    updated_at = now()
		WHERE id = $1
		  AND user_id = $2
		  AND course_draft_id IS NULL
		  AND is_active = true
	`, goalID, userID, courseDraftID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetSystemLLMSetting(ctx context.Context, feature string) (provider, model string, err error) {
	err = r.pool.QueryRow(ctx, `
		SELECT provider, model FROM system_llm_settings WHERE feature = $1
	`, feature).Scan(&provider, &model)
	return
}

func (r *Repository) GetSystemAPIKey(ctx context.Context, provider string) (string, error) {
	var encryptedKey *string
	err := r.pool.QueryRow(ctx, `
		SELECT api_key_encrypted FROM system_api_keys WHERE provider = $1
	`, provider).Scan(&encryptedKey)
	if err == nil && encryptedKey != nil && strings.TrimSpace(*encryptedKey) != "" {
		decrypted, decErr := secretmanager.Decrypt(*encryptedKey)
		if decErr != nil {
			return *encryptedKey, nil // backward compat: plaintext
		}
		return decrypted, nil
	}
	// fallback to environment variable
	envKey := ""
	switch strings.ToLower(provider) {
	case "openai":
		envKey = os.Getenv("OPENAI_API_KEY")
	case "anthropic":
		envKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	return envKey, nil
}

func (r *Repository) GetUserRuntimeAIConfig(ctx context.Context, userID uuid.UUID) (*UserRuntimeAIConfig, error) {
	var provider string
	var apiKeyEncrypted *string
	var endpointURL *string
	var isEnabled bool

	err := r.pool.QueryRow(ctx, `
		SELECT provider, api_key_encrypted, endpoint_url, COALESCE(is_enabled, true)
		FROM user_api_keys
		WHERE user_id = $1
	`, userID).Scan(&provider, &apiKeyEncrypted, &endpointURL, &isEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	config := &UserRuntimeAIConfig{
		Mode:        "managed_credit",
		Provider:    strings.TrimSpace(provider),
		EndpointURL: endpointURL,
	}
	if isEnabled && apiKeyEncrypted != nil && strings.TrimSpace(*apiKeyEncrypted) != "" {
		decrypted, err := secretmanager.Decrypt(*apiKeyEncrypted)
		if err != nil {
			return nil, err
		}
		config.Mode = "byok"
		config.APIKey = strings.TrimSpace(decrypted)
	}
	return config, nil
}

func (r *Repository) RecordAIUsageEvent(ctx context.Context, input AIUsageEventInput) error {
	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = "byok"
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO ai_usage_events (
			id, user_id, source, provider, model, feature, billing_status,
			input_tokens, output_tokens, estimated_cost_usd, success, error_code, created_at
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, NULLIF($7, ''),
		        $8, $9, $10, $11, NULLIF($12, ''), NOW())
	`, uuid.New(), input.UserID, source, strings.TrimSpace(input.Provider), strings.TrimSpace(input.Model),
		strings.TrimSpace(input.Feature), strings.TrimSpace(input.BillingStatus), input.InputTokens, input.OutputTokens,
		input.EstimatedCostUSD, input.Success, strings.TrimSpace(input.ErrorCode))
	return err
}

func (r *Repository) RecordBYOKValidation(ctx context.Context, userID uuid.UUID, provider string, valid bool, message string) error {
	status := "valid"
	var errorText *string
	var lastFailedAt *time.Time
	var nextRetryAt *time.Time
	now := time.Now()

	if !valid {
		status = "invalid"
		trimmed := logsafe.ProviderErrorCode(message)
		errorText = &trimmed
		lastFailedAt = &now
		retry := now.Add(15 * time.Minute)
		nextRetryAt = &retry
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_api_keys (id, user_id, provider, last_validation_status, last_validated_at, last_validation_error, last_failed_at, next_retry_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), $5, $6, $7, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET last_validation_status = EXCLUDED.last_validation_status,
		    last_validated_at = NOW(),
		    last_validation_error = EXCLUDED.last_validation_error,
		    last_failed_at = EXCLUDED.last_failed_at,
		    next_retry_at = EXCLUDED.next_retry_at,
		    updated_at = NOW()
	`, uuid.New(), userID, strings.TrimSpace(strings.ToLower(provider)), status, errorText, lastFailedAt, nextRetryAt)
	return err
}

func (r *Repository) VerifyCourseDraftOwner(ctx context.Context, courseDraftID, userID uuid.UUID) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM course_drafts WHERE id = $1 AND user_id = $2
	`, courseDraftID, userID).Scan(&count)
	return count > 0, err
}

func scanGoalProfile(row pgx.Row) (*GoalProfile, error) {
	var p GoalProfile
	var courseDraftID *uuid.UUID
	var msgsRaw []byte
	var learningIntentRaw []byte
	var rebuildDecision *string
	var revisionSnapshot []byte
	err := row.Scan(
		&p.ID, &courseDraftID, &p.UserID,
		&p.UserIntent, &p.Motivation, &p.UsageContext,
		&p.ConfirmedGoal, &p.GoalType, &p.OutputType, &p.DifficultyLevel, &p.TimeHorizon,
		&p.Language, &p.SummarizedContext, &learningIntentRaw, &rebuildDecision, &revisionSnapshot,
		&p.InterviewState, &msgsRaw,
		&p.Version, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.CourseDraftID = courseDraftID
	p.Language = normalizeLearningLanguage(p.Language)
	if len(learningIntentRaw) > 0 {
		_ = json.Unmarshal(learningIntentRaw, &p.LearningIntent)
	}
	p.LearningIntent = NormalizeLearningIntentProfile(p.LearningIntent)
	if rebuildDecision != nil {
		d := RebuildDecision(*rebuildDecision)
		p.RebuildDecision = &d
	}
	if len(revisionSnapshot) > 0 {
		var snapshot GoalRevisionState
		if err := json.Unmarshal(revisionSnapshot, &snapshot); err == nil {
			p.RevisionSnapshot = &snapshot
		}
	}
	if len(msgsRaw) > 0 {
		_ = json.Unmarshal(msgsRaw, &p.Messages)
	}
	if p.Messages == nil {
		p.Messages = []InterviewMessage{}
	}
	return &p, nil
}
