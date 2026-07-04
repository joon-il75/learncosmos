import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'
import type { DashboardSortKey } from '@/lib/world-ui-engine/types'
import type { PlanetStatus } from '@/components/dashboard/planetMapTokens'

type DashboardTodayTaskKind = 'create_course' | 'start_planet' | 'continue_point' | 'complete_planet' | 'review_records'

export type DashboardMainCopy = {
  loading: string
  heroTitle: string
  heroSubtitle: string
  cta: {
    title: string
    placeholder: string
    submitLabel: string
    submittingLabel: string
    slogan: { default: string; completedOnly: string }
    description: { default: string; completedOnly: string }
    inputHelp: string
    topicChips: string[]
    errors: {
      required: string
      createFailed: string
    }
  }
  todayTask: {
    ariaLabel: string
    eyebrow: string
    kindLabels: Record<DashboardTodayTaskKind, string>
  }
  galaxy: {
    eyebrow: string
    courseCount: (count: number) => string
    statusFilter: (label: string) => string
    imageAlt: string
    missingSystem: string
    selectedSystemFallback: (index: number) => string
    sortOptions: Array<{ key: DashboardSortKey; label: string }>
    statusLabels: Record<PlanetStatus, string> & { inactive: string }
    statusLine: {
      prefix: string
      countSuffix: string
      empty: string
    }
    updatedAt: {
      prefix: string
      empty: string
    }
    mobile: {
      title: { galaxy: string; starSystem: string }
      primary: { galaxy: string; starSystem: string }
      secondary: { galaxy: string; starSystem: string }
      current: (summaryTitle: string) => string
    }
    hover: {
      title: (pageNumber: number) => string
      courseCount: (count: number) => string
      emptySystem: (title: string) => string
      learningOnly: (title: string) => string
      inactiveOnly: (title: string) => string
      completedOnly: (title: string) => string
      readyOnly: (title: string) => string
      draftOnly: (title: string) => string
      mixed: (title: string) => string
      courseList: string
    }
  }
  planetMap: {
    systemTitle: (title: string | null | undefined) => string
    starSystemMode: string
    returnToGalaxy: string
    empty: string
    ariaLabel: (title: string, status: string) => string
    lessonCount: (count: number) => string
    progress: (percent: number) => string
    regionCount: (count: number) => string
    completedCount: (done: number, total: number) => string
    inactiveMessage: string
    completedMessage: string
    learningMessage: string
    readyMessage: string
    draftMessage: string
  }
  list: {
    modeLabel: { galaxy: string; starSystem: (index: string | number) => string }
    pageLabel: { galaxy: (page: number, total: number) => string; starSystem: (index: string | number) => string }
    totalPages: (total: number) => string
    statusFilters: {
      all: string
      learning: string
      completed: string
    }
    inactive: {
      show: string
      hide: string
      badge: string
      showTitle: string
      hideTitle: string
    }
    searchPlaceholder: string
    searchAria: string
    pagination: { previous: string; next: string }
    systemNav: {
      previous: string
      next: string
      label: string
      previousAria: string
      nextAria: string
    }
    emptySearch: (query: string) => string
    emptyFilter: string
    progress: {
      inactive: string
      learningWithProgress: (percent: number, done: number, total: number) => string
      learning: string
      draftWithProgress: (percent: number, planned: number, total: number) => string
      draft: string
      readyWithRegions: (count: number) => string
      readyWithLessons: (count: number) => string
      ready: string
      completedWithLessons: (count: number) => string
      completed: string
    }
  }
}

