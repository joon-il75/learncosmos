package curriculum

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentSearchCandidate struct {
	ContentID    *uuid.UUID
	ResourceType ResourceType
	ContentType  string
	Title        string
	Description  *string
	ThumbnailURL *string
	ExternalURL  *string
	PriceType    PriceType
	Language     string
	QualityScore float64
	CreatedAt    time.Time
	RankScore    float64
}

var (
	errDraftResourceNotFound                = errors.New("draft resource not found")
	errDraftResourceInvalidState            = errors.New("draft resource invalid state")
	errInactivePlanetAccessDenied           = errors.New("inactive planet access denied")
	errDraftStartNotAllowed                 = errors.New("draft start not allowed")
	errLessonRuntimeInvalidInput            = errors.New("lesson runtime invalid input")
	errLearningLessonNotFound               = errors.New("learning lesson not found")
	errPointRuntimeInvalidInput             = errors.New("point runtime invalid input")
	errLearningPointNotFound                = errors.New("learning point not found")
	errLearningPointCompletionNotReady      = errors.New("learning point completion not ready")
	errLearningPointMaterialNotReady        = errors.New("learning point material not ready")
	errDraftLevelNotFound                   = errors.New("draft level not found")
	errDraftLevelInvalidInput               = errors.New("draft level invalid input")
	errDraftLessonNotFound                  = errors.New("draft lesson not found")
	errDraftLessonInvalidInput              = errors.New("draft lesson invalid input")
	errDraftMemoInvalidInput                = errors.New("draft memo invalid input")
	errDraftJournalInvalidInput             = errors.New("draft journal invalid input")
	errDraftRecordInvalidInput              = errors.New("draft record invalid input")
	errDraftArtifactInvalidInput            = errors.New("draft artifact invalid input")
	errDraftStructureNotFound               = errors.New("draft structure not found")
	errDraftStructureInvalidInput           = errors.New("draft structure invalid input")
	errDraftCompleteNotAllowed              = errors.New("draft complete not allowed")
	errDraftContentNotFound                 = errors.New("draft content not found")
	errDraftResourceDuplicate               = errors.New("draft resource duplicate")
	errLearningPointJournalInvalidInput     = errors.New("learning point journal invalid input")
	errLearningPointObservationInvalidInput = errors.New("learning point observation invalid input")
	errLearningPointRecordInvalidInput      = errors.New("learning point record invalid input")
	errLearningPointArtifactInvalidInput    = errors.New("learning point artifact invalid input")
	errLearningPointPracticeLogInvalidInput = errors.New("learning point practice log invalid input")
	errLearningPointBlockInvalidInput       = errors.New("learning point block invalid input")
	errLearningPointBlockLimitExceeded      = errors.New("learning point block limit exceeded")
	errLearningPointMaterialLocked          = errors.New("learning point material locked")
	errLearningPointAISummaryInvalidInput   = errors.New("learning point ai summary invalid input")
	errLearningPointAttachmentLimitExceeded = errors.New("learning point attachment limit exceeded")
	errLearningPointQuestionInvalidInput    = errors.New("learning point question invalid input")
	errLearningPointSelfEvalInvalidInput    = errors.New("learning point self evaluation invalid input")
)

type LessonCandidateSearchOptions struct {
	PreferredFormat   *string
	PreferredLanguage string
	Now               time.Time
}

type LessonCandidateSearchDiagnostics struct {
	LexicalCandidates []ContentSearchCandidate
	VectorCandidates  []ContentSearchCandidate
	MergedCandidates  []ContentSearchCandidate
	RetrievalLimit    int
}

type Repository struct {
	pool *pgxpool.Pool
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

func hashCurriculumSeed(value string) int {
	hash := 0
	for _, ch := range value {
		hash = ((hash << 5) - hash) + int(ch)
	}
	if hash < 0 {
		return -hash
	}
	return hash
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func valueOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func valueOrDefault(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func isUndefinedTableError(err error, tableName string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "42P01" && strings.Contains(pgErr.Message, tableName)
}

func learningPointCompletionReady(questions []CoursePointQuestion, selfEvaluation *CoursePointSelfEvaluation) bool {
	answeredQuestions := 0
	for _, question := range questions {
		if strings.TrimSpace(valueOrEmpty(question.Answer)) == "" {
			continue
		}
		answeredQuestions++
	}
	if answeredQuestions == 0 || selfEvaluation == nil {
		return false
	}

	return learningPointSelfEvaluationQualityReady(selfEvaluation)
}

func learningPointSelfEvaluationQualityReady(selfEvaluation *CoursePointSelfEvaluation) bool {
	if selfEvaluation == nil {
		return false
	}
	if !validSelfEvaluationScore(selfEvaluation.UnderstandingScore) ||
		!validSelfEvaluationScore(selfEvaluation.ApplicationScore) ||
		!validSelfEvaluationScore(selfEvaluation.ProficiencyScore) ||
		!validSelfEvaluationScore(selfEvaluation.ProblemSolvingScore) ||
		!validSelfEvaluationScore(selfEvaluation.ExpressionScore) {
		return false
	}
	requiredNotes := []string{
		selfEvaluation.UnderstandingReason,
		selfEvaluation.ApplicationReason,
		selfEvaluation.ProficiencyReason,
		selfEvaluation.ProblemSolvingReason,
		selfEvaluation.ExpressionReason,
		selfEvaluation.GoalAlignmentNote,
	}
	for _, note := range requiredNotes {
		if strings.TrimSpace(note) == "" {
			return false
		}
	}
	return true
}

func validSelfEvaluationScore(score *int) bool {
	if score == nil {
		return false
	}
	return *score >= 1 && *score <= 5
}
