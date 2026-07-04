import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DiaryPointStatusCopy = {
  completed: string
  learning: string
  ready: string
  draft: string
}

export type PlanetDiaryCopy = {
  loading: string
  notFound: string
  backToGalaxy: string
  errors: {
    inactiveDiaryOnlyPlanning: string
    loadFailed: string
    startNotReady: string
    loginRequired: string
    startFailed: string
    startSuccess: string
    startFallback: string
    completeNeedsAllPoints: string
    completeFailed: string
    invalidPlanetPath: string
  }
  header: {
    meta: string
    learning: string
    shared: string
  }
  section: {
    planetPrefix: string
    sharedPrefix: string
    tabs: {
      journal: string
      records: string
      results: string
    }
    editPlan: string
  }
  summary: {
    goalEyebrow: string
    goalTitle: string
    fallbackGoal: string
    completionRate: (percent: number) => string
    completedLessons: (completed: number, total: number) => string
    completedPoints: (completed: number, total: number) => string
    progressAction: {
      readyBadge: string
      starting: string
      start: string
      learningBadge: string
      complete: string
      completeBlocked: string
    }
    guide: {
      shared: string
      ready: string
      learning: string
      allPointsDone: string
      inProgress: string
    }
    sharedNotice: {
      title: string
      message: string
    }
  }
  completionModal: {
    aria: string
    title: string
    guideTitle: string
    guideMessage: (completed: number, total: number) => string
    shareTitle: string
    shareMessage: string
    actionPrompt: string
    cancel: string
    busy: string
    confirm: string
  }
  planningConfirm: {
    aria: string
    title: string
    message: string
    cancel: string
    confirm: string
  }
  stage: {
    alt: string
    tabs: {
      planning: { label: string; shortLabel: string }
      journal: { label: string; shortLabel: string }
      records: { label: string; shortLabel: string }
      results: { label: string; shortLabel: string }
      community: { label: string; shortLabel: string }
      civilization: { label: string; shortLabel: string }
    }
    readyTitle: (label: string) => string
    save: string
    cancel: string
    collapseMap: string
    expandMap: string
  }
  guidePanel: {
    unsavedPointNotice: string
    recommendationAdded: string
    action: string
    addExplorationElement: string
  }
  firstElementOnboarding: {
    nowTitle: string
    featureTitle: string
    skip: string
    lessonTask: string
    lessonFeature: string
    lessonBubble: string
    addTask: string
    addFeature: string
    addBubble: string
    modalTask: string
    modalFeature: string
    saveTask: string
    saveFeature: string
    saveBubble: string
    openPointTask: string
    openPointFeature: string
    openPointBubble: string
    doneTask: string
    doneFeature: string
    emptyTask: string
    emptyFeature: string
    defaultTask: string
    defaultFeature: string
  }
  mapPanel: {
    backToStarSystemAria: string
    backToStarSystemLines: [string, string]
  }
  tree: {
    courseTitle: (title: string) => string
    emptyCourseTitle: string
    unnamed: string
    completing: string
    complete: string
  }
  modals: {
    exitTitle: string
    exitMessage: string
    exitCancel: string
    exitConfirm: string
    planningTitle: string
    planningMessage: string
    cancel: string
    confirm: string
    pendingPoint: {
      region: (name: string) => string
      subRegion: (name: string) => string
      type: (type: string) => string
      status: (status: string) => string
      explorationType: string
      researchType: string
      prompt: string
      action: string
    }
    plannedBookmark: {
      communityTitle: string
      communityMessage: string
      civilizationTitle: string
      civilizationMessage: string
    }
    blockedPoint: {
      title: string
      unsaved: string
      message: string
    }
  }
  statusLabels: DiaryPointStatusCopy
  statusSummaries: DiaryPointStatusCopy
}

