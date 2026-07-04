package curriculum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/llm"
)

type Service struct{}

var errDraftGenerationRetryable = errors.New("draft generation retryable")

type DraftBuildOptions struct {
	SkipReview        bool
	GenerationTimeout time.Duration
	ReviewTimeout     time.Duration
}

func NewService() *Service {
	return &Service{}
}

type generatedDraftDocument struct {
	Title              string                `json:"title"`
	Description        string                `json:"description"`
	MainLessons        []generatedMainLesson `json:"main_lessons"`
	Lessons            []generatedMainLesson `json:"lessons"`
	CompletionCriteria []string              `json:"completion_criteria"`
}

type reviewedDraftDocument struct {
	Strengths                 []string              `json:"strengths"`
	Issues                    []string              `json:"issues"`
	MissingAxes               []string              `json:"missing_axes"`
	BoundaryProblems          []string              `json:"boundary_problems"`
	RevisionInstructions      []string              `json:"revision_instructions"`
	RevisedMainLessons        []generatedMainLesson `json:"revised_main_lessons"`
	RevisedCompletionCriteria []string              `json:"revised_completion_criteria"`
}

func (s *Service) BuildInitialDraft(userID uuid.UUID, req CreateCourseDraftRequest) (*DraftAggregate, error) {
	sourceQuery := strings.TrimSpace(req.SourceQuery)
	if sourceQuery == "" {
		return nil, fmt.Errorf("source_query is required")
	}
	learningGoal := strings.TrimSpace(derefString(req.LearningGoal))
	if learningGoal == "" {
		return nil, fmt.Errorf("learning_goal is required")
	}

	title := sourceQuery + " 학습 코스"
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		title = sourceQuery + " Learning Course"
	}
	fallbackPlan := buildFallbackDraftPlan(req)
	description := fallbackPlan.Description

	draftID := uuid.New()
	lessons := make([]DraftLessonTree, 0, len(fallbackPlan.MainLessons))
	for idx, lessonTemplate := range fallbackPlan.MainLessons {
		objective := strings.TrimSpace(lessonTemplate.Objective)
		lessons = append(lessons, DraftLessonTree{
			Lesson: CourseDraftLesson{
				ID:                       uuid.New(),
				CourseDraftID:            draftID,
				Title:                    strings.TrimSpace(lessonTemplate.Title),
				Objective:                &objective,
				RecommendationSearchSpec: BuildFallbackLessonRecommendationSearchSpec(req, lessonTemplate, idx),
				LessonRole:               LessonRoleCore,
				SourceType:               LessonSourceAI,
				OrderIndex:               idx,
			},
		})
	}

	draft := &DraftAggregate{
		Draft: CourseDraft{
			ID:                 draftID,
			UserID:             userID,
			SourceQuery:        sourceQuery,
			LearningGoal:       req.LearningGoal,
			GoalProfileID:      req.GoalProfileID,
			GoalProfileVersion: req.GoalProfileVersion,
			CurrentLevel:       req.CurrentLevel,
			DurationWeeks:      req.DurationWeeks,
			StudyHoursPerWeek:  req.StudyHoursPerWeek,
			PreferredFormat:    req.PreferredFormat,
			GenerationLanguage: normalizeLearningLanguage(req.LearningLanguage),
			Title:              title,
			Description:        &description,
			CompletionCriteria: fallbackPlan.CompletionCriteria,
			Status:             DraftStatusDraft,
		},
		Lessons: lessons,
	}

	if err := s.NormalizeDraftAggregate(draft); err != nil {
		return nil, err
	}
	return draft, nil
}

func (s *Service) BuildDraftWithLLM(ctx context.Context, client llm.Client, userID uuid.UUID, req CreateCourseDraftRequest, opts DraftBuildOptions) (*DraftAggregate, error) {
	if client == nil {
		return s.BuildInitialDraft(userID, req)
	}

	generationPrompt := buildCurriculumGenerationPrompt(req)
	raw, err := completeWithOptionalTimeout(ctx, client, generationPrompt, opts.GenerationTimeout)
	if err != nil {
		return nil, err
	}

	var doc generatedDraftDocument
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("%w: parse llm curriculum json: %v", errDraftGenerationRetryable, err)
	}
	doc = refineGeneratedDraftDocument(req, doc)

	if strings.TrimSpace(doc.Title) == "" || len(doc.MainLessons) == 0 {
		return nil, fmt.Errorf("%w: llm returned incomplete curriculum", errDraftGenerationRetryable)
	}

	if !opts.SkipReview {
		reviewPrompt, err := buildCurriculumReviewPrompt(req, doc)
		if err == nil {
			if reviewRaw, reviewErr := completeWithOptionalTimeout(ctx, client, reviewPrompt, opts.ReviewTimeout); reviewErr == nil {
				var review reviewedDraftDocument
				if err := json.Unmarshal([]byte(reviewRaw), &review); err == nil {
					doc = applyCurriculumReview(doc, review)
				}
			}
		}
	}
	doc = refineGeneratedDraftDocument(req, doc)
	draft, err := s.buildDraftAggregateFromGeneratedDocument(userID, req, doc)
	if err != nil {
		return nil, err
	}

	return draft, nil
}

