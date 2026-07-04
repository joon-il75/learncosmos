package curriculum

import "strings"

func containsLeathercraftContext(context string) bool {
	return containsAny(context, "가죽공예", "가죽 공예", "가죽", "leathercraft", "leather craft", "leatherworking", "leather working", "leather wallet", "leather card", "leather")
}

func containsCalligraphyContext(context string) bool {
	return containsAny(context, "캘리그라피", "손글씨", "붓펜", "레터링", "calligraphy", "lettering", "brush pen")
}

func containsWoodworkingSmallProjectContext(context string) bool {
	return containsAny(context,
		"목공예", "목공", "woodworking", "wood work", "woodcraft",
		"작은 선반", "벽 선반", "선반 만들기", "도마 만들기", "주방용 도마",
		"birdhouse", "bird house",
		"목재 재단", "나무 재단", "샌딩", "사포질", "오일 마감", "목재 마감",
	)
}

func containsKnittingBasicWearableContext(context string) bool {
	return containsAny(context,
		"대바늘", "뜨개", "뜨개질", "knitting", "knit",
		"목도리", "스카프", "니트 소품", "착용 소품",
		"겉뜨기", "안뜨기", "코 잡기", "단수", "폭 유지",
	)
}

func containsCrochetSmallDollProjectContext(context string) bool {
	return containsAny(context,
		"코바늘", "crochet", "amigurumi", "아미구루미",
		"granny square", "crochet square", "coaster", "pouch",
		"코바늘 인형", "인형 만들기", "작은 인형",
		"원형뜨기", "짧은뜨기", "증감", "솜 넣기", "부품 연결",
	)
}

func containsPotteryHandbuildingCupProjectContext(context string) bool {
	return containsAny(context,
		"도자기", "pottery", "ceramic", "ceramics",
		"머그컵", "컵 만들기", "도자기 컵", "핸드빌딩",
		"small bowl", "bowl", "glaze", "glazing",
		"흙 성형", "손잡이 붙이기", "손잡이 접합", "유약", "소성",
	)
}

func containsMIDIProductionContext(context string) bool {
	return containsAny(context,
		"미디", "midi", "daw", "bandlab", "밴드랩", "lmms", "tunepad", "soundtrap",
		"음악제작", "음악 제작", "뮤직 프로덕션", "music production", "beat making", "비트메이킹", "비트 만들기",
		"piano roll", "피아노롤", "시퀀싱", "sequencing",
	)
}

func containsInstrumentPerformanceContext(context string) bool {
	return containsAny(context,
		"기타", "통기타", "일렉기타", "일렉 기타", "베이스기타", "베이스 기타",
		"피아노", "드럼", "바이올린", "비올라", "첼로", "콘트라베이스",
		"플룻", "플루트", "색소폰", "리코더", "하모니카",
		"거문고", "가야금", "해금", "대금", "소금", "장구",
		"보컬", "악기", "연주",
		"guitar", "electric guitar", "acoustic guitar", "bass", "bass guitar", "piano", "drums", "drum",
		"violin", "viola", "cello", "flute", "saxophone", "recorder", "harmonica", "instrument", "play through",
	)
}

func containsInstrumentSongCompletionContext(context string) bool {
	return containsAny(context,
		"한 곡", "한곡", "완주", "끝까지", "좋아하는 곡",
		"곡 한 곡", "곡을 연주", "곡 연주", "곡을 끝까지",
		"클래식 곡", "록 곡", "재즈 곡", "재즈 스탠더드",
		"민요 한 곡", "전통 가락", "산조", "리프와 솔로",
		"one song", "complete a song", "finish one song", "full song", "play through", "favorite song",
		"complete one piece", "full run-through",
	)
}

func containsVibeCodingContext(context string) bool {
	return containsAny(context,
		"바이브 코딩", "바이브코딩", "vibe coding", "ai coding", "ai 코딩", "ai 개발",
		"cursor", "lovable", "replit", "bolt.new", "bolt new", "windsurf", "claude code", "클로드 코드",
		"github copilot", "copilot", "mcp", "agentic coding", "agentic development",
	)
}

func containsDanceCoverRoutineContext(context string) bool {
	return containsAny(context,
		"케이팝 댄스", "kpop dance", "k-pop dance", "댄스 커버", "dance cover",
		"안무 커버", "안무 한 곡", "안무 한곡", "한 곡 안무", "한곡 안무",
		"포인트 안무", "후렴 안무", "동작 구간", "구간 연습", "리듬과 동작",
	)
}

func containsYogaDailyRoutineContext(context string) bool {
	return containsAny(context,
		"아침 요가", "요가 루틴", "데일리 요가", "매일 요가", "15분 요가", "10분 요가",
		"요가 자세", "요가 호흡", "스트레칭 자세", "루틴을 꾸준히",
		"routine yoga", "daily yoga", "beginner yoga", "morning yoga", "yoga routine", "flexibility", "breathing",
		"home workout", "workout routine", "exercise habit", "weekly exercise", "fitness routine",
	)
}

