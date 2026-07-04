import type { Locale } from '@/lib/i18n/locales'

export type LandingPageCopy = {
  locale: Locale
  homeHref: string
  loginHref: string
  navLinks: Array<{ name: string; href: string }>
  alphaNotice: {
    eyebrow: string
    title: string
    bodyLines: string[]
    applicationHref: string
    applicationLabel: string
    videoLabel: string
    confirm: string
  }
  introVideo: {
    eyebrow: string
    title: string
    descriptionLines: string[]
    lumiSpeechLines: string[]
    videoTitle: string
    placeholderLabel: string
    simulationFlowLabel: string
    youtubeEmbedUrl: string
    alphaCta: string
  }
  hero: {
    lumiMessages: string[]
    descriptionLines: string[]
    inputLabel: string
    placeholder: string
    submitLabel: string
    inputHelp: string
    chips: string[]
    lumiName: string
    onlineLabel: string
  }
  lumiScrollGuide: {
    eyebrow: string
    title: string
    steps: Array<{
      badge: string
      title: string
      description: string
      demoTitle: string
      demoItems: string[]
      demoAction: string
    }>
  }
  features: {
    eyebrow: string
    title: string
    description: string
    items: Array<{ img: string; title: string; desc: string; delay: string }>
  }
  howItWorks: {
    eyebrow: string
    title: string
    descriptionLines: string[]
    steps: Array<{ step: string; title: string; desc: string }>
  }
  aiUsage: {
    eyebrow: string
    titlePrefix: string
    titleHighlight: string
    descriptionLines: string[]
    byokBadge: string
    byokTitle: string
    byokDescriptionLines: string[]
    byokBullets: string[]
    pointsTitle: string
    pointsDescription: (welcomePoints: number) => string[]
    pointExampleLabel: string
    pointExamples: Array<{ key: 'course' | 'lesson'; label: string; icon: string }>
    pointBullets: string[]
    sensitiveNoticeLines: string[]
    futureProvidersLabel: string
  }
  categories: {
    eyebrow: string
    title: string
    descriptionLines: string[]
    items: Array<{ name: string; icon: string; query: string; direct?: boolean }>
  }
  policies: {
    eyebrow: string
    title: string
    descriptionLines: string[]
    docs: Array<{ href: string; badge: string; title: string; description: string }>
    linkLabel: string
    noteLines: string[]
  }
  socialLogin: {
    eyebrow: string
    title: string
    loggedInMessage: string
    loggedInHighlight: string
    dashboardCta: string
    loggedInHint: string
    signupDescriptionLines: string[]
    socialDescriptionLines: string[]
    aiNoticeLines: string[]
    providers: {
      google: string
      kakao: string
      naver: string
    }
    policyNotice: {
      before: string
      terms: string
      middle: string
      privacy: string
      after: string
    }
    naverNotice: {
      title: string
      bodyLines: string[]
      confirm: string
    }
  }
  footer: {
    taglineLines: string[]
    columns: Array<{ title: string; links: Array<{ label: string; href?: string }> }>
  }
}

const sharedFeatureImages = [
  '/images/ai_roadmap_gold.webp',
  '/images/search_gold_network.webp',
  '/images/streak_gold.webp',
  '/images/book_gold_glow.webp',
]

