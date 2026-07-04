package curriculum

import "strings"

func refineWorshipGuitarLessons(lessons []generatedMainLesson) []generatedMainLesson {
	if len(lessons) == 0 {
		return lessons
	}

	refined := make([]generatedMainLesson, 0, maxInt(len(lessons), 4))
	refined = append(refined, generatedMainLesson{
		Title:     "첫 소리를 안정적으로 낸다",
		Objective: "기타를 안정적으로 잡고 튜닝하며, 찬양 반주에 필요한 오픈 코드를 눌러 첫 반주를 시작할 수 있게 한다.",
	})
	refined = append(refined, generatedMainLesson{
		Title:     "찬양 반주에 필요한 코드를 끊기지 않고 전환한다",
		Objective: "찬양곡에서 자주 쓰는 오픈 코드를 끊기지 않게 전환하며 반주 흐름을 유지할 수 있게 한다.",
	})
	refined = append(refined, generatedMainLesson{
		Title:     "기본 스트로크로 찬양곡 한 곡을 끝까지 이어간다",
		Objective: "기본 스트로크를 유지하며 쉬운 찬양곡 한 곡을 처음부터 끝까지 안정적으로 반주할 수 있게 한다.",
	})

	hasPrepLesson := false
	for _, lesson := range lessons {
		text := strings.ToLower(strings.TrimSpace(lesson.Title + " " + lesson.Objective))
		if strings.Contains(text, "합주") || strings.Contains(text, "흐름") || strings.Contains(text, "준비") {
			hasPrepLesson = true
			break
		}
	}
	if hasPrepLesson || len(lessons) >= 3 {
		refined = append(refined, generatedMainLesson{
			Title:     "찬양단 합주 흐름에 맞춰 반주를 준비한다",
			Objective: "곡의 시작, 진행, 마침 타이밍을 맞추며 찬양단 합주 흐름에 맞춰 반주를 연결할 수 있게 한다.",
		})
	}

	return refined
}

func refineCrochetTeachingLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "기초 뜨개법으로 작은 결과물을 완성한다",
			Objective: "사슬뜨기와 짧은뜨기 같은 기초 뜨개법을 익혀 작은 코바늘 결과물을 직접 완성할 수 있게 한다.",
		},
		{
			Title:     "반복 패턴을 안정적으로 익혀 작품 형태를 만든다",
			Objective: "반복되는 뜨개 패턴을 안정적으로 이어가며 원하는 작품 형태를 만들 수 있게 한다.",
		},
		{
			Title:     "나만의 작품을 구상하고 도안 흐름을 정리한다",
			Objective: "작품 아이디어를 구체화하고, 실과 바늘 선택, 반복 패턴, 도안 표기, 단수와 코수 흐름을 정리할 수 있게 한다.",
		},
		{
			Title:     "직접 만든 작품 한 점을 완성한다",
			Objective: "구상한 도안 흐름을 바탕으로 나만의 코바늘 작품 한 점을 완성할 수 있게 한다.",
		},
		{
			Title:     "작품 제작 과정을 설명하고 시연할 준비를 한다",
			Objective: "작품 제작 과정을 단계별로 설명하고, 유튜브 강의를 위한 시연 흐름을 정리할 수 있게 한다.",
		},
	}
}

func refineWorshipTeamSupportLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "첫 소리를 안정적으로 낸다",
			Objective: "기타를 안정적으로 잡고 튜닝하며, 찬양 반주에 필요한 오픈 코드를 눌러 첫 반주를 시작할 수 있게 한다.",
		},
		{
			Title:     "찬양 반주에 필요한 코드를 끊기지 않고 전환한다",
			Objective: "찬양곡에서 자주 쓰는 오픈 코드를 끊기지 않게 전환하며 반주 흐름을 유지할 수 있게 한다.",
		},
		{
			Title:     "기본 스트로크로 찬양곡 한 곡을 끝까지 이어간다",
			Objective: "기본 스트로크를 유지하며 쉬운 찬양곡 한 곡을 처음부터 끝까지 안정적으로 반주할 수 있게 한다.",
		},
		{
			Title:     "찬양단 합주 흐름에 맞춰 반주를 준비한다",
			Objective: "곡의 시작, 진행, 마침 타이밍을 맞추며 찬양단 합주 흐름에 맞춰 반주를 연결할 수 있게 한다.",
		},
	}
}

func refineSongCompletionLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "첫 소리와 기본 패턴으로 곡을 시작한다",
			Objective: "목표 곡에 필요한 기본 소리와 패턴을 연결해 곡을 실제로 시작할 수 있게 한다.",
		},
		{
			Title:     "자주 나오는 구간 전환을 끊기지 않게 잇는다",
			Objective: "자주 막히는 구간과 전환을 반복해 끊기지 않게 이어갈 수 있게 한다.",
		},
		{
			Title:     "목표 곡을 처음부터 끝까지 이어간다",
			Objective: "목표 곡을 처음부터 끝까지 멈추지 않고 이어가며 완주 흐름을 만들 수 있게 한다.",
		},
		{
			Title:     "목표 연주를 끝까지 이어갈 준비를 한다",
			Objective: "실제 연주 직전까지 약한 구간, 연결 흐름, 마무리 타이밍을 점검하고 정리할 수 있게 한다.",
		},
	}
}

func refineBrunchSerialPublishLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "짧은 글 한 편을 끝까지 쓴다",
			Objective: "짧지만 완결된 글 한 편을 직접 끝까지 써 보며 공개 가능한 글의 최소 단위를 만든다.",
		},
		{
			Title:     "전달하려는 주제와 독자 흐름을 잡는다",
			Objective: "무엇을 누구에게 어떻게 전달할지 정리하고 글의 중심 흐름을 잡을 수 있게 한다.",
		},
		{
			Title:     "공개 가능한 글 묶음의 구조를 정리한다",
			Objective: "연재나 공개를 염두에 두고 글 묶음의 순서와 연결 구조를 정리할 수 있게 한다.",
		},
		{
			Title:     "연재 또는 공개 직전까지 글을 다듬고 정리한다",
			Objective: "공개 직전까지 문장, 흐름, 제목, 구성 요소를 다듬어 밖에 내놓을 준비를 마칠 수 있게 한다.",
		},
	}
}

func refineDailyWritingHabitLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "짧은 기록 한 편을 바로 남긴다",
			Objective: "부담 없는 짧은 글 한 편을 바로 써 보며 기록을 시작할 수 있게 한다.",
		},
		{
			Title:     "반복 가능한 글쓰기 단위를 만든다",
			Objective: "매일 반복해도 부담이 크지 않은 글쓰기 단위와 길이를 정할 수 있게 한다.",
		},
		{
			Title:     "부담 없이 이어가는 쓰기 흐름을 만든다",
			Objective: "소재 찾기, 초안 쓰기, 마무리하기를 부담 없이 이어가는 흐름을 만들 수 있게 한다.",
		},
		{
			Title:     "지속 가능한 글쓰기 루틴을 정리한다",
			Objective: "무리 없이 이어갈 수 있는 빈도와 흐름을 정리해 일상 루틴으로 만들 수 있게 한다.",
		},
	}
}

func refineCraftTeachYoutubeLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "기초 뜨개법으로 작은 결과물을 완성한다",
			Objective: "사슬뜨기와 짧은뜨기 같은 기초 뜨개법을 익혀 작은 코바늘 결과물을 직접 완성할 수 있게 한다.",
		},
		{
			Title:     "반복 패턴을 안정적으로 익혀 작품 형태를 만든다",
			Objective: "반복되는 뜨개 패턴을 안정적으로 이어가며 원하는 작품 형태를 만들 수 있게 한다.",
		},
		{
			Title:     "나만의 작품을 구상하고 도안 흐름을 정리한다",
			Objective: "작품 아이디어를 구체화하고, 실과 바늘 선택, 반복 패턴, 도안 표기, 단수와 코수 흐름을 정리할 수 있게 한다.",
		},
		{
			Title:     "직접 만든 작품 한 점을 완성한다",
			Objective: "구상한 도안 흐름을 바탕으로 나만의 코바늘 작품 한 점을 완성할 수 있게 한다.",
		},
		{
			Title:     "작품 제작 과정을 설명하고 시연할 준비를 한다",
			Objective: "작품 제작 과정을 단계별로 설명하고, 유튜브 강의를 위한 시연 흐름을 정리할 수 있게 한다.",
		},
	}
}

func refineTravelConversationLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "핵심 표현으로 짧은 여행 대화를 시작한다",
			Objective: "공항, 식당, 길 묻기처럼 자주 마주치는 여행 상황에서 핵심 표현으로 짧은 대화를 시작할 수 있게 한다.",
		},
		{
			Title:     "기본 문장 패턴으로 질문과 응답을 이어간다",
			Objective: "자주 쓰는 기본 문장 패턴을 연결해 짧은 질문과 응답을 끊기지 않게 이어갈 수 있게 한다.",
		},
		{
			Title:     "여행 상황에 맞는 대화를 자연스럽게 연결한다",
			Objective: "공항, 식당, 길 묻기 같은 목표 상황에 맞는 질문과 응답 흐름을 자연스럽게 연결할 수 있게 한다.",
		},
		{
			Title:     "여행 직전까지 필요한 표현 흐름을 정리한다",
			Objective: "실제 여행 직전까지 필요한 표현, 반응, 이어가기 흐름을 스스로 정리하고 점검할 수 있게 한다.",
		},
	}
}

func refineArtPublishShowcaseLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "기본 표현으로 작은 그림을 완성한다",
			Objective: "선, 형태, 명암, 색 같은 기본 표현을 사용해 작은 그림 한 점을 직접 완성할 수 있게 한다.",
		},
		{
			Title:     "색과 형태 표현을 안정적으로 이어간다",
			Objective: "색과 형태를 다루는 핵심 패턴을 반복해 표현 흐름을 안정적으로 이어갈 수 있게 한다.",
		},
		{
			Title:     "공개할 작품 묶음의 흐름을 정리한다",
			Objective: "전시나 업로드를 염두에 두고 공개할 작품 묶음의 순서와 연결 흐름을 정리할 수 있게 한다.",
		},
		{
			Title:     "공개 직전까지 작품과 소개 흐름을 정리한다",
			Objective: "전시나 업로드 직전까지 작품 묶음, 제목, 소개 흐름을 정리할 수 있게 한다.",
		},
	}
}

func refineCalligraphyBasicLetteringLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "붓펜과 기본 선 흐름을 익힌다",
			Objective: "붓펜 잡기, 압력 조절, 굵고 얇은 선의 변화를 익혀 캘리그라피의 기본 손 움직임을 만들 수 있게 한다.",
		},
		{
			Title:     "자모와 한 단어를 안정적으로 쓴다",
			Objective: "자음과 모음의 형태를 반복하고 한 글자와 한 단어를 균형 있게 쓰며 기본 글씨체를 잡을 수 있게 한다.",
		},
		{
			Title:     "짧은 문장의 리듬과 여백을 맞춘다",
			Objective: "짧은 문장을 쓰며 글자 크기, 간격, 행 흐름, 여백을 조절해 읽기 좋은 문장 구성을 만들 수 있게 한다.",
		},
		{
			Title:     "짧은 캘리그라피 결과물을 완성한다",
			Objective: "연습한 선과 글씨체를 바탕으로 짧은 문장 캘리그라피 결과물 한 점을 완성하고 흔들림과 균형을 점검할 수 있게 한다.",
		},
	}
}

func refineCalligraphyQuoteArtProjectLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "문구와 작품 방향을 정한다",
			Objective: "작품에 담을 명언이나 문장을 고르고, 전달하려는 분위기와 사용할 종이·도구·색의 방향을 정할 수 있게 한다.",
		},
		{
			Title:     "문장 구도와 글자 리듬을 설계한다",
			Objective: "글자 크기, 줄 나눔, 중심축, 여백을 설계해 문구가 잘 읽히는 캘리그라피 구도를 만들 수 있게 한다.",
		},
		{
			Title:     "배경과 장식을 더해 작품성을 높인다",
			Objective: "간단한 수채 배경, 포인트 장식, 색 대비를 적용해 글씨가 묻히지 않는 작품 흐름을 만들 수 있게 한다.",
		},
		{
			Title:     "캘리그라피 문구 작품을 마무리한다",
			Objective: "최종 문구를 다시 쓰고 번짐, 여백, 장식, 마감 상태를 점검해 선물하거나 보관할 수 있는 작품 한 점을 완성할 수 있게 한다.",
		},
	}
}

func refineCalligraphyPublishShowcaseLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "공개할 캘리그라피 작품의 방향을 정한다",
			Objective: "어떤 주제와 분위기의 작품을 공개할지 정하고 첫 작품 한 점을 공개 가능한 형태로 시작할 수 있게 한다.",
		},
		{
			Title:     "작품 묶음의 글씨체와 구도 흐름을 맞춘다",
			Objective: "여러 작품을 만들며 글씨체, 문구 길이, 여백, 색감의 일관성을 맞춰 공개용 묶음으로 정리할 수 있게 한다.",
		},
		{
			Title:     "촬영과 소개 문구로 작품 전달력을 높인다",
			Objective: "작품 사진을 정리하고 제목, 설명, 제작 의도를 써서 보는 사람이 작품 흐름을 이해할 수 있게 한다.",
		},
		{
			Title:     "공개 직전까지 작품 묶음과 게시 흐름을 정리한다",
			Objective: "SNS나 공개 플랫폼에 올리기 직전까지 작품 순서, 설명 문구, 게시 형식을 점검하고 정리할 수 있게 한다.",
		},
	}
}

func refinePortfolioPublishLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "공개할 결과물의 방향과 포맷을 정한다",
			Objective: "어떤 결과물을 어떤 플랫폼과 포맷으로 공개할지 방향을 정할 수 있게 한다.",
		},
		{
			Title:     "핵심 제작 패턴을 반복해 공개 품질을 만든다",
			Objective: "편집, 구성, 시각 요소 같은 핵심 제작 패턴을 반복해 공개 가능한 품질을 만들 수 있게 한다.",
		},
		{
			Title:     "표현 요소를 통합해 완성 흐름을 만든다",
			Objective: "배운 요소를 통합해 처음부터 끝까지 이어지는 완성 흐름을 만들 수 있게 한다.",
		},
		{
			Title:     "공개 직전까지 결과물과 게시 흐름을 정리한다",
			Objective: "공개 직전까지 결과물 품질과 게시 흐름을 점검하고 정리할 수 있게 한다.",
		},
	}
}

func refineCreatorTutorialPublishLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "도구에 익숙해지며 첫 결과물을 시작한다",
			Objective: "편집기나 제작 도구의 기본 흐름에 익숙해지며 첫 디지털 결과물을 바로 시작할 수 있게 한다.",
		},
		{
			Title:     "핵심 제작 패턴을 반복해 익힌다",
			Objective: "컷 편집, 배치, 레이어, 효과 같은 핵심 제작 패턴을 반복해 익힐 수 있게 한다.",
		},
		{
			Title:     "표현 요소를 통합해 완성 흐름을 만든다",
			Objective: "배운 제작 요소를 통합해 처음부터 끝까지 이어지는 완성 흐름을 만들 수 있게 한다.",
		},
		{
			Title:     "설명 가능한 결과물 한 편을 완성한다",
			Objective: "다른 사람에게 설명할 수 있을 만큼 구조가 분명한 결과물 한 편 또는 한 개를 완성할 수 있게 한다.",
		},
		{
			Title:     "튜토리얼 공개 직전까지 설명 흐름을 정리한다",
			Objective: "튜토리얼 공개 직전까지 단계 설명, 예시 화면, 전달 순서를 정리할 수 있게 한다.",
		},
	}
}

func refineFreewareMIDIBeatmakingLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "무료 DAW에서 첫 프로젝트를 준비한다",
			Objective: "BandLab, LMMS, TunePad 같은 무료 도구의 기본 화면과 트랙 구조를 이해하고 첫 프로젝트를 바로 시작할 수 있게 한다.",
		},
		{
			Title:     "드럼 패턴과 첫 루프를 만든다",
			Objective: "드럼 머신이나 step sequencer로 기본 리듬을 만들고 반복 재생되는 첫 loop를 구성할 수 있게 한다.",
		},
		{
			Title:     "MIDI 노트로 베이스와 멜로디를 얹는다",
			Objective: "Piano Roll이나 MIDI editor에서 노트 길이, 위치, velocity를 조절해 beat 위에 베이스와 짧은 멜로디를 얹을 수 있게 한다.",
		},
		{
			Title:     "짧은 MIDI beat를 완성하고 저장한다",
			Objective: "드럼, 베이스, 멜로디 균형을 맞춰 짧은 beat 한 개를 완성하고 재생·저장 상태를 점검할 수 있게 한다.",
		},
	}
}

func refineFreewareMIDIFullTrackLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "무료 DAW와 MIDI 제작 흐름을 잡는다",
			Objective: "무료 DAW의 트랙, 클립, Piano Roll, mixer 위치를 파악하고 짧은 곡 제작에 필요한 기본 흐름을 잡을 수 있게 한다.",
		},
		{
			Title:     "드럼과 코드로 곡의 스케치를 만든다",
			Objective: "드럼 패턴과 코드 진행을 MIDI로 입력해 곡의 기본 groove와 분위기를 잡는 첫 sketch를 만들 수 있게 한다.",
		},
		{
			Title:     "베이스와 멜로디를 더해 편곡을 확장한다",
			Objective: "베이스, 멜로디, 보조 악기를 배치하고 구간별 변화를 만들어 loop 반복을 짧은 arrangement로 확장할 수 있게 한다.",
		},
		{
			Title:     "효과와 오토메이션으로 흐름을 다듬는다",
			Objective: "EQ, reverb, delay, filter sweep, volume automation 같은 기본 처리를 적용해 구간 전환과 소리 균형을 다듬을 수 있게 한다.",
		},
		{
			Title:     "짧은 트랙을 export 직전까지 정리한다",
			Objective: "인트로, 전개, 마무리 구조를 점검하고 전체 볼륨과 믹스 상태를 정리해 export 직전 상태까지 만들 수 있게 한다.",
		},
	}
}

func refineMIDIFoundationWorkflowLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "MIDI가 소리를 움직이는 방식을 이해한다",
			Objective: "MIDI가 실제 오디오가 아니라 note, timing, velocity 같은 연주 정보를 전달한다는 점을 이해하고 제작 흐름과 연결할 수 있게 한다.",
		},
		{
			Title:     "MIDI 입력과 송수신 흐름을 연결한다",
			Objective: "키보드 입력, MIDI controller, software instrument 사이의 sending과 receiving 흐름을 이해하고 기본 연결을 점검할 수 있게 한다.",
		},
		{
			Title:     "Piano Roll에서 짧은 phrase를 만든다",
			Objective: "Piano Roll이나 MIDI editor에서 note 위치, 길이, velocity를 조절해 짧은 rhythm 또는 melody phrase를 만들 수 있게 한다.",
		},
		{
			Title:     "MIDI phrase를 악기 소리로 재생한다",
			Objective: "만든 MIDI phrase를 가상 악기와 연결하고 소리, timing, 반복 상태를 점검해 제작에 쓸 수 있는 기본 phrase로 정리할 수 있게 한다.",
		},
	}
}

func refineVibeCodingMVPAppLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "AI 코딩 도구와 만들 앱의 범위를 정한다",
			Objective: "Cursor, Replit, Lovable 같은 AI 코딩 도구의 역할을 이해하고 만들 앱의 사용자, 핵심 기능, 제외할 범위를 정리할 수 있게 한다.",
		},
		{
			Title:     "아이디어를 PRD와 첫 프롬프트로 바꾼다",
			Objective: "앱 아이디어를 화면, 데이터, 사용자 흐름이 보이는 간단한 PRD로 정리하고 AI가 이해할 첫 생성 프롬프트를 만들 수 있게 한다.",
		},
		{
			Title:     "AI로 첫 작동 버전을 만든다",
			Objective: "AI 코딩 도구에 PRD를 전달해 첫 앱 버전을 만들고, 생성된 화면과 핵심 흐름이 의도와 맞는지 확인할 수 있게 한다.",
		},
		{
			Title:     "프롬프트 반복으로 기능과 오류를 고친다",
			Objective: "원하는 기능, 깨진 화면, 동작 오류를 구체적으로 설명하며 AI와 반복 수정해 앱의 핵심 흐름을 안정화할 수 있게 한다.",
		},
		{
			Title:     "작은 MVP를 완성하고 사용자 흐름을 점검한다",
			Objective: "작은 앱 또는 MVP의 시작, 입력, 결과, 저장 흐름을 직접 점검하고 남은 개선 목록을 정리할 수 있게 한다.",
		},
	}
}

func refineVibeCodingFullstackShipLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "AI coding 환경과 앱 구조를 준비한다",
			Objective: "Cursor, Replit, GitHub, Supabase, Vercel 같은 도구의 역할을 나누고 full-stack 앱의 화면, 데이터, 배포 구조를 준비할 수 있게 한다.",
		},
		{
			Title:     "프론트엔드 화면과 사용자 흐름을 만든다",
			Objective: "AI 도구로 주요 화면과 컴포넌트를 만들고 사용자가 이동하는 흐름이 끊기지 않도록 구조를 정리할 수 있게 한다.",
		},
		{
			Title:     "데이터베이스와 인증을 연결한다",
			Objective: "Supabase 같은 backend 서비스를 사용해 데이터 모델, 저장 흐름, 로그인·권한 흐름을 앱에 연결할 수 있게 한다.",
		},
		{
			Title:     "AI가 만든 코드를 검토하고 오류를 정리한다",
			Objective: "생성된 코드의 중복, 깨진 상태, 보안상 위험한 처리, 콘솔 오류를 찾아 AI와 함께 수정하고 안정화할 수 있게 한다.",
		},
		{
			Title:     "배포 직전 상태로 앱을 점검한다",
			Objective: "Vercel 같은 배포 흐름을 기준으로 환경 변수, 인증, 데이터 접근, 핵심 사용자 시나리오를 점검해 배포 직전 상태까지 정리할 수 있게 한다.",
		},
	}
}

func refineClaudeCodeAgenticWorkflowLessons(_ []generatedMainLesson) []generatedMainLesson {
	return []generatedMainLesson{
		{
			Title:     "Claude Code 작업 흐름과 안전 기준을 잡는다",
			Objective: "Claude Code가 맡을 작업 범위, 파일 읽기, 계획 수립, 코드 변경, 테스트 실행 흐름을 이해하고 안전한 작업 기준을 세울 수 있게 한다.",
		},
		{
			Title:     "context engineering으로 작업 지시를 정리한다",
			Objective: "목표, 제약, 관련 파일, 테스트 기준을 명확히 전달해 Claude Code가 작은 기능 변경을 일관되게 수행할 수 있게 한다.",
		},
		{
			Title:     "작은 기능을 agentic workflow로 구현한다",
			Objective: "계획, 코드 수정, 테스트, 오류 수정의 반복 흐름으로 작은 기능 하나를 완성하며 AI 개발 cycle을 체득할 수 있게 한다.",
		},
		{
			Title:     "hooks와 MCP로 반복 작업을 보강한다",
			Objective: "반복 검증, 외부 도구 연결, 문서 조회 같은 작업을 hooks나 MCP 개념으로 보강해 개발 흐름을 더 안정적으로 만들 수 있게 한다.",
		},
		{
			Title:     "테스트와 정리까지 포함한 개발 workflow를 완성한다",
			Objective: "작은 프로젝트에 Claude Code 기반 workflow를 적용하고 테스트, 문서 정리, 남은 위험 확인까지 마칠 수 있게 한다.",
		},
	}
}