func containsHomeStrengthRoutineContext(context string) bool {
	return containsAny(context,
		"홈트", "홈 트레이닝", "맨몸운동", "맨몸 운동", "근력운동", "근력 운동", "덤벨", "스쿼트", "푸쉬업", "푸시업", "팔굽혀펴기", "코어 운동", "근력 루틴",
		"home workout", "home strength", "bodyweight workout", "strength routine", "dumbbell workout", "push-up", "pushup", "squat", "core workout",
	)
}

func containsRunning5KBeginnerRoutineContext(context string) bool {
	return containsAny(context,
		"러닝", "달리기", "조깅", "5km", "5키로", "5킬로", "초보 러너", "러닝 루틴", "인터벌 러닝", "페이스 조절", "달리기 자세",
		"running beginner", "beginner running", "5k", "5 km", "couch to 5k", "jogging", "running plan", "running form", "interval running",
	)
}

func containsVocalSongPracticeContext(context string) bool {
	return containsAny(context,
		"보컬", "노래 연습", "노래 한 곡", "노래 커버", "커버곡", "가창",
		"vocal practice", "singing practice", "sing a song", "vocal warmup", "vocal warm-up", "cover song",
	) || (containsAny(context, "노래", "singing", "song") && containsAny(context, "발성", "음정", "pitch"))
}

func containsDailyEnglishConversationRoutineContext(context string) bool {
	return containsAny(context,
		"일상 영어", "생활 영어", "매일 영어", "영어 회화 루틴", "영어 말하기 루틴", "쉐도잉", "섀도잉", "하루 10문장", "영어 표현 연습", "상황별 영어",
		"daily english", "everyday english", "english conversation routine", "daily conversation", "shadowing", "speaking practice", "english phrases",
	)
}

func containsVideoEditingShortformProjectContext(context string) bool {
	return containsAny(context,
		"영상 편집", "영상편집", "숏폼", "쇼츠", "릴스", "캡컷", "프리미어", "다빈치 리졸브", "자막", "컷편집", "컷 편집",
		"video editing", "shortform video", "short-form video", "youtube shorts", "shorts", "reels", "capcut", "premiere pro", "davinci resolve", "subtitles", "cut editing",
	)
}

func containsAIToolProductivityWorkflowContext(context string) bool {
	return containsAny(context,
		"chatgpt", "챗gpt", "챗지피티", "생성형 ai", "ai 도구", "업무 자동화", "문서 요약", "프롬프트", "노션 ai", "구글시트 자동화", "AI 생산성",
		"ai productivity", "ai tools", "chatgpt workflow", "prompt workflow", "document summary", "spreadsheet automation", "notion ai", "work automation",
	)
}

func containsCandleSoapResinProjectContext(context string) bool {
	return containsAny(context,
		"캔들", "향초", "비누 만들기", "수제비누", "레진", "레진아트", "석고방향제", "몰드", "향료", "왁스",
		"candle making", "soy candle", "soap making", "handmade soap", "resin art", "epoxy resin", "mold", "fragrance oil",
	)
}

func containsKoreanHomeCookingBasicsContext(context string) bool {
	return containsAny(context,
		"집밥", "생활요리", "한식 기초", "한식 기본", "반찬 만들기", "밑반찬", "된장찌개", "김치찌개", "계란말이", "나물무침", "초보 요리", "자취요리",
		"korean home cooking", "korean cooking basics", "banchan", "doenjang jjigae", "kimchi jjigae", "home cooking basics",
	)
}

func containsArduinoSensorProjectContext(context string) bool {
	return containsAny(context,
		"아두이노", "arduino", "아두이노 센서", "아두이노 프로젝트", "아두이노 키트", "아두이노 LED", "아두이노 led",
		"arduino project", "arduino sensor", "arduino starter kit", "arduino led",
	)
}

func containsEnglishMOOCProgrammingContext(context string) bool {
	return containsAny(context,
		"mit ocw 6.0001", "6.0001", "introduction to computer science", "programming in python",
		"computer science and programming", "computational thinking", "problem set", "problem sets",
		"saylor cs102", "cs102", "object oriented programming", "java programming", "c++ programming",
		"freecodecamp", "free code camp", "interactive lessons", "coding curriculum", "programming curriculum",
	)
}

func containsProblemSetFinalProjectContext(context string) bool {
	return containsAny(context, "cs50", "cs50x", "harvard cs50") ||
		(containsAny(context, "final project") && containsAny(context, "problem set", "problem sets"))
}