const ko: DashboardMainCopy = {
  loading: '대시보드를 불러오는 중...',
  heroTitle: '내 학습지도 · Galaxy',
  heroSubtitle: '오늘 이어갈 일, 새 학습탐험 만들기, 현재 행성 지도를 한 화면에서 확인하세요.',
  cta: {
    title: '새 학습탐험 경로 만들기',
    placeholder: '예: 기타 코드 3개월 안에 마스터하기',
    submitLabel: '학습탐험 시작하기',
    submittingLabel: '행성 생성 중...',
    slogan: {
      default: '배우고 싶은 것을 입력하면, 학습탐험 경로가 열립니다',
      completedOnly: '다음에 배우고 싶은 것을 입력해 보세요',
    },
    description: {
      default: '루미가 목표 정리를 돕고, 지금 따라갈 수 있는 학습탐험 경로를 준비합니다.',
      completedOnly: '최근 탐험을 돌아봤다면, 새 주제를 입력해 다음 학습탐험 경로를 이어갈 수 있습니다.',
    },
    inputHelp: '배우고 싶은 것을 편하게 적어 주세요.',
    topicChips: ['기타 입문', '수채화 기초', '라떼아트', 'DSLR 사진촬영'],
    errors: {
      required: '탐험할 행성 주제를 먼저 입력해 주세요.',
      createFailed: '행성탐험계획 생성에 실패했습니다.',
    },
  },
  todayTask: {
    ariaLabel: '오늘 이어갈 학습',
    eyebrow: 'Today Route',
    kindLabels: {
      create_course: '새 여정',
      start_planet: '탐험 시작',
      continue_point: '다음 지점',
      complete_planet: '마무리',
      review_records: '기록 회고',
    },
  },
  galaxy: {
    eyebrow: 'Space Map',
    courseCount: (count) => `현재 조건에 맞는 코스 ${count}개`,
    statusFilter: (label) => `상태 필터 ${label}`,
    imageAlt: 'Galaxy map',
    missingSystem: '선택한 항성계 정보를 다시 맞추는 중입니다. Galaxy Mode로 돌아가 항성계를 다시 선택해 주세요.',
    selectedSystemFallback: (index) => `항성계 ${index}`,
    sortOptions: [
      { key: 'updated_desc', label: '최신순' },
      { key: 'updated_asc', label: '오래된순' },
      { key: 'title_asc', label: '제목 오름차순' },
      { key: 'title_desc', label: '제목 내림차순' },
    ],
    statusLabels: {
      draft: '이전 버전',
      ready: '준비중',
      learning: '학습중',
      completed: '완료',
      inactive: '비활성',
    },
    statusLine: { prefix: '상태:', countSuffix: '개', empty: '상태 정보 없음' },
    updatedAt: { prefix: '최근 업데이트:', empty: '최근 업데이트 정보 없음' },
    mobile: {
      title: { galaxy: 'Galaxy Dashboard 안내', starSystem: 'Star System 모드 안내' },
      primary: {
        galaxy: '항성을 탭하여 해당 항성계의 행성을 확인하고, 하단 리스트에서 행성을 스크롤하세요.',
        starSystem: '행성을 눌러 학습을 이어가거나, Lumi 버튼을 눌러 설명을 확인하세요.',
      },
      secondary: {
        galaxy: '필터·정렬 버튼을 사용해 원하는 행성을 빠르게 찾을 수 있습니다.',
        starSystem: '현재 항성계 정보를 확인한 뒤, 행성을 눌러 즉시 학습을 이어가세요.',
      },
      current: (summaryTitle) => `현재 ${summaryTitle}`,
    },
    hover: {
      title: (pageNumber) => `Lumi · 항성계 ${pageNumber}`,
      courseCount: (count) => `코스 ${count}개`,
      emptySystem: (title) => `${title}는 아직 비어 있어요. 새 학습 행성을 만들면 이 구역이 먼저 채워집니다.`,
      learningOnly: (title) => `${title}는 지금 학습이 활발하게 진행되고 있어요.`,
      inactiveOnly: (title) => `${title}에는 비활성화한 행성만 모여 있어요.`,
      completedOnly: (title) => `${title}에는 학습을 완료한 행성만 모여 있어요.`,
      readyOnly: (title) => `${title}의 행성들이 학습 준비 상태예요.`,
      draftOnly: (title) => `${title}의 행성들은 이전 버전 계획 데이터예요.`,
      mixed: (title) => `${title}는 다양한 상태가 섞여 있어요.`,
      courseList: 'Course List',
    },
  },
  planetMap: {
    systemTitle: (title) => title ? `${title} 항성계 · Star System Mode` : 'Star System Mode',
    starSystemMode: 'Star System Mode',
    returnToGalaxy: 'Galaxy Mode로 복귀',
    empty: '아직 생성된 행성이 없습니다.',
    ariaLabel: (title, status) => `${title} - ${status}`,
    lessonCount: (count) => `레슨 ${count}개`,
    progress: (percent) => `진행률 ${percent}%`,
    regionCount: (count) => `지역 ${count}개`,
    completedCount: (done, total) => `완료 ${done}/${total}`,
    inactiveMessage: '이 행성은 지금 목록에서 쉬고 있는 상태예요. 필요하면 다시 활성화해서 이어갈 수 있습니다.',
    completedMessage: '이 행성은 탐험을 마친 행성이에요. 기록과 결과물을 다시 살펴볼 수 있습니다.',
    learningMessage: '이 행성은 지금 탐험이 진행 중이에요. 탐험일지에서 다음 지점을 이어가면 됩니다.',
    readyMessage: '이 행성은 학습을 시작할 준비가 되어 있어요. 탐험일지를 열어 첫 지점부터 시작해 보세요.',
    draftMessage: '이 행성은 아직 탐험계획을 다듬는 중이에요. 계획을 정리한 뒤 학습을 시작할 수 있습니다.',
  },
  list: {
    modeLabel: {
      galaxy: 'Galaxy Course List',
      starSystem: (index) => `Star System ${index} 코스`,
    },
    pageLabel: {
      galaxy: (page, total) => `페이지 ${page}/${total}`,
      starSystem: (index) => `항성계 ${index}`,
    },
    totalPages: (total) => `전체 ${total}페이지`,
    statusFilters: { all: '전체', learning: '학습중', completed: '완료' },
    inactive: {
      show: '비활성 보기',
      hide: '비활성 숨김',
      badge: '비활성',
      showTitle: '비활성화된 코스도 함께 봅니다',
      hideTitle: '비활성화된 코스를 숨깁니다',
    },
    searchPlaceholder: '코스 제목 검색...',
    searchAria: '코스 검색',
    pagination: { previous: '이전', next: '다음' },
    systemNav: {
      previous: '← 이전',
      next: '다음 →',
      label: '항성계 이동',
      previousAria: '이전 항성계',
      nextAria: '다음 항성계',
    },
    emptySearch: (query) => `"${query}"에 해당하는 행성이 없습니다.`,
    emptyFilter: '현재 필터 조건에 맞는 행성이 없습니다.',
    progress: {
      inactive: '학습중 비활성',
      learningWithProgress: (percent, done, total) => `학습진행 ${percent}% · ${done}/${total}개 지역 완료`,
      learning: '지금 학습을 이어갈 수 있어요.',
      draftWithProgress: (percent, planned, total) => `이전 버전 계획 ${percent}% · ${planned}/${total}개 지역 준비`,
      draft: '이전 버전 계획 데이터입니다.',
      readyWithRegions: (count) => `${count}개 지역 학습 준비`,
      readyWithLessons: (count) => `${count}개 지역 학습 준비`,
      ready: '학습 준비 상태입니다.',
      completedWithLessons: (count) => `${count}개 지역 학습을 마쳤어요. 이제 내 행성을 공유해 볼 수 있어요.`,
      completed: '학습을 완료했어요.',
    },
  },
}

