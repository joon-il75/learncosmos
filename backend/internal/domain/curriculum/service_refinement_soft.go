package curriculum

func softenInstrumentPerformanceExecution(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) == 0 {
		return doc
	}
	lastIdx := len(doc.MainLessons) - 1
	doc.MainLessons[lastIdx] = generatedMainLesson{
		Title:     "목표 연주를 끝까지 이어갈 준비를 한다",
		Objective: "배운 기술을 연결해 목표 연주나 과제를 처음부터 끝까지 이어갈 준비를 할 수 있게 한다.",
	}
	return doc
}

func softenWritingHabitLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) == 0 {
		return doc
	}
	lastIdx := len(doc.MainLessons) - 1
	doc.MainLessons[lastIdx] = generatedMainLesson{
		Title:     "지속 가능한 글쓰기 루틴을 정리한다",
		Objective: "무리 없이 이어갈 수 있는 글쓰기 빈도와 흐름을 정리해 일상 루틴으로 만들 수 있게 한다.",
	}
	return doc
}

func softenVisualArtPublishLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) == 0 {
		return doc
	}
	lastIdx := len(doc.MainLessons) - 1
	doc.MainLessons[lastIdx] = generatedMainLesson{
		Title:     "공개 직전까지 작품과 소개 흐름을 정리한다",
		Objective: "전시나 업로드 직전까지 작품 묶음, 제목, 소개 흐름을 정리할 수 있게 한다.",
	}
	return doc
}

func softenCookingTeachingLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) == 0 {
		return doc
	}
	lastIdx := len(doc.MainLessons) - 1
	doc.MainLessons[lastIdx] = generatedMainLesson{
		Title:     "레시피와 시연 흐름을 설명할 준비를 한다",
		Objective: "배운 메뉴의 레시피, 순서, 시연 포인트를 정리해 설명할 준비를 할 수 있게 한다.",
	}
	return doc
}

func softenCookingHabitLessons(doc generatedDraftDocument) generatedDraftDocument {
	if len(doc.MainLessons) < 3 {
		doc.MainLessons = []generatedMainLesson{
			{
				Title:     "기본 재료와 조리 흐름으로 한 끼를 완성한다",
				Objective: "부담 없는 재료와 기본 조리 흐름으로 가족과 함께 먹을 한 끼를 직접 완성할 수 있게 한다.",
			},
			{
				Title:     "반복 가능한 건강 식사 패턴을 만든다",
				Objective: "자주 쓰는 재료 조합과 조리 패턴을 익혀 건강한 식사를 반복 가능하게 준비할 수 있게 한다.",
			},
			{
				Title:     "가족과 일상에 맞는 식사 루틴을 정리한다",
				Objective: "가족의 식사 흐름과 생활 리듬에 맞춰 무리 없이 이어갈 식사 준비 루틴을 정리할 수 있게 한다.",
			},
		}
		return doc
	}

	doc.MainLessons[0] = generatedMainLesson{
		Title:     "기본 재료와 조리 흐름으로 한 끼를 완성한다",
		Objective: "부담 없는 재료와 기본 조리 흐름으로 가족과 함께 먹을 한 끼를 직접 완성할 수 있게 한다.",
	}
	lastIdx := len(doc.MainLessons) - 1
	doc.MainLessons[lastIdx] = generatedMainLesson{
		Title:     "가족과 일상에 맞는 식사 루틴을 정리한다",
		Objective: "가족의 식사 흐름과 생활 리듬에 맞춰 무리 없이 이어갈 식사 준비 루틴을 정리할 수 있게 한다.",
	}
	return doc
}