func containsWebDevProjectFoundationContext(context string) bool {
	return containsAny(context,
		"the odin project", "odin project", "web development foundations", "foundations course",
		"recipes project", "landing page", "rock paper scissors", "etch-a-sketch", "calculator project",
		"web app mvp", "simple web app", "full-stack web app", "full stack web app", "authentication and database",
		"build and deploy", "crud app",
	)
}

func containsWebPlatformCoreCurriculumContext(context string) bool {
	return containsAny(context,
		"mdn learn web development", "learn web development", "web platform core", "html css javascript",
		"semantic html", "accessibility", "responsive design", "web tools",
	)
}

func containsFrontendCompetencyMapContext(context string) bool {
	return containsAny(context,
		"mdn curriculum", "front-end curriculum", "frontend curriculum", "front-end developer",
		"frontend developer", "competency map", "job-ready", "job ready",
	)
}

func containsCreativeCodingVisualProjectContext(context string) bool {
	return containsAny(context,
		"khan academy", "drawing and animation", "intro to js", "javascript drawing",
		"creative coding", "talk-through", "challenge", "shapes", "coloring",
	)
}

func containsGame3DGuidedPathwayContext(context string) bool {
	return containsAny(context,
		"unity learn", "unity essentials", "junior programmer", "real-time 3d", "realtime 3d",
		"game prototype", "gameplay prototype", "unity pathway",
		"3d game", "tiny 3d game", "3d game scene", "guided path", "guided player path",
		"player controller", "simple interaction",
	)
}

func containsOpenTextbookChapterPracticeContext(context string) bool {
	return containsAny(context,
		"openstax", "open textbook", "open textbooks", "textbook chapter", "chapter practice",
	)
}

func containsOERLearningObjectClusterContext(context string) bool {
	return containsAny(context,
		"libretexts", "learning object", "learning objects", "worksheet", "visualization", "simulation",
	)
}

func containsInteractiveTutorConceptPracticeContext(context string) bool {
	return containsAny(context,
		"carnegie mellon open learning initiative", "cmu oli", "oli course", "open learning initiative",
		"interactive tutor", "computer tutor", "tutor feedback", "introduction to biology open",
	)
}

func containsCloudLabSkillPathContext(context string) bool {
	return containsAny(context,
		"google cloud skills boost", "cloud skills boost", "skill badge", "hands-on lab", "hands on lab",
		"cloud lab", "cloud labs", "role-based cloud", "google cloud learning path",
	)
}

func containsEnglishMOOCHumanitiesSurveyContext(context string) bool {
	return containsAny(context,
		"introduction to western music", "humanities survey", "readings and discussion",
		"listening assignment", "listening assignments",
	) || (containsAny(context, "mit ocw", "ocw", "open courseware", "openlearn", "open university") &&
		containsAny(context, "western music", "music history", "humanities", "essay", "paper"))
}

func containsDigitalLiteracyBadgedCourseContext(context string) bool {
	return containsAny(context,
		"digital skills", "succeeding in a digital world", "digital literacy", "online identity",
		"digital safety", "information overload", "digital wellbeing", "digital well-being", "openlearn digital",
	)
}

func containsRoleSkillLearningPathContext(context string) bool {
	return containsAny(context,
		"microsoft learn", "learning path", "learning paths", "role based", "role-based",
		"career path", "career paths", "skill path", "training module", "training modules",
	)
}

func containsStructuredDrawingFoundationContext(context string) bool {
	return containsAny(context,
		"drawspace", "drawing lessons", "drawing foundation", "learn to draw", "contour drawing",
		"shading", "perspective", "portrait drawing", "color drawing",
		"watercolor", "watercolour", "landscape postcard", "simple washes", "sky depth",
		"어반스케치", "어반 스케치", "urban sketch", "urban sketching", "풍경 스케치",
		"구도", "선 표현", "관찰 드로잉", "관찰 스케치",
	)
}

func containsWatercolorBeginnerLandscapeContext(context string) bool {
	return containsAny(context,
		"수채화", "물감", "수채화 풍경", "수채화 엽서", "번짐", "워시", "그라데이션", "하늘 그리기", "나무 그리기", "풍경화", "watercolor", "watercolour", "watercolor landscape", "watercolor postcard", "wash technique", "gradient wash",
	)
}

func containsCreativeToolBeginnerOutputContext(context string) bool {
	return containsAny(context,
		"envato tuts", "tutsplus", "tuts+", "photoshop tutorial", "illustrator tutorial",
		"after effects tutorial", "premiere pro tutorial", "creative tool", "beginner tutorial",
	)
}

func containsScienceSimulationInquiryPathContext(context string) bool {
	return containsAny(context,
		"phet interactive simulations", "phet simulations", "phet simulation", "interactive science simulation",
		"interactive math simulation", "cause-effect exploration", "multiple representations", "variable experiment",
	)
}

