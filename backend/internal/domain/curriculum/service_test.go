package curriculum

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type stubLLMClient struct {
	responses []string
	errs      []error
	callCount int
}

func (s *stubLLMClient) Complete(_ context.Context, _ string) (string, error) {
	idx := s.callCount
	s.callCount++
	if idx < len(s.errs) && s.errs[idx] != nil {
		return "", s.errs[idx]
	}
	if idx < len(s.responses) {
		return s.responses[idx], nil
	}
	return "", fmt.Errorf("unexpected llm call %d", idx)
}

func TestNormalizeDraftAggregate(t *testing.T) {
	svc := NewService()
	userID := uuid.New()

	draft := DraftAggregate{
		Draft: CourseDraft{
			UserID:      userID,
			SourceQuery: "  통기타 코드 배우기  ",
			Title:       "  첫 기타 코스  ",
		},
		Lessons: []DraftLessonTree{
			{
				Lesson: CourseDraftLesson{Title: " 기초 ", OrderIndex: 5},
				SubLessons: []DraftLessonTree{
					{
						Lesson: CourseDraftLesson{Title: " 코드 잡기 ", OrderIndex: 9},
						Points: []DraftPointAggregate{
							{Point: CourseDraftPoint{Title: "  기본 코드 영상  ", OrderIndex: 7}},
						},
					},
				},
			},
			{
				Lesson: CourseDraftLesson{Title: " 입문 ", OrderIndex: 1},
			},
		},
	}

	if err := svc.NormalizeDraftAggregate(&draft); err != nil {
		t.Fatalf("NormalizeDraftAggregate() error = %v", err)
	}

	if got := draft.Draft.SourceQuery; got != "통기타 코드 배우기" {
		t.Fatalf("trimmed source_query mismatch: %q", got)
	}
	if got := draft.Draft.Title; got != "첫 기타 코스" {
		t.Fatalf("trimmed title mismatch: %q", got)
	}
	if got := draft.Draft.Status; got != DraftStatusDraft {
		t.Fatalf("default status mismatch: %q", got)
	}
	// sorted by OrderIndex: 입문(1) first, 기초(5) second
	if got := draft.Lessons[0].Lesson.Title; got != "입문" {
		t.Fatalf("lessons not normalized in order: %q", got)
	}
	if got := draft.Lessons[0].Lesson.OrderIndex; got != 0 {
		t.Fatalf("first lesson order_index = %d, want 0", got)
	}
	if got := draft.Lessons[1].SubLessons[0].Lesson.OrderIndex; got != 0 {
		t.Fatalf("sub-lesson order_index = %d, want 0", got)
	}
	point := draft.Lessons[1].SubLessons[0].Points[0]
	if point.Point.OrderIndex != 0 {
		t.Fatalf("point order_index = %d, want 0", point.Point.OrderIndex)
	}
}

func TestBuildConfirmedCourseFromDraft(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	draftID := uuid.New()
	contentID := uuid.New()
	mainLessonID := uuid.New()
	subLessonID := uuid.New()

	draft := DraftAggregate{
		Draft: CourseDraft{
			ID:                 draftID,
			UserID:             userID,
			SourceQuery:        "수채화 기초",
			Title:              "수채화 시작 코스",
			CompletionCriteria: []string{"작은 수채화 작품 한 점을 완성할 수 있다."},
			Status:             DraftStatusDraft,
		},
		Lessons: []DraftLessonTree{
			{
				Lesson: CourseDraftLesson{
					ID:            mainLessonID,
					CourseDraftID: draftID,
					Title:         "입문",
					LessonRole:    LessonRoleCore,
					SourceType:    LessonSourceRecommended,
					OrderIndex:    0,
				},
				SubLessons: []DraftLessonTree{
					{
						Lesson: CourseDraftLesson{
							ID:             subLessonID,
							CourseDraftID:  draftID,
							ParentLessonID: &mainLessonID,
							Title:          "붓과 종이 이해",
							LessonRole:     LessonRoleCore,
							SourceType:     LessonSourceRecommended,
							OrderIndex:     0,
						},
						Points: []DraftPointAggregate{
							{
								Point: CourseDraftPoint{
									ID:                  uuid.New(),
									PointType:           PointTypeExploration,
									Title:               "재료 소개 영상",
									ContentID:           &contentID,
									CourseDraftLessonID: subLessonID,
									OrderIndex:          0,
								},
							},
						},
					},
				},
			},
		},
	}

	confirmed, err := svc.BuildConfirmedCourseFromDraft(draft)
	if err != nil {
		t.Fatalf("BuildConfirmedCourseFromDraft() error = %v", err)
	}

	if confirmed.Course.ID == uuid.Nil {
		t.Fatal("confirmed course id not generated")
	}
	if confirmed.Course.SourceDraftID == nil || *confirmed.Course.SourceDraftID != draftID {
		t.Fatal("source draft id not preserved")
	}
	if len(confirmed.Course.CompletionCriteria) != 1 || confirmed.Course.CompletionCriteria[0] != "작은 수채화 작품 한 점을 완성할 수 있다." {
		t.Fatalf("completion criteria not preserved: %#v", confirmed.Course.CompletionCriteria)
	}
	if len(confirmed.Lessons) != 1 {
		t.Fatalf("lessons len = %d, want 1", len(confirmed.Lessons))
	}
	mainTree := confirmed.Lessons[0]
	if mainTree.Lesson.CourseID != confirmed.Course.ID {
		t.Fatal("main lesson course_id mismatch")
	}
	if len(mainTree.SubLessons) != 1 {
		t.Fatalf("sub-lessons len = %d, want 1", len(mainTree.SubLessons))
	}
	subTree := mainTree.SubLessons[0]
	if subTree.Lesson.CourseID != confirmed.Course.ID {
		t.Fatal("sub-lesson course_id mismatch")
	}
	if len(subTree.Points) != 1 {
		t.Fatalf("points len = %d, want 1", len(subTree.Points))
	}
	point := subTree.Points[0]
	if point.Point.CourseLessonID != subTree.Lesson.ID {
		t.Fatal("point course_lesson_id mismatch")
	}
	if point.Point.ContentID == nil || *point.Point.ContentID != contentID {
		t.Fatal("point content_id not preserved")
	}
}