func completeWithOptionalTimeout(ctx context.Context, client llm.Client, prompt string, timeout time.Duration) (string, error) {
	if timeout <= 0 {
		return client.Complete(ctx, prompt)
	}
	childCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return client.Complete(childCtx, prompt)
}

func buildCurriculumGenerationPrompt(req CreateCourseDraftRequest) string {
	var b strings.Builder
	learningGoal := strings.TrimSpace(derefString(req.LearningGoal))
	language := normalizeLearningLanguage(req.LearningLanguage)
	b.WriteString("당신은 목표 기반 커리큘럼 설계자입니다. JSON 객체만 반환하세요.\n")
	if language == "en" {
		b.WriteString("Learner-facing output language must be English. title, description, lesson title/objective, and completion_criteria must be natural English.\n")
		b.WriteString("Do not translate JSON keys, enum-like values, provider names, or internal identifiers.\n")
	} else {
		b.WriteString("출력 언어: 한국어. 사용자 표시 문구(title, description, lesson title/objective, completion_criteria)는 자연스러운 한국어로 작성하세요.\n")
		b.WriteString("영문 제목이나 영문 설명을 기본값처럼 사용하지 마세요.\n")
	}
	b.WriteString("주제: " + strings.TrimSpace(req.SourceQuery) + "\n")
	b.WriteString("확정 목표: " + learningGoal + "\n")
	for _, detail := range buildGoalPromptDetails(req) {
		b.WriteString(detail + "\n")
	}
	if req.CurrentLevel != nil && strings.TrimSpace(*req.CurrentLevel) != "" {
		b.WriteString("현재 수준: " + strings.TrimSpace(*req.CurrentLevel) + "\n")
	}
	if req.PreferredFormat != nil && strings.TrimSpace(*req.PreferredFormat) != "" {
		b.WriteString("선호 형식: " + strings.TrimSpace(*req.PreferredFormat) + "\n")
	}
	if req.DurationWeeks != nil {
		b.WriteString(fmt.Sprintf("학습 기간(주): %d\n", *req.DurationWeeks))
	}
	if req.StudyHoursPerWeek != nil {
		b.WriteString(fmt.Sprintf("주당 학습 시간: %d\n", *req.StudyHoursPerWeek))
	}
	b.WriteString("출력 JSON:\n")
	b.WriteString("{\n")
	b.WriteString("  \"title\": string,\n")
	b.WriteString("  \"description\": string,\n")
	b.WriteString("  \"main_lessons\": [{ \"title\": string, \"objective\": string }],\n")
	b.WriteString("  \"completion_criteria\": [string]\n")
	b.WriteString("}\n")
	b.WriteString("규칙:\n")
	b.WriteString("- main_lessons는 보통 4~6개의 관문 단계로 구성하세요.\n")
	b.WriteString("- completion_criteria는 2~4개의 짧은 체크리스트로 작성하세요.\n")
	if language == "en" {
		b.WriteString("- completion_criteria should cover only output quality, completion, and usability. Do not add course attendance, sales, gifts, public release, certification, or other outcomes absent from the learner goal.\n")
		b.WriteString("- Lessons must be a goal-achievement flow, not a topic list.\n")
		b.WriteString("- Required functions, situations, or outcomes explicitly named in the confirmed goal must appear in lessons or completion_criteria.\n")
		b.WriteString("- The first lesson should combine beginner setup with understanding the target artifact or goal structure, so the learner can start doing.\n")
		b.WriteString("- Each objective must describe what the learner can do at the end of that stage.\n")
		b.WriteString("- Do not make tools, materials, basic concepts, posture, tuning, or other tiny elements standalone lessons.\n")
		b.WriteString("- Put the final output/event itself in completion_criteria.\n")
		b.WriteString("- If the confirmed goal explicitly includes deployment, upload, or shipping, include it as deployment check, release readiness, or operational verification in the last lesson or completion_criteria.\n")
		b.WriteString("- The last lesson must be integration practice, finishing adjustment, usability check, or performance prep just before the final result.\n")
		b.WriteString("- The last lesson title must not use completion-event words such as complete, finished, sell, gift, participate, pass, or publish.\n")
		b.WriteString("- lesson titles should be short, clear English action phrases.\n")
		b.WriteString("- Do not copy internal flow labels such as setup and entry, first output, core pattern practice, integration practice, or performance prep as lesson titles. Convert them into concrete learner actions for the goal.\n")
		b.WriteString("- description should be concise natural English.\n")
	} else {
		b.WriteString("- completion_criteria는 결과물 완성도, 품질, 사용성 확인만 다루고 수업 참여/판매/선물/공개/자격 취득을 쓰지 마세요.\n")
		b.WriteString("- 리슨은 주제 목록이 아니라 목표 달성 흐름이어야 합니다.\n")
		b.WriteString("- 확정 목표에 명시된 필수 기능, 상황, 결과는 리슨 또는 완료 기준에 반드시 포함하세요.\n")
		b.WriteString("- 첫 리슨은 초보자가 바로 수행을 시작할 수 있도록 기초 준비와 목표물 구조 이해를 묶은 단계로 만드세요.\n")
		b.WriteString("- 각 objective는 단계 끝에서 학습자가 할 수 있는 행동으로 쓰세요.\n")
		b.WriteString("- 도구, 재료, 기초 개념, 자세, 튜닝 같은 작은 요소는 독립 리슨으로 만들지 마세요.\n")
		b.WriteString("- 최종 결과물 완성, 발표, 시험, 실전 참여 자체는 completion_criteria로 두세요.\n")
		b.WriteString("- 확정 목표에 배포, 업로드, 출시가 명시되어 있으면 마지막 리슨 또는 완료 기준에 배포 점검, 출시 준비, 운영 확인으로 반드시 포함하세요.\n")
		b.WriteString("- 마지막 리슨은 최종 결과 직전의 통합 연습, 마감 보정, 사용성 점검, 실전 준비 단계여야 합니다.\n")
		b.WriteString("- 마지막 리슨 title에는 완성/완성된/판매/선물/참여/합격/공개를 쓰지 마세요.\n")
		b.WriteString("- 마지막 리슨 title 예: 마감 보정과 사용성 점검하기, 실전 사용 전 품질 점검하기.\n")
		b.WriteString("- lesson 제목은 짧고 명확한 한국어로 작성하세요.\n")
		b.WriteString("- description은 불필요한 영어 표현 없이 한국어로 작성하세요.\n")
	}
	b.WriteString("- sub-lesson은 만들지 마세요.\n")
	b.WriteString("- recommendation_search_spec은 만들지 마세요. 서버가 패턴 템플릿으로 생성합니다.\n")
	patternMatchGuidance := strings.TrimSpace(req.PatternMatchGuidance)
	if patternMatchGuidance == "" {
		if patternGuidance := buildGoalPatternGuidance(req); patternGuidance != "" {
			if normalizeLearningLanguage(req.LearningLanguage) == "en" {
				b.WriteString("- Goal-pattern structure guide:\n")
			} else {
				b.WriteString("- 목표 패턴 기반 구조 가이드:\n")
			}
			b.WriteString(patternGuidance)
		}
		if subpatternGuidance := buildGoalSubpatternGuidance(req); subpatternGuidance != "" {
			if normalizeLearningLanguage(req.LearningLanguage) == "en" {
				b.WriteString("- Recurring real-goal subtype guide:\n")
			} else {
				b.WriteString("- 자주 등장하는 실제 목표형 세부 가이드:\n")
			}
			b.WriteString(subpatternGuidance)
		}
		if domainGuidance := buildGoalTypeGuidance(req.SourceQuery, buildGoalNarrativeContext(req)); domainGuidance != "" {
			if normalizeLearningLanguage(req.LearningLanguage) != "en" {
				b.WriteString("- 목표 유형별 필수 분해 규칙:\n")
				b.WriteString(domainGuidance)
			}
		}
	} else {
		if normalizeLearningLanguage(req.LearningLanguage) == "en" {
			b.WriteString("- Data-based pattern retrieval guide:\n")
		} else {
			b.WriteString("- 데이터 기반 패턴 검색 가이드:\n")
		}
		b.WriteString(patternMatchGuidance)
		b.WriteString("\n")
		if normalizeLearningLanguage(req.LearningLanguage) == "en" {
			b.WriteString("- The data-based guide is only the closest reference pattern. If it conflicts with the learner confirmed goal, prioritize the confirmed goal.\n")
		} else {
			b.WriteString("- 위 데이터 기반 가이드는 목표와 가장 가까운 참고 패턴입니다. 사용자 확정 목표와 충돌하면 확정 목표를 우선하세요.\n")
		}
	}
	return b.String()
}