func containsGuidedSimulationActivitySheetContext(context string) bool {
	return containsAny(context,
		"teaching with phet", "phet small group", "phet activity sheet", "activity sheet design",
		"guided simulation activity", "simulation activity sheet", "guided inquiry activity sheet",
	)
}

func containsFoodPreservationSafetyPathContext(context string) bool {
	return containsAny(context,
		"penn state extension home food preservation", "water bath canning", "atmospheric steam canning",
		"high acid foods", "acidifying tomato products", "research-tested recipes", "food preservation safety",
		"home food preservation",
		"김치 담그기", "김치", "발효 관리", "발효 보관", "배추 절이기", "양념 비율", "식품 보존", "저장 발효",
	)
}

func containsHorticultureDiagnosticFoundationContext(context string) bool {
	return containsAny(context,
		"university of minnesota extension prohort", "prohort", "horticulture education",
		"botany and horticulture", "integrated pest management", "plant pathology", "plant diagnostics",
		"horticulture diagnostics",
		"반려식물", "분갈이 시기", "흙 배합", "물주기 기준", "식물 진단", "식물 관리", "화분 관리",
	)
}

func containsComputationalMusicAnalysisProjectContext(context string) bool {
	return containsAny(context,
		"computational music theory and analysis", "21m.383", "music21", "symbolic score",
		"score-based", "corpus studies", "algorithmic composition", "computational musicology",
	)
}

func containsDesignThinkingProblemFramingPathContext(context string) bool {
	return containsAny(context,
		"design thinking bootleg", "stanford d.school design thinking", "empathize define ideate prototype test",
		"human-centered design", "human centered design", "design thinking methods", "problem framing",
	)
}

func containsNoBuildPrototypeTestContext(context string) bool {
	return containsAny(context,
		"no-build hack", "no build hack", "no-build prototyping", "no build prototyping",
		"scrappy prototype", "moveable furniture", "simple prototyping supplies", "prototype test",
	)
}

func containsClimateDataGraphingExplanationContext(context string) bool {
	return containsAny(context,
		"graphing global temperature trends", "global temperature trends", "temperature anomaly",
		"climate data graphing", "climate data graph", "climate data", "data graph", "data trend",
		"short-term trends", "long-term trends", "nasa science graphing",
	)
}

func containsCitizenScienceObservationRecordContext(context string) bool {
	return containsAny(context,
		"inaturalist training", "bioblitz", "biodiversity record", "document organisms",
		"field guide", "citizen science observation", "participatory science observation",
		"biodiversity observation", "local plants", "local plant observation", "inaturalist",
		"observation log", "biodiversity notes",
	)
}

func containsPrimarySourceInquiryNoteContext(context string) bool {
	return containsAny(context,
		"library of congress primary source", "primary source analysis tool", "observe reflect question",
		"investigate further", "primary source inquiry", "analyzing primary sources",
		"primary source analysis", "historical primary source", "inquiry note", "observe-reflect-question",
	)
}

func containsSewingProjectFoundationContext(context string) bool {
	return containsAny(context,
		"alison basics of sewing", "basics of sewing", "learn to sew", "sewing machine",
		"hand-stitching", "hand stitching", "pattern envelope", "seam allowance", "sewing project",
		"beginner sewing", "tote bag", "straight stitches", "straight stitch", "attach handles", "clean finishing",
		"재봉틀", "재봉", "봉제", "에코백", "파우치", "쿠션커버", "원단 재단", "직선 박기", "박음질",
	)
}

func containsCardboardCircuitInventionPathContext(context string) bool {
	return containsAny(context,
		"cardboard circuits", "invent with cardboard", "circuit playground", "cardboard invention",
		"electronics robotics programming", "copper tape", "alligator clips",
		"cardboard circuit", "circuit invention", "leds", "led", "switches", "working prototype",
	)
}

func containsMusicProductionFullPipelineContext(context string) bool {
	return containsAny(context,
		"mixxin academy", "music production certificate", "music production 101", "audio processing",
		"sound design", "composition and sketching", "arrangements and structure", "mixing and mastering",
	)
}

func containsPhotographyFundamentalsProjectContext(context string) bool {
	return containsAny(context,
		"icon photography school", "photography fundamentals", "camera simulator", "exposure triangle",
		"aperture shutter iso", "photography project", "photography basics",
		"smartphone photo", "smartphone photos", "mobile photo", "mobile photos", "take better photos", "better smartphone photos",
		"portrait photo", "portrait photos", "natural light", "photo composition", "photo editing",
		"사진", "스마트폰 사진", "휴대폰 사진", "폰 사진", "일상 사진", "인물 사진",
		"사진 잘 찍기", "촬영 구도", "사진 구도", "빛 활용", "자연광", "초점", "노출",
		"보정", "사진 보정", "배경 정리", "포즈 유도", "portrait photography", "mobile photography",
	)
}