const ko: LandingPageCopy = {
  locale: 'ko',
  homeHref: '/',
  loginHref: '/login',
  navLinks: [
    { name: '홈', href: '/' },
    { name: '개인학습', href: '/dashboard' },
    { name: '커뮤니티', href: '/dashboard/community' },
    { name: '크리에이터', href: '/dashboard/creator' },
    { name: '플랫폼', href: '/platform' },
  ],
  alphaNotice: {
    eyebrow: 'Alpha Test',
    title: 'LearnCosmos는 현재 테스트 운영 중입니다',
    bodyLines: [
      '이 서비스는 알파테스트 단계의 학습 도구입니다.',
      '일부 기능, 화면, 추천 결과, AI 사용 정책은 운영 안정화 과정에서 변경될 수 있습니다.',
      '테스트 중 발견한 문제나 의견은 운영자에게 알려 주세요.',
    ],
    applicationHref: '/alpha',
    applicationLabel: '알파 테스터 신청하기',
    videoLabel: '소개 영상 보기',
    confirm: '확인하고 둘러보기',
  },
  introVideo: {
    eyebrow: 'Learning Guide',
    title: '학습 안내',
    descriptionLines: [
      '목표 입력에서 목표 채팅, 코스 생성, 탐험계획, Explorer Diary까지 이어지는 흐름을 영상으로 확인해 보세요.',
      '알파 테스트에서는 이 여정을 직접 체험하고 피드백을 남길 수 있습니다.',
    ],
    lumiSpeechLines: [
      '이제 실제 학습 흐름을 영상으로 보여드릴게요.',
      '주제를 입력하고, 목표를 정리한 뒤',
      '코스를 만들고 학습을 진행하는 과정을',
      '한 번에 따라가 볼 수 있어요.',
    ],
    videoTitle: 'LearnCosmos 학습 흐름 안내',
    placeholderLabel: '영상 준비 중 · 학습 안내 슬롯',
    simulationFlowLabel: '주제 입력 -> 목표 정리 -> 코스 생성 -> 학습 진행',
    youtubeEmbedUrl: 'https://www.youtube-nocookie.com/embed/76bq5ixUzmE',
    alphaCta: '알파 테스터 2차 신청하기',
  },
  hero: {
    lumiMessages: [
      '안녕하세요, 저는 Lumi예요.\n배우고 싶은 주제를 편하게 적어 주세요.',
      '제가 목표를 함께 정리하고,',
      '첫 번째 학습탐험 경로를 만들어 드릴게요.',
      '오늘은 어떤 배움부터 시작해 볼까요?',
    ],
    descriptionLines: ['목표 정리부터 콘텐츠 탐색, 기록과 완료까지', '끝까지 배우도록 곁에서 도와드립니다.'],
    inputLabel: '✦ 무엇을 배우고 싶으세요?',
    placeholder: '예: 목공예 간단한 수납함 만들기',
    submitLabel: '나만의 코스 만들기',
    inputHelp: '입력하면 학습 경로를 만들고, 오늘 할 일을 안내합니다.',
    chips: ['통기타 입문', '수채화 기초', '코바늘 기초'],
    lumiName: 'Lumi · 루미',
    onlineLabel: '온라인',
  },
  lumiScrollGuide: {
    eyebrow: 'Lumi Guide',
    title: '루미가 학습탐험을 함께 안내해요',
    steps: [
      {
        badge: '01',
        title: '배우고 싶은 것을 입력하세요',
        description: '루미가 첫 주제를 받아 학습탐험의 출발점을 열어 줍니다.',
        demoTitle: '기타 입문',
        demoItems: ['주제를 입력하면 바로 시작점이 열려요', '목표 채팅으로 이어집니다'],
        demoAction: '배우고 싶은 것을 입력해 보세요',
      },
      {
        badge: '02',
        title: '왜 배우고 싶은지 함께 정리해요',
        description: '목표와 이유를 묻고 답하며 지금 필요한 방향을 또렷하게 만듭니다.',
        demoTitle: '목표 채팅',
        demoItems: ['왜 배우고 싶은가요?', '3개월 안에 좋아하는 곡을 연주하고 싶어요', '목표가 더 선명해졌어요'],
        demoAction: '목표 확정',
      },
      {
        badge: '03',
        title: '경로를 만들고 콘텐츠를 찾아요',
        description: '지역과 지점을 만들고, 필요한 콘텐츠를 살펴보며 경로를 조정할 수 있습니다.',
        demoTitle: '학습탐험 경로',
        demoItems: ['기초 코드 익히기', '추천 콘텐츠 2개', '필요하면 지점을 편집해요'],
        demoAction: '콘텐츠 모아보기',
      },
      {
        badge: '04',
        title: '기록하며 끝까지 진행해요',
        description: '배운 내용과 진행 상황을 남기며 나만의 행성을 완성해 갑니다.',
        demoTitle: '오늘의 진행',
        demoItems: ['학습 대상 확인', '기록 2개 저장', '완료까지 68%'],
        demoAction: '기록/완료',
      },
    ],
  },
  features: {
    eyebrow: '핵심 기능',
    title: '시작은 빠르게, 학습은 끝까지',
    description: '배우고 싶은 주제를 입력하면 목표 정리부터 콘텐츠 탐색, 기록과 완료까지 이어집니다.',
    items: [
      { img: sharedFeatureImages[0], title: 'AI 학습탐험 경로 생성', desc: '주제를 입력하면 루미가 목표 정리를 돕고, 바로 시작할 수 있는 학습탐험 경로를 만들어 줍니다.', delay: '0ms' },
      { img: sharedFeatureImages[1], title: '콘텐츠 모아보기', desc: 'YouTube, 블로그, 외부 링크를 한 흐름 안에 모아 학습할 수 있습니다.', delay: '150ms' },
      { img: sharedFeatureImages[2], title: '진도와 연속 학습', desc: '오늘 배운 것을 남기고, 꾸준히 이어갈 수 있도록 진행 상황을 보여줍니다.', delay: '300ms' },
      { img: sharedFeatureImages[3], title: '완료와 성장 기록', desc: '완료한 학습이 기록으로 남고, 나의 성장을 확인할 수 있습니다.', delay: '450ms' },
    ],
  },
  howItWorks: {
    eyebrow: '이용 방법',
    title: '4단계로 시작하는 학습탐험',
    descriptionLines: ['배우고 싶은 주제를 입력하고, 목표를 정리한 뒤, 바로 학습을 시작하세요.', '학습 중에는 콘텐츠와 기록을 더해가며 나만의 경로를 완성할 수 있습니다.'],
    steps: [
      { step: '01', title: '배우고 싶은 것 입력', desc: '기타, 수채화, 코딩처럼 배우고 싶은 주제를 편하게 적어 주세요.' },
      { step: '02', title: '시작 방식 선택', desc: '바로 시작하거나, 루미와 목표를 먼저 정리한 뒤 시작할 수 있습니다.' },
      { step: '03', title: '학습탐험 경로 생성', desc: '정리된 목표를 바탕으로 지금 따라갈 수 있는 학습탐험 경로가 만들어집니다.' },
      { step: '04', title: '배우며 완성', desc: '학습하면서 콘텐츠를 추가하고 기록을 남기며, 끝까지 완료해 갑니다.' },
    ],
  },
  aiUsage: {
    eyebrow: 'AI 사용 방식',
    titlePrefix: 'AI 모델은 선택하고,',
    titleHighlight: '추천은 설치형 임베딩으로',
    descriptionLines: ['콘텐츠 검색과 추천은 서비스 서버의 설치형 임베딩 모델을 기반으로 처리합니다.', '목표 채팅, 코스 생성, 학습 코칭 같은 생성형 AI 기능은 사용자가 연결한 AI 키로 단계적으로 선택할 수 있게 준비하고 있습니다.'],
    byokBadge: 'BYOK',
    byokTitle: '내 AI 키 연결하기',
    byokDescriptionLines: ['OpenAI를 시작으로 Claude, Gemini 등 다양한 LLM 제공자를 단계적으로 연결할 예정입니다.', 'API 키는 민감 정보로 다루며, 생성형 AI 호출에 필요한 순간에만 사용합니다.'],
    byokBullets: ['생성형 AI 기능에 개인 키 사용', 'API 키 변경·삭제 가능', '사용량과 호출 기록 확인 방향', '보안 운영 기준에 따라 민감 정보로 관리'],
    pointsTitle: 'AI 포인트로 먼저 시작하기',
    pointsDescription: (welcomePoints) => ['API 키가 없어도 괜찮습니다.', `가입 시 지급되는 ${welcomePoints} AI 포인트로 먼저 학습탐험을 시작할 수 있어요.`],
    pointExampleLabel: '포인트 사용 예시',
    pointExamples: [
      { key: 'course', label: '학습탐험 경로 생성', icon: '🎯' },
      { key: 'lesson', label: '콘텐츠 추천', icon: '💡' },
    ],
    pointBullets: ['가입 시 AI 포인트 지급', 'API 키 없이도 먼저 시작 가능', '나중에 BYOK LLM으로 전환 가능'],
    sensitiveNoticeLines: ['API 키는 사용자의 비용과 권한이 연결된 민감 정보입니다.', '콘텐츠 추천과 유사 자료 검색은 별도의 설치형 임베딩 기반으로 처리해 BYOK LLM 사용 범위를 명확히 분리합니다.'],
    futureProvidersLabel: '단계적 연결 예정 AI 제공자',
  },
  categories: {
    eyebrow: '카테고리',
    title: '무엇이든 학습탐험으로 시작할 수 있어요',
    descriptionLines: ['관심 있는 분야를 골라 첫 경로를 열어 보세요.', '목록에 없어도 배우고 싶은 주제를 직접 입력할 수 있습니다.'],
    items: [
      { name: '기타 입문', icon: '🎸', query: '기타 입문' },
      { name: '수채화 기초', icon: '🎨', query: '수채화 기초' },
      { name: '사진 촬영', icon: '📷', query: '사진 촬영' },
      { name: '코딩 배우기', icon: '💻', query: '코딩 배우기' },
      { name: '홈카페', icon: '☕', query: '홈카페' },
      { name: '운동·건강', icon: '🧘', query: '운동·건강' },
      { name: '글쓰기', icon: '📝', query: '글쓰기' },
      { name: '외국어', icon: '🌍', query: '외국어' },
      { name: '요리', icon: '🍳', query: '요리' },
      { name: '자격증 공부', icon: '📚', query: '자격증 공부' },
      { name: '직접 입력하기', icon: '✨', query: '', direct: true },
    ],
  },
  policies: {
    eyebrow: '공개 문서',
    title: '가입 전에도 운영 기준을 확인할 수 있어요',
    descriptionLines: ['LearnCosmos는 서비스 이용약관과 개인정보처리방침을 공개 문서로 제공합니다.', '가입 전에도 설치형 임베딩 기반 추천, BYOK LLM 운영 기준, 개인정보 처리 원칙을 확인할 수 있습니다.'],
    docs: [
      { href: '/terms', badge: 'Terms', title: '서비스 이용약관', description: '서비스 이용 조건, AI 포인트, BYOK LLM, 콘텐츠 이용 기준을 확인할 수 있습니다.' },
      { href: '/privacy', badge: 'Privacy', title: '개인정보처리방침', description: '개인정보 수집·이용, 보관, 보호 조치, API 키 등 민감 설정 정보 처리 기준을 확인할 수 있습니다.' },
      { href: '/open-source', badge: 'Open Source', title: '오픈소스 라이선스', description: '서비스에 포함된 주요 오픈소스 소프트웨어와 라이선스 고지 기준을 확인할 수 있습니다.' },
    ],
    linkLabel: '문서 보기',
    noteLines: ['AI 포인트와 BYOK는 사용자의 비용과 권한이 연결될 수 있는 기능입니다.', '검색·추천용 설치형 임베딩과 생성형 AI용 BYOK LLM은 사용 범위를 분리해 안내합니다.'],
  },
  socialLogin: {
    eyebrow: '지금 시작하기',
    title: '첫 학습탐험, 무료로 시작하세요',
    loggedInMessage: '반가워요! 준비된 학습이 기다리고 있어요.',
    loggedInHighlight: '',
    dashboardCta: '나만의 대시보드로 이동하기',
    loggedInHint: '현재 Login 상태입니다. 학습 설정을 변경하려면 대시보드를 방문하세요.',
    signupDescriptionLines: ['가입하면 AI 포인트가 지급됩니다.', '배우고 싶은 주제를 입력하고, 바로 첫 학습탐험 경로를 만들어 보세요.'],
    socialDescriptionLines: ['LearnCosmos 계정이 없어도 괜찮습니다.', '소셜 계정으로 시작하면 바로 회원가입으로 이어집니다.'],
    aiNoticeLines: ['API 키가 없어도 먼저 시작할 수 있고,', '나중에 내 AI 키를 연결해 사용할 수도 있습니다.'],
    providers: {
      google: 'Continue with Google',
      kakao: 'Continue with Kakao',
      naver: 'Continue with Naver',
    },
    policyNotice: {
      before: '가입 후 서비스 이용 전에 필수 약관 동의가 진행됩니다.',
      terms: '이용약관',
      middle: '과',
      privacy: '개인정보처리방침',
      after: '을 확인할 수 있습니다.',
    },
    naverNotice: {
      title: 'Naver login is not ready',
      bodyLines: ['현재 개발 중인 기능이라 Naver Login은 아직 사용할 수 없습니다.', 'Google 또는 Kakao Login을 이용해 주세요.'],
      confirm: 'OK',
    },
  },
  footer: {
    taglineLines: ['배우고 싶은 것을 입력하면,', '학습탐험 경로가 바로 열립니다.'],
    columns: [
      { title: '서비스', links: [{ label: '코스생성', href: '#course-create-ready' }, { label: '학습안내', href: '#intro-video-ready' }, { label: '운영기준', href: '#public-documents-ready' }, { label: '새소식(공지사항)', href: '#platform-news-ready' }] },
      { title: '연결', links: [{ label: '블로그', href: 'https://blog.naver.com/learnweaver' }, { label: 'YouTube', href: 'https://www.youtube.com/@햇볕냥이-g9d' }] },
      { title: '지원', links: [{ label: '도움말', href: '#intro-video-ready' }, { label: '문의하기: learnweavr@gmail.com' }, { label: '개인정보처리방침', href: '/privacy' }, { label: '이용약관', href: '/terms' }, { label: '오픈소스 라이선스', href: '/open-source' }] },
    ],
  },
}

