package curriculum

import "strings"

func buildIntroStageTitle(sourceQuery, learningGoal, language string) string {
	context := strings.ToLower(strings.TrimSpace(sourceQuery + " " + learningGoal))
	if normalizeLearningLanguage(language) == "en" {
		switch {
		case containsCreatorTutorialPublishContext(context):
			return "Plan the tutorial structure and demo flow"
		case containsMIDIProductionContext(context):
			return "Set up the DAW and first track idea"
		case containsCalligraphyContext(context):
			return "Plan the quote layout and lettering tools"
		case containsLeathercraftContext(context):
			return "Prepare the pattern, leather, and stitching tools"
		case containsSewingProjectFoundationContext(context):
			return "Prepare fabric, pattern pieces, and stitch settings"
		case containsAny(context, "arduino", "sensor demo", "maker workshop", "prototype demo"):
			return "Set up the prototype and demo components"
		case containsCardboardCircuitInventionPathContext(context):
			return "Plan the circuit layout and cardboard build"
		case containsDesignLabPrototypeProjectContext(context):
			return "Frame the user problem and prototype goal"
		case containsGame3DGuidedPathwayContext(context):
			return "Set up the 3D scene and player path goal"
		case containsVibeCodingContext(context) && containsAny(context, "claude code", "mcp", "hook", "hooks", "agentic"):
			return "Set up the Claude Code workflow workspace"
		case containsClimateDataGraphingExplanationContext(context):
			return "Choose the dataset and trend question"
		case containsCitizenScienceObservationRecordContext(context):
			return "Choose the observation area and recording tools"
		case containsPrimarySourceInquiryNoteContext(context):
			return "Choose the source and inquiry focus"
		case containsShortCourseWeeklyDiscussionContext(context):
			return "Map the course schedule and discussion routine"
		case containsOrganicGrowingCyclePlanContext(context):
			return "Prepare the bed, seed choices, and care calendar"
		case containsAIDigitalLiteracyLifePracticeContext(context):
			return "Set safe boundaries for daily AI use"
		case containsAny(context, "essay series", "weekly essay", "publish a short essay", "online publication"):
			return "Define the essay theme and publishing rhythm"
		case strings.Contains(context, "yoga") || strings.Contains(context, "flexibility") || strings.Contains(context, "breathing routine") || strings.Contains(context, "morning routine"):
			return "Start with breath and safe beginner poses"
		case strings.Contains(context, "dance") || strings.Contains(context, "choreography") || strings.Contains(context, "k-pop") || strings.Contains(context, "cover routine"):
			return "Map the rhythm and point moves"
		case strings.Contains(context, "guitar") || strings.Contains(context, "piano") || strings.Contains(context, "drum") || strings.Contains(context, "violin") || strings.Contains(context, "cello") || strings.Contains(context, "performance") || strings.Contains(context, "play "):
			return "Make the first stable sound"
		case strings.Contains(context, "english") || strings.Contains(context, "japanese") || strings.Contains(context, "conversation") || strings.Contains(context, "language"):
			return "Prepare to start a short conversation"
		case strings.Contains(context, "drawing") || strings.Contains(context, "watercolor") || strings.Contains(context, "sketch") || strings.Contains(context, "painting") || strings.Contains(context, "photo"):
			return "Start a small output with the core tools"
		default:
			return "Combine the basics and begin the first practice"
		}
	}
	switch {
	case strings.Contains(context, "기타") || strings.Contains(context, "피아노") || strings.Contains(context, "드럼") || strings.Contains(context, "연주"):
		return "첫 소리를 안정적으로 낸다"
	case strings.Contains(context, "영어") || strings.Contains(context, "일본어") || strings.Contains(context, "회화") || strings.Contains(context, "언어"):
		return "짧은 대화를 시작할 준비를 한다"
	case strings.Contains(context, "그림") || strings.Contains(context, "수채화") || strings.Contains(context, "드로잉") || strings.Contains(context, "스케치"):
		return "도구를 써서 작은 결과물을 시작한다"
	default:
		return "기본 요소를 묶어 첫 수행을 시작한다"
	}
}