func containsGardenDesignMaintenancePathContext(context string) bool {
	return containsAny(context,
		"garden tutor", "free online gardening course", "garden design", "garden installation",
		"garden maintenance", "learn to garden", "balcony herb garden", "container garden", "grow herbs",
		"herb garden", "plant care routine", "watering schedule",
		"베란다 허브", "허브 키우기", "허브 화분", "베란다 화분", "햇빛과 물주기", "물주기 루틴",
		"작은 정원", "컨테이너 정원", "화분 가꾸기",
	)
}

func containsMealPrepWeekdayLunchContext(context string) bool {
	return containsAny(context,
		"meal prep", "weekday meal", "healthy lunch", "healthy lunches", "lunch prep",
		"meal prep routine",
	) && containsAny(context, "lunch", "meal", "cook", "store", "weekday")
}

func containsHouseplantRepottingCareContext(context string) bool {
	return (containsAny(context, "houseplant", "house plant", "반려식물") &&
		containsAny(context, "repot", "repotting", "soil mix", "potting mix", "plant recovery", "monitor recovery", "분갈이", "흙 배합")) ||
		containsAny(context, "houseplant repotting", "repot a small houseplant")
}

func containsPersonalFinanceBudgetRoutineContext(context string) bool {
	return containsAny(context,
		"personal finance", "monthly budget", "budget beginner", "budgeting", "spending categories",
		"track spending", "weekly expenses", "expense review",
	)
}

func containsDogBasicTrainingRoutineContext(context string) bool {
	return strings.Contains(context, "dog") && containsAny(context,
		"basic training", "sit", "stay", "recall", "come", "positive reinforcement",
		"training sessions", "command",
	)
}

func containsNotionStudyDashboardContext(context string) bool {
	return containsAny(context,
		"notion study dashboard", "notion dashboard", "study dashboard", "notion workspace",
		"notion productivity dashboard", "productivity dashboard", "notion template",
		"task tracker", "tasks database", "task database", "notes database",
		"weekly review pages", "weekly review template", "study planner",
		"노션", "학습 대시보드", "공부 대시보드", "생산성 대시보드", "노션 템플릿", "태스크 트래커",
	)
}

func containsCookingBasicsFoundationContext(context string) bool {
	return containsAny(context,
		"alison cooking basics", "cooking basics", "food hygiene", "kitchen utensils",
		"knife skills", "stock making", "sauces", "herbs and spices",
		"홈베이킹", "케이크", "제누와즈", "아이싱", "생크림", "디저트", "반죽", "오븐", "굽기",
		"home baking", "beginner cake", "sponge cake", "whipped cream cake", "cake batter", "bake a cake",
		"simple dessert",
	)
}

func containsHomeCafeLatteFoundationContext(context string) bool {
	return containsAny(context,
		"홈카페", "home cafe", "라떼", "latte",
		"에스프레소 추출", "espresso", "우유 스티밍", "milk steaming",
		"우유 질감", "라떼아트", "커피", "coffee",
	)
}

func containsCreatorTutorialPublishContext(context string) bool {
	return containsAny(context,
		"강의 영상", "튜토리얼 영상", "설명 영상", "시연 영상", "교육 영상",
		"유튜브 강의", "촬영하고 편집", "촬영하고 게시", "편집해 게시",
		"creator tutorial", "tutorial video", "how-to video",
	) || (containsAny(context, "영상", "콘텐츠", "촬영", "편집", "유튜브", "youtube") &&
		containsAny(context, "강의", "튜토리얼", "시연", "가르치"))
}

func containsCriticalReadingToWritingDraftContext(context string) bool {
	return containsAny(context,
		"creative writing and critical reading", "critical reading", "reading as a writer",
		"fiction creative nonfiction poetry scriptwriting", "writing draft",
	)
}

func containsPhotographyStudioProjectContext(context string) bool {
	return containsAny(context,
		"introduction to photography and related media", "photography and related media",
		"darkroom techniques", "digital imaging", "studio lighting", "photography studio",
	)
}

func containsDesignLabPrototypeProjectContext(context string) bool {
	return containsAny(context,
		"d-lab ii", "d-lab: design", "d-lab design", "design lab", "design packet",
		"hands-on prototyping", "hands on prototyping", "term-long project", "community partner",
		"design problem", "simple prototype", "user feedback", "prototype for user feedback",
	)
}

func containsOrganicGrowingCyclePlanContext(context string) bool {
	return containsAny(context,
		"growing organic food sustainably", "organic food gardening", "organic gardening",
		"crop seasons", "sowing and growing", "growing cycle",
		"vegetable growing cycle", "grow vegetables", "organic growing",
		"organic vegetable", "seed to harvest", "small organic vegetable bed", "care cycle",
	)
}

func containsAdobeCreativeLearningPathContext(context string) bool {
	return containsAny(context,
		"adobe learn", "adobe learning path", "adobe creative cloud", "photoshop learning path",
		"illustrator learning path", "adobe express", "adobe creative tools",
	)
}