func buildCurriculumReviewPrompt(req CreateCourseDraftRequest, doc generatedDraftDocument) (string, error) {
	payload, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("marshal generated curriculum for review: %w", err)
	}

	var b strings.Builder
	b.WriteString("당신은 학습자의 확정 목표를 기준으로 단계형 커리큘럼의 품질을 검토하고 교정하는 교육설계 리뷰어입니다.\n")
	b.WriteString("아래 학습 커리큘럼 초안을 검토하고, 목표를 이루기 위한 관문 단계 구조로 설계되어 있는지 판단한 뒤 더 나은 메인 리슨 구조로 교정하세요.\n")
	b.WriteString("출력은 JSON 객체만 반환하세요.\n")
	if normalizeLearningLanguage(req.LearningLanguage) == "en" {
		b.WriteString("Keep all learner-facing revised titles, objectives, and completion criteria in English.\n")
	} else {
		b.WriteString("수정된 사용자 표시 제목, 목표, 완료 기준은 모두 한국어로 유지하세요.\n")
	}
	b.WriteString("최상위 키는 strengths, issues, missing_axes, boundary_problems, revision_instructions, revised_main_lessons, revised_completion_criteria 입니다.\n")
	b.WriteString("revised_main_lessons의 각 항목은 title, objective를 가져야 합니다.\n")
	b.WriteString("revised_main_lessons는 보통 4~6개, 필요 시 3~8개 범위 안에서만 수정하세요.\n")
	b.WriteString("main lesson은 단순 주제 목록이 아니라 목표 달성을 위한 관문 단계여야 합니다.\n")
	b.WriteString("각 revised main lesson은 학습자가 그 단계 끝에서 무엇을 할 수 있게 되는지를 보여줘야 합니다.\n")
	b.WriteString("기본자세, 도구 이해, 튜닝, 첫 소리, 기초 개념, 재료 소개 같은 작은 요소는 독립 리슨이 아니라 탐험지점 수준의 세부 활동으로 간주하세요.\n")
	b.WriteString("너무 작은 lesson은 삭제하거나 인접한 더 큰 단계로 흡수하는 방향으로 수정하세요.\n")
	b.WriteString("경계가 겹치거나 거의 같은 의미인 lesson은 병합하세요.\n")
	b.WriteString("최종 목표 이벤트(실전 투입, 공식 참여, 발표, 시험 응시)는 revised_main_lessons에 넣지 말고 revised_completion_criteria로 분리하세요.\n")
	b.WriteString("마지막 revised main lesson은 완료 이벤트가 아니라 실전 직전 준비, 적용, 리허설, 통합 수행 단계여야 합니다.\n")
	b.WriteString("좋은 lesson 예: 첫 소리를 안정적으로 낸다, 자주 쓰는 코드를 끊기지 않고 전환한다.\n")
	b.WriteString("나쁜 lesson 예: 기본자세 익히기, 도구 이해하기, 기초 이론 알아보기.\n")
	b.WriteString("사용자 주제(도메인 맥락): " + strings.TrimSpace(req.SourceQuery) + "\n")
	b.WriteString("확정 목표: " + strings.TrimSpace(derefString(req.LearningGoal)) + "\n")
	for _, detail := range buildGoalPromptDetails(req) {
		b.WriteString(detail + "\n")
	}
	if req.CurrentLevel != nil && strings.TrimSpace(*req.CurrentLevel) != "" {
		b.WriteString("현재 수준: " + strings.TrimSpace(*req.CurrentLevel) + "\n")
	}
	if patternGuidance := buildGoalPatternGuidance(req); patternGuidance != "" {
		b.WriteString("목표 패턴 기반 검토 기준:\n")
		b.WriteString(patternGuidance)
	}
	if subpatternGuidance := buildGoalSubpatternGuidance(req); subpatternGuidance != "" {
		b.WriteString("실사용 목표형 세부 검토 기준:\n")
		b.WriteString(subpatternGuidance)
	}
	b.WriteString("검토할 초안 JSON:\n")
	b.Write(payload)
	return b.String(), nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type generatedMainLesson struct {
	Title                    string                         `json:"title"`
	Objective                string                         `json:"objective"`
	RecommendationSearchSpec LessonRecommendationSearchSpec `json:"recommendation_search_spec"`
}