const ko: PlanetDiaryCopy = {
  loading: '행성 기록을 준비하는 중...',
  notFound: '행성 기록을 찾지 못했습니다.',
  backToGalaxy: 'Galaxy로 돌아가기',
  errors: {
    inactiveDiaryOnlyPlanning: '비활성 탐험 다이어리는 탐험계획 페이지만 사용할 수 있습니다.',
    loadFailed: '행성 데이터를 불러오지 못했습니다.',
    startNotReady: '이 행성은 아직 학습을 시작할 수 없습니다. 준비 상태를 다시 확인해 주세요.',
    loginRequired: '로그인 상태를 확인한 뒤 다시 시도해 주세요.',
    startFailed: '학습 시작 처리에 실패했습니다.',
    startSuccess: '학습을 시작했습니다. 이제 탐험일지에서 지점을 눌러 바로 학습을 이어갈 수 있습니다.',
    startFallback: '탐험시작 처리에 실패했습니다.',
    completeNeedsAllPoints: '모든 탐험지점을 완료해야 탐험완료할 수 있습니다.',
    completeFailed: '탐험완료 처리에 실패했습니다.',
    invalidPlanetPath: '유효하지 않은 행성 ID입니다.',
  },
  header: {
    meta: 'Explorer Diary',
    learning: 'Learning Explorer Diary',
    shared: 'Shared Read-only Diary',
  },
  section: {
    planetPrefix: '행성 ',
    sharedPrefix: ' 공유 ',
    tabs: { journal: '탐험일지', records: '탐험기록', results: '탐험결과물' },
    editPlan: '탐험계획 편집',
  },
  summary: {
    goalEyebrow: 'Goal Context',
    goalTitle: '탐험 목표',
    fallbackGoal: '현재 탐험일지에 연결된 학습 목표를 불러오지 못했습니다.',
    completionRate: (percent) => `완료율 ${percent}%`,
    completedLessons: (completed, total) => `완료지역 ${completed}/${total}`,
    completedPoints: (completed, total) => `완료지점 ${completed}/${total}`,
    progressAction: {
      readyBadge: '학습 준비 상태',
      starting: '탐험시작 처리 중...',
      start: '탐험시작',
      learningBadge: 'learning 상태',
      complete: '탐험완료',
      completeBlocked: '탐험완료 조건 미충족',
    },
    guide: {
      shared: '완료된 탐험일지예요. 일지에서는 경로를 다시 살펴보고, 기록과 결과 책갈피에서 남긴 내용을 모아볼 수 있습니다.',
      ready: '탐험을 시작하면 지점별 학습 화면에서 일지, 기록, 결과물을 남길 수 있어요. 준비가 되면 진행률 영역의 탐험시작을 눌러주세요.',
      learning: '지역(리슨)을 선택해서 지점(학습콘텐츠)을 추가할 수 있어요. 지점을 선택하면 학습화면으로 이동해요. 진행 중인 지점은 파란색으로 표시돼요.',
      allPointsDone: '모든 지점을 완료했어요. 기록과 결과를 한 번 확인한 뒤 진행률 영역에서 탐험완료를 확정할 수 있습니다.',
      inProgress: '지점을 선택하면 학습 화면으로 이동해요. 기록과 결과는 오른쪽 책갈피에서 모아보고, 남은 지점은 진행률 영역에서 확인할 수 있습니다.',
    },
    sharedNotice: {
      title: '공유 탐험일지',
      message: '이 화면은 완료된 탐험 구조와 각 지점 기록을 다시 읽는 용도입니다. 상태 전이나 구조 편집 없이 지점 회고만 확인합니다.',
    },
  },
  completionModal: {
    aria: '코스 완료 확인',
    title: '코스를 완료 처리할까요?',
    guideTitle: '완료 처리 안내',
    guideMessage: (completed, total) => `지역 ${completed}/${total}곳의 지점을 모두 완료했습니다. 완료 처리 후에는 읽기 전용 탐험일지로 전환됩니다.`,
    shareTitle: '앞으로의 공유 안내',
    shareMessage: '완료된 행성의 코스와 학습콘텐츠는 학습기록 같은 개인 학습데이터를 제외하고, 이후 사용자의 허락에 따라 다른 학습자에게 공유될 수 있게 구현될 예정입니다. 내 코스를 다른 학습자가 많이 활용하면 생성자에게 혜택을 주는 방향도 준비하고 있습니다.',
    actionPrompt: '완료 여부를 선택해 주세요.',
    cancel: '취소',
    busy: '처리 중...',
    confirm: '최종확인 · 코스 완료',
  },
  planningConfirm: {
    aria: '탐험계획 편집 이동 확인',
    title: '탐험계획 편집 화면으로 이동할까요?',
    message: '탐험일지는 학습 실행 화면이고 탐험계획은 구조 편집 화면입니다. 편집 화면으로 이동하시겠습니까?',
    cancel: '취소',
    confirm: '확인',
  },
  stage: {
    alt: 'Explorer Diary',
    tabs: {
      planning: { label: '탐험계획', shortLabel: '계획' },
      journal: { label: '탐험일지', shortLabel: '일지' },
      records: { label: '탐험기록', shortLabel: '기록' },
      results: { label: '탐험결과물', shortLabel: '결과' },
      community: { label: '커뮤니티', shortLabel: '커뮤' },
      civilization: { label: '문명발전현황', shortLabel: '문명' },
    },
    readyTitle: (label) => `${label} 준비중`,
    save: '저장',
    cancel: '취소',
    collapseMap: '지도 접기',
    expandMap: '지도 펼치기',
  },
  guidePanel: {
    unsavedPointNotice: '저장 전까지는 학습페이지로 이동할 수 없습니다.',
    recommendationAdded: '추천 지점을 추가했습니다. 저장 후 학습페이지로 이동할 수 있습니다.',
    action: '작업',
    addExplorationElement: '➕ 탐험요소',
  },
  firstElementOnboarding: {
    nowTitle: '지금 할 일',
    featureTitle: '기능 소개',
    skip: '나중에 보기',
    lessonTask: '첫 리슨을 선택해 주세요.',
    lessonFeature: '탐험일지는 학습을 바로 시작하는 화면입니다. 먼저 리슨을 선택한 뒤 탐험요소를 하나 추가하면 됩니다.',
    lessonBubble: '먼저 이 리슨을 선택하세요.',
    addTask: '+탐험요소를 눌러 학습할 지점이나 연구 과제를 추가해 주세요.',
    addFeature: '탐험지점은 자료를 열어 학습하는 지점이고, 연구지점은 직접 조사할 과제입니다. 서브리슨은 흐름을 더 작게 나눌 때 씁니다.',
    addBubble: '이 버튼으로 탐험요소를 추가하세요.',
    modalTask: '탭 설명을 확인하고 탐험지점, 연구지점, 서브리슨 중 하나를 추가해 주세요.',
    modalFeature: '추천을 쓰면 현재 목표와 리슨 문맥에 맞는 자료 후보를 찾아볼 수 있습니다.',
    saveTask: '추가한 탐험요소를 저장해 주세요.',
    saveFeature: '저장 후에는 지점을 선택해 학습페이지로 이동할 수 있습니다.',
    saveBubble: '추가한 탐험요소를 저장해 주세요.',
    openPointTask: '저장된 첫 탐험요소를 선택해 학습페이지로 이동해 주세요.',
    openPointFeature: '지점을 열면 학습 대상 확인, 기록 작성, 결과물 제출을 이어갈 수 있습니다.',
    openPointBubble: '이 탐험요소를 선택해 학습을 시작하세요.',
    doneTask: '이제 지점을 선택해 학습페이지로 이동할 수 있어요.',
    doneFeature: '기록과 결과는 오른쪽 책갈피에서 모아볼 수 있습니다.',
    emptyTask: '첫 리슨을 선택하고 +탐험요소를 눌러 학습을 시작할 지점이나 연구 과제를 하나 추가해 주세요.',
    emptyFeature: '탐험일지는 학습을 진행하는 화면입니다. 탐험요소 추가와 저장은 여기서 바로 할 수 있어요.',
    defaultTask: '지점을 선택해 학습화면으로 이동해 주세요. 진행 중인 지점은 파란색으로 표시돼요.',
    defaultFeature: '기록과 결과는 오른쪽 책갈피에서 모아볼 수 있어요. 삭제, 위치 조정, 리슨 구조 변경은 탐험계획에서 할 수 있습니다.',
  },
  mapPanel: {
    backToStarSystemAria: '스타 시스템으로 돌아가기',
    backToStarSystemLines: ['← 스타', '시스템'],
  },
  tree: {
    courseTitle: (title) => `행성 ${title} 탐험 계획`,
    emptyCourseTitle: '행성 탐험 계획',
    unnamed: '이름 없음',
    completing: '완료 처리 중...',
    complete: '코스 완료',
  },
  modals: {
    exitTitle: '스타 시스템으로 돌아갈까요?',
    exitMessage: '지금 탐험일지 화면을 나가면 스타 시스템으로 돌아갑니다.\n현재 보고 있던 지점과 위치는 다시 선택해 이어볼 수 있습니다.',
    exitCancel: '계속 보기',
    exitConfirm: '스타 시스템으로 이동',
    planningTitle: '탐험계획 편집 화면으로 이동할까요?',
    planningMessage: '책갈피로 이동해도 탐험일지에서 탐험계획 편집 화면으로 전환됩니다. 계속하시겠습니까?',
    cancel: '취소',
    confirm: '확인',
    pendingPoint: {
      region: (name) => `지역 ${name}`,
      subRegion: (name) => `서브 지역 ${name}`,
      type: (type) => `유형 ${type}`,
      status: (status) => `상태 ${status}`,
      explorationType: '탐험지점',
      researchType: '연구지점',
      prompt: '이 지점의 탐험(학습)을 시작하거나 이어서 진행할까요?',
      action: '탐험(학습) 하기',
    },
    plannedBookmark: {
      communityTitle: '커뮤니티는 정식 오픈 때 제공할 예정입니다',
      communityMessage: '베타 기간에는 개인 탐험 흐름을 먼저 안정화하고, 이후 같은 목표나 비슷한 행성을 탐험하는 학습자들이 지점별 질문, 참고자료, 응원 기록을 나눌 수 있는 공간으로 열 예정입니다.',
      civilizationTitle: '문명발전현황은 정식 오픈 때 제공할 예정입니다',
      civilizationMessage: '정식 오픈되면 완료한 코스의 리슨과 학습콘텐츠를 공유하고, 다른 학습자들이 내 코스를 얼마나 활용하는지 확인할 수 있게 준비하겠습니다. 공유율이 높아지면 코스 생성자에게 주어질 혜택도 함께 연결할 예정입니다.',
    },
    blockedPoint: {
      title: '저장 후 학습페이지로 이동할 수 있습니다',
      unsaved: '탐험일지에 저장되지 않은 변경이 있습니다.',
      message: '추가하거나 수정한 탐험요소를 먼저 저장하면 이 지점의 학습페이지로 이동할 수 있습니다.',
    },
  },
  statusLabels: { completed: '완료', learning: '진행 중', ready: '시작 전', draft: '준비 중' },
  statusSummaries: {
    completed: '탐험이 끝난 지점입니다.',
    learning: '현재 학습 중인 지점입니다.',
    ready: '학습페이지로 이동해 탐험을 시작할 수 있습니다.',
    draft: '아직 준비 단계에 있는 지점입니다.',
  },
}