func containsTemplateDesignQuickOutputContext(context string) bool {
	return containsAny(context,
		"canva design school", "canva", "template design", "social media design", "brand kit",
		"quick design output", "design template",
	)
}

func containsSequentialCreativeClassPathContext(context string) bool {
	return containsAny(context,
		"skillshare learning path", "skillshare learning paths", "creative class", "creative classes",
		"illustration learning path", "project gallery", "class project", "sequential class",
	)
}

func containsAnimationFundamentalsShotProgressionContext(context string) bool {
	return containsAny(context,
		"blender studio", "animation fundamentals", "bouncing ball", "body mechanics",
		"character animation", "walk cycle", "pantomime", "animation shot",
		"beginner blender animation", "blender animation", "keyframes", "timing and spacing",
	)
}

func containsInteractiveMusicProductionFoundationContext(context string) bool {
	return containsAny(context,
		"ableton learning music", "learning music ableton", "music production foundation",
		"beat loop", "song section", "bassline", "melody", "browser music production",
	)
}

func containsSongDrivenInstrumentPathContext(context string) bool {
	if containsAny(context,
		"fender play", "song-driven", "song driven", "guitar lessons", "guitar path",
		"riff", "chords", "song practice", "instrument path",
	) {
		return true
	}
	if !containsInstrumentPerformanceContext(context) {
		return false
	}
	return containsAny(context,
		"곡 카피", "원곡 카피", "원곡에 맞춰", "원곡 느낌", "커버",
		"리프", "솔로", "베이스 라인", "베이스라인", "그루브", "필인",
		"팝송", "밴드곡", "반주", "피아노 반주", "팝송 반주", "코드 반주", "왼손 반주", "반주 패턴", "코드 진행",
		"song cover", "cover song", "favorite song", "play through", "riff", "solo", "bass line",
		"groove", "fill", "accompaniment", "accompaniment pattern", "chord accompaniment", "left hand accompaniment", "chord progression",
	)
}

func containsMakerProjectClassContext(context string) bool {
	return containsAny(context,
		"instructables", "instructables classes", "maker project", "diy project", "classroom project",
		"electronics project", "woodworking project", "project class",
	)
}

func containsShortDiplomaAssessmentPathContext(context string) bool {
	return containsAny(context,
		"alison diploma", "alison certificate", "diploma course", "short diploma",
		"learner record", "self-paced module", "self paced module", "career skills assessment",
	)
}

func containsShortCourseWeeklyDiscussionContext(context string) bool {
	return containsAny(context,
		"futurelearn", "futurelearn short course", "short course", "weekly activities",
		"weekly course", "course discussion", "discussion activity", "case study discussion",
	)
}

func containsAuditCoursePracticePathContext(context string) bool {
	return containsAny(context,
		"edx free audit", "edx audit", "audit course", "audit track", "free audit course",
		"ungraded practice", "verified upgrade",
	)
}

func containsPublicLifelongOnlineCourseContext(context string) bool {
	return containsAny(context, "gseek 온라인학습", "gseek online", "경기도 평생학습포털 지식 온라인", "지식 온라인학습") ||
		(containsAny(context, "gseek", "경기도 평생학습포털", "지식") &&
			containsAny(context, "온라인학습", "온라인 강좌", "평생학습") &&
			containsAny(context, "취미", "건강", "생활상식", "디지털역량", "자격증"))
}

func containsBlendedLifelongParticipationContext(context string) bool {
	return containsAny(context, "gseek 오프라인", "gseek 화상학습", "지식 오프라인", "지식 화상학습") ||
		(containsAny(context, "경기도 평생학습포털", "gseek", "평생학습") &&
			containsAny(context, "오프라인", "화상학습", "화상 학습", "지역 강좌") &&
			containsAny(context, "참여", "수강", "학습이력"))
}

func containsMidlifeTransitionLearningSupportContext(context string) bool {
	return containsAny(context, "서울런4050 직업교육경비", "서울런 4050 직업교육경비", "중장년 직업교육경비") ||
		(containsAny(context, "서울런4050", "서울런 4050", "중장년") &&
			containsAny(context, "직업역량", "직업 전환", "직업전환", "취창업", "취·창업", "학습경비"))
}

func containsCompletionRefundOnlineCourseContext(context string) bool {
	return containsAny(context, "서울런4050 중장년 특화 온라인", "서울런 4050 중장년 특화 온라인") ||
		(containsAny(context, "서울런4050", "서울런 4050", "중장년") &&
			containsAny(context, "온라인과정", "온라인 과정", "진도율", "환급", "수료"))
}