func (s *Service) buildDraftAggregateFromGeneratedDocument(userID uuid.UUID, req CreateCourseDraftRequest, doc generatedDraftDocument) (*DraftAggregate, error) {
	mainLessons := normalizeGeneratedMainLessons(doc.MainLessons)
	if strings.TrimSpace(doc.Title) == "" || len(mainLessons) == 0 {
		return nil, fmt.Errorf("%w: llm returned incomplete curriculum", errDraftGenerationRetryable)
	}

	description := strings.TrimSpace(doc.Description)
	completionCriteria := normalizeCompletionCriteria(doc.CompletionCriteria)
	if len(completionCriteria) == 0 {
		completionCriteria = buildFallbackDraftPlan(req).CompletionCriteria
	}

	draftID := uuid.New()
	draft := &DraftAggregate{
		Draft: CourseDraft{
			ID:                 draftID,
			UserID:             userID,
			SourceQuery:        strings.TrimSpace(req.SourceQuery),
			LearningGoal:       req.LearningGoal,
			GoalProfileID:      req.GoalProfileID,
			GoalProfileVersion: req.GoalProfileVersion,
			CurrentLevel:       req.CurrentLevel,
			DurationWeeks:      req.DurationWeeks,
			StudyHoursPerWeek:  req.StudyHoursPerWeek,
			PreferredFormat:    req.PreferredFormat,
			GenerationLanguage: normalizeLearningLanguage(req.LearningLanguage),
			Title:              strings.TrimSpace(doc.Title),
			Description:        &description,
			CompletionCriteria: completionCriteria,
			Status:             DraftStatusDraft,
		},
		Lessons: make([]DraftLessonTree, 0, len(mainLessons)),
	}

	for lessonIdx, lessonDoc := range mainLessons {
		mainObjective := strings.TrimSpace(lessonDoc.Objective)
		draft.Lessons = append(draft.Lessons, DraftLessonTree{
			Lesson: CourseDraftLesson{
				ID:                       uuid.New(),
				CourseDraftID:            draftID,
				Title:                    strings.TrimSpace(lessonDoc.Title),
				Objective:                &mainObjective,
				RecommendationSearchSpec: normalizeGeneratedLessonRecommendationSearchSpec(req, lessonDoc, lessonIdx),
				LessonRole:               LessonRoleCore,
				SourceType:               LessonSourceAI,
				OrderIndex:               lessonIdx,
			},
		})
	}

	if err := s.NormalizeDraftAggregate(draft); err != nil {
		return nil, err
	}
	return draft, nil
}