const en: LandingPageCopy = {
  ...ko,
  locale: 'en',
  homeHref: '/en',
  loginHref: '/en/login',
  navLinks: [
    { name: 'Home', href: '/en' },
    { name: 'Personal Learning', href: '/dashboard' },
    { name: 'Community', href: '/dashboard/community' },
    { name: 'Creator', href: '/dashboard/creator' },
    { name: 'Platform', href: '/platform' },
  ],
  alphaNotice: {
    eyebrow: 'Alpha Test',
    title: 'LearnCosmos is currently in alpha testing',
    bodyLines: [
      'This service is a learning tool under active alpha testing.',
      'Some features, screens, recommendations, and AI usage policies may change as we stabilize the service.',
      'If you find an issue or have feedback during testing, please contact the operator.',
    ],
    applicationHref: '/alpha',
    applicationLabel: 'Apply for Alpha Test',
    videoLabel: 'Watch Intro Video',
    confirm: 'Continue',
  },
  introVideo: {
    eyebrow: 'Learning Guide',
    title: 'Learning Guide',
    descriptionLines: [
      'Watch how LearnCosmos connects goal input, goal chat, course creation, exploration planning, and Explorer Diary into one guided learning flow.',
      'Alpha testers can experience this flow directly and share feedback.',
    ],
    lumiSpeechLines: [
      'Now I will show you the real learning flow in a video.',
      'You will enter a topic, refine your goal,',
      'create a course, and move through learning',
      'all in one guided simulation.',
    ],
    videoTitle: 'LearnCosmos Learning Flow Guide',
    placeholderLabel: 'Video coming soon · Learning guide slot',
    simulationFlowLabel: 'Topic input -> Goal setup -> Course creation -> Learning flow',
    youtubeEmbedUrl: 'https://www.youtube-nocookie.com/embed/76bq5ixUzmE',
    alphaCta: 'Apply for Alpha Test',
  },
  hero: {
    lumiMessages: [
      'Hi, I am Lumi.\nTell me what you want to learn.',
      'I will help clarify your goal,',
      'then build your first learning journey.',
      'What would you like to start learning today?',
    ],
    descriptionLines: ['From goal setup to content discovery, records, and completion,', 'LearnCosmos helps you keep learning through the end.'],
    inputLabel: '✦ What learning journey do you want to start?',
    placeholder: 'Example: Master guitar chords in three months',
    submitLabel: 'Start Learning →',
    inputHelp: 'Enter a topic and your learning journey opens.',
    chips: ['🎸 Beginner Guitar', '🎨 Watercolor Basics', '☕ Home Cafe Latte Art'],
    lumiName: 'Lumi',
    onlineLabel: 'Online',
  },
  lumiScrollGuide: {
    eyebrow: 'Lumi Guide',
    title: 'Lumi walks through the learning journey with you',
    steps: [
      {
        badge: '01',
        title: 'Enter what you want to learn',
        description: 'Lumi turns your first topic into the starting point of a learning journey.',
        demoTitle: 'Beginner Guitar',
        demoItems: ['A topic opens your starting point', 'Then the goal chat begins'],
        demoAction: 'Start Learning',
      },
      {
        badge: '02',
        title: 'Clarify why it matters',
        description: 'A short goal chat helps shape the purpose and direction of your journey.',
        demoTitle: 'Goal Chat',
        demoItems: ['Why do you want to learn this?', 'I want to play favorite songs in three months', 'Your goal is clearer now'],
        demoAction: 'Confirm Goal',
      },
      {
        badge: '03',
        title: 'Build the path and explore content',
        description: 'Create regions and points, review useful content, and adjust the journey as needed.',
        demoTitle: 'Learning Journey',
        demoItems: ['Learn basic chords', '2 recommended materials', 'Edit points when needed'],
        demoAction: 'Collect Content',
      },
      {
        badge: '04',
        title: 'Learn, record, and keep going',
        description: 'Your progress, notes, and completion records help your planet take shape.',
        demoTitle: 'Today’s Progress',
        demoItems: ['Review learning material', '2 records saved', '68% to completion'],
        demoAction: 'Record/Complete',
      },
    ],
  },
  features: {
    eyebrow: 'Core Features',
    title: 'Start quickly, keep learning to the end',
    description: 'Enter a topic and move from goal setup to content discovery, records, and completion.',
    items: [
      { img: sharedFeatureImages[0], title: 'AI Learning Journey Builder', desc: 'Lumi helps clarify your goal and creates a learning journey you can start right away.', delay: '0ms' },
      { img: sharedFeatureImages[1], title: 'Content in One Flow', desc: 'Bring YouTube, blogs, and external links into one learning path.', delay: '150ms' },
      { img: sharedFeatureImages[2], title: 'Progress and Streaks', desc: 'Record what you learned today and see your progress as you continue.', delay: '300ms' },
      { img: sharedFeatureImages[3], title: 'Completion and Growth Records', desc: 'Completed learning stays in your records so you can review your growth.', delay: '450ms' },
    ],
  },
  howItWorks: {
    eyebrow: 'How It Works',
    title: 'Start a learning journey in four steps',
    descriptionLines: ['Enter what you want to learn, refine your goal, and begin right away.', 'As you learn, you can add content and records to complete your own path.'],
    steps: [
      { step: '01', title: 'Enter what you want to learn', desc: 'Write a topic naturally, such as guitar, watercolor, or coding.' },
      { step: '02', title: 'Choose how to start', desc: 'Start right away, or organize your goal with Lumi first.' },
      { step: '03', title: 'Create your learning journey', desc: 'A practical path is generated from your refined goal.' },
      { step: '04', title: 'Learn and complete it', desc: 'Add content, leave records, and keep moving until completion.' },
    ],
  },
  aiUsage: {
    eyebrow: 'AI Use',
    titlePrefix: 'Choose your AI model,',
    titleHighlight: 'while recommendations use installed embeddings',
    descriptionLines: ['Content search and recommendations are handled through an installed embedding model on the LearnCosmos server.', 'Generative AI features such as goal chat, course creation, and learning coaching are being prepared to work with user-connected AI keys step by step.'],
    byokBadge: 'BYOK',
    byokTitle: 'Connect your AI key',
    byokDescriptionLines: ['OpenAI starts first, with Claude, Gemini, and other LLM providers planned for staged connection.', 'API keys are treated as sensitive information and used only when generative AI calls need them.'],
    byokBullets: ['Use your key for generative AI features', 'Change or delete your API key', 'Usage and call history visibility planned', 'Managed as sensitive information under security rules'],
    pointsTitle: 'Start with AI points first',
    pointsDescription: (welcomePoints) => ['You can start without an API key.', `New accounts receive ${welcomePoints} AI points to begin a learning journey.`],
    pointExampleLabel: 'Point usage examples',
    pointExamples: [
      { key: 'course', label: 'Learning journey creation', icon: '🎯' },
      { key: 'lesson', label: 'Content recommendation', icon: '💡' },
    ],
    pointBullets: ['AI points are granted at sign-up', 'Start before connecting an API key', 'Switch to a BYOK LLM later'],
    sensitiveNoticeLines: ['An API key is sensitive because it can connect to your cost and permissions.', 'Content recommendations and similar-resource search use a separate installed embedding path, so BYOK LLM usage remains clearly scoped.'],
    futureProvidersLabel: 'AI providers planned for staged connection',
  },
  categories: {
    eyebrow: 'Categories',
    title: 'Anything can become a learning journey',
    descriptionLines: ['Pick an area of interest and open your first path.', 'If it is not on the list, you can enter your own topic.'],
    items: [
      { name: 'Beginner Guitar', icon: '🎸', query: 'beginner guitar' },
      { name: 'Watercolor Basics', icon: '🎨', query: 'watercolor basics' },
      { name: 'Photography', icon: '📷', query: 'photography basics' },
      { name: 'Learn Coding', icon: '💻', query: 'learn coding' },
      { name: 'Home Cafe', icon: '☕', query: 'home cafe latte art' },
      { name: 'Fitness & Health', icon: '🧘', query: 'fitness and health' },
      { name: 'Writing', icon: '📝', query: 'writing practice' },
      { name: 'Languages', icon: '🌍', query: 'learn a foreign language' },
      { name: 'Cooking', icon: '🍳', query: 'cooking basics' },
      { name: 'Certification Study', icon: '📚', query: 'certification study plan' },
      { name: 'Enter My Own Topic', icon: '✨', query: '', direct: true },
    ],
  },
  policies: {
    eyebrow: 'Public Documents',
    title: 'Review operating rules before signing up',
    descriptionLines: ['LearnCosmos provides public Terms and Privacy Policy documents.', 'Before signing up, you can review installed-embedding recommendations, BYOK LLM rules, and privacy principles.'],
    docs: [
      { href: '/en/terms', badge: 'Terms', title: 'Terms of Service', description: 'Review service terms, AI points, BYOK LLM usage, and content usage rules.' },
      { href: '/en/privacy', badge: 'Privacy', title: 'Privacy Policy', description: 'Review personal data collection, storage, protection, and sensitive settings such as API keys.' },
      { href: '/en/open-source', badge: 'Open Source', title: 'Open Source Licenses', description: 'Review the major open source software and license notices included in the service.' },
    ],
    linkLabel: 'View document',
    noteLines: ['AI points and BYOK may connect to user cost and permissions.', 'Installed embeddings for search/recommendation and BYOK LLMs for generative AI are described as separate usage paths.'],
  },
  socialLogin: {
    eyebrow: 'Start Now',
    title: 'Start your first learning journey for free',
    loggedInMessage: 'Welcome back! Your learning is waiting for you.',
    loggedInHighlight: '',
    dashboardCta: 'Go to My Dashboard',
    loggedInHint: 'You are currently logged in. Visit the dashboard to change learning settings.',
    signupDescriptionLines: ['AI points are granted when you sign up.', 'Enter what you want to learn and create your first learning journey right away.'],
    socialDescriptionLines: ['You do not need a LearnCosmos account yet.', 'Starting with a social account will continue directly to sign-up.'],
    aiNoticeLines: ['You can start without an API key,', 'and connect your own AI key later.'],
    providers: {
      google: 'Continue with Google',
      kakao: 'Continue with Kakao',
      naver: 'Continue with Naver',
    },
    policyNotice: {
      before: 'Required policy consent is requested before using the service.',
      terms: 'Terms',
      middle: 'and',
      privacy: 'Privacy Policy',
      after: 'are available before sign-up.',
    },
    naverNotice: {
      title: 'Naver login is not ready',
      bodyLines: ['Naver login is still in development.', 'Please use Google or Kakao for now.'],
      confirm: 'OK',
    },
  },
  footer: {
    taglineLines: ['Enter what you want to learn,', 'and a learning journey opens right away.'],
    columns: [
      { title: 'Service', links: [{ label: 'Course Creation', href: '#course-create-ready' }, { label: 'Learning Guide', href: '#intro-video-ready' }, { label: 'Operation Standards', href: '#public-documents-ready' }, { label: 'News & Notices', href: '#platform-news-ready' }] },
      { title: 'Connect', links: [{ label: 'Blog', href: 'https://blog.naver.com/learnweaver' }, { label: 'YouTube', href: 'https://www.youtube.com/@햇볕냥이-g9d' }] },
      { title: 'Support', links: [{ label: 'Help', href: '#intro-video-ready' }, { label: 'Contact: learnweavr@gmail.com' }, { label: 'Privacy Policy', href: '/en/privacy' }, { label: 'Terms', href: '/en/terms' }, { label: 'Open Source Licenses', href: '/en/open-source' }] },
    ],
  },
}

export function getLandingPageCopy(locale: Locale): LandingPageCopy {
  return locale === 'en' ? en : ko
}