const en: DashboardMainCopy = {
  loading: 'Loading dashboard...',
  heroTitle: 'My Learning Map · Galaxy',
  heroSubtitle: 'See today’s next step, create a new learning route, and review your planet map in one place.',
  cta: {
    title: 'Create a New Learning Route',
    placeholder: 'Example: Master guitar chords in 3 months',
    submitLabel: 'Start Learning Route',
    submittingLabel: 'Creating planet...',
    slogan: {
      default: 'Enter what you want to learn, and a route will open.',
      completedOnly: 'What would you like to learn next?',
    },
    description: {
      default: 'Lumi helps shape your goal and prepares a learning route you can follow now.',
      completedOnly: 'After reviewing your recent exploration, enter a new topic to continue with another route.',
    },
    inputHelp: 'Write what you want to learn in your own words.',
    topicChips: ['Beginner Guitar', 'Watercolor Basics', 'Latte Art', 'DSLR Photography'],
    errors: {
      required: 'Enter a planet topic to explore first.',
      createFailed: 'Failed to create the planet learning plan.',
    },
  },
  todayTask: {
    ariaLabel: 'Today learning task',
    eyebrow: 'Today Route',
    kindLabels: {
      create_course: 'New route',
      start_planet: 'Start',
      continue_point: 'Next point',
      complete_planet: 'Wrap up',
      review_records: 'Review',
    },
  },
  galaxy: {
    eyebrow: 'Space Map',
    courseCount: (count) => `${count} courses match the current filters`,
    statusFilter: (label) => `Status filter ${label}`,
    imageAlt: 'Galaxy map',
    missingSystem: 'Resyncing the selected star system. Return to Galaxy Mode and choose the star system again.',
    selectedSystemFallback: (index) => `Star System ${index}`,
    sortOptions: [
      { key: 'updated_desc', label: 'Latest' },
      { key: 'updated_asc', label: 'Oldest' },
      { key: 'title_asc', label: 'Title A-Z' },
      { key: 'title_desc', label: 'Title Z-A' },
    ],
    statusLabels: {
      draft: 'Legacy',
      ready: 'Ready',
      learning: 'Learning',
      completed: 'Completed',
      inactive: 'Inactive',
    },
    statusLine: { prefix: 'Status:', countSuffix: '', empty: 'No status info' },
    updatedAt: { prefix: 'Updated:', empty: 'No recent update info' },
    mobile: {
      title: { galaxy: 'Galaxy Dashboard Guide', starSystem: 'Star System Mode Guide' },
      primary: {
        galaxy: 'Tap a star to view planets in that star system, then scroll planets in the list below.',
        starSystem: 'Tap a planet to continue learning, or tap Lumi for context.',
      },
      secondary: {
        galaxy: 'Use filters and sorting to find the planet you need quickly.',
        starSystem: 'Review this star system, then tap a planet to continue learning.',
      },
      current: (summaryTitle) => `Current: ${summaryTitle}`,
    },
    hover: {
      title: (pageNumber) => `Lumi · Star System ${pageNumber}`,
      courseCount: (count) => `${count} courses`,
      emptySystem: (title) => `${title} is empty. New learning planets will fill this area first.`,
      learningOnly: (title) => `${title} is actively in progress.`,
      inactiveOnly: (title) => `${title} contains only inactive planets.`,
      completedOnly: (title) => `${title} contains only completed planets.`,
      readyOnly: (title) => `${title} planets are ready for learning.`,
      draftOnly: (title) => `${title} contains legacy plan data.`,
      mixed: (title) => `${title} has a mix of planet statuses.`,
      courseList: 'Course List',
    },
  },
  planetMap: {
    systemTitle: (title) => title ? `${title} · Star System Mode` : 'Star System Mode',
    starSystemMode: 'Star System Mode',
    returnToGalaxy: 'Return to Galaxy Mode',
    empty: 'No planets have been created yet.',
    ariaLabel: (title, status) => `${title} - ${status}`,
    lessonCount: (count) => `${count} lessons`,
    progress: (percent) => `${percent}% progress`,
    regionCount: (count) => `${count} regions`,
    completedCount: (done, total) => `Completed ${done}/${total}`,
    inactiveMessage: 'This planet is currently resting in the list. You can reactivate it when needed.',
    completedMessage: 'This planet has been explored. You can revisit its records and outcomes.',
    learningMessage: 'This planet is in progress. Continue from the next point in the exploration diary.',
    readyMessage: 'This planet is ready for learning. Open the diary and start from the first point.',
    draftMessage: 'This planet plan is still being refined. You can start learning after the plan is ready.',
  },
  list: {
    modeLabel: {
      galaxy: 'Galaxy Course List',
      starSystem: (index) => `Star System ${index} Courses`,
    },
    pageLabel: {
      galaxy: (page, total) => `Page ${page}/${total}`,
      starSystem: (index) => `Star System ${index}`,
    },
    totalPages: (total) => `${total} pages total`,
    statusFilters: { all: 'All', learning: 'Learning', completed: 'Completed' },
    inactive: {
      show: 'Show inactive',
      hide: 'Hide inactive',
      badge: 'Inactive',
      showTitle: 'Include inactive courses',
      hideTitle: 'Hide inactive courses',
    },
    searchPlaceholder: 'Search course titles...',
    searchAria: 'Search courses',
    pagination: { previous: 'Previous', next: 'Next' },
    systemNav: {
      previous: '← Previous',
      next: 'Next →',
      label: 'Star System Navigation',
      previousAria: 'Previous star system',
      nextAria: 'Next star system',
    },
    emptySearch: (query) => `No planets match "${query}".`,
    emptyFilter: 'No planets match the current filters.',
    progress: {
      inactive: 'Learning inactive',
      learningWithProgress: (percent, done, total) => `Learning ${percent}% · ${done}/${total} regions complete`,
      learning: 'You can continue learning now.',
      draftWithProgress: (percent, planned, total) => `Legacy plan ${percent}% · ${planned}/${total} regions ready`,
      draft: 'This is legacy plan data.',
      readyWithRegions: (count) => `${count} regions ready`,
      readyWithLessons: (count) => `${count} regions ready`,
      ready: 'Ready to learn.',
      completedWithLessons: (count) => `${count} regions completed. You can share your planet now.`,
      completed: 'Learning completed.',
    },
  },
}

export function getDashboardMainCopy(locale: Locale | string | null | undefined): DashboardMainCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