const en: PlanetDiaryCopy = {
  loading: 'Preparing planet diary...',
  notFound: 'Could not find this planet diary.',
  backToGalaxy: 'Back to Galaxy',
  errors: {
    inactiveDiaryOnlyPlanning: 'Inactive explorer diaries can only use the planning page.',
    loadFailed: 'Could not load planet data.',
    startNotReady: 'This planet is not ready to start yet. Check its readiness and try again.',
    loginRequired: 'Check your login session and try again.',
    startFailed: 'Could not start learning.',
    startSuccess: 'Learning started. You can now continue from points in the Explorer Diary.',
    startFallback: 'Could not start exploration.',
    completeNeedsAllPoints: 'Complete every exploration point before finishing this exploration.',
    completeFailed: 'Could not complete this exploration.',
    invalidPlanetPath: 'Invalid planet ID.',
  },
  header: {
    meta: 'Explorer Diary',
    learning: 'Learning Explorer Diary',
    shared: 'Shared Read-only Diary',
  },
  section: {
    planetPrefix: 'Planet ',
    sharedPrefix: ' Shared ',
    tabs: { journal: 'Explorer Diary', records: 'Exploration Records', results: 'Exploration Artifacts' },
    editPlan: 'Edit Plan',
  },
  summary: {
    goalEyebrow: 'Goal Context',
    goalTitle: 'Exploration Goal',
    fallbackGoal: 'Could not load the learning goal connected to this diary.',
    completionRate: (percent) => `Completion ${percent}%`,
    completedLessons: (completed, total) => `Regions ${completed}/${total}`,
    completedPoints: (completed, total) => `Points ${completed}/${total}`,
    progressAction: {
      readyBadge: 'Ready to learn',
      starting: 'Starting...',
      start: 'Start Exploration',
      learningBadge: 'Learning',
      complete: 'Complete Exploration',
      completeBlocked: 'Completion requirements unmet',
    },
    guide: {
      shared: 'This is a completed Explorer Diary. Review the path, records, and artifacts from the bookmarks.',
      ready: 'After you start, each point lets you save diary notes, records, and artifacts. Use Start Exploration when you are ready.',
      learning: 'Select a region to add points. Select a point to move to its learning page. Points in progress are shown in blue.',
      allPointsDone: 'All points are complete. Review records and artifacts, then confirm completion from the progress area.',
      inProgress: 'Select a point to move to learning. Records and artifacts are collected in the right bookmarks, and remaining points appear in the progress area.',
    },
    sharedNotice: {
      title: 'Shared Explorer Diary',
      message: 'This read-only view is for reviewing a completed exploration structure and point records. No status changes or structure edits are available here.',
    },
  },
  completionModal: {
    aria: 'Confirm course completion',
    title: 'Complete this course?',
    guideTitle: 'Completion guide',
    guideMessage: (completed, total) => `You completed the points in ${completed}/${total} regions. After completion, this diary becomes read-only.`,
    shareTitle: 'Future sharing',
    shareMessage: 'Completed course and learning content may later be shared with other learners, excluding personal learning data such as private records and only with your permission. We are also preparing benefits for creators when their courses are used often.',
    actionPrompt: 'Choose whether to complete it.',
    cancel: 'Cancel',
    busy: 'Processing...',
    confirm: 'Confirm · Complete Course',
  },
  planningConfirm: {
    aria: 'Confirm plan editor navigation',
    title: 'Move to the plan editor?',
    message: 'Explorer Diary is the learning screen, and Exploration Plan is the structure editor. Move to the editor?',
    cancel: 'Cancel',
    confirm: 'Confirm',
  },
  stage: {
    alt: 'Explorer Diary',
    tabs: {
      planning: { label: 'Exploration Plan', shortLabel: 'Plan' },
      journal: { label: 'Explorer Diary', shortLabel: 'Diary' },
      records: { label: 'Exploration Records', shortLabel: 'Records' },
      results: { label: 'Exploration Artifacts', shortLabel: 'Artifacts' },
      community: { label: 'Community', shortLabel: 'Community' },
      civilization: { label: 'Civilization Progress', shortLabel: 'Civ' },
    },
    readyTitle: (label) => `${label} coming soon`,
    save: 'Save',
    cancel: 'Cancel',
    collapseMap: 'Collapse Map',
    expandMap: 'Expand Map',
  },
  guidePanel: {
    unsavedPointNotice: 'You can move to the learning page after saving.',
    recommendationAdded: 'Recommended point added. Save it before moving to the learning page.',
    action: 'Action',
    addExplorationElement: '+ Add Element',
  },
  firstElementOnboarding: {
    nowTitle: 'Next Action',
    featureTitle: 'Feature',
    skip: 'Later',
    lessonTask: 'Select the first lesson.',
    lessonFeature: 'Explorer Diary is where learning starts. Select a lesson first, then add one exploration element.',
    lessonBubble: 'Select this lesson first.',
    addTask: 'Use + Add Element to add a point or research task.',
    addFeature: 'Exploration points open learning resources, research points hold investigation tasks, and sub-lessons split the flow.',
    addBubble: 'Add an exploration element here.',
    modalTask: 'Review the tab descriptions and add an exploration point, research point, or sub-lesson.',
    modalFeature: 'Recommendations can find candidate resources from the current goal and lesson context.',
    saveTask: 'Save the exploration element you added.',
    saveFeature: 'After saving, select a point to move to its learning page.',
    saveBubble: 'Save the added element.',
    openPointTask: 'Select the first saved element to move to its learning page.',
    openPointFeature: 'Opening a point lets you review the learning target, write records, and submit artifacts.',
    openPointBubble: 'Select this element to start learning.',
    doneTask: 'Now select a point to move to the learning page.',
    doneFeature: 'Records and artifacts are collected from the right bookmarks.',
    emptyTask: 'Select the first lesson and use + Add Element to add one point or research task.',
    emptyFeature: 'Explorer Diary is the learning screen. You can add elements and save them here.',
    defaultTask: 'Select a point to move to the learning page. Points in progress are shown in blue.',
    defaultFeature: 'Records and artifacts are collected from the right bookmarks. Delete, reorder, and lesson structure edits are handled in Exploration Plan.',
  },
  mapPanel: {
    backToStarSystemAria: 'Back to Star System',
    backToStarSystemLines: ['← Star', 'System'],
  },
  tree: {
    courseTitle: (title) => `Planet ${title} Exploration Plan`,
    emptyCourseTitle: 'Planet Exploration Plan',
    unnamed: 'Untitled',
    completing: 'Completing...',
    complete: 'Complete Course',
  },
  modals: {
    exitTitle: 'Return to Star System?',
    exitMessage: 'You will leave the Explorer Diary and return to the Star System.\nYou can select this point and location again later.',
    exitCancel: 'Keep Viewing',
    exitConfirm: 'Go to Star System',
    planningTitle: 'Move to the plan editor?',
    planningMessage: 'This bookmark switches from Explorer Diary to the Exploration Plan editor. Continue?',
    cancel: 'Cancel',
    confirm: 'Confirm',
    pendingPoint: {
      region: (name) => `Region ${name}`,
      subRegion: (name) => `Subregion ${name}`,
      type: (type) => `Type ${type}`,
      status: (status) => `Status ${status}`,
      explorationType: 'Exploration Point',
      researchType: 'Research Point',
      prompt: 'Start or continue learning at this point?',
      action: 'Start Learning',
    },
    plannedBookmark: {
      communityTitle: 'Community will open at the official release',
      communityMessage: 'During beta, we are stabilizing personal exploration first. Later this area will let learners with similar goals or planets share point questions, references, and encouragement.',
      civilizationTitle: 'Civilization Progress will open at the official release',
      civilizationMessage: 'After the official release, you will be able to share completed course lessons and learning content and see how other learners use your course. Creator benefits are also planned for highly used courses.',
    },
    blockedPoint: {
      title: 'Save before moving to the learning page',
      unsaved: 'There are unsaved changes in this Explorer Diary.',
      message: 'Save the added or edited exploration elements first, then you can move to this point learning page.',
    },
  },
  statusLabels: { completed: 'Completed', learning: 'In Progress', ready: 'Not Started', draft: 'Preparing' },
  statusSummaries: {
    completed: 'This point is complete.',
    learning: 'This point is currently in progress.',
    ready: 'Move to the learning page to start this exploration.',
    draft: 'This point is still being prepared.',
  },
}

export function getPlanetDiaryCopy(locale: Locale | string | null | undefined): PlanetDiaryCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