func containsFederatedLifelongOpenResourceContext(context string) bool {
	return containsAny(context, "국가평생학습포털 늘배움", "늘배움", "lifelongedu") ||
		(containsAny(context, "평생교육정보", "평생학습 포털", "공개 평생교육") &&
			containsAny(context, "민·관·학", "민관학", "강좌 묶음", "학습경로"))
}

func containsKoreanContentCreationPipelineContext(context string) bool {
	return containsAny(context, "한국콘텐츠아카데미", "kocca", "콘텐츠아카데미") ||
		(containsAny(context, "방송영상", "게임", "만화", "애니메이션", "캐릭터", "공연") &&
			containsAny(context, "콘텐츠 제작", "콘텐츠 분야", "콘텐츠 포트폴리오"))
}

func containsKoreanDesignToolPracticalCertificationContext(context string) bool {
	return containsAny(context, "gtq", "그래픽기술자격", "gtq포토샵", "gtq 포토샵", "포토샵 자격증", "포토샵 2급", "photoshop certification") ||
		(containsAny(context, "포토샵", "photoshop", "웹디자인기능사", "컴퓨터그래픽스운용기능사") &&
			containsAny(context, "자격증", "자격", "취득", "시험", "실기", "모의고사", "기출"))
}

func containsKoreanNCSWebPublisherTrainingContext(context string) bool {
	return containsAny(context, "ncs 디지털디자인", "ncs 디지털 디자인", "hrd-net 웹퍼블리셔", "work24 웹퍼블리셔", "국민내일배움카드 웹퍼블리셔") ||
		(containsAny(context, "웹퍼블리셔", "웹 퍼블리셔", "프론트엔드", "frontend", "반응형 웹", "ui/ux", "uiux") &&
			containsAny(context, "ncs", "hrd-net", "work24", "국민내일배움카드", "직업훈련", "취업", "포트폴리오"))
}

func containsRuralReturnAgricultureFoundationContext(context string) bool {
	return containsAny(context, "농업교육포털 귀농귀촌", "귀농귀촌 온라인 교육", "귀농 귀촌 온라인 교육") ||
		(containsAny(context, "귀농", "귀촌") &&
			containsAny(context, "농업", "품목기술", "정착", "교육시간", "수료"))
}

func containsUrbanAgriculturePracticeSeriesContext(context string) bool {
	return containsAny(context, "도시농업교육", "도시농업 교육", "인천광역시농업기술센터 도시농업") ||
		(containsAny(context, "도시농업", "텃밭", "가정원예", "원예식물") &&
			containsAny(context, "재배", "실습", "관리", "화분"))
}

func containsAIDigitalLiteracyLifePracticeContext(context string) bool {
	return containsAny(context, "ai·디지털배움터", "ai 디지털배움터", "디지털배움터 생활밀착형") ||
		(containsAny(context, "스마트폰", "키오스크", "보이스피싱", "생활앱", "생활 앱") &&
			containsAny(context, "디지털 문해", "디지털배움터", "ai", "디지털 역량")) ||
		(containsAny(context, "ai tools", "ai assistant", "ai assistants", "artificial intelligence") &&
			containsAny(context, "digital literacy", "safe", "safely", "privacy", "daily life", "email", "search", "planning"))
}

func containsLocalHobbyWorkshopSeriesContext(context string) bool {
	return containsAny(context, "계룡시 평생학습포털", "계룡 평생학습", "생활취미 강좌") ||
		(containsAny(context, "지역 평생학습", "평생학습센터", "지역센터") &&
			containsAny(context, "디지털 드로잉", "바이올린 초급", "가죽공예 입문", "생활취미"))
}

func containsCreditBankPracticumContext(context string) bool {
	return containsAny(context,
		"사회복지현장실습", "보육실습", "평생교육실습", "교육실습", "한국어교육실습",
		"실습세미나", "현장실습", "기관실습",
	)
}

func containsKMOOCOnlineCourseContext(context string) bool {
	return containsAny(context,
		"k-mooc", "kmooc", "케이무크", "한국형 온라인 공개강좌", "온라인 공개강좌",
		"k-mooc 강좌", "kmooc 강좌", "공개강좌형", "공개 강좌형",
		"소셜 네트워크 분석", "비정형 데이터 분석", "대학생활 101", "외국인 학습자를 위한 대학생활",
		"클래식 친구 만들기", "클래식 음악 감상", "3d 모션캡쳐", "3d 모션 캡쳐", "모션캡쳐", "모션 캡쳐",
		"ai와 윤리", "ai 와 윤리", "인공지능 윤리",
	)
}

func containsKoreanOpenLectureWeeklySurveyContext(context string) bool {
	return containsAny(context, "kocw", "한국교육학술정보원 공개강의", "대학 공개강의", "공개 강의계획서") ||
		(containsAny(context, "k-mooc", "kmooc", "케이무크", "한국형 온라인 공개강좌", "온라인 공개강좌") &&
			containsAny(context, "강의계획서", "강좌운영", "학습목표", "중간고사", "기말고사")) ||
		(containsAny(context, "강의계획서", "주차별 강의", "주차별 학습", "syllabus") &&
			containsAny(context, "공개강의", "공개 강의", "대학 강의", "온라인 공개강좌"))
}