func buildIntroStageObjective(sourceQuery, learningGoal, originalObjective, language string) string {
	base := strings.TrimSpace(originalObjective)
	if base == "" {
		base = buildPreMilestoneObjective(learningGoal, language)
	}
	if normalizeLearningLanguage(language) == "en" {
		context := strings.ToLower(strings.TrimSpace(sourceQuery + " " + learningGoal))
		switch {
		case containsCreatorTutorialPublishContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect the teaching outline, recording setup, and final explanation flow."))
		case containsMIDIProductionContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect beat, arrangement, mixing, and export decisions."))
		case containsCalligraphyContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect lettering practice, layout, decoration, and final card review."))
		case containsLeathercraftContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect patterning, cutting, stitching, edge finishing, and final quality checks."))
		case containsSewingProjectFoundationContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect fabric cutting, straight stitching, handle attachment, edge finishing, and final inspection."))
		case containsAny(context, "arduino", "sensor demo", "maker workshop", "prototype demo"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect the circuit, code, sensor output, and workshop explanation."))
		case containsCardboardCircuitInventionPathContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect circuit layout, LED and switch wiring, cardboard assembly, testing, and prototype revision."))
		case containsDesignLabPrototypeProjectContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect problem framing, prototype building, user feedback, and revision decisions."))
		case containsGame3DGuidedPathwayContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect scene setup, guided player movement, simple interaction, and playtesting."))
		case containsVibeCodingContext(context) && containsAny(context, "claude code", "mcp", "hook", "hooks", "agentic"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect Claude Code setup, hooks, MCP context, feature implementation, and workflow review."))
		case containsClimateDataGraphingExplanationContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect data reading, graph creation, trend comparison, and explanation."))
		case containsCitizenScienceObservationRecordContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect field observation, iNaturalist recording, plant identification, and biodiversity note review."))
		case containsPrimarySourceInquiryNoteContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect source observation, reflection, question writing, and inquiry note revision."))
		case containsShortCourseWeeklyDiscussionContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect weekly course viewing, note taking, discussion posting, and peer response."))
		case containsOrganicGrowingCyclePlanContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect soil preparation, sowing, watering, pest checks, and harvest review."))
		case containsAIDigitalLiteracyLifePracticeContext(context):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect safe prompting, source checking, privacy checks, and daily task routines."))
		case containsAny(context, "essay series", "weekly essay", "publish a short essay", "online publication"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect topic planning, drafting, revision, and publication preparation."))
		case strings.Contains(context, "yoga") || strings.Contains(context, "flexibility") || strings.Contains(context, "breathing routine") || strings.Contains(context, "morning routine"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to connect breath, comfort checks, and a short repeatable flow."))
		case strings.Contains(context, "dance") || strings.Contains(context, "choreography") || strings.Contains(context, "k-pop") || strings.Contains(context, "cover routine"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to practice counts, point moves, and transitions before recording."))
		case strings.Contains(context, "guitar") || strings.Contains(context, "piano") || strings.Contains(context, "drum") || strings.Contains(context, "violin") || strings.Contains(context, "cello") || strings.Contains(context, "performance") || strings.Contains(context, "play "):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to start the first playable performance."))
		case strings.Contains(context, "english") || strings.Contains(context, "japanese") || strings.Contains(context, "conversation") || strings.Contains(context, "language"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to begin short questions and answers."))
		case strings.Contains(context, "drawing") || strings.Contains(context, "watercolor") || strings.Contains(context, "sketch") || strings.Contains(context, "painting") || strings.Contains(context, "photo"):
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to start a small practical output."))
		default:
			return dedupeObjectiveSentences(strings.TrimSpace(base + " Use this process to enter the first practice stage."))
		}
	}
	context := strings.ToLower(strings.TrimSpace(sourceQuery + " " + learningGoal))
	switch {
	case strings.Contains(context, "기타") || strings.Contains(context, "피아노") || strings.Contains(context, "드럼") || strings.Contains(context, "연주"):
		return dedupeObjectiveSentences(strings.TrimSpace(base + " 이 과정을 바탕으로 첫 연주 수행을 시작할 수 있게 한다."))
	case strings.Contains(context, "영어") || strings.Contains(context, "일본어") || strings.Contains(context, "회화") || strings.Contains(context, "언어"):
		return dedupeObjectiveSentences(strings.TrimSpace(base + " 이 과정을 바탕으로 짧은 질문과 응답을 시작할 수 있게 한다."))
	case strings.Contains(context, "그림") || strings.Contains(context, "수채화") || strings.Contains(context, "드로잉") || strings.Contains(context, "스케치"):
		return dedupeObjectiveSentences(strings.TrimSpace(base + " 이 과정을 바탕으로 작은 작업을 실제로 시작할 수 있게 한다."))
	default:
		return dedupeObjectiveSentences(strings.TrimSpace(base + " 이 과정을 바탕으로 첫 수행 단계에 들어갈 수 있게 한다."))
	}
}

func dedupeObjectiveSentences(objective string) string {
	trimmed := strings.TrimSpace(objective)
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, ".")
	unique := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		sentence := strings.Join(strings.Fields(strings.TrimSpace(part)), " ")
		if sentence == "" {
			continue
		}
		if _, exists := seen[sentence]; exists {
			continue
		}
		seen[sentence] = struct{}{}
		unique = append(unique, sentence)
	}
	if len(unique) == 0 {
		return ""
	}
	return strings.Join(unique, ". ") + "."
}