func TestValidateRecommendationEvent(t *testing.T) {
	svc := NewService()

	if err := svc.ValidateRecommendationEvent(RecommendationEvent{}); err == nil {
		t.Fatal("expected validation error for empty event")
	}

	err := svc.ValidateRecommendationEvent(RecommendationEvent{
		UserID:    uuid.New(),
		EventType: EventLessonSelected,
	})
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestBuildInitialDraft(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "기본 코드 진행을 이해하고 간단한 반주를 할 수 있다"
	currentLevel := "완전 초보"
	durationWeeks := 6
	studyHours := 3
	preferredFormat := "영상"

	draft, err := svc.BuildInitialDraft(userID, CreateCourseDraftRequest{
		SourceQuery:       "통기타",
		LearningGoal:      &goal,
		CurrentLevel:      &currentLevel,
		DurationWeeks:     &durationWeeks,
		StudyHoursPerWeek: &studyHours,
		PreferredFormat:   &preferredFormat,
	})
	if err != nil {
		t.Fatalf("BuildInitialDraft() error = %v", err)
	}

	if draft.Draft.UserID != userID {
		t.Fatal("user_id mismatch")
	}
	if draft.Draft.Title != "통기타 학습 코스" {
		t.Fatalf("draft title = %q", draft.Draft.Title)
	}
	if len(draft.Lessons) != 5 {
		t.Fatalf("lessons len = %d, want 5", len(draft.Lessons))
	}
	if draft.Lessons[0].Lesson.Title != "기초 셋팅과 첫 소리 내기" {
		t.Fatalf("first lesson title = %q", draft.Lessons[0].Lesson.Title)
	}
	if draft.Draft.LearningGoal == nil || *draft.Draft.LearningGoal != goal {
		t.Fatalf("learning goal mismatch: %#v", draft.Draft.LearningGoal)
	}
	if len(draft.Lessons[0].SubLessons) != 0 {
		t.Fatalf("first lesson sub-lessons len = %d, want 0", len(draft.Lessons[0].SubLessons))
	}
	if draft.Lessons[0].Lesson.SourceType != LessonSourceAI {
		t.Fatalf("lesson source type = %q", draft.Lessons[0].Lesson.SourceType)
	}
	if len(draft.Draft.CompletionCriteria) == 0 {
		t.Fatal("completion criteria should not be empty")
	}
}

func TestBuildInitialDraftUsesEnglishFallbackWhenLearningLanguageIsEnglish(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "Hold a short travel conversation at the airport"

	draft, err := svc.BuildInitialDraft(userID, CreateCourseDraftRequest{
		SourceQuery:      "travel English conversation",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	})
	if err != nil {
		t.Fatalf("BuildInitialDraft() error = %v", err)
	}

	if draft.Draft.GenerationLanguage != "en" {
		t.Fatalf("generation language = %q, want en", draft.Draft.GenerationLanguage)
	}
	if !strings.Contains(draft.Draft.Title, "Learning Course") {
		t.Fatalf("draft title = %q, want English course title", draft.Draft.Title)
	}
	if strings.Contains(draft.Lessons[0].Lesson.Title, "핵심") || strings.Contains(draft.Lessons[0].Lesson.Title, "정리") {
		t.Fatalf("first fallback lesson should be English, got %q", draft.Lessons[0].Lesson.Title)
	}
	if !strings.Contains(draft.Lessons[0].Lesson.Title, "Organize") {
		t.Fatalf("first fallback lesson title = %q, want English travel-conversation template", draft.Lessons[0].Lesson.Title)
	}
	if len(draft.Draft.CompletionCriteria) == 0 || strings.Contains(draft.Draft.CompletionCriteria[0], "목표") {
		t.Fatalf("completion criteria should be English, got %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMDoesNotOverwriteEnglishLessonsWithKoreanPatternRefinement(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "Hold a short travel conversation at the airport"
	client := &stubLLMClient{responses: []string{`{
		"title":"Airport Conversation Course",
		"description":"A staged course for short travel conversations.",
		"main_lessons":[
			{"title":"Map airport survival phrases","objective":"Identify the first phrases needed at check-in, security, and boarding."},
			{"title":"Build short question and answer patterns","objective":"Practice concise questions and answers for airport situations."},
			{"title":"Connect phrases in a travel dialogue","objective":"Link the phrases into a realistic airport conversation."},
			{"title":"Run a final conversation rehearsal","objective":"Rehearse the full flow before using it in the airport."}
		],
		"completion_criteria":["Handle a short airport conversation using prepared phrases."]
	}`}}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:      "영어 travel conversation",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if draft.Lessons[0].Lesson.Title != "Map airport survival phrases" {
		t.Fatalf("English LLM lesson was overwritten by refinement, got %q", draft.Lessons[0].Lesson.Title)
	}
	for _, lesson := range draft.Lessons {
		if containsHangul(lesson.Lesson.Title) || containsHangul(derefString(lesson.Lesson.Objective)) {
			t.Fatalf("English lesson should not contain Korean refinement text: %#v", lesson.Lesson)
		}
	}
	last := draft.Lessons[len(draft.Lessons)-1].Lesson
	if last.Title != "Prepare to sustain the target conversation" {
		t.Fatalf("English pre-milestone title = %q", last.Title)
	}
	if !strings.Contains(derefString(last.Objective), "Connect the preparation") {
		t.Fatalf("English pre-milestone objective = %q", derefString(last.Objective))
	}
}

func TestBuildPreMilestoneTitleUsesSpecificEnglishFallbacks(t *testing.T) {
	tests := []struct {
		name        string
		sourceQuery string
		goal        string
		want        string
	}{
		{
			name:        "woodworking birdhouse",
			sourceQuery: "beginner woodworking birdhouse",
			goal:        "I want to build a simple wooden birdhouse with safe measuring, cutting, sanding, and assembly",
			want:        "Inspect and finish the assembled birdhouse",
		},
		{
			name:        "watercolor postcard",
			sourceQuery: "beginner watercolor landscape postcard",
			goal:        "I want to paint a small watercolor landscape postcard with simple washes, trees, and sky depth",
			want:        "Refine the watercolor postcard before completion",
		},
		{
			name:        "ukulele first song",
			sourceQuery: "ukulele first song strumming",
			goal:        "I want to learn ukulele chords and strumming so I can play one simple song from start to finish",
			want:        "Play the first song from start to finish",
		},
		{
			name:        "harmonica blues riff",
			sourceQuery: "harmonica blues riff beginner",
			goal:        "I want to play a short blues riff on harmonica with clean single notes and simple bends",
			want:        "Play the blues riff with clean notes and bends",
		},
		{
			name:        "bass groove",
			sourceQuery: "bass groove foundation",
			goal:        "I want to play steady bass grooves with basic rhythm and clean note transitions",
			want:        "Play the target bass groove steadily",
		},
		{
			name:        "crochet granny square",
			sourceQuery: "crochet granny square beginner",
			goal:        "I want to crochet a neat granny square and join several squares into a small coaster or pouch",
			want:        "Join and finish the crochet square project",
		},
		{
			name:        "pottery glazing bowl",
			sourceQuery: "pottery glazing small bowl",
			goal:        "I want to handbuild a small bowl and test simple glazing choices before finishing it",
			want:        "Review the glazed bowl before firing",
		},
		{
			name:        "short video publishing",
			sourceQuery: "shorts video editing workflow",
			goal:        "I want to edit a short vertical video with cuts, captions, music timing, and a publish-ready thumbnail",
			want:        "Prepare the short video for publishing",
		},
		{
			name:        "podcast publishing",
			sourceQuery: "beginner podcast first episode",
			goal:        "I want to plan, record, edit, and publish a short first podcast episode with clear audio",
			want:        "Prepare the first podcast episode for publishing",
		},
		{
			name:        "home workout routine",
			sourceQuery: "home workout beginner routine",
			goal:        "I want to build a simple home workout routine and keep a safe weekly exercise habit",
			want:        "Review the weekly workout routine safely",
		},
		{
			name:        "sourdough starter bread",
			sourceQuery: "beginner sourdough starter bread",
			goal:        "I want to create and maintain a sourdough starter and bake one simple loaf with clear fermentation steps",
			want:        "Bake and review the first sourdough loaf",
		},
		{
			name:        "meal prep healthy lunch",
			sourceQuery: "healthy meal prep beginner lunch",
			goal:        "I want to plan, cook, and store simple healthy lunches for a weekday meal prep routine",
			want:        "Review the weekday meal prep routine",
		},
		{
			name:        "indoor herb garden",
			sourceQuery: "indoor herb garden beginner",
			goal:        "I want to grow basil and mint indoors, manage watering and light, and harvest herbs for simple cooking",
			want:        "Review the herb care and harvest routine",
		},
		{
			name:        "houseplant repotting",
			sourceQuery: "houseplant repotting soil mix",
			goal:        "I want to repot a small houseplant safely, choose a basic soil mix, and monitor recovery after repotting",
			want:        "Monitor the repotted plant recovery",
		},
		{
			name:        "chess opening practice",
			sourceQuery: "beginner chess opening practice",
			goal:        "I want to learn one beginner chess opening, understand the first few moves, and practice simple middlegame plans",
			want:        "Play and review the opening practice game",
		},
		{
			name:        "personal finance budget",
			sourceQuery: "personal finance monthly budget beginner",
			goal:        "I want to build a simple monthly budget, track spending categories, and review one week of expenses",
			want:        "Review the first weekly budget check",
		},
		{
			name:        "mindfulness breathing routine",
			sourceQuery: "mindfulness breathing routine beginner",
			goal:        "I want to practice a short mindfulness breathing routine every day and record how it affects my focus",
			want:        "Review the daily breathing practice log",
		},
		{
			name:        "dog basic training",
			sourceQuery: "dog basic training sit stay recall",
			goal:        "I want to train my dog to sit, stay, and come using short positive reinforcement sessions",
			want:        "Practice the commands in a short real session",
		},
		{
			name:        "notion study dashboard",
			sourceQuery: "Notion study dashboard beginner",
			goal:        "I want to build a simple Notion study dashboard with tasks, notes, and weekly review pages",
			want:        "Use the study dashboard for a weekly review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPreMilestoneTitle(tt.sourceQuery, tt.goal, "en")
			if got != tt.want {
				t.Fatalf("buildPreMilestoneTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildDraftWithLLMReplacesKoreanTitleDescriptionInEnglishMode(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "I want to improve my spoken answers for an English speaking interview test"
	client := &stubLLMClient{responses: []string{`{
		"title":"영어 면접 답변 향상",
		"description":"이 커리큘럼은 영어 면접 시험의 구술 답변을 향상시키기 위한 단계별 접근 방식을 제공합니다.",
		"main_lessons":[
			{"title":"Craft clear spoken answers","objective":"Build concise answers for common interview prompts."},
			{"title":"Practice mock interviews","objective":"Use timed practice to improve fluency and confidence."},
			{"title":"Review feedback and refine responses","objective":"Adjust response structure based on feedback."}
		],
		"completion_criteria":["Complete a final mock speaking interview."]
	}`}}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:      "English speaking interview test",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if containsHangul(draft.Draft.Title) || containsHangul(derefString(draft.Draft.Description)) {
		t.Fatalf("English draft title/description should not contain Korean: title=%q description=%q", draft.Draft.Title, derefString(draft.Draft.Description))
	}
	if !strings.Contains(draft.Draft.Title, "English speaking interview test") {
		t.Fatalf("draft title = %q, want English source fallback", draft.Draft.Title)
	}
}

func TestBuildInitialDraftRequiresLearningGoal(t *testing.T) {
	svc := NewService()
	_, err := svc.BuildInitialDraft(uuid.New(), CreateCourseDraftRequest{
		SourceQuery: "통기타",
	})
	if err == nil || err.Error() != "learning_goal is required" {
		t.Fatalf("expected learning_goal is required, got %v", err)
	}
}

func TestBuildDraftWithLLMUsesReviewedLessons(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "찬양단에서 기타를 연주하기"
	currentLevel := "완전 초보"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"기타 입문 학습 코스",
				"description":"기타 입문 초안",
				"main_lessons":[
					{"title":"기타 기본 잡기","objective":"기타를 안정적으로 잡는다."},
					{"title":"기초 코드 배우기","objective":"기본 코드를 익힌다."},
					{"title":"찬양단에서 실제 연주하기","objective":"찬양단에서 실제 연주한다."}
				],
				"completion_criteria":["기본 코드로 간단한 반주를 할 수 있다."]
			}`,
			`{
				"strengths":["목표와 관련은 있다."],
				"issues":["최종 목표가 lesson에 들어가 있다."],
				"missing_axes":["코드 전환","스트로크","곡 적용"],
				"boundary_problems":["실전 참여가 lesson으로 들어감"],
				"revision_instructions":["실전 참여는 완료 기준으로 분리"],
				"revised_main_lessons":[
					{"title":"기타 기본 자세와 장비 익히기","objective":"기타를 안정적으로 잡고 기본 장비를 다룬다."},
					{"title":"찬양 반주 기본 코드 익히기","objective":"찬양곡에 자주 쓰는 코드를 익힌다."},
					{"title":"코드 전환과 스트로크 안정화하기","objective":"코드 전환과 기본 스트로크를 끊기지 않게 연결한다."},
					{"title":"쉬운 찬양곡 한 곡 완주하기","objective":"쉬운 찬양곡 한 곡을 끝까지 반주한다."},
					{"title":"찬양단 합주 흐름에 맞춰 준비하기","objective":"합주 흐름과 시작, 마침 타이밍에 맞춘다."}
				],
				"revised_completion_criteria":[
					"찬양곡 1곡을 기본 코드와 스트로크로 완주할 수 있다.",
					"찬양단 연습 흐름에 맞춰 반주에 참여할 준비가 되어 있다."
				]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "기타 입문",
		LearningGoal: &goal,
		CurrentLevel: &currentLevel,
	}, DraftBuildOptions{})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if client.callCount != 2 {
		t.Fatalf("llm call count = %d, want 2", client.callCount)
	}
	if len(draft.Lessons) != 4 {
		t.Fatalf("lessons len = %d, want 4", len(draft.Lessons))
	}
	if draft.Lessons[0].Lesson.Title != "첫 소리를 안정적으로 낸다" {
		t.Fatalf("first reviewed lesson title = %q", draft.Lessons[0].Lesson.Title)
	}
	if len(draft.Draft.CompletionCriteria) < 2 {
		t.Fatalf("completion criteria len = %d, want at least 2", len(draft.Draft.CompletionCriteria))
	}
}

func TestBuildDraftWithLLMFallsBackWhenReviewFails(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "여행에서 간단히 영어로 대화하기"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영어 회화 입문",
				"description":"여행 영어 초안",
				"main_lessons":[
					{"title":"핵심 표현과 상황 정리하기","objective":"여행 상황에 필요한 표현을 정리한다."},
					{"title":"기본 문장 패턴 익히기","objective":"짧은 기본 문장을 익힌다."},
					{"title":"상황별 대화 연습하기","objective":"상황별 대화를 연습한다."}
				],
				"completion_criteria":["여행 상황에서 기본 질문과 응답을 이어갈 수 있다."]
			}`,
			`{invalid json`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영어 회화 입문",
		LearningGoal: &goal,
	}, DraftBuildOptions{})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if len(draft.Lessons) != 4 {
		t.Fatalf("lessons len = %d, want 4", len(draft.Lessons))
	}
	if draft.Lessons[0].Lesson.Title != "핵심 표현으로 짧은 여행 대화를 시작한다" {
		t.Fatalf("fallback refine failed, got %q", draft.Lessons[0].Lesson.Title)
	}
	if len(draft.Draft.CompletionCriteria) < 1 {
		t.Fatalf("completion criteria len = %d, want at least 1", len(draft.Draft.CompletionCriteria))
	}
}

func TestBuildDraftWithLLMSkipsReviewWhenOptionEnabled(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "작은 수채화 엽서 한 장 완성하기"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"수채화 입문",
				"description":"수채화 초안",
				"main_lessons":[
					{"title":"도구와 물 조절 익히기","objective":"기본 도구를 다루고 물 농도를 맞춘다."},
					{"title":"기본 번짐과 색 혼합 연습하기","objective":"번짐과 색 혼합을 제어한다."},
					{"title":"작은 엽서 구도 잡기","objective":"간단한 엽서 구도를 정한다."}
				],
				"completion_criteria":["작은 수채화 엽서 한 장을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "수채화 입문",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if client.callCount != 1 {
		t.Fatalf("llm call count = %d, want 1", client.callCount)
	}
	if draft.Lessons[0].Lesson.Title != "기본 표현으로 작은 그림을 완성한다" {
		t.Fatalf("expected refined first lesson, got %q", draft.Lessons[0].Lesson.Title)
	}
}

func TestBuildDraftWithLLMReturnsRetryableErrorOnInvalidGenerationJSON(t *testing.T) {
	svc := NewService()
	client := &stubLLMClient{
		responses: []string{
			`{invalid json`,
		},
	}
	goal := "기타 반주를 배우고 싶다"

	_, err := svc.BuildDraftWithLLM(context.Background(), client, uuid.New(), CreateCourseDraftRequest{
		SourceQuery:  "통기타 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err == nil {
		t.Fatal("expected retryable generation error")
	}
	if !errors.Is(err, errDraftGenerationRetryable) {
		t.Fatalf("expected retryable error, got %v", err)
	}
}

func TestBuildDraftWithLLMReturnsRetryableErrorOnEmptyLessons(t *testing.T) {
	svc := NewService()
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"기타 입문",
				"description":"초안",
				"main_lessons":[],
				"completion_criteria":["기본 반주를 시작할 수 있다."]
			}`,
		},
	}
	goal := "기타 반주를 배우고 싶다"

	_, err := svc.BuildDraftWithLLM(context.Background(), client, uuid.New(), CreateCourseDraftRequest{
		SourceQuery:  "통기타 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err == nil {
		t.Fatal("expected retryable generation error")
	}
	if !errors.Is(err, errDraftGenerationRetryable) {
		t.Fatalf("expected retryable error, got %v", err)
	}
}

func TestRefineGeneratedDraftDocumentMergesAtomicLessonsAndMovesMilestones(t *testing.T) {
	goal := "찬양단에서 기타 반주를 자연스럽게 이어가기"
	doc := refineGeneratedDraftDocument(CreateCourseDraftRequest{
		SourceQuery:  "통기타",
		LearningGoal: &goal,
	}, generatedDraftDocument{
		Title:       "기타 코스",
		Description: "초안",
		MainLessons: []generatedMainLesson{
			{Title: "기본자세 익히기", Objective: "기타를 안정적으로 잡고 튜닝한다."},
			{Title: "기초 코드 배우기", Objective: "기본 코드를 익힌다."},
			{Title: "코드 전환 연습하기", Objective: "코드 전환을 끊기지 않게 이어간다."},
			{Title: "찬양단에서 실제 연주하기", Objective: "찬양단에서 실제 연주한다."},
		},
		CompletionCriteria: []string{"기본 코드로 간단히 반주할 수 있다."},
	})

	if len(doc.MainLessons) != 4 {
		t.Fatalf("main lessons len = %d, want 4", len(doc.MainLessons))
	}
	for _, lesson := range doc.MainLessons {
		if strings.Contains(lesson.Title, "기본자세") {
			t.Fatalf("atomic lesson should be merged, got %q", lesson.Title)
		}
		if strings.Contains(lesson.Title, "실제 연주") {
			t.Fatalf("milestone lesson should be moved out of lessons, got %q", lesson.Title)
		}
	}
	if len(doc.CompletionCriteria) < 2 {
		t.Fatalf("completion criteria len = %d, want at least 2", len(doc.CompletionCriteria))
	}
	foundMilestone := false
	for _, item := range doc.CompletionCriteria {
		if strings.Contains(item, "실제 연주") {
			foundMilestone = true
			break
		}
	}
	if !foundMilestone {
		t.Fatalf("milestone objective not moved to completion criteria: %#v", doc.CompletionCriteria)
	}
}

func TestIsMilestoneLikeLessonKeepsPreparationStage(t *testing.T) {
	lesson := generatedMainLesson{
		Title:     "모임 흐름에 맞춰 참여를 준비한다",
		Objective: "소셜댄스 모임의 시작, 연결, 마무리 흐름에 맞춰 자연스럽게 들어갈 준비를 한다.",
	}

	if isMilestoneLikeLesson(lesson) {
		t.Fatalf("preparation lesson should not be treated as milestone: %#v", lesson)
	}
}

func TestIsMilestoneLikeLessonOnlyTreatsExplicitEventAsMilestone(t *testing.T) {
	lesson := generatedMainLesson{
		Title:     "졸업식에서 친구들 앞에서 발표하기",
		Objective: "졸업식에서 실제로 무대에 올라 친구들 앞에서 춤을 발표한다.",
	}

	if !isMilestoneLikeLesson(lesson) {
		t.Fatalf("explicit completion event should be treated as milestone: %#v", lesson)
	}
}

func TestIsAtomicLessonDoesNotCollapseBroadFirstStage(t *testing.T) {
	lesson := generatedMainLesson{
		Title:     "재료와 도구 흐름 익히기",
		Objective: "도자기 재료와 기본 도구를 다루며 작은 머그컵 형태를 바로 만들기 시작한다.",
	}

	if isAtomicLesson(lesson) {
		t.Fatalf("broad first stage should not be treated as atomic: %#v", lesson)
	}
}

func TestIsAtomicLessonKeepsMicroIntroOnly(t *testing.T) {
	lesson := generatedMainLesson{
		Title:     "재료 소개",
		Objective: "기본 재료의 차이를 이해한다.",
	}

	if !isAtomicLesson(lesson) {
		t.Fatalf("micro intro should still be treated as atomic: %#v", lesson)
	}
}

func TestIsMilestoneLikeLessonWithRealUserDancePhrasing(t *testing.T) {
	prep := generatedMainLesson{
		Title:     "졸업식 무대 흐름에 맞춰 참여를 준비한다",
		Objective: "친구들 앞에서 춤추기 전 동선과 시작 타이밍을 정리한다.",
	}
	event := generatedMainLesson{
		Title:     "졸업식에서 친구들 앞에서 발표하기",
		Objective: "졸업식에서 실제로 춤을 발표한다.",
	}

	if isMilestoneLikeLesson(prep) {
		t.Fatalf("real-user prep phrasing should not be treated as milestone: %#v", prep)
	}
	if !isMilestoneLikeLesson(event) {
		t.Fatalf("real-user completion event should be treated as milestone: %#v", event)
	}
}

func TestIsAtomicLessonWithRealUserUrbanSketchPhrasing(t *testing.T) {
	broadFirstStage := generatedMainLesson{
		Title:     "재료와 도구 흐름 익히기",
		Objective: "여행에서 쓸 재료와 도구를 다루며 작은 장면을 바로 그리기 시작한다.",
	}
	microStage := generatedMainLesson{
		Title:     "도구 소개",
		Objective: "펜과 종이의 차이를 이해한다.",
	}

	if isAtomicLesson(broadFirstStage) {
		t.Fatalf("real-user broad first stage should not be treated as atomic: %#v", broadFirstStage)
	}
	if !isAtomicLesson(microStage) {
		t.Fatalf("real-user micro intro should still be treated as atomic: %#v", microStage)
	}
}

func TestBuildDraftWithLLMSkipsReviewAndRefinesAtomicLesson(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "찬양단에서 기타 반주를 자연스럽게 이어가기"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"기타 입문",
				"description":"기타 초안",
				"main_lessons":[
					{"title":"기본자세 익히기","objective":"기타를 안정적으로 잡고 튜닝한다."},
					{"title":"기초 코드 배우기","objective":"기본 코드를 익힌다."},
					{"title":"코드 전환 연습하기","objective":"코드 전환을 끊기지 않게 이어간다."},
					{"title":"찬양단에서 실제 연주하기","objective":"찬양단에서 실제 연주한다."}
				],
				"completion_criteria":["기본 코드로 간단히 반주할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "통기타",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if len(draft.Lessons) != 4 {
		t.Fatalf("lessons len = %d, want 4", len(draft.Lessons))
	}
	for _, lesson := range draft.Lessons {
		if strings.Contains(lesson.Lesson.Title, "기본자세") {
			t.Fatalf("atomic lesson should be merged after refine, got %q", lesson.Lesson.Title)
		}
		if strings.Contains(lesson.Lesson.Title, "실제 연주") {
			t.Fatalf("milestone lesson should not remain in lessons, got %q", lesson.Lesson.Title)
		}
	}
	if len(draft.Draft.CompletionCriteria) < 2 {
		t.Fatalf("completion criteria len = %d, want at least 2", len(draft.Draft.CompletionCriteria))
	}
}

func TestDedupeObjectiveSentences(t *testing.T) {
	got := dedupeObjectiveSentences("기타를 안정적으로 잡고 튜닝한다. 이 과정을 바탕으로 첫 연주 수행을 시작할 수 있게 한다. 이 과정을 바탕으로 첫 연주 수행을 시작할 수 있게 한다.")
	want := "기타를 안정적으로 잡고 튜닝한다. 이 과정을 바탕으로 첫 연주 수행을 시작할 수 있게 한다."
	if got != want {
		t.Fatalf("dedupeObjectiveSentences() = %q, want %q", got, want)
	}
}

func TestBuildDraftWithLLMRefinesWorshipGuitarProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "찬양단 봉사"
	currentLevel := "완전 초보"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"통기타 찬양 반주 코스",
				"description":"찬양단 봉사를 위한 통기타 입문 초안",
				"main_lessons":[
					{"title":"기본자세 익히기","objective":"기타를 안정적으로 잡고 튜닝한다."},
					{"title":"찬양 반주 오픈 코드 익히기","objective":"찬양곡에서 자주 쓰는 오픈 코드를 익힌다."},
					{"title":"코드 전환과 스트로크 연습하기","objective":"코드 전환과 기본 스트로크를 끊기지 않게 이어간다."},
					{"title":"쉬운 찬양곡 한 곡 완주하기","objective":"쉬운 찬양곡 한 곡을 끝까지 반주한다."},
					{"title":"찬양단에서 실제 연주하기","objective":"찬양단에서 실제로 반주한다."}
				],
				"completion_criteria":["기본 코드로 찬양 반주를 이어갈 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "통기타 배우기",
		LearningGoal: &goal,
		CurrentLevel: &currentLevel,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"첫 소리를 안정적으로 낸다",
		"찬양 반주에 필요한 코드를 끊기지 않고 전환한다",
		"기본 스트로크로 찬양곡 한 곡을 끝까지 이어간다",
		"찬양단 합주 흐름에 맞춰 반주를 준비한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	found := false
	for _, item := range draft.Draft.CompletionCriteria {
		if item == "찬양단 합주 흐름에 맞춰 반주를 준비하고 연결할 수 있다." {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("completion criteria missing worship-team prep criterion: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesCrochetTeachingProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "나만의 작품을 만들고 유튜브를 통해 강의하고 싶다"
	currentLevel := "완전 초보"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"코바늘 입문 코스",
				"description":"작품 제작과 유튜브 강의를 목표로 하는 코바늘 입문 초안",
				"main_lessons":[
					{"title":"코바늘 도구 이해하기","objective":"코바늘과 실, 기본 준비물을 이해한다."},
					{"title":"기초 뜨개법 익히기","objective":"사슬뜨기와 짧은뜨기 같은 기초 뜨개법을 익힌다."},
					{"title":"작은 소품 만들어보기","objective":"간단한 소품 하나를 완성한다."},
					{"title":"나만의 작품 구상하기","objective":"내 스타일의 작품을 구상하고 도안을 정리한다."},
					{"title":"유튜브 강의 시작하기","objective":"작품 만드는 과정을 촬영하고 유튜브 강의를 시작한다."}
				],
				"completion_criteria":["작은 코바늘 소품을 스스로 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "코바늘 배우기",
		LearningGoal: &goal,
		CurrentLevel: &currentLevel,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기초 뜨개법으로 작은 결과물을 완성한다",
		"반복 패턴을 안정적으로 익혀 작품 형태를 만든다",
		"나만의 작품을 구상하고 도안 흐름을 정리한다",
		"직접 만든 작품 한 점을 완성한다",
		"작품 제작 과정을 설명하고 시연할 준비를 한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	found := false
	for _, item := range draft.Draft.CompletionCriteria {
		if item == "직접 만든 작품의 제작 과정을 설명하고 시연할 준비가 되어 있다." {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("completion criteria missing crochet-teaching prep criterion: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesConversationPerformanceProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "여행에서 간단히 영어로 대화하기"
	currentLevel := "완전 초보"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영어 회화 입문",
				"description":"여행 영어 초안",
				"main_lessons":[
					{"title":"핵심 표현과 상황 정리하기","objective":"여행 상황에 필요한 표현을 정리한다."},
					{"title":"기본 문장 패턴 익히기","objective":"짧은 기본 문장을 익힌다."},
					{"title":"상황별 대화 연습하기","objective":"상황별 대화를 연습한다."},
					{"title":"실전 대화 시작하기","objective":"여행에서 실제로 대화를 시작한다."}
				],
				"completion_criteria":["여행 상황에서 기본 질문과 응답을 이어갈 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영어 회화 입문",
		LearningGoal: &goal,
		CurrentLevel: &currentLevel,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"핵심 표현으로 짧은 여행 대화를 시작한다",
		"기본 문장 패턴으로 질문과 응답을 이어간다",
		"여행 상황에 맞는 대화를 자연스럽게 연결한다",
		"여행 직전까지 필요한 표현 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	found := false
	for _, item := range draft.Draft.CompletionCriteria {
		if item == "여행 상황에서 짧은 질문과 응답을 이어가며 실제 대화를 시작할 수 있다." {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("completion criteria missing conversation-performance criterion: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesWritingPresentationProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "브런치에 글을 연재하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"글쓰기 코스",
				"description":"브런치 연재 준비 초안",
				"main_lessons":[
					{"title":"글쓰기 기본 다지기","objective":"문장을 다듬는 법을 익힌다."},
					{"title":"소재 고르기","objective":"무슨 글을 쓸지 정한다."},
					{"title":"브런치에 바로 올리기","objective":"브런치에 글을 올린다."}
				],
				"completion_criteria":["한 편의 글을 공개할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "에세이 글쓰기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"짧은 글 한 편을 끝까지 쓴다",
		"전달하려는 주제와 독자 흐름을 잡는다",
		"공개 가능한 글 묶음의 구조를 정리한다",
		"연재 또는 공개 직전까지 글을 다듬고 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesVisualArtArtifactProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "수채화 작품 한 점을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"수채화 코스",
				"description":"작품 완성 초안",
				"main_lessons":[
					{"title":"붓과 종이 이해하기","objective":"도구를 이해한다."},
					{"title":"기초 워시 익히기","objective":"기초 워시를 연습한다."},
					{"title":"작품 전시하기","objective":"완성한 그림을 전시한다."}
				],
				"completion_criteria":["작품 한 점을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "수채화 입문",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 표현으로 작은 그림을 완성한다",
		"색과 형태 표현을 안정적으로 이어간다",
		"표현 요소를 통합해 한 장면을 구성한다",
		"완성도 있는 작품 한 점을 마무리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesMakerProfessionalProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "3D 프린팅 결과물을 판매하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"3D 프린팅 코스",
				"description":"판매 준비 초안",
				"main_lessons":[
					{"title":"프린터 구조 이해하기","objective":"장비를 이해한다."},
					{"title":"첫 결과물 만들어보기","objective":"작은 출력물을 만든다."},
					{"title":"온라인에서 실제 판매하기","objective":"온라인에서 실제 판매한다."}
				],
				"completion_criteria":["작은 출력물을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "3D 프린팅 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"만들고 싶은 결과물과 판매 방향을 정한다",
		"상품이 되는 제작 흐름을 설계한다",
		"작동하는 결과물과 샘플을 완성한다",
		"판매 또는 의뢰 직전까지 소개와 채널을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesVisualArtHabitProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "매일 드로잉하는 습관을 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"드로잉 루틴 코스",
				"description":"매일 그리기 초안",
				"main_lessons":[
					{"title":"연필과 도구 알아보기","objective":"도구를 이해한다."},
					{"title":"선 긋기 연습","objective":"선을 연습한다."},
					{"title":"매일 올리기","objective":"그림을 매일 올린다."}
				],
				"completion_criteria":["매일 10분씩 그릴 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "드로잉 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"빈 화면 부담 없이 바로 그리기 시작한다",
		"짧은 스케치를 반복해 손을 풀고 익숙해진다",
		"일상에서 자주 그리는 루틴을 만든다",
		"꾸준히 이어갈 수 있는 개인 드로잉 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesDigitalCreationArtifactProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "영상 한 편을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영상 편집 코스",
				"description":"영상 완성 초안",
				"main_lessons":[
					{"title":"프리미어 인터페이스 이해하기","objective":"도구를 이해한다."},
					{"title":"컷 편집 익히기","objective":"컷 편집을 연습한다."},
					{"title":"실제 업로드하기","objective":"유튜브에 업로드한다."}
				],
				"completion_criteria":["짧은 영상 한 편을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"도구에 익숙해지며 첫 결과물을 시작한다",
		"핵심 제작 패턴을 반복해 익힌다",
		"표현 요소를 통합해 완성 흐름을 만든다",
		"완성도 있는 디지털 결과물 한 편을 마무리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesDigitalCreationProfessionalProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "영상편집으로 프리랜서 의뢰를 받고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영상 편집 프리랜서 코스",
				"description":"프리랜서 준비 초안",
				"main_lessons":[
					{"title":"프리미어 기본 익히기","objective":"기본 도구를 익힌다."},
					{"title":"짧은 영상 만들어보기","objective":"짧은 결과물을 만든다."},
					{"title":"실제 의뢰 받기","objective":"실제 의뢰를 받는다."}
				],
				"completion_criteria":["짧은 영상 샘플을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"판매하거나 제안할 결과물 방향을 정한다",
		"결과물 구조와 작업 흐름을 설계한다",
		"소개 가능한 결과물과 샘플을 완성한다",
		"판매 또는 의뢰 직전까지 채널과 소개 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesVisualArtProfessionalProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "그림을 판매하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"그림 판매 코스",
				"description":"판매 준비 초안",
				"main_lessons":[
					{"title":"재료 알아보기","objective":"재료를 이해한다."},
					{"title":"작품 한 점 그리기","objective":"작품을 완성한다."},
					{"title":"실제 판매하기","objective":"그림을 판매한다."}
				],
				"completion_criteria":["작품 한 점을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "수채화 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"판매하거나 소개할 작품 방향을 정한다",
		"작품을 상품과 포트폴리오 구조로 정리한다",
		"소개 가능한 작품과 샘플을 완성한다",
		"판매 또는 의뢰 직전까지 채널과 소개 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCookingArtifactProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "케이크 한 점을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"베이킹 코스",
				"description":"결과물 완성 초안",
				"main_lessons":[
					{"title":"재료 이해하기","objective":"재료를 이해한다."},
					{"title":"반죽 만들기","objective":"반죽을 만든다."},
					{"title":"작품 전시하기","objective":"완성한 케이크를 전시한다."}
				],
				"completion_criteria":["케이크를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "베이킹 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 재료와 도구로 첫 결과물을 시작한다",
		"핵심 조리와 제과 패턴을 안정적으로 익힌다",
		"표현 요소를 통합해 완성 흐름을 만든다",
		"완성도 있는 대표 결과물 한 점을 마무리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCookingProfessionalProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "베이킹으로 주문을 받고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"베이킹 판매 코스",
				"description":"주문 준비 초안",
				"main_lessons":[
					{"title":"재료 이해하기","objective":"재료를 이해한다."},
					{"title":"케이크 만들기","objective":"케이크를 만든다."},
					{"title":"실제 주문받기","objective":"실제로 주문을 받는다."}
				],
				"completion_criteria":["케이크를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "베이킹 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"판매하거나 주문받을 메뉴 방향을 정한다",
		"메뉴 구성과 작업 흐름을 설계한다",
		"소개 가능한 대표 메뉴와 샘플을 완성한다",
		"판매 또는 주문 직전까지 채널과 소개 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesDigitalCreationPublishProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "유튜브에 영상 한 편을 공개하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"유튜브 편집 코스",
				"description":"공개 준비 초안",
				"main_lessons":[
					{"title":"프리미어 기본 익히기","objective":"도구를 익힌다."},
					{"title":"컷 편집 하기","objective":"컷 편집을 한다."},
					{"title":"유튜브에 바로 올리기","objective":"유튜브에 바로 올린다."}
				],
				"completion_criteria":["영상 한 편을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"공개할 결과물의 방향과 포맷을 정한다",
		"핵심 제작 패턴을 반복해 공개 품질을 만든다",
		"표현 요소를 통합해 완성 흐름을 만든다",
		"공개 직전까지 결과물과 게시 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesUniversitySurveySubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "대학 강의처럼 현대미술 흐름을 넓게 이해하고 비교해서 정리하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"현대미술사 학습 코스",
				"description":"survey 초안",
				"main_lessons":[
					{"title":"작가 이름 익히기","objective":"대표 작가를 외운다."},
					{"title":"사조 이해하기","objective":"사조를 이해한다."},
					{"title":"레포트 제출하기","objective":"최종 레포트를 제출한다."}
				],
				"completion_criteria":["현대미술의 큰 흐름을 설명할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "현대미술사",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"과목 범위와 문제의식을 큰 흐름으로 잡는다",
		"핵심 텍스트와 사례군을 따라 읽고 비교한다",
		"주제와 사례를 엮어 해석 관점을 만든다",
		"비교·비평 관점으로 전체 흐름을 통합 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesUniversityStudioSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "희곡 워크숍처럼 초안부터 써 보고 수정해 제출 가능한 형태까지 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"희곡 쓰기 코스",
				"description":"studio 초안",
				"main_lessons":[
					{"title":"희곡 이론 이해하기","objective":"희곡 이론을 이해한다."},
					{"title":"장면 하나 써보기","objective":"장면 하나를 써본다."},
					{"title":"최종 제출하기","objective":"작품을 제출한다."}
				],
				"completion_criteria":["짧은 희곡 초안을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "희곡 쓰기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"작은 초안이나 첫 결과물을 빠르게 만든다",
		"핵심 기법과 제작 패턴을 반복해 익힌다",
		"피드백을 반영해 구조와 표현을 다듬는다",
		"목표 작업물을 완성 가능한 형태로 발전시킨다",
		"발표나 제출 직전까지 작업물을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesWritingTeachingProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "글쓰기를 가르치는 워크숍을 열고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"글쓰기 강의 코스",
				"description":"워크숍 준비 초안",
				"main_lessons":[
					{"title":"좋은 문장 알아보기","objective":"좋은 문장을 이해한다."},
					{"title":"예시 모으기","objective":"예시를 모은다."},
					{"title":"워크숍 바로 열기","objective":"워크숍을 바로 연다."}
				],
				"completion_criteria":["짧은 글을 쓸 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "글쓰기 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"짧은 글로 자신의 글쓰기 흐름을 정리한다",
		"글쓰기 과정을 단계별로 나눠 설명한다",
		"예시와 비교를 통해 이해 흐름을 만든다",
		"강의나 워크숍 직전까지 설명 구조를 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCraftArtifactProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "코바늘 작품 한 점을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"코바늘 작품 코스",
				"description":"작품 완성 초안",
				"main_lessons":[
					{"title":"도구 이해하기","objective":"도구를 이해한다."},
					{"title":"기본 포인트 익히기","objective":"기본 포인트를 익힌다."},
					{"title":"작품 전시하기","objective":"작품을 전시한다."}
				],
				"completion_criteria":["작품 하나를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "코바늘 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기초 기법으로 작은 결과물을 완성한다",
		"반복 패턴과 형태 만들기를 안정적으로 익힌다",
		"작품 구조를 구상하고 설계 흐름을 정리한다",
		"직접 만든 작품 한 점을 완성한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesBodyMovementParticipationProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "소셜댄스 모임에 참여하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"댄스 참여 코스",
				"description":"참여 준비 초안",
				"main_lessons":[
					{"title":"기본 스텝 배우기","objective":"기본 스텝을 익힌다."},
					{"title":"턴 익히기","objective":"턴을 익힌다."},
					{"title":"모임에서 바로 참여하기","objective":"모임에서 바로 참여한다."}
				],
				"completion_criteria":["기본 스텝을 할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "살사 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 리듬과 흐름으로 첫 참여를 시작한다",
		"핵심 동작 패턴을 끊기지 않게 이어간다",
		"상대나 팀의 흐름과 연결해 움직인다",
		"실제 참여 직전까지 흐름과 타이밍을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesKnowledgeHabitProgression(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "별을 관찰하고 기록하는 취미를 꾸준히 하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"천문 취미 코스",
				"description":"루틴 초안",
				"main_lessons":[
					{"title":"별자리 알아보기","objective":"별자리를 알아본다."},
					{"title":"망원경 알아보기","objective":"도구를 이해한다."},
					{"title":"매주 기록하기","objective":"매주 기록한다."}
				],
				"completion_criteria":["별자리를 볼 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "천문 관찰 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"관심 있는 대상을 일상에서 바로 관찰하기 시작한다",
		"핵심 기준으로 구분하고 선택하는 힘을 기른다",
		"반복 가능한 관찰과 기록 루틴을 만든다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMSoftRefinesInstrumentPerformanceExecution(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "피아노로 한 곡을 완주하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"피아노 코스",
				"description":"연주 초안",
				"main_lessons":[
					{"title":"기본 자세 익히기","objective":"자세를 익힌다."},
					{"title":"양손 패턴 익히기","objective":"양손 패턴을 연습한다."},
					{"title":"무대에서 실제 연주하기","objective":"무대에서 실제 연주한다."}
				],
				"completion_criteria":["간단한 곡을 연주할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "피아노 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	last := draft.Lessons[len(draft.Lessons)-1].Lesson.Title
	if last != "목표 연주를 끝까지 이어갈 준비를 한다" {
		t.Fatalf("last lesson title = %q, want soft-refined pre-milestone title", last)
	}
}

func TestBuildDraftWithLLMSoftRefinesWritingHabit(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "매일 짧게 글을 쓰는 습관을 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"글쓰기 루틴 코스",
				"description":"루틴 초안",
				"main_lessons":[
					{"title":"소재 찾기","objective":"소재를 찾는다."},
					{"title":"짧은 글 쓰기","objective":"짧은 글을 쓴다."},
					{"title":"매일 발행하기","objective":"매일 발행한다."}
				],
				"completion_criteria":["짧은 글을 쓸 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "글쓰기 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	last := draft.Lessons[len(draft.Lessons)-1].Lesson.Title
	if last != "지속 가능한 글쓰기 루틴을 정리한다" {
		t.Fatalf("last lesson title = %q, want habit soft-refinement title", last)
	}
}

func TestBuildDraftWithLLMRefinesSongCompletionSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "좋아하는 곡 한 곡을 끝까지 완주하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"피아노 완주 코스",
				"description":"완주 초안",
				"main_lessons":[
					{"title":"기본 자세 익히기","objective":"자세를 익힌다."},
					{"title":"양손 패턴 익히기","objective":"양손 패턴을 연습한다."},
					{"title":"좋아하는 곡 연주하기","objective":"좋아하는 곡을 연주한다."}
				],
				"completion_criteria":["간단한 곡을 연주할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "피아노 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"첫 소리와 기본 패턴으로 곡을 시작한다",
		"자주 나오는 구간 전환을 끊기지 않게 잇는다",
		"목표 곡을 처음부터 끝까지 이어간다",
		"목표 연주를 끝까지 이어갈 준비를 한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesWorshipTeamSupportSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "찬양단 봉사에서 자연스럽게 반주하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"워십 기타 코스",
				"description":"찬양단 반주 초안",
				"main_lessons":[
					{"title":"기본 자세 익히기","objective":"자세를 익힌다."},
					{"title":"오픈 코드 배우기","objective":"오픈 코드를 배운다."},
					{"title":"찬양단에서 실제 반주하기","objective":"찬양단에서 실제 반주한다."}
				],
				"completion_criteria":["간단한 찬양 반주를 할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "통기타 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"첫 소리를 안정적으로 낸다",
		"찬양 반주에 필요한 코드를 끊기지 않고 전환한다",
		"기본 스트로크로 찬양곡 한 곡을 끝까지 이어간다",
		"찬양단 합주 흐름에 맞춰 반주를 준비한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesBrunchSerialPublishSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "브런치에 글을 연재하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"브런치 글쓰기 코스",
				"description":"연재 초안",
				"main_lessons":[
					{"title":"소재 찾기","objective":"소재를 찾는다."},
					{"title":"한 편 쓰기","objective":"한 편을 쓴다."},
					{"title":"브런치에 바로 올리기","objective":"브런치에 바로 올린다."}
				],
				"completion_criteria":["짧은 글을 공개할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "에세이 글쓰기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"짧은 글 한 편을 끝까지 쓴다",
		"전달하려는 주제와 독자 흐름을 잡는다",
		"공개 가능한 글 묶음의 구조를 정리한다",
		"연재 또는 공개 직전까지 글을 다듬고 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCreditBankPracticumSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "학점은행제로 사회복지현장실습을 이수하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"사회복지현장실습 코스",
				"description":"실습 초안",
				"main_lessons":[
					{"title":"오리엔테이션 보기","objective":"실습 OT를 본다."},
					{"title":"기관 찾기","objective":"실습 기관을 찾는다."},
					{"title":"현장실습 하기","objective":"현장실습을 한다."}
				],
				"completion_criteria":["현장실습을 마칠 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "사회복지현장실습",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"실습 목표와 선수요건을 확인한다",
		"기관 섭외와 실습 준비 흐름을 정리한다",
		"현장 실습과 중간점검 흐름을 수행한다",
		"실습일지와 사후평가를 정리해 인정 직전까지 마무리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	if len(draft.Draft.CompletionCriteria) != 2 {
		t.Fatalf("completion criteria len = %d, want 2", len(draft.Draft.CompletionCriteria))
	}
}

func TestBuildDraftWithLLMRefinesCreditBankInstructionalDesignSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "학점은행제로 영유아교수방법론을 통해 수업 설계를 배우고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영유아교수방법론 코스",
				"description":"교수설계 초안",
				"main_lessons":[
					{"title":"교수법 이론 읽기","objective":"교수법 이론을 읽는다."},
					{"title":"교안 써보기","objective":"교안을 쓴다."},
					{"title":"수업하기","objective":"수업을 한다."}
				],
				"completion_criteria":["수업안을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영유아교수방법론",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"교수 목표와 설계 원리를 이해한다",
		"교안과 활동계획 구조를 설계한다",
		"평가와 피드백 흐름을 적용해 다듬는다",
		"수업 운영 직전까지 설계를 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCreditBankStandardTheorySubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "학점은행제로 인간행동과사회환경을 배우고 사회복지 기초를 다지고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"인간행동과사회환경 코스",
				"description":"표준 이론 초안",
				"main_lessons":[
					{"title":"개념 읽기","objective":"개념을 읽는다."},
					{"title":"이론 이해하기","objective":"이론을 이해한다."},
					{"title":"실무에 적용하기","objective":"실무에 적용한다."}
				],
				"completion_criteria":["핵심 이론을 이해할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "인간행동과사회환경",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"과목 범위와 핵심 관점을 이해한다",
		"주요 이론과 개념 축을 체계적으로 익힌다",
		"사례와 제도 관점으로 이론을 적용한다",
		"실무 관점으로 핵심 이론을 통합 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesDailyWritingHabitSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "매일 짧게 기록하는 글쓰기 습관을 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"글쓰기 루틴 코스",
				"description":"루틴 초안",
				"main_lessons":[
					{"title":"소재 찾기","objective":"소재를 찾는다."},
					{"title":"짧은 글 쓰기","objective":"짧은 글을 쓴다."},
					{"title":"매일 발행하기","objective":"매일 발행한다."}
				],
				"completion_criteria":["짧은 글을 쓸 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "글쓰기 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"짧은 기록 한 편을 바로 남긴다",
		"반복 가능한 글쓰기 단위를 만든다",
		"부담 없이 이어가는 쓰기 흐름을 만든다",
		"지속 가능한 글쓰기 루틴을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCraftTeachYoutubeSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "나만의 작품을 만들고 유튜브를 통해 강의하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"코바늘 강의 코스",
				"description":"강의 초안",
				"main_lessons":[
					{"title":"코바늘 도구 이해하기","objective":"코바늘과 실, 기본 준비물을 이해한다."},
					{"title":"기초 뜨개법 익히기","objective":"사슬뜨기와 짧은뜨기 같은 기초 뜨개법을 익힌다."},
					{"title":"유튜브 강의 시작하기","objective":"작품 만드는 과정을 촬영하고 유튜브 강의를 시작한다."}
				],
				"completion_criteria":["작은 코바늘 소품을 스스로 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "코바늘 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기초 뜨개법으로 작은 결과물을 완성한다",
		"반복 패턴을 안정적으로 익혀 작품 형태를 만든다",
		"나만의 작품을 구상하고 도안 흐름을 정리한다",
		"직접 만든 작품 한 점을 완성한다",
		"작품 제작 과정을 설명하고 시연할 준비를 한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesLeathercraftWalletSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "가죽공예로 카드지갑을 직접 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"가죽 카드지갑 코스",
				"description":"카드지갑 제작 초안",
				"main_lessons":[
					{"title":"가죽 도구 이해하기","objective":"도구를 이해한다."},
					{"title":"카드지갑 재단하기","objective":"재단한다."},
					{"title":"카드지갑 완성하기","objective":"완성한다."}
				],
				"completion_criteria":["카드지갑 하나를 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "가죽공예 카드지갑 만들기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"지갑 구조와 패턴을 이해한다",
		"외피와 내피를 재단하고 부품을 준비한다",
		"타공과 새들스티치로 지갑을 조립한다",
		"엣지 마감과 수납 품질을 점검한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	if !containsString(draft.Draft.CompletionCriteria, "실제로 사용할 수 있는 가죽 지갑 한 점을 완성하고 수납과 마감 품질을 점검할 수 있다.") {
		t.Fatalf("completion criteria missing leathercraft wallet caption: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesLeathercraftDimensionalBagSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "가죽공예로 미니백과 클러치 같은 입체 소품을 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"가죽 미니백 코스",
				"description":"입체 가방 제작 초안",
				"main_lessons":[
					{"title":"도구 익히기","objective":"도구를 익힌다."},
					{"title":"가방 만들기","objective":"가방을 만든다."}
				],
				"completion_criteria":["미니백 하나를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "가죽공예 가방 만들기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"평면 소품 기술을 점검하고 입체 구조를 이해한다",
		"입체 패턴과 안감·거싯 구조를 설계한다",
		"포켓·여밈·하드웨어를 순서대로 조립한다",
		"다층 구조를 통합하고 곡선 엣지를 마감한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCalligraphyQuoteArtSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "캘리그라피로 명언 엽서 작품을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"캘리그라피 엽서 코스",
				"description":"문구 작품 제작 초안",
				"main_lessons":[
					{"title":"붓펜 익히기","objective":"붓펜을 익힌다."},
					{"title":"명언 고르기","objective":"명언을 고른다."},
					{"title":"엽서 완성하기","objective":"엽서를 완성한다."}
				],
				"completion_criteria":["캘리그라피 엽서 한 장을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "붓펜 캘리그라피 엽서",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"문구와 작품 방향을 정한다",
		"문장 구도와 글자 리듬을 설계한다",
		"배경과 장식을 더해 작품성을 높인다",
		"캘리그라피 문구 작품을 마무리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	if !containsString(draft.Draft.CompletionCriteria, "문구와 구도를 정해 선물하거나 보관할 수 있는 캘리그라피 작품 한 점을 완성할 수 있다.") {
		t.Fatalf("completion criteria missing calligraphy quote caption: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesCalligraphyPublishSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "캘리그라피 작품을 인스타에 공개하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"캘리그라피 공개 코스",
				"description":"공개 초안",
				"main_lessons":[
					{"title":"문구 쓰기","objective":"문구를 쓴다."},
					{"title":"작품 올리기","objective":"인스타에 작품을 올린다."}
				],
				"completion_criteria":["작품을 공개할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "캘리그라피 작품 만들기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"공개할 캘리그라피 작품의 방향을 정한다",
		"작품 묶음의 글씨체와 구도 흐름을 맞춘다",
		"촬영과 소개 문구로 작품 전달력을 높인다",
		"공개 직전까지 작품 묶음과 게시 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesArtPublishShowcaseSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "그림을 전시와 인스타 업로드용으로 공개하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"작품 공개 코스",
				"description":"공개 초안",
				"main_lessons":[
					{"title":"재료 알아보기","objective":"재료를 이해한다."},
					{"title":"작품 한 점 완성하기","objective":"작품을 완성한다."},
					{"title":"전시하고 업로드하기","objective":"전시하고 업로드한다."}
				],
				"completion_criteria":["작품을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "드로잉 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 표현으로 작은 그림을 완성한다",
		"색과 형태 표현을 안정적으로 이어간다",
		"공개할 작품 묶음의 흐름을 정리한다",
		"공개 직전까지 작품과 소개 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesFreewareMIDIBeatmakingSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "무료 BandLab으로 미디 비트를 처음 만들어보고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"무료 MIDI beat 코스",
				"description":"BandLab beat 제작 초안",
				"main_lessons":[
					{"title":"BandLab 가입하기","objective":"가입한다."},
					{"title":"드럼 샘플 고르기","objective":"샘플을 고른다."},
					{"title":"비트 만들기","objective":"비트를 만든다."}
				],
				"completion_criteria":["짧은 비트를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "무료 미디 음악 제작",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"무료 DAW에서 첫 프로젝트를 준비한다",
		"드럼 패턴과 첫 루프를 만든다",
		"MIDI 노트로 베이스와 멜로디를 얹는다",
		"짧은 MIDI beat를 완성하고 저장한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	if !containsString(draft.Draft.CompletionCriteria, "무료 도구에서 짧은 MIDI beat 한 개를 완성하고 재생·저장 상태를 점검할 수 있다.") {
		t.Fatalf("completion criteria missing freeware MIDI beat caption: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesFreewareMIDIFullTrackSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "무료 DAW로 미디 드럼과 멜로디를 편곡해 짧은 트랙을 완성하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"무료 DAW 트랙 제작 코스",
				"description":"전체 트랙 제작 초안",
				"main_lessons":[
					{"title":"DAW 화면 보기","objective":"화면을 본다."},
					{"title":"루프 배치하기","objective":"루프를 배치한다."},
					{"title":"트랙 완성하기","objective":"트랙을 완성한다."}
				],
				"completion_criteria":["짧은 트랙을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "BandLab MIDI track 만들기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"무료 DAW와 MIDI 제작 흐름을 잡는다",
		"드럼과 코드로 곡의 스케치를 만든다",
		"베이스와 멜로디를 더해 편곡을 확장한다",
		"효과와 오토메이션으로 흐름을 다듬는다",
		"짧은 트랙을 export 직전까지 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesMIDIFoundationWorkflowSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "MIDI 기초 개념과 컨트롤러 송수신 흐름을 이해하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"MIDI 101 코스",
				"description":"MIDI 기초 초안",
				"main_lessons":[
					{"title":"MIDI 정의 보기","objective":"정의를 본다."},
					{"title":"컨트롤러 비교하기","objective":"제품을 비교한다."}
				],
				"completion_criteria":["MIDI를 이해할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "MIDI 101",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"MIDI가 소리를 움직이는 방식을 이해한다",
		"MIDI 입력과 송수신 흐름을 연결한다",
		"Piano Roll에서 짧은 phrase를 만든다",
		"MIDI phrase를 악기 소리로 재생한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesVibeCodingMVPAppSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "바이브 코딩으로 작은 앱 MVP를 처음 만들어보고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"바이브 코딩 입문 코스",
				"description":"AI 앱 제작 초안",
				"main_lessons":[
					{"title":"도구 가입하기","objective":"AI 코딩 도구에 가입한다."},
					{"title":"프롬프트 읽기","objective":"예시 프롬프트를 읽는다."},
					{"title":"앱 만들기","objective":"앱을 만든다."}
				],
				"completion_criteria":["작은 앱을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "vibe coding app 만들기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"AI 코딩 도구와 만들 앱의 범위를 정한다",
		"아이디어를 PRD와 첫 프롬프트로 바꾼다",
		"AI로 첫 작동 버전을 만든다",
		"프롬프트 반복으로 기능과 오류를 고친다",
		"작은 MVP를 완성하고 사용자 흐름을 점검한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
	if !containsString(draft.Draft.CompletionCriteria, "AI 도구로 작은 앱 또는 MVP를 만들고 핵심 사용자 흐름을 직접 점검할 수 있다.") {
		t.Fatalf("completion criteria missing vibe coding MVP caption: %#v", draft.Draft.CompletionCriteria)
	}
}

func TestBuildDraftWithLLMRefinesVibeCodingFullstackShipSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "Cursor와 Supabase, Vercel로 AI 코딩 full-stack 앱을 배포 직전까지 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"AI full-stack 앱 코스",
				"description":"배포형 앱 초안",
				"main_lessons":[
					{"title":"React 컴포넌트 하나 만들기","objective":"컴포넌트를 만든다."},
					{"title":"DB 테이블명 정하기","objective":"테이블명을 정한다."},
					{"title":"배포 버튼 누르기","objective":"배포한다."}
				],
				"completion_criteria":["앱을 배포할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "vibe coding full-stack",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"AI coding 환경과 앱 구조를 준비한다",
		"프론트엔드 화면과 사용자 흐름을 만든다",
		"데이터베이스와 인증을 연결한다",
		"AI가 만든 코드를 검토하고 오류를 정리한다",
		"배포 직전 상태로 앱을 점검한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesClaudeCodeAgenticWorkflowSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "Claude Code와 MCP를 사용해 agentic coding workflow를 만들고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"Claude Code workflow 코스",
				"description":"agentic workflow 초안",
				"main_lessons":[
					{"title":"Claude Code 설치하기","objective":"설치한다."},
					{"title":"slash command 목록 보기","objective":"명령어를 본다."},
					{"title":"MCP 이름 외우기","objective":"MCP 이름을 외운다."}
				],
				"completion_criteria":["Claude Code를 사용할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "Claude Code vibe coding",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"Claude Code 작업 흐름과 안전 기준을 잡는다",
		"context engineering으로 작업 지시를 정리한다",
		"작은 기능을 agentic workflow로 구현한다",
		"hooks와 MCP로 반복 작업을 보강한다",
		"테스트와 정리까지 포함한 개발 workflow를 완성한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesTravelConversationSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "여행 가서 영어로 자연스럽게 말하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"여행 영어 코스",
				"description":"여행 회화 초안",
				"main_lessons":[
					{"title":"기본 문법 보기","objective":"기본 문법을 본다."},
					{"title":"표현 외우기","objective":"표현을 외운다."},
					{"title":"해외에서 바로 말하기","objective":"해외에서 바로 말한다."}
				],
				"completion_criteria":["기본 표현을 말할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영어 회화 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"핵심 표현으로 짧은 여행 대화를 시작한다",
		"기본 문장 패턴으로 질문과 응답을 이어간다",
		"여행 상황에 맞는 대화를 자연스럽게 연결한다",
		"여행 직전까지 필요한 표현 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesPortfolioPublishSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "편집한 영상을 포트폴리오와 유튜브에 공개하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영상 공개 코스",
				"description":"포트폴리오 공개 초안",
				"main_lessons":[
					{"title":"툴 익히기","objective":"툴을 익힌다."},
					{"title":"영상 만들기","objective":"영상을 만든다."},
					{"title":"바로 업로드하기","objective":"바로 업로드한다."}
				],
				"completion_criteria":["영상 한 편을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"공개할 결과물의 방향과 포맷을 정한다",
		"핵심 제작 패턴을 반복해 공개 품질을 만든다",
		"표현 요소를 통합해 완성 흐름을 만든다",
		"공개 직전까지 결과물과 게시 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesCreatorTutorialPublishSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "영상편집 과정을 튜토리얼로 설명하고 강의하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"영상편집 튜토리얼 코스",
				"description":"튜토리얼 초안",
				"main_lessons":[
					{"title":"툴 익히기","objective":"툴을 익힌다."},
					{"title":"영상 만들기","objective":"영상을 만든다."},
					{"title":"강의 올리기","objective":"강의를 올린다."}
				],
				"completion_criteria":["영상 한 편을 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "영상편집 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"도구에 익숙해지며 첫 결과물을 시작한다",
		"핵심 제작 패턴을 반복해 익힌다",
		"표현 요소를 통합해 완성 흐름을 만든다",
		"설명 가능한 결과물 한 편을 완성한다",
		"튜토리얼 공개 직전까지 설명 흐름을 정리한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesMakerWorkshopDemoSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "아두이노 프로젝트 만드는 법을 워크숍에서 설명하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"아두이노 워크숍 코스",
				"description":"워크숍 초안",
				"main_lessons":[
					{"title":"부품 이해하기","objective":"부품을 이해한다."},
					{"title":"회로 연결하기","objective":"회로를 연결한다."},
					{"title":"워크숍 열기","objective":"워크숍을 연다."}
				],
				"completion_criteria":["간단한 회로를 연결할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "아두이노 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 부품과 도구로 첫 작동을 만든다",
		"핵심 제작 패턴을 반복해 익힌다",
		"작동 흐름이 분명한 프로젝트를 완성한다",
		"제작 과정과 원리를 설명할 준비를 한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMRefinesBakingClassDemoSubpattern(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "베이킹 레시피를 설명하며 클래스를 해보고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"베이킹 클래스 코스",
				"description":"강의 초안",
				"main_lessons":[
					{"title":"재료 이해하기","objective":"재료를 이해한다."},
					{"title":"케이크 만들기","objective":"케이크를 만든다."},
					{"title":"바로 클래스 열기","objective":"바로 클래스를 연다."}
				],
				"completion_criteria":["케이크를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "베이킹 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	wantTitles := []string{
		"기본 재료와 도구로 첫 결과물을 시작한다",
		"핵심 조리와 제과 패턴을 안정적으로 익힌다",
		"설명 가능한 대표 메뉴를 완성한다",
		"레시피와 시연 흐름을 설명할 준비를 한다",
	}
	if len(draft.Lessons) != len(wantTitles) {
		t.Fatalf("lessons len = %d, want %d", len(draft.Lessons), len(wantTitles))
	}
	for idx, want := range wantTitles {
		if got := draft.Lessons[idx].Lesson.Title; got != want {
			t.Fatalf("lesson[%d] title = %q, want %q", idx, got, want)
		}
	}
}

func TestBuildDraftWithLLMSoftRefinesVisualArtPublish(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "그림을 전시와 인스타 업로드용으로 공개하고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"작품 공개 코스",
				"description":"공개 초안",
				"main_lessons":[
					{"title":"재료 알아보기","objective":"재료를 이해한다."},
					{"title":"작품 한 점 완성하기","objective":"작품을 완성한다."},
					{"title":"전시하고 업로드하기","objective":"전시하고 업로드한다."}
				],
				"completion_criteria":["작품을 완성할 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "드로잉 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	last := draft.Lessons[len(draft.Lessons)-1].Lesson.Title
	if last != "공개 직전까지 작품과 소개 흐름을 정리한다" {
		t.Fatalf("last lesson title = %q, want publish soft-refinement title", last)
	}
}

func TestBuildDraftWithLLMSoftRefinesCookingTeaching(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "베이킹 레시피를 설명하며 클래스를 해보고 싶다"
	client := &stubLLMClient{
		responses: []string{
			`{
				"title":"베이킹 클래스 코스",
				"description":"강의 초안",
				"main_lessons":[
					{"title":"재료 이해하기","objective":"재료를 이해한다."},
					{"title":"케이크 만들기","objective":"케이크를 만든다."},
					{"title":"바로 클래스 열기","objective":"바로 클래스를 연다."}
				],
				"completion_criteria":["케이크를 만들 수 있다."]
			}`,
		},
	}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:  "베이킹 배우기",
		LearningGoal: &goal,
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	last := draft.Lessons[len(draft.Lessons)-1].Lesson.Title
	if last != "레시피와 시연 흐름을 설명할 준비를 한다" {
		t.Fatalf("last lesson title = %q, want cooking teaching soft-refinement title", last)
	}
}

func TestBuildDraftWithLLMPreservesGeneratedRecommendationSearchSpec(t *testing.T) {
	svc := NewService()
	userID := uuid.New()
	goal := "Facilitate a short remote team meeting with a clear agenda and follow-up"
	client := &stubLLMClient{responses: []string{`{
		"title":"Remote Meeting Facilitation Course",
		"description":"A short course for planning and leading remote meetings.",
		"main_lessons":[
			{
				"title":"Set the meeting goal and agenda",
				"objective":"Define the meeting outcome, agenda items, and participant roles.",
				"recommendation_search_spec":{
					"primary_query":"remote meeting facilitation agenda tutorial",
					"intent":"tutorial",
					"must_include":["remote meeting", "agenda", "facilitation"],
					"nice_to_have":["beginner", "team meeting"],
					"avoid":["sales", "software ad"],
					"content_types":["video"],
					"language":"en",
					"stage_role":"setup_intro",
					"source":"generated"
				}
			},
			{"title":"Practice opening and guiding discussion","objective":"Use short prompts to open the meeting and keep discussion focused."},
			{"title":"Handle decisions and action items","objective":"Summarize decisions and assign clear next steps."},
			{"title":"Prepare to run the full meeting","objective":"Connect the preparation into a complete remote meeting run-through."}
		],
		"completion_criteria":["Run a 15-minute remote meeting with agenda, discussion, and follow-up."]
	}`}}

	draft, err := svc.BuildDraftWithLLM(context.Background(), client, userID, CreateCourseDraftRequest{
		SourceQuery:      "remote meeting facilitation",
		LearningGoal:     &goal,
		LearningLanguage: "en",
	}, DraftBuildOptions{SkipReview: true})
	if err != nil {
		t.Fatalf("BuildDraftWithLLM() error = %v", err)
	}

	if len(draft.Lessons) == 0 {
		t.Fatal("expected generated lessons")
	}
	spec := draft.Lessons[0].Lesson.RecommendationSearchSpec
	if spec.Source != SearchSpecSourceGenerated {
		t.Fatalf("search spec source = %q, want %q", spec.Source, SearchSpecSourceGenerated)
	}
	if spec.PrimaryQuery != "remote meeting facilitation agenda tutorial" {
		t.Fatalf("primary query = %q", spec.PrimaryQuery)
	}
	if !searchSpecPrefersContentType(spec.ContentTypes, "video") {
		t.Fatalf("expected video content type, got %#v", spec.ContentTypes)
	}
}