func containsUniversitySurveyContext(context string) bool {
	return containsAny(context,
		"현대미술사", "미술사", "건축사", "로마건축", "사회학", "비평", "비평입문",
		"현대시", "문학사", "survey course", "history of art", "art since", "roman architecture",
	)
}

func containsUniversityStudioWritingContext(context string) bool {
	return containsAny(context,
		"희곡", "희곡쓰기", "플레이라이팅", "playwriting", "극작", "시나리오 워크숍",
	)
}

func containsUniversityStudioVisualContext(context string) bool {
	return containsAny(context,
		"디자인 스튜디오", "design studio", "principles of design", "디자인 원리", "스튜디오 수업",
	)
}

func containsUniversityStudioCreativeTechContext(context string) bool {
	return containsAny(context,
		"사운드 디자인", "sound design", "작곡", "입문 작곡", "musical composition",
		"생성 음악", "algorithmic", "generative music", "music and technology",
	)
}

func containsCreditBankStandardTheoryContext(context string) bool {
	return containsAny(context,
		"인간행동과사회환경", "정신건강론", "사회복지실천론", "사회복지실천기술론",
		"사회복지조사론", "사회복지행정론", "사회복지정책론", "사회복지법제",
		"사회복지정책론", "사회복지행정론", "인간행동", "정신건강",
	)
}

func containsCreditBankInstructionalDesignContext(context string) bool {
	return containsAny(context,
		"영유아교수방법론", "교수방법론", "교수설계", "교육공학", "외국어로서의한국어능력평가론",
		"능력평가론", "평가론", "교안", "활동계획안", "평가문항", "수업설계",
	)
}

func containsLanguageCertificationContext(context string) bool {
	if containsEnglishSpeakingInterviewCertificationContext(context) {
		return true
	}
	if containsAny(context, "jlpt", "toeic", "toefl", "opic", "ielts", "teps", "n1", "n2", "n3", "n4", "n5") {
		return true
	}
	return containsAny(context, "영어", "일본어", "언어") &&
		containsAny(context, "시험", "인증", "자격", "자격증", "취득", "test", "certification", "interview")
}

func containsEnglishScoreCertificationContext(context string) bool {
	return containsAny(context, "toeic", "teps")
}

func containsEnglishSpeakingInterviewCertificationContext(context string) bool {
	return containsAny(context,
		"opic", "speaking interview test", "english speaking interview", "speaking interview",
		"spoken answers", "interview test", "mock interview",
	)
}

func containsEnglishIntegratedFourSkillsCertificationContext(context string) bool {
	return containsAny(context, "ielts", "toefl")
}

func containsCloudCertificationContext(context string) bool {
	return containsAny(context,
		"aws", "cloud practitioner", "solutions architect", "ncp", "naver cloud", "google cloud", "gcp", "azure",
		"az-900", "az 900", "az-104", "az 104", "associate cloud engineer", "digital leader", "certification", "certified",
	)
}

func containsTheoryCertificationContext(context string) bool {
	return containsAny(context,
		"시험", "자격", "자격증", "취득", "기사", "산업기사", "기능사", "직업상담사", "정보처리기사", "컴퓨터활용능력", "컴활", "대기환경기사",
		"건축기사", "건축산업기사", "토목기사", "전기기사", "전기산업기사", "전기기능사", "전기공사산업기사",
		"certification", "certified", "certificate", "exam",
		"counseling certification", "counselor certification", "environmental engineering certification",
		"environmental engineer exam", "office spreadsheet certification", "spreadsheet certification", "excel certification",
	)
}

func containsPracticalCertificationContext(context string) bool {
	return containsAny(context,
		"드론", "초경량비행장치", "지게차", "용접", "피복아크", "자동차정비", "정비기능사",
		"gtq", "그래픽기술자격", "웹디자인기능사", "컴퓨터그래픽스운용기능사",
		"제과기능사", "제빵기능사", "한식조리기능사", "조리기능사",
		"drone license", "drone operator", "beginner drone license", "basic drone operation",
	) || containsBeautyServicePracticalContext(context)
}

func containsBeautyServicePracticalContext(context string) bool {
	return containsAny(context,
		"미용사", "메이크업", "피부관리", "피부 관리", "피부미용", "피부 미용",
		"네일아트", "네일 아트", "네일미용", "네일 미용", "손톱관리", "손톱 관리",
		"헤어미용", "헤어 미용", "헤어디자인", "헤어 디자인", "이발", "커트 실기",
	)
}

func containsAny(text string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(text, strings.ToLower(term)) {
			return true
		}
	}
	return false
}