func normalizeGeneratedLessonRecommendationSearchSpec(req CreateCourseDraftRequest, lessonDoc generatedMainLesson, lessonIdx int) LessonRecommendationSearchSpec {
	fallback := BuildFallbackLessonRecommendationSearchSpec(req, lessonDoc, lessonIdx)
	generated := lessonDoc.RecommendationSearchSpec
	if fallback.Source == SearchSpecSourcePatternTemplate && fallback.StageRole == "artifact_finish" && generated.StageRole != fallback.StageRole {
		generated.StageRole = fallback.StageRole
		generated.MustInclude = nil
		generated.NiceToHave = nil
		generated.PrimaryQuery = ""
	}
	return NormalizeLessonRecommendationSearchSpec(generated, fallback)
}

func (s *Service) ValidateRecommendationEvent(event RecommendationEvent) error {
	if event.UserID == uuid.Nil {
		return fmt.Errorf("event user_id is required")
	}
	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	return nil
}

func (s *Service) ApplyDraftUpdate(existing DraftAggregate, req UpdateCourseDraftRequest) (*DraftAggregate, error) {
	updated := existing

	if req.Title != nil {
		updated.Draft.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		updated.Draft.Description = &trimmed
	}
	if req.LearningGoal != nil {
		trimmed := strings.TrimSpace(*req.LearningGoal)
		updated.Draft.LearningGoal = &trimmed
	}
	if req.CurrentLevel != nil {
		trimmed := strings.TrimSpace(*req.CurrentLevel)
		updated.Draft.CurrentLevel = &trimmed
	}
	if req.DurationWeeks != nil {
		updated.Draft.DurationWeeks = req.DurationWeeks
	}
	if req.StudyHoursPerWeek != nil {
		updated.Draft.StudyHoursPerWeek = req.StudyHoursPerWeek
	}
	if req.PreferredFormat != nil {
		trimmed := strings.TrimSpace(*req.PreferredFormat)
		updated.Draft.PreferredFormat = &trimmed
	}
	if req.PlanetTypeID != nil {
		updated.Draft.PlanetTypeID = req.PlanetTypeID
	}
	if req.PlanetTextureMapID != nil {
		updated.Draft.PlanetTextureMapID = req.PlanetTextureMapID
	}

	if err := s.NormalizeDraftAggregate(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}
