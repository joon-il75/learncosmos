import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DashboardCourseDraftCopy = {
  page: {
    loading: string
    notFound: string
    invalidDraftId: string
    loadFailed: string
    backToDashboard: string
    headerTitle: string
    headerSubtitle: string
    sectionTitle: (title: string) => string
    fallbackCourseName: string
    returnToDiary: string
    goal: {
      planningEyebrow: string
      journalEyebrow: string
      planningTitle: string
      journalTitle: string
      edit: string
      set: string
      loading: string
      empty: string
      currentRegion: (name: string) => string
      currentSubRegion: (name: string) => string
      journalPreLearning: string
      journalPointReady: string
      journalPointPending: string
      journalGuideCanOpen: string
      journalGuideExplore: string
      usage: (value: string) => string
      motivation: (value: string) => string
      learnerEditDecision: string
      rebuildAllDecision: string
      statusLabel: (status?: string) => string
      planningBeforeStart: string
      planningChangedBeforeStart: string
      planningChangedLearning: string
      status: (value: string) => string
      completedRegions: (completed: number, total: number) => string
      completedPoints: (completed: number, total: number, explorationCompleted: number, explorationTotal: number, researchCompleted: number, researchTotal: number) => string
      completionRate: (percent: number) => string
      regionProgress: (completed: number, total: number) => string
      learningRegions: (count: number) => string
      remainingRegions: (count: number) => string
    }
    goalRevision: {
      aria: string
      eyebrow: string
      title: string
      description: string
      noInterviewTitle: string
      noInterviewDescription: string
      start: string
      starting: string
      reviseSeed: string
      rebuildTitle: string
      rebuildDescriptionBeforeStart: string
      rebuildDescriptionLearning: string
      rebuildOptions: {
        rebuildAll: { label: string; desc: string }
        keepStructure: { label: string; desc: string }
      }
      confirmedAction: string
      confirmedActionBusy: string
      confirmedMessage: string
      applySuccess: string
      applyFailed: string
      decisionSuccess: string
      cancelSuccess: string
      close: string
    }
    rebuildLoading: {
      aria: string
      eyebrow: string
      title: string
      message: string
      steps: string[]
    }
    returnConfirm: {
      aria: string
      eyebrow: string
      title: string
      message: string
      cancel: string
      confirm: string
    }
    planned: {
      community: { title: string; message: string }
      civilization: { title: string; message: string }
      eyebrow: string
      confirm: string
    }
    overview: {
      descriptionFallback: (sourceQuery: string) => string
      sourceLabel: string
      levelLabel: string
      durationLabel: string
      weeklyHoursLabel: string
      unset: string
      weeks: (count: number) => string
      hours: (count: number) => string
      lumiTitle: string
      editEyebrow: string
      editTitle: string
      editHint: string
      titleLabel: string
      descriptionLabel: string
      saving: string
      save: string
    }
  }
  hook: {
    saveTitleRequired: string
    saveFailed: string
    saveSuccess: string
    loginPath: (draftId: string) => string
    agreementsPath: (draftId: string) => string
    entries: {
      authRequired: string
      targetNotFound: string
      journalContentRequired: string
      journalSaveFailed: string
      journalResponseMissing: string
      journalSaveSuccess: string
      recordInvalid: string
      recordSaveFailed: string
      recordResponseMissing: string
      recordSaveSuccess: string
      artifactTargetNotFound: string
      artifactInvalid: string
      artifactSaveFailed: string
      artifactResponseMissing: string
      artifactSaveSuccess: string
    }
    detail: {
      savePoint: string
      savingPoint: string
      saveResearch: string
      savingResearch: string
      saveRegion: string
      savingRegion: string
      draftMissing: string
      authRequired: string
      memoPointNotFound: string
      memoRegionNotFound: string
      memoInvalid: string
      memoSaveFailed: string
      memoResponseMissing: string
      pointTitleRequired: string
      pointNotFound: string
      pointSaveFailed: string
      pointResponseMissing: string
      regionTitleRequired: string
      regionNotFound: string
      regionSaveFailed: string
      regionResponseMissing: string
      researchTitleRequired: string
      researchNotFound: string
      researchInvalid: string
      researchSaveFailed: string
      researchResponseMissing: string
      structureSaveFirstAdd: string
      structureSaveFirstDelete: string
      researchAddNotFound: string
      researchAddFailed: string
      researchAddResponseMissing: string
      researchAddSuccess: string
      researchDeleteNotFound: string
      researchDeleteFailed: string
      researchDeleteSuccess: string
      subregionTitleRequired: string
      subregionAddNotFound: string
      subregionInvalid: string
      subregionAddFailed: string
      subregionResponseMissing: string
      subregionAddSuccess: string
      savedTitleDescriptionDifficulty: string
      savedTitleDescription: string
      savedMemo: string
      pointSaved: (parts: string) => string
      regionSaved: (parts: string) => string
      researchSaved: string
      templateLabel: (templateType?: string | null) => string
      templateDescription: (templateType?: string | null) => string
      defaultResearchTitle: (templateType: string, index: number) => string
    }
    resources: {
      titleLabel: string
      titleValid: string
      titleRequired: string
      urlLabel: string
      urlValid: string
      urlInvalid: string
      memoLabel: string
      memoValid: string
      memoOptional: string
      authRequired: string
      parseFailed: string
      parseFailedKeep: string
      parsePartial: string
      parseSuccess: string
      contentInvalid: string
      contentSaveFailed: string
      contentSaveFailedKeep: string
      contentResponseMissing: string
      contentCreated: string
      attachNotFound: string
      attachConflict: string
      attachInvalidState: string
      attachFailed: string
      attachFailedKeep: string
      attachResponseMissing: string
      contentCreatedAndAttached: string
      attachSuccess: string
      attachCandidateMessage: string
      insufficientPoints: string
      candidateRefreshFailed: string
      candidateRefreshFailedKeep: string
      candidateResponseMissing: string
      selectNotFound: string
      selectCandidateOnly: string
      selectFailed: string
      selectFailedKeep: string
      selectResponseMissing: string
      selectSuccess: string
      costPreview: (cost: number) => string
      selectedRegionBasis: string
      allRegionBasis: string
      candidateRefreshSelected: (basis: string, count: number) => string
      candidateRefreshAll: (basis: string, count: number) => string
    }
  }
  planning: {
    backgroundAlt: {
      smallPlanning: string
      smallJournal: string
      expandedPlanning: string
      expandedJournal: string
      collapsedPlanning: string
      collapsedJournal: string
    }
    inactiveBanner: string
    collapseMap: string
    expandMap: string
    pointModal: {
      aria: string
      region: (name: string) => string
      subRegion: (name: string) => string
      type: (value: 'exploration' | 'research') => string
      question: string
      cancel: string
      open: string
    }
    transition: {
      aria: string
      toPlanningTitle: string
      toJournalTitle: string
      toPlanningMessage: string
      toJournalMessage: string
      cancel: string
      confirm: string
    }
    rightPanel: {
      bookmarks: Array<{ key: 'planning' | 'journal' | 'records' | 'artifacts' | 'community' | 'civilization'; label: string; shortLabel: string }>
      inactiveTitle: string
      lockedTitle: (label: string) => string
      planned: {
        community: { title: string; message: string }
        civilization: { title: string; message: string }
        eyebrow: string
        confirm: string
      }
    }
    lumiGuide: {
      eyebrow: string
      inactive: string
      unsaved: string
      course: string
      region: string
      subregion: string
      explorationNode: string
      researchNode: string
      default: string
    }
    map: {
      headerBadge: string
      courseMapTitle: string
      regionMapFallback: string
      subregionMapFallback: string
      returnToStarSystemAria: string
      returnToStarSystemLine1: string
      returnToStarSystemLine2: string
      returnToCourseMapAria: string
      returnToCourseMapLine1: string
      returnToCourseMapLine2: string
      returnToRegionMapAria: string
      returnToRegionMapLine1: string
      returnToRegionMapLine2: string
      emptyCourseLine1: string
      emptyCourseLine2: string
      missingSubregion: string
      missingRegion: string
      exitTitle: string
      exitEyebrow: string
      exitMessage: string
      exitCancel: string
      exitConfirm: string
      previousPageAria: string
      nextPageAria: string
      subregionLegend: string
      explorationLegend: string
      researchLegend: string
      formatRegionTitle: (index: number, name: string) => string
      formatRegionMeta: (subregions: number, exploration: number, research: number) => string
      formatRegionAria: (title: string, subregions: number, exploration: number, research: number) => string
      regionEmptyLine1: string
      regionEmptyLine2: string
      regionEmptyLine3: string
      subregionEmptyLine1: string
      subregionEmptyLine2: string
      subregionEmptyLine3: string
      formatSubregionNodeCounts: (exploration: number, research: number) => string
      formatRegionSummary: (subregions: number, exploration: number, research: number) => string[]
      formatSubregionSummary: (exploration: number, research: number) => string[]
    }
  }
  supportSections: {
    common: {
      currentRegion: string
      currentPoint: string
      currentStage: string
      regionNotSelected: string
      pointNotSelected: string
      selectRegionFirst: string
      reset: string
      draftPreview: string
    }
    records: {
      eyebrow: string
      title: string
      statusChanged: string
      statusConnected: string
      statusNeedRegion: string
      regionDescriptionFallback: string
      pointDescriptionFallback: string
      focusMinutes: string
      practiceCount: string
      confidence: string
      applicationNote: string
      applicationPlaceholder: string
      applicationDisabledPlaceholder: string
      saving: string
      save: string
      noticeReady: string
      noticeNeedRegion: string
    }
    artifacts: {
      eyebrow: string
      title: string
      statusChanged: string
      statusConnected: string
      statusNeedRegion: string
      regionDescriptionFallback: string
      pointDescriptionFallback: string
      type: string
      typeOptions: Array<{ value: string; label: string }>
      artifactTitle: string
      link: string
      description: string
      descriptionPlaceholder: string
      descriptionDisabledPlaceholder: string
      saving: string
      save: string
      noticeReady: string
      noticeNeedRegion: string
    }
    community: {
      eyebrow: string
      title: string
      status: string
      pointDescription: string
      stageDescription: string
      items: Array<{ title: string; body: string }>
    }
    civilization: {
      eyebrow: string
      title: string
      status: string
      stageDescription: string
      metricLabel: string
      metricValue: string
      metricDescription: string
      items: Array<{ title: string; body: string }>
    }
    selectedPoint: {
      researchLabel: string
      researchTemplateDescription: (template: string) => string
      researchTemplateLabel: (templateType?: string | null) => string | null
      researchDescription: string
      explorationLabel: string
      selectedExploration: string
      candidateExploration: string
      genericExploration: string
      formatExplorationDescription: (label: string, resourceType: string) => string
    }
  }
  tree: {
    loading: string
    emptyRegions: string
    formatPlanTitle: (title?: string | null) => string
    formatPlanetName: (title?: string | null) => string
    planetPrefix: string
    actionLabel: string
    pendingDelete: string
    pendingCreate: string
    openLinkTitle: string
    regionLabel: (index: number, name: string) => string
    regionPlaceholder: (index: number) => string
    collapseRegion: string
    expandRegion: string
    subRegionLabel: (index: number, name: string) => string
    subRegionPlaceholder: (index: number) => string
    collapseSubRegion: string
    expandSubRegion: string
    actionBar: {
      save: string
      cancel: string
      delete: string
      undoDelete: string
      activatePlanet: string
      activate: string
      hideInactive: string
      showInactive: string
      inactiveActivateTitle: string
      activateTitle: string
      undoDeleteTitle: string
      deleteCourseTitle: string
      deleteWithChildrenTitle: string
      inactiveToggleTitle: string
      reactivateTitle: string
    }
    editPanel: {
      action: string
      pendingDeleteNotice: string
      pendingCreateNotice: string
      renamePlanet: string
      addRegion: string
      renameRegion: string
      renameSubRegion: string
      addChild: string
      editNode: string
      moveUp: string
      moveDown: string
      modals: {
        close: string
        cancel: string
        apply: string
        edit: string
        add: string
        change: string
        check: string
        checking: string
        valid: string
        error: string
        openLink: string
        internalContent: string
        urlMustPass: string
        parentPath: (value: string) => string
        recommendation: {
          aria: string
          title: string
          queryLabel: string
          queryPlaceholder: string
          search: string
          searching: string
          recommendAgain: string
          empty: string
          pointErrorEmpty: string
          loading: string
          noResults: string
          loaded: (count: number) => string
          missingCandidate: string
          loginRequired: string
          insufficientPoints: (balance: number, cost: number) => string
          loadFailed: string
        }
        addExploration: {
          aria: string
          title: string
          recommendationTab: string
          urlTab: string
          titleLabel: string
          titlePlaceholder: string
          urlLabel: string
          addButton: string
        }
        addChild: {
          aria: string
          title: string
          tabAria: string
          subregionTab: string
          explorationTab: string
          researchTab: string
          subregionDescription: string
          subregionLimit: string
          subregionNameLabel: string
          subregionNamePlaceholder: string
          addSubregion: string
          explorationDescription: string
          explorationMethodAria: string
          recommendationDescription: string
          urlDescription: string
          urlDirectTab: string
          researchDescription: string
          researchTitleLabel: string
          researchTitlePlaceholder: string
          addResearch: string
        }
        addUrl: {
          explorationTitle: string
          researchTitle: string
        }
        nodeEdit: {
          explorationAria: string
          researchAria: string
          explorationTitle: string
          researchTitle: string
          titleLabel: string
          explorationPlaceholder: string
          researchPlaceholder: string
          linkLabel: string
          recommend: string
        }
        deleteConfirm: {
          deleteAria: string
          activateAria: string
          alertAria: string
          warningEyebrow: string
          alertTitle: string
          alertSuffix: string
          deleteTitle: string
          activateTitle: string
          targetLabel: string
          courseConfirmLabel: string
          exactTitleHelp: (action: 'delete' | 'activate') => string
          delete: string
          activate: string
        }
        courseRename: {
          aria: string
          title: string
          label: string
          placeholder: string
          help: string
        }
        regionRename: {
          regionTitle: string
          subregionTitle: string
          regionLabel: string
          subregionLabel: string
          aria: string
          placeholder: (label: string) => string
          help: string
        }
        regionAdd: {
          aria: string
          title: string
          label: string
          description: string
          placeholder: string
          help: string
          addButton: string
        }
        messages: {
          urlValid: string
          urlInvalid: string
          urlCheckFailed: string
          contentIdBased: string
          deleteCourse: string
          activateCourse: string
          lastRegionTitle: string
          lastRegionBlock: string
          deleteRegionWithChildren: string
          deleteRegion: string
          deleteSubregionWithChildren: string
          deleteSubregion: string
          deleteNode: string
          fallbackRegion: string
          fallbackSubregion: string
        }
      }
    }
  }
}

const plannedKo = {
  community: {
    title: '커뮤니티는 정식 오픈 때 제공할 예정입니다',
    message:
      '베타 기간에는 개인 탐험 흐름을 먼저 안정화하고, 이후 같은 목표나 비슷한 행성을 탐험하는 학습자들이 지점별 질문, 참고자료, 응원 기록을 나눌 수 있는 공간으로 열 예정입니다.',
  },
  civilization: {
    title: '문명발전현황은 정식 오픈 때 제공할 예정입니다',
    message:
      '정식 오픈되면 완료한 코스의 리슨과 학습콘텐츠를 공유하고, 다른 학습자들이 내 코스를 얼마나 활용하는지 확인할 수 있게 준비하겠습니다. 공유율이 높아지면 코스 생성자에게 주어질 혜택도 함께 연결할 예정입니다.',
  },
}

const plannedEn = {
  community: {
    title: 'Community opens at launch',
    message:
      'During beta, LearnCosmos is stabilizing the personal exploration flow first. Later, learners exploring similar goals or planets will be able to share point-level questions, references, and encouragement.',
  },
  civilization: {
    title: 'Civilization View opens at launch',
    message:
      'After launch, completed courses, lessons, and learning content can be shared, and you will be able to see how other learners use your course. Creator benefits may be connected as sharing grows.',
  },
}

function researchTemplateLabelKo(templateType?: string | null) {
  switch (templateType) {
    case 'concept_summary':
      return '개념 요약'
    case 'practice_strategy':
      return '실습 전략'
    case 'problem_solving':
      return '문제 해결'
    case 'free_research':
      return '자유 연구'
    default:
      return '자유 연구'
  }
}

function researchTemplateLabelEn(templateType?: string | null) {
  switch (templateType) {
    case 'concept_summary':
      return 'Concept Summary'
    case 'practice_strategy':
      return 'Practice Strategy'
    case 'problem_solving':
      return 'Problem Solving'
    case 'free_research':
      return 'Free Research'
    default:
      return 'Free Research'
  }
}

const ko: DashboardCourseDraftCopy = {
  page: {
    loading: 'Explorer Diary를 준비하는 중...',
    notFound: '행성탐험계획을 찾지 못했습니다.',
    invalidDraftId: '유효하지 않은 draft ID입니다.',
    loadFailed: '행성탐험계획을 불러오지 못했습니다.',
    backToDashboard: 'Galaxy로 돌아가기',
    headerTitle: 'Explorer Diary',
    headerSubtitle: 'Shared Explorer Layout',
    sectionTitle: (title) => `행성 ${title || '탐험'} 탐험 계획`,
    fallbackCourseName: '탐험',
    returnToDiary: '← 탐험일지로 돌아가기',
    goal: {
      planningEyebrow: 'Goal Context',
      journalEyebrow: 'Journal Context',
      planningTitle: '탐험 목표',
      journalTitle: '이 탐험일지가 따라가는 목표와 현재 상태',
      edit: '목표 수정',
      set: '목표 설정',
      loading: '목표 문맥을 불러오는 중입니다.',
      empty: '아직 탐험계획에 연결된 목표가 없습니다. 목표를 먼저 정하면 리슨 생성과 추천이 더 선명해집니다.',
      currentRegion: (name) => `현재 지역 ${name}`,
      currentSubRegion: (name) => `현재 서브 지역 ${name}`,
      journalPreLearning: '학습 진입 전 편집 경로',
      journalPointReady: '지점을 선택하면 학습 화면으로 이동',
      journalPointPending: '지점 선택 상태를 따라 학습 화면 연결을 준비중',
      journalGuideCanOpen: '지점을 클릭하면 Lumi가 학습 진입을 한 번 더 확인해드려요.',
      journalGuideExplore: '트리와 지도를 둘러보며 현재 탐험 구조를 확인해보세요.',
      usage: (value) => `상황 ${value}`,
      motivation: (value) => `이유 ${value}`,
      learnerEditDecision: '선택: 학습자 편집',
      rebuildAllDecision: '선택: 모두 재구성',
      statusLabel: (status) => {
        if (status === 'learning') return '탐험중(learning)'
        if (status === 'confirmed') return 'confirmed'
        if (status === 'archived') return '탐험완료(archived)'
        return '탐험계획중(draft)'
      },
      planningBeforeStart: '학습 시작 전에는 목표를 다시 조정한 뒤 모두 재구성할지, 현재 구조를 학습자 편집으로 이어갈지 선택할 수 있습니다. 모두 재구성은 코스 생성과 같은 포인트가 차감됩니다.',
      planningChangedBeforeStart: '목표가 바뀌었습니다. 학습 시작 전에는 모두 재구성할지 학습자 편집으로 갈지 학습자가 선택합니다. 모두 재구성은 코스 생성과 같은 포인트가 차감됩니다.',
      planningChangedLearning: '목표가 바뀌었습니다. 학습 중에도 모두 재구성할지 학습자 편집으로 갈지 학습자가 선택합니다. 모두 재구성은 코스 생성과 같은 포인트가 차감됩니다.',
      status: (value) => `행성상태 ${value}`,
      completedRegions: (completed, total) => `완료지역 ${completed}/${total}`,
      completedPoints: (completed, total, ec, et, rc, rt) => `완료지점 ${completed}/${total}  🎯 ${ec}/${et}  🔬 ${rc}/${rt}`,
      completionRate: (percent) => `완료율 ${percent}%`,
      regionProgress: (completed, total) => `완료 지역 ${completed} / 전체 지역 ${total}`,
      learningRegions: (count) => `진행 중 지역 ${count}개`,
      remainingRegions: (count) => `남은 지역 ${count}개`,
    },
    goalRevision: {
      aria: '탐험계획 목표 조정',
      eyebrow: 'Goal Revision',
      title: '탐험계획 목표 조정',
      description: '목표를 다시 조정하면 이후 리슨 구성과 콘텐츠 추천이 새 목표 문맥을 따라가게 됩니다.',
      noInterviewTitle: '아직 연결된 목표 인터뷰가 없습니다.',
      noInterviewDescription: '현재 탐험계획 제목이나 검색 문구를 시작점으로 목표 인터뷰를 열 수 있습니다.',
      start: '목표 인터뷰 시작',
      starting: '시작 중...',
      reviseSeed: '목표를 수정하고 싶어요.',
      rebuildTitle: '목표 수정 후 어떻게 이어갈까요?',
      rebuildDescriptionBeforeStart: '학습 시작 전이므로 새 목표 기준으로 새 리슨을 생성할지, 현재 구조를 학습자가 직접 편집할지 선택해주세요. 모두 재구성은 기존 지역/지점을 비활성화하고 코스 생성과 같은 5포인트가 차감됩니다. BYOK LLM 활성 사용자는 개인 키로 생성합니다.',
      rebuildDescriptionLearning: '학습 중에도 새 목표에 맞춰 새 리슨을 생성할지, 현재 구조를 유지한 채 학습자가 직접 편집할지 선택해주세요. 모두 재구성은 기존 지역/지점을 비활성화하고 코스 생성과 같은 5포인트가 차감됩니다. BYOK LLM 활성 사용자는 개인 키로 생성합니다.',
      rebuildOptions: {
        rebuildAll: { label: '모두 재구성', desc: '새 리슨을 생성하고 기존 지역/지점은 비활성화합니다. 코스 생성과 같은 5포인트가 차감되며, BYOK LLM 활성 사용자는 개인 키로 생성합니다' },
        keepStructure: { label: '학습자 편집', desc: '현재 구조를 유지한 채 학습자가 직접 리슨과 자료를 수정합니다' },
      },
      confirmedAction: '목표 다시 조정',
      confirmedActionBusy: '목표 반영 중...',
      confirmedMessage: '새 목표를 확정했습니다. 탐험계획을 어떻게 이어갈지 선택해 주세요.',
      applySuccess: '새 목표를 탐험계획에 반영했습니다.',
      applyFailed: '목표 반영에 실패했습니다.',
      decisionSuccess: '목표 변경 기준을 탐험계획에 반영했습니다.',
      cancelSuccess: '목표 변경을 취소하고 이전 목표로 돌아왔습니다.',
      close: '닫기',
    },
    rebuildLoading: {
      aria: '탐험계획 재구성 진행 중',
      eyebrow: 'Lumi is rebuilding',
      title: '새 목표로 리슨을 재구성하고 있어요',
      message: 'AI가 새 목표 기준으로 탐험계획을 반영하는 중입니다. 기존 구조는 비활성화하고 새 리슨을 준비하므로 잠시만 기다려 주세요.',
      steps: ['목표 반영', '리슨 생성', '탐험계획 갱신'],
    },
    returnConfirm: {
      aria: '탐험일지 복귀 확인',
      eyebrow: 'Lumi Confirm',
      title: '탐험일지로 돌아갈까요?',
      message: '탐험계획 편집 화면을 나가고 탐험일지 학습 화면으로 돌아갑니다. 계속하시겠습니까?',
      cancel: '취소',
      confirm: '확인',
    },
    planned: { ...plannedKo, eyebrow: 'Coming Soon', confirm: '확인' },
    overview: {
      descriptionFallback: (sourceQuery) => `${sourceQuery} 목표를 기준으로 AI가 만든 첫 행성탐험계획입니다.`,
      sourceLabel: '탐험 주제',
      levelLabel: '현재 수준',
      durationLabel: '권장 기간',
      weeklyHoursLabel: '주간 시간',
      unset: '미설정',
      weeks: (count) => `${count}주`,
      hours: (count) => `${count}시간`,
      lumiTitle: '현재 초안 상태 해석',
      editEyebrow: 'Safe Edit',
      editTitle: '탐험계획 기본 정보',
      editHint: '제목과 설명만 저장합니다. 탐험 구조는 변경하지 않습니다.',
      titleLabel: '탐험계획 제목',
      descriptionLabel: '탐험계획 설명',
      saving: '저장 중...',
      save: '기본 정보 저장',
    },
  },
  hook: {
    saveTitleRequired: '탐험계획 제목을 입력해 주세요.',
    saveFailed: '탐험계획 저장에 실패했습니다.',
    saveSuccess: '탐험계획 기본 정보가 저장되었습니다.',
    loginPath: (draftId) => `/login?redirect_after=/dashboard/course-drafts/${draftId}`,
    agreementsPath: (draftId) => `/agreements?redirect_after=/dashboard/course-drafts/${draftId}`,
    entries: {
      authRequired: '로그인 상태를 확인한 뒤 다시 시도해 주세요.',
      targetNotFound: '선택한 지역 기록 대상을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      journalContentRequired: '탐험일지에는 최소 한 칸 이상 기록이 필요합니다.',
      journalSaveFailed: '탐험일지 저장에 실패했습니다. 기존 입력은 유지됩니다.',
      journalResponseMissing: '탐험일지 저장 응답을 확인하지 못했습니다.',
      journalSaveSuccess: '탐험일지를 저장했습니다. 다음 진입에서도 같은 지점 학습 맥락의 기록을 이어서 볼 수 있습니다.',
      recordInvalid: '탐험기록 값 범위를 다시 확인해 주세요.',
      recordSaveFailed: '탐험기록 저장에 실패했습니다. 기존 입력은 유지됩니다.',
      recordResponseMissing: '탐험기록 저장 응답을 확인하지 못했습니다.',
      recordSaveSuccess: '탐험기록을 저장했습니다. 지점 학습의 성장 데이터로 다시 확인할 수 있습니다.',
      artifactTargetNotFound: '선택한 지역 결과물 대상을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      artifactInvalid: '결과물 제목과 URL 형식을 다시 확인해 주세요.',
      artifactSaveFailed: '탐험결과물 저장에 실패했습니다. 기존 입력은 유지됩니다.',
      artifactResponseMissing: '탐험결과물 저장 응답을 확인하지 못했습니다.',
      artifactSaveSuccess: '탐험결과물을 저장했습니다. 지점 학습 결과물 메타데이터로 다시 확인할 수 있습니다.',
    },
    detail: {
      savePoint: '지역 상세 저장',
      savingPoint: '지역 저장 중...',
      saveResearch: '연구지점 저장',
      savingResearch: '연구지점 저장 중...',
      saveRegion: '그룹 상세 저장',
      savingRegion: '그룹 저장 중...',
      draftMissing: '탐험계획 정보를 확인하지 못했습니다.',
      authRequired: '로그인 상태를 확인한 뒤 다시 시도해 주세요.',
      memoPointNotFound: '선택한 지역 메모 대상을 찾지 못했습니다.',
      memoRegionNotFound: '선택한 그룹 메모 대상을 찾지 못했습니다.',
      memoInvalid: '운영 메모 값을 확인해 주세요.',
      memoSaveFailed: '운영 메모 저장에 실패했습니다. 기존 입력은 유지됩니다.',
      memoResponseMissing: '운영 메모 저장 응답을 확인하지 못했습니다.',
      pointTitleRequired: '지역 이름을 입력해 주세요.',
      pointNotFound: '선택한 지역을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      pointSaveFailed: '지역 상세 저장에 실패했습니다. 기존 구조는 유지됩니다.',
      pointResponseMissing: '지역 상세 저장 응답을 확인하지 못했습니다.',
      regionTitleRequired: '그룹 이름을 입력해 주세요.',
      regionNotFound: '선택한 그룹을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      regionSaveFailed: '그룹 상세 저장에 실패했습니다. 기존 구조는 유지됩니다.',
      regionResponseMissing: '그룹 상세 저장 응답을 확인하지 못했습니다.',
      researchTitleRequired: '연구지점 제목을 입력해 주세요.',
      researchNotFound: '선택한 연구지점을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      researchInvalid: '연구지점 제목과 템플릿을 확인해 주세요.',
      researchSaveFailed: '연구지점 저장에 실패했습니다. 기존 구조는 유지됩니다.',
      researchResponseMissing: '연구지점 저장 응답을 확인하지 못했습니다.',
      structureSaveFirstAdd: '탐험 구조 변경을 먼저 저장한 뒤 연구지점을 추가해 주세요.',
      structureSaveFirstDelete: '탐험 구조 변경을 먼저 저장한 뒤 연구지점을 삭제해 주세요.',
      researchAddNotFound: '연구지점을 추가할 지역을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      researchAddFailed: '연구지점 추가에 실패했습니다. 기존 구조는 유지됩니다.',
      researchAddResponseMissing: '연구지점 추가 응답을 확인하지 못했습니다.',
      researchAddSuccess: '연구지점을 추가했습니다. 제목과 템플릿을 확인해 주세요.',
      researchDeleteNotFound: '삭제할 연구지점을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      researchDeleteFailed: '연구지점 삭제에 실패했습니다. 기존 구조는 유지됩니다.',
      researchDeleteSuccess: '연구지점을 삭제했습니다.',
      subregionTitleRequired: '서브지역 제목을 입력해 주세요.',
      subregionAddNotFound: '서브지역을 추가할 지역을 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      subregionInvalid: '서브지역 제목을 확인해 주세요.',
      subregionAddFailed: '서브지역 추가에 실패했습니다.',
      subregionResponseMissing: '서브지역 추가 응답을 확인하지 못했습니다.',
      subregionAddSuccess: '서브지역을 추가했습니다.',
      savedTitleDescriptionDifficulty: '제목/설명/난이도',
      savedTitleDescription: '제목/설명',
      savedMemo: '운영 메모',
      pointSaved: (parts) => `지역 ${parts}를 저장했습니다.`,
      regionSaved: (parts) => `그룹 ${parts}를 저장했습니다.`,
      researchSaved: '연구지점 제목과 템플릿을 저장했습니다.',
      templateLabel: researchTemplateLabelKo,
      templateDescription: (templateType) => `${researchTemplateLabelKo(templateType)} 템플릿으로 기록하는 연구지점입니다.`,
      defaultResearchTitle: (templateType, index) => `${researchTemplateLabelKo(templateType)} 연구지점 ${index}`,
    },
    resources: {
      titleLabel: '자료 제목',
      titleValid: '제목이 입력되었습니다.',
      titleRequired: '자료를 구분할 제목이 필요합니다.',
      urlLabel: 'URL 형식',
      urlValid: 'URL 형식이 안전합니다.',
      urlInvalid: 'http:// 또는 https://로 시작해야 합니다.',
      memoLabel: '활용 메모',
      memoValid: '학습 맥락 메모가 있습니다.',
      memoOptional: '선택 사항입니다. 어느 단계에서 볼지 적어두면 좋습니다.',
      authRequired: '로그인 상태를 확인한 뒤 다시 시도해 주세요.',
      parseFailed: 'URL 메타 확인에 실패했습니다. URL을 다시 확인해 주세요.',
      parseFailedKeep: 'URL 메타 확인에 실패했습니다. 기존 입력은 유지됩니다.',
      parsePartial: 'URL은 확인했지만 일부 메타 정보를 읽지 못했습니다. 저장은 아직 하지 않습니다.',
      parseSuccess: 'URL 메타 정보를 확인했습니다. 저장은 아직 하지 않습니다.',
      contentInvalid: '제목과 URL을 다시 확인해 주세요. 콘텐츠는 아직 저장되지 않았습니다.',
      contentSaveFailed: '콘텐츠 저장에 실패했습니다. 잠시 후 다시 시도해 주세요.',
      contentSaveFailedKeep: '콘텐츠 저장에 실패했습니다. 기존 입력은 유지됩니다.',
      contentResponseMissing: '콘텐츠 저장 응답을 확인하지 못했습니다.',
      contentCreated: '콘텐츠로 저장했습니다. 아직 선택 지역에는 연결하지 않았습니다.',
      attachNotFound: '선택한 지역 또는 콘텐츠를 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      attachConflict: '이미 이 지역에 연결된 자료입니다. 기존 노트북 상태는 유지됩니다.',
      attachInvalidState: '연결 상태 값을 확인하지 못했습니다. 기존 노트북 상태는 유지됩니다.',
      attachFailed: '지역 연결에 실패했습니다. 기존 노트북 상태는 유지됩니다.',
      attachFailedKeep: '지역 연결에 실패했습니다. 기존 노트북 상태는 유지됩니다.',
      attachResponseMissing: '지역 연결 응답을 확인하지 못했습니다. 기존 노트북 상태는 유지됩니다.',
      contentCreatedAndAttached: '콘텐츠 저장과 지역 연결을 완료했습니다.',
      attachSuccess: '선택 지역에 selected 자료로 연결했습니다.',
      attachCandidateMessage: '직접 추가한 content를 선택 지역에 selected 자료로 연결했습니다.',
      insufficientPoints: '포인트가 부족해 추천 후보를 새로 찾지 못했습니다.',
      candidateRefreshFailed: '추천 후보 교체에 실패했습니다. 잠시 후 다시 시도해 주세요.',
      candidateRefreshFailedKeep: '추천 후보 교체에 실패했습니다. 기존 노트북 상태는 유지됩니다.',
      candidateResponseMissing: '추천 후보 응답을 확인하지 못했습니다. 기존 노트북 상태는 유지됩니다.',
      selectNotFound: '선택할 자료를 찾지 못했습니다. 새로고침 후 다시 확인해 주세요.',
      selectCandidateOnly: 'candidate 자료만 선택 확정할 수 있습니다.',
      selectFailed: '선택 확정에 실패했습니다. 기존 노트북 상태는 유지됩니다.',
      selectFailedKeep: '선택 확정에 실패했습니다. 기존 노트북 상태는 유지됩니다.',
      selectResponseMissing: '선택 확정 응답을 확인하지 못했습니다. 기존 노트북 상태는 유지됩니다.',
      selectSuccess: 'candidate 자료 1개를 selected로 확정했습니다.',
      costPreview: (cost) => `예상 비용 ${cost}pt 기준으로`,
      selectedRegionBasis: '선택 지역 기준으로',
      allRegionBasis: '전체 지역 기준으로',
      candidateRefreshSelected: (basis, count) => `${basis} candidate ${count}개를 다시 불러왔습니다.`,
      candidateRefreshAll: (basis, count) => `${basis} candidate ${count}개를 전체 지역에 다시 불러왔습니다.`,
    },
  },
  planning: {
    backgroundAlt: {
      smallPlanning: '탐험계획 - 작은 디바이스',
      smallJournal: '탐험일지 - 작은 디바이스',
      expandedPlanning: '탐험계획 - 지도 펼침',
      expandedJournal: '탐험일지 - 지도 펼침',
      collapsedPlanning: '탐험계획 - 지도 접힘',
      collapsedJournal: '탐험일지 - 지도 접힘',
    },
    inactiveBanner: '비활성화된 행성',
    collapseMap: '◀ 지도 접기',
    expandMap: '지도 펼치기 ▶',
    pointModal: {
      aria: '지점 학습 진입 확인',
      region: (name) => `지역 ${name}`,
      subRegion: (name) => `서브 지역 ${name}`,
      type: (value) => `유형 ${value === 'research' ? '연구지점' : '탐험지점'}`,
      question: '이 지점의 학습페이지로 이동할까요?',
      cancel: '취소',
      open: '첫 학습 지점 열기',
    },
    transition: {
      aria: '화면 전환 확인',
      toPlanningTitle: '탐험계획 편집 화면으로 이동할까요?',
      toJournalTitle: '첫 학습을 시작할까요?',
      toPlanningMessage: '탐험일지는 학습 실행 화면이고 탐험계획은 구조 편집 화면입니다. 편집 화면으로 이동하시겠습니까?',
      toJournalMessage: '탐험계획을 확인했다면 탐험일지에서 첫 지점을 열어 바로 학습을 시작할 수 있습니다. 계속할까요?',
      cancel: '취소',
      confirm: '확인',
    },
    rightPanel: {
      bookmarks: [
        { key: 'planning', label: '탐험계획', shortLabel: '계획' },
        { key: 'journal', label: '탐험일지', shortLabel: '일지' },
        { key: 'records', label: '탐험기록', shortLabel: '기록' },
        { key: 'artifacts', label: '탐험결과물', shortLabel: '결과' },
        { key: 'community', label: '커뮤니티', shortLabel: '커뮤' },
        { key: 'civilization', label: '문명발전현황', shortLabel: '문명' },
      ],
      inactiveTitle: '비활성 행성은 탐험계획만 사용할 수 있습니다',
      lockedTitle: (label) => `${label} 준비중`,
      planned: { ...plannedKo, eyebrow: 'Coming Soon', confirm: '확인' },
    },
    lumiGuide: {
      eyebrow: 'Lumi Guide',
      inactive: '이 행성(코스)은 비활성화되었습니다. 활성화 버튼으로 다시 탐험일지를 이어갈 수 있어요.',
      unsaved: '첫 학습으로 넘어가기 전, 바뀐 탐험계획을 먼저 저장해 주세요.',
      course: '행성(코스)명을 수정하거나 비활성화하고, 지역(레슨)을 추가할 수 있어요.',
      region: '지역(레슨)의 제목이나 위치를 바꾸고, 탐험요소[서브 지역(서브 레슨), 지점(학습콘텐츠)]를 추가할 수 있어요.',
      subregion: '서브 지역(서브 레슨)의 제목이나 위치를 바꾸고, 탐험요소[지점(학습콘텐츠)]를 추가할 수 있어요.',
      explorationNode: '첫 학습에서 볼 콘텐츠 지점입니다. 제목과 위치를 확인해 시작 흐름을 준비해 주세요.',
      researchNode: '직접 기록하며 배울 연구지점입니다. 첫 학습 전에 제목과 위치를 가볍게 정리해 주세요.',
      default: '탐험계획은 첫 학습을 여는 지도입니다. 지역과 지점을 확인한 뒤 준비되면 탐험일지에서 시작하세요.',
    },
    map: {
      headerBadge: 'EXPLORER DIARY',
      courseMapTitle: '탐험 계획 지도',
      regionMapFallback: '지역 내부 지도',
      subregionMapFallback: '서브지역 내부 지도',
      returnToStarSystemAria: '스타 시스템으로 돌아가기',
      returnToStarSystemLine1: '← 스타',
      returnToStarSystemLine2: '시스템',
      returnToCourseMapAria: '전체 지역 지도로 돌아가기',
      returnToCourseMapLine1: '← 전체',
      returnToCourseMapLine2: '지역 지도',
      returnToRegionMapAria: '지역 지도로 돌아가기',
      returnToRegionMapLine1: '← 지역',
      returnToRegionMapLine2: '지도',
      emptyCourseLine1: '아직 구성된 지역이 없습니다.',
      emptyCourseLine2: '지역을 추가하면 여기에 탐험 경계가 표시됩니다.',
      missingSubregion: '선택한 서브지역 정보를 찾지 못했습니다.',
      missingRegion: '선택한 지역 정보를 찾지 못했습니다.',
      exitTitle: '스타 시스템으로 돌아갈까요?',
      exitEyebrow: 'Explorer Exit',
      exitMessage: '지금 탐험계획 화면을 나가면 스타 시스템으로 돌아갑니다.\n저장되지 않은 편집 내용이 있다면 먼저 확인하세요.',
      exitCancel: '계속 보기',
      exitConfirm: '스타 시스템으로 이동',
      previousPageAria: '이전 지역 페이지',
      nextPageAria: '다음 지역 페이지',
      subregionLegend: '서브지역',
      explorationLegend: '탐험지점',
      researchLegend: '연구지점',
      formatRegionTitle: (index, name) => `${index + 1}. "${name}" 지역`,
      formatRegionMeta: (subregions, exploration, research) => `📁 ${subregions}개   🎯 ${exploration}개   🔬 ${research}개`,
      formatRegionAria: (title, subregions, exploration, research) => `${title}, 서브지역 ${subregions}개, 탐험지점 ${exploration}개, 연구지점 ${research}개`,
      regionEmptyLine1: '아직 구성된 서브지역이 없습니다.',
      regionEmptyLine2: '오른쪽 트리에서 서브지역이나 지점을 추가하면',
      regionEmptyLine3: '이 지역 지도에 표시됩니다.',
      subregionEmptyLine1: '아직 구성된 지점이 없습니다.',
      subregionEmptyLine2: '오른쪽 트리에서 탐험지점이나 연구지점을 추가하면',
      subregionEmptyLine3: '이 서브지역 지도에 표시됩니다.',
      formatSubregionNodeCounts: (exploration, research) => `🎯 ${exploration}개  🔬 ${research}개`,
      formatRegionSummary: (subregions, exploration, research) => [`📁 서브지역 ${subregions}개`, `🎯 탐험지점 ${exploration}개`, `🔬 연구지점 ${research}개`],
      formatSubregionSummary: (exploration, research) => [`🎯 탐험지점 ${exploration}개`, `🔬 연구지점 ${research}개`],
    },
  },
  supportSections: {
    common: {
      currentRegion: '현재 지역',
      currentPoint: '현재 지점',
      currentStage: '현재 단계',
      regionNotSelected: '지역 미선택',
      pointNotSelected: '지점 미선택',
      selectRegionFirst: '먼저 지역을 선택해 주세요.',
      reset: '원본 복원',
      draftPreview: 'draft 프리뷰',
    },
    records: {
      eyebrow: 'Record Summary',
      title: '탐험기록',
      statusChanged: '로컬 기록 수정 있음',
      statusConnected: '지점 기록 연결됨',
      statusNeedRegion: '지역 선택 필요',
      regionDescriptionFallback: '먼저 지역을 선택하면 그 안의 지점 학습 데이터를 정리할 수 있습니다.',
      pointDescriptionFallback: '탐험지점이나 연구지점이 선택되면 해당 지점의 성장 데이터를 중심으로 기록합니다.',
      focusMinutes: '집중 시간 (분)',
      practiceCount: '연습 횟수',
      confidence: '이해도',
      applicationNote: '적용 메모',
      applicationPlaceholder: '직접 적용해 본 예시, 막힌 지점, 다음 반복 지점을 적어 주세요.',
      applicationDisabledPlaceholder: '먼저 지역을 선택해 주세요.',
      saving: '기록 저장 중...',
      save: '탐험기록 저장',
      noticeReady: '탐험기록은 현재 지점 학습의 성장 데이터를 정리합니다. 현재 구현 저장소는 지역 단위를 사용하지만, 기록 의미는 지점 기준입니다.',
      noticeNeedRegion: '지역을 선택하면 해당 지점의 성장 데이터를 기록할 수 있습니다.',
    },
    artifacts: {
      eyebrow: 'Artifact Entry',
      title: '탐험결과물',
      statusChanged: '로컬 결과물 수정 있음',
      statusConnected: '지점 결과물 연결됨',
      statusNeedRegion: '지역 선택 필요',
      regionDescriptionFallback: '먼저 지역을 선택하면 그 안의 지점 학습 결과를 정리할 수 있습니다.',
      pointDescriptionFallback: '탐험지점이나 연구지점이 선택되면 그 학습 결과물을 중심으로 정리합니다.',
      type: '결과물 유형',
      typeOptions: [
        { value: 'note', label: '노트' },
        { value: 'code', label: '코드' },
        { value: 'project', label: '프로젝트' },
        { value: 'writing', label: '글' },
        { value: 'design', label: '디자인' },
        { value: 'link', label: '외부 링크' },
      ],
      artifactTitle: '결과물 제목',
      link: '결과물 링크',
      description: '결과물 설명',
      descriptionPlaceholder: '이 지점 학습 결과물이 무엇을 보여주는지 적어 주세요.',
      descriptionDisabledPlaceholder: '먼저 지역을 선택해 주세요.',
      saving: '결과물 저장 중...',
      save: '탐험결과물 저장',
      noticeReady: '탐험결과물은 현재 지점 학습의 산출물 메타데이터입니다. 현재 구현 저장소는 지역 단위를 사용하지만, 결과 의미는 지점 기준입니다.',
      noticeNeedRegion: '지역을 선택하면 해당 지점의 결과물 메타데이터를 저장할 수 있습니다.',
    },
    community: {
      eyebrow: 'Community Preview',
      title: '커뮤니티',
      status: 'learning 이후 본격 활성',
      pointDescription: '해당 지점 관련 질문, 답변, 팁, 결과물 피드백을 한곳에서 찾는 지식 허브입니다.',
      stageDescription: 'learning diary에 들어가면 질문 작성과 피드백 요청이 열립니다. 지금은 구조와 진입 안내만 먼저 연결합니다.',
      items: [
        { title: '열릴 기능', body: '지점 기준 질문/답변 검색, creator 도움 요청, 결과물 피드백 축' },
        { title: '현재 프리뷰 이유', body: '커뮤니티는 learning diary 진입과 함께 활성하는 정책이라 draft 단계에서는 쓰기 기능을 열지 않습니다.' },
        { title: '후속 연결 기준', body: 'learning diary 런타임 연결 뒤 질문 작성과 지점 검색 결과를 붙입니다.' },
      ],
    },
    civilization: {
      eyebrow: 'Civilization Preview',
      title: '문명발전현황',
      status: 'shared 이후 본격 활성',
      stageDescription: '문명발전현황은 완료 후 `공유적용`이 끝난 shared 상태에서 확산 데이터를 집계하는 화면입니다.',
      metricLabel: '예정 지표',
      metricValue: '복제 수 · 반응 수 · 확산 단계',
      metricDescription: '지금은 진입 조건과 집계 항목만 미리 보여주고, 실제 수치는 shared API 분리 이후 연결합니다.',
      items: [
        { title: '활성 조건', body: '`completed + shared` 이후부터 문명 확산 단계와 공개 반응 지표를 보여줍니다.' },
        { title: '현재 프리뷰 이유', body: 'draft 단계에서는 집계할 공개 데이터가 없으므로 구조와 목적만 먼저 연결합니다.' },
        { title: '후속 연결 기준', body: '완료 화면의 `공유적용` 흐름과 shared 상태 집계 API가 준비되면 실제 데이터로 전환합니다.' },
      ],
    },
    selectedPoint: {
      researchLabel: '연구지점',
      researchTemplateDescription: (template) => `${template} 템플릿으로 기록하는 학습자 노트 지점입니다.`,
      researchTemplateLabel: (templateType) => {
        switch (templateType) {
          case 'concept_summary':
            return '개념 요약'
          case 'practice_strategy':
            return '실습 전략'
          case 'problem_solving':
            return '문제 해결'
          case 'free_research':
            return '자유 연구'
          default:
            return null
        }
      },
      researchDescription: '학습자가 직접 기록을 남기는 연구지점입니다.',
      explorationLabel: '탐험지점',
      selectedExploration: '선택된 탐험지점',
      candidateExploration: '추천된 탐험지점 후보',
      genericExploration: '탐험지점',
      formatExplorationDescription: (label, resourceType) => `${label} · ${resourceType}`,
    },
  },
  tree: {
    loading: '불러오는 중...',
    emptyRegions: '아직 탐험 지역이 없습니다.',
    formatPlanTitle: (title) => {
      const trimmed = title?.trim()
      return trimmed ? `행성 ${trimmed} 탐험 계획` : '행성 탐험 계획'
    },
    formatPlanetName: (title) => title?.trim() || '이름 없음',
    planetPrefix: '행성',
    actionLabel: '작업',
    pendingDelete: '삭제 예정',
    pendingCreate: '추가 예정',
    openLinkTitle: '새 창에서 바로가기',
    regionLabel: (index, name) => (/^지역\s*\d+$/u.test(name) ? `지역 ${index}` : `지역 ${index} · ${name}`),
    regionPlaceholder: (index) => `지역 ${index}`,
    collapseRegion: '지역 접기',
    expandRegion: '지역 펼치기',
    subRegionLabel: (index, name) => (/^(서브|서브지역)\s*\d+$/u.test(name) ? `서브 ${index}` : `서브 ${index} · ${name}`),
    subRegionPlaceholder: (index) => `서브 ${index}`,
    collapseSubRegion: '서브지역 접기',
    expandSubRegion: '서브지역 펼치기',
    actionBar: {
      save: '저장',
      cancel: '취소',
      delete: '삭제',
      undoDelete: '삭제 취소',
      activatePlanet: '행성 활성화',
      activate: '활성화',
      hideInactive: '비활성 숨김',
      showInactive: '비활성 보기',
      inactiveActivateTitle: '비활성화된 행성을 다시 활성화하고 탐험일지로 이동합니다',
      activateTitle: '코스를 활성화하고 탐험일지로 이동합니다',
      undoDeleteTitle: '선택 항목의 삭제 예정 표시를 취소합니다',
      deleteCourseTitle: '코스를 비활성화하고 갤럭시 대시보드로 이동합니다',
      deleteWithChildrenTitle: '하위 항목도 함께 삭제 예정으로 표시됩니다. 학습 시작 전 코스는 저장 시 완전히 삭제됩니다',
      inactiveToggleTitle: '비활성 항목 표시를 전환합니다',
      reactivateTitle: '선택 항목을 다시 활성화합니다',
    },
    editPanel: {
      action: '작업',
      pendingDeleteNotice: '삭제 예정 항목입니다. 저장하면 반영되고, 삭제 취소로 되돌릴 수 있습니다.',
      pendingCreateNotice: '추가 예정 항목입니다. 저장하면 탐험계획에 반영됩니다.',
      renamePlanet: '행성명 변경',
      addRegion: '➕ 지역',
      renameRegion: '지역명 수정',
      renameSubRegion: '서브지역명 수정',
      addChild: '➕ 탐험요소',
      editNode: '지역 수정',
      moveUp: '한 칸 위로 이동',
      moveDown: '한 칸 아래로 이동',
      modals: {
        close: '닫기',
        cancel: '취소',
        apply: '적용',
        edit: '수정',
        add: '추가',
        change: '변경',
        check: '점검',
        checking: '확인중',
        valid: '✓ 유효',
        error: '✗ 오류',
        openLink: '새 창에서 바로가기',
        internalContent: '내부 콘텐츠',
        urlMustPass: 'URL 점검을 통과해야 추가할 수 있습니다.',
        parentPath: (value) => `상위 지역(리슨): ${value}`,
        recommendation: {
          aria: '콘텐츠 추천',
          title: '콘텐츠 추천',
          queryLabel: '추천 기준',
          queryPlaceholder: '콘텐츠 제목, 설명, 작성자',
          search: '콘텐츠 추천',
          searching: '추천중',
          recommendAgain: '재추천',
          empty: '콘텐츠 추천 버튼을 눌러 후보를 불러오세요.',
          pointErrorEmpty: '포인트를 충전 후 다시 시도해 주세요.',
          loading: '후보를 불러오고 있는 중입니다.',
          noResults: '일치하는 콘텐츠가 없습니다. 추천 기준을 바꿔보세요.',
          loaded: (count) => `${count}개 후보를 불러왔습니다.`,
          missingCandidate: '제목 또는 콘텐츠 정보가 없습니다. 다른 콘텐츠를 선택해주세요.',
          loginRequired: '로그인 상태를 확인한 뒤 다시 시도해 주세요.',
          insufficientPoints: (balance, cost) => `포인트가 부족합니다. 현재 ${balance}pt (필요: ${cost}pt)`,
          loadFailed: '콘텐츠 후보를 불러오지 못했습니다.',
        },
        addExploration: {
          aria: '탐험지점 추가',
          title: '탐험지점 추가',
          recommendationTab: '콘텐츠 추천',
          urlTab: 'URL 추가',
          titleLabel: '탐험지점 제목',
          titlePlaceholder: '탐험지점 제목을 입력하세요',
          urlLabel: '링크 URL',
          addButton: '탐험지점 추가',
        },
        addChild: {
          aria: '탐험요소 추가',
          title: '탐험요소 추가',
          tabAria: '탐험요소 종류 선택',
          subregionTab: '서브리슨',
          explorationTab: '탐험지점',
          researchTab: '연구지점',
          subregionDescription: '서브리슨은 선택한 리슨을 더 작은 학습 흐름으로 나누는 하위 리슨입니다. 관련 탐험지점과 연구지점을 묶어 학습 순서를 정리할 때 사용합니다.',
          subregionLimit: '서브리슨은 리슨에서만 추가할 수 있고 최대 3개까지 만들 수 있습니다.',
          subregionNameLabel: '서브리슨 이름',
          subregionNamePlaceholder: '서브리슨 이름을 입력하세요',
          addSubregion: '서브리슨 추가',
          explorationDescription: '탐험지점은 학습자가 실제로 열어 보고 학습할 외부 콘텐츠입니다. 추천으로 후보를 고르거나, 직접 찾은 URL을 점검한 뒤 추가할 수 있습니다.',
          explorationMethodAria: '탐험지점 추가 방식 선택',
          recommendationDescription: '콘텐츠 추천은 현재 지역/목표 문맥과 추천 기준을 바탕으로 학습에 쓸 후보 자료를 찾아오는 방식입니다. 후보를 검토한 뒤 `적용`하면 탐험지점으로 추가됩니다.',
          urlDescription: 'URL 직접 추가는 학습자가 직접 찾은 자료를 탐험지점으로 넣는 방식입니다. 링크 점검을 통과해야 저장 전 목록에 추가할 수 있습니다.',
          urlDirectTab: 'URL 직접 추가',
          researchDescription: '연구지점은 학습자가 직접 조사하거나 정리할 학습 콘텐츠 자리입니다. 자료를 나중에 직접 입력하고 확정한 뒤 학습을 진행할 수 있습니다.',
          researchTitleLabel: '연구지점 제목',
          researchTitlePlaceholder: '연구지점 제목을 입력하세요',
          addResearch: '연구지점 추가',
        },
        addUrl: {
          explorationTitle: '탐험지점 추가 (URL)',
          researchTitle: '연구지점 추가',
        },
        nodeEdit: {
          explorationAria: '탐험지점 수정',
          researchAria: '연구지점 수정',
          explorationTitle: '탐험지점 수정',
          researchTitle: '연구지점 수정',
          titleLabel: '제목',
          explorationPlaceholder: '탐험지점 제목',
          researchPlaceholder: '연구지점 제목',
          linkLabel: '링크',
          recommend: '추천',
        },
        deleteConfirm: {
          deleteAria: '삭제 확인',
          activateAria: '활성화 확인',
          alertAria: '삭제 불가 안내',
          warningEyebrow: 'Lumi Warning',
          alertTitle: '삭제할 수 없어요',
          alertSuffix: '리슨 하나는 남겨두고 계속 구조를 다듬어 주세요.',
          deleteTitle: '삭제 확인',
          activateTitle: '활성화 확인',
          targetLabel: '대상',
          courseConfirmLabel: '코스명 확인',
          exactTitleHelp: (action) => `행성명을 정확히 입력해야 ${action === 'activate' ? '활성화할' : '삭제할'} 수 있습니다.`,
          delete: '삭제',
          activate: '활성화',
        },
        courseRename: {
          aria: '행성명 변경',
          title: '행성명 변경',
          label: '행성명',
          placeholder: '새 행성명을 입력하세요',
          help: '행성명은 탐험계획 제목과 지도 상단 이름에 함께 반영됩니다.',
        },
        regionRename: {
          regionTitle: '지역명 수정',
          subregionTitle: '서브지역명 수정',
          regionLabel: '지역명',
          subregionLabel: '서브지역명',
          aria: '지역명 수정',
          placeholder: (label) => `${label}을 입력하세요`,
          help: '변경한 지역명은 탐험계획 저장 전까지 화면에만 표시됩니다.',
        },
        regionAdd: {
          aria: '지역 추가',
          title: '지역 추가',
          label: '지역 이름',
          description: '지역은 코스의 큰 학습 단계인 리슨입니다. 목표를 향해 순서대로 탐험할 주요 주제를 나눌 때 사용합니다.',
          placeholder: '새 지역 이름',
          help: '추가한 지역은 탐험계획 저장 전까지 화면에만 표시됩니다.',
          addButton: '지역 추가',
        },
        messages: {
          urlValid: '유효한 링크입니다.',
          urlInvalid: '유효하지 않은 링크입니다.',
          urlCheckFailed: '링크 확인 실패',
          contentIdBased: '콘텐츠 ID 기반 탐험지점입니다.',
          deleteCourse: '코스를 비활성화하고 갤럭시 대시보드로 이동합니다. 진행하려면 코스명을 그대로 입력해 주세요.',
          activateCourse: '비활성화된 코스를 다시 활성화하고 탐험일지로 이동합니다. 진행하려면 행성명을 그대로 입력해 주세요.',
          lastRegionTitle: '지역',
          lastRegionBlock: '마지막 리슨은 삭제 예정으로 표시할 수 없습니다. 빈 리슨 상태를 막기 위해 최소 1개의 리슨은 남겨두어야 합니다.',
          deleteRegionWithChildren: '이 지역과 하위 항목을 삭제 예정으로 표시합니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.',
          deleteRegion: '이 지역을 삭제 예정으로 표시합니다.',
          deleteSubregionWithChildren: '이 서브지역과 하위 지점을 삭제 예정으로 표시합니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.',
          deleteSubregion: '이 서브지역을 삭제 예정으로 표시합니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.',
          deleteNode: '이 지점을 삭제 예정으로 표시합니다. 학습 시작 전 코스는 최종 저장 시 완전히 삭제됩니다.',
          fallbackRegion: '지역',
          fallbackSubregion: '서브지역',
        },
      },
    },
  },
}

const en: DashboardCourseDraftCopy = {
  page: {
    loading: 'Preparing Explorer Diary...',
    notFound: 'Could not find this learning plan.',
    invalidDraftId: 'Invalid draft ID.',
    loadFailed: 'Could not load the learning plan.',
    backToDashboard: 'Back to Galaxy',
    headerTitle: 'Explorer Diary',
    headerSubtitle: 'Shared Explorer Layout',
    sectionTitle: (title) => `${title || 'Learning'} Exploration Plan`,
    fallbackCourseName: 'Learning',
    returnToDiary: 'Back to Diary',
    goal: {
      planningEyebrow: 'Goal Context',
      journalEyebrow: 'Journal Context',
      planningTitle: 'Exploration Goal',
      journalTitle: 'Goal and Current Status',
      edit: 'Edit Goal',
      set: 'Set Goal',
      loading: 'Loading goal context.',
      empty: 'No goal is connected to this plan yet. Set a goal first to make lessons and recommendations clearer.',
      currentRegion: (name) => `Current region ${name}`,
      currentSubRegion: (name) => `Current sub-region ${name}`,
      journalPreLearning: 'Pre-learning edit path',
      journalPointReady: 'Select a point to open its learning screen',
      journalPointPending: 'Preparing the point learning screen link',
      journalGuideCanOpen: 'Click a point and Lumi will confirm before opening the learning screen.',
      journalGuideExplore: 'Browse the tree and map to review the current exploration structure.',
      usage: (value) => `Context ${value}`,
      motivation: (value) => `Reason ${value}`,
      learnerEditDecision: 'Choice: Learner Edit',
      rebuildAllDecision: 'Choice: Rebuild All',
      statusLabel: (status) => {
        if (status === 'learning') return 'Learning'
        if (status === 'confirmed') return 'Ready'
        if (status === 'archived') return 'Completed'
        return 'Planning'
      },
      planningBeforeStart: 'Before learning starts, you can adjust the goal and choose whether to rebuild everything or keep editing the current structure. Rebuild all costs the same points as course creation.',
      planningChangedBeforeStart: 'The goal changed. Before learning starts, choose whether to rebuild all or keep learner editing. Rebuild all costs the same points as course creation.',
      planningChangedLearning: 'The goal changed. While learning, choose whether to rebuild all or keep learner editing. Rebuild all costs the same points as course creation.',
      status: (value) => `Planet status ${value}`,
      completedRegions: (completed, total) => `Regions ${completed}/${total}`,
      completedPoints: (completed, total, ec, et, rc, rt) => `Points ${completed}/${total}  🎯 ${ec}/${et}  🔬 ${rc}/${rt}`,
      completionRate: (percent) => `Completion ${percent}%`,
      regionProgress: (completed, total) => `Completed regions ${completed} / ${total}`,
      learningRegions: (count) => `In progress ${count}`,
      remainingRegions: (count) => `Remaining ${count}`,
    },
    goalRevision: {
      aria: 'Adjust exploration goal',
      eyebrow: 'Goal Revision',
      title: 'Adjust Exploration Goal',
      description: 'If you adjust the goal, later lesson structure and content recommendations will follow the new goal context.',
      noInterviewTitle: 'No goal interview is connected yet.',
      noInterviewDescription: 'You can start a goal interview from the current plan title or search phrase.',
      start: 'Start Goal Interview',
      starting: 'Starting...',
      reviseSeed: 'I want to revise my goal.',
      rebuildTitle: 'How should this plan continue?',
      rebuildDescriptionBeforeStart: 'Before learning starts, choose whether to generate new lessons from the new goal or keep editing the current structure yourself. Rebuild all deactivates existing regions and points and costs 5 points. BYOK LLM users generate with their personal key.',
      rebuildDescriptionLearning: 'While learning, choose whether to generate new lessons from the new goal or keep editing the current structure yourself. Rebuild all deactivates existing regions and points and costs 5 points. BYOK LLM users generate with their personal key.',
      rebuildOptions: {
        rebuildAll: { label: 'Rebuild All', desc: 'Generate new lessons and deactivate existing regions and points. This costs 5 points unless BYOK LLM is active.' },
        keepStructure: { label: 'Learner Edit', desc: 'Keep the current structure and edit lessons and resources yourself.' },
      },
      confirmedAction: 'Revise Goal',
      confirmedActionBusy: 'Applying goal...',
      confirmedMessage: 'The new goal is confirmed. Choose how to continue the exploration plan.',
      applySuccess: 'The new goal was applied to the exploration plan.',
      applyFailed: 'Could not apply the goal.',
      decisionSuccess: 'The goal change was applied to the exploration plan.',
      cancelSuccess: 'Goal change was cancelled and the previous goal was restored.',
      close: 'Close',
    },
    rebuildLoading: {
      aria: 'Rebuilding exploration plan',
      eyebrow: 'Lumi is rebuilding',
      title: 'Rebuilding lessons for the new goal',
      message: 'AI is applying the new goal to the exploration plan. Existing structure will be deactivated while new lessons are prepared.',
      steps: ['Apply goal', 'Generate lessons', 'Update plan'],
    },
    returnConfirm: {
      aria: 'Confirm return to diary',
      eyebrow: 'Lumi Confirm',
      title: 'Return to Explorer Diary?',
      message: 'You will leave the plan editor and return to the learning diary. Continue?',
      cancel: 'Cancel',
      confirm: 'Confirm',
    },
    planned: { ...plannedEn, eyebrow: 'Coming Soon', confirm: 'OK' },
    overview: {
      descriptionFallback: (sourceQuery) => `This first planet exploration plan was generated from the ${sourceQuery} goal.`,
      sourceLabel: 'Exploration Topic',
      levelLabel: 'Current Level',
      durationLabel: 'Recommended Duration',
      weeklyHoursLabel: 'Weekly Hours',
      unset: 'Not set',
      weeks: (count) => `${count} weeks`,
      hours: (count) => `${count} hours`,
      lumiTitle: 'Current Draft State',
      editEyebrow: 'Safe Edit',
      editTitle: 'Plan Basics',
      editHint: 'Only the title and description are saved. The exploration structure is not changed.',
      titleLabel: 'Plan Title',
      descriptionLabel: 'Plan Description',
      saving: 'Saving...',
      save: 'Save Basics',
    },
  },
  hook: {
    saveTitleRequired: 'Enter a learning plan title.',
    saveFailed: 'Could not save the learning plan.',
    saveSuccess: 'Plan basics saved.',
    loginPath: (draftId) => `/en/login?redirect_after=/dashboard/course-drafts/${draftId}`,
    agreementsPath: (draftId) => `/en/agreements?redirect_after=/dashboard/course-drafts/${draftId}`,
    entries: {
      authRequired: 'Check your login status and try again.',
      targetNotFound: 'Could not find the selected record target. Refresh and check again.',
      journalContentRequired: 'Add at least one journal field before saving.',
      journalSaveFailed: 'Could not save the diary entry. Your current input is preserved.',
      journalResponseMissing: 'Could not verify the diary save response.',
      journalSaveSuccess: 'Diary entry saved. You can continue this point context the next time you return.',
      recordInvalid: 'Check the record value range again.',
      recordSaveFailed: 'Could not save the record. Your current input is preserved.',
      recordResponseMissing: 'Could not verify the record save response.',
      recordSaveSuccess: 'Record saved. You can review it again as growth data for this point.',
      artifactTargetNotFound: 'Could not find the selected artifact target. Refresh and check again.',
      artifactInvalid: 'Check the artifact title and URL format again.',
      artifactSaveFailed: 'Could not save the artifact. Your current input is preserved.',
      artifactResponseMissing: 'Could not verify the artifact save response.',
      artifactSaveSuccess: 'Artifact saved. You can review it again as output metadata for this point.',
    },
    detail: {
      savePoint: 'Save Region Details',
      savingPoint: 'Saving Region...',
      saveResearch: 'Save Research Point',
      savingResearch: 'Saving Research Point...',
      saveRegion: 'Save Group Details',
      savingRegion: 'Saving Group...',
      draftMissing: 'Could not verify the exploration plan.',
      authRequired: 'Check your login status and try again.',
      memoPointNotFound: 'Could not find the selected region memo target.',
      memoRegionNotFound: 'Could not find the selected group memo target.',
      memoInvalid: 'Check the operation memo value.',
      memoSaveFailed: 'Could not save the operation memo. Your current input is preserved.',
      memoResponseMissing: 'Could not verify the memo save response.',
      pointTitleRequired: 'Enter a region name.',
      pointNotFound: 'Could not find the selected region. Refresh and check again.',
      pointSaveFailed: 'Could not save region details. The current structure is preserved.',
      pointResponseMissing: 'Could not verify the region save response.',
      regionTitleRequired: 'Enter a group name.',
      regionNotFound: 'Could not find the selected group. Refresh and check again.',
      regionSaveFailed: 'Could not save group details. The current structure is preserved.',
      regionResponseMissing: 'Could not verify the group save response.',
      researchTitleRequired: 'Enter a research point title.',
      researchNotFound: 'Could not find the selected research point. Refresh and check again.',
      researchInvalid: 'Check the research point title and template.',
      researchSaveFailed: 'Could not save the research point. The current structure is preserved.',
      researchResponseMissing: 'Could not verify the research point save response.',
      structureSaveFirstAdd: 'Save the exploration structure before adding a research point.',
      structureSaveFirstDelete: 'Save the exploration structure before deleting a research point.',
      researchAddNotFound: 'Could not find the region for the new research point. Refresh and check again.',
      researchAddFailed: 'Could not add the research point. The current structure is preserved.',
      researchAddResponseMissing: 'Could not verify the research point add response.',
      researchAddSuccess: 'Research point added. Check its title and template.',
      researchDeleteNotFound: 'Could not find the research point to delete. Refresh and check again.',
      researchDeleteFailed: 'Could not delete the research point. The current structure is preserved.',
      researchDeleteSuccess: 'Research point deleted.',
      subregionTitleRequired: 'Enter a sub-region title.',
      subregionAddNotFound: 'Could not find the region for the new sub-region. Refresh and check again.',
      subregionInvalid: 'Check the sub-region title.',
      subregionAddFailed: 'Could not add the sub-region.',
      subregionResponseMissing: 'Could not verify the sub-region add response.',
      subregionAddSuccess: 'Sub-region added.',
      savedTitleDescriptionDifficulty: 'title/description/difficulty',
      savedTitleDescription: 'title/description',
      savedMemo: 'operation memo',
      pointSaved: (parts) => `Saved region ${parts}.`,
      regionSaved: (parts) => `Saved group ${parts}.`,
      researchSaved: 'Research point title and template saved.',
      templateLabel: researchTemplateLabelEn,
      templateDescription: (templateType) => `${researchTemplateLabelEn(templateType)} template research point.`,
      defaultResearchTitle: (templateType, index) => `${researchTemplateLabelEn(templateType)} Research Point ${index}`,
    },
    resources: {
      titleLabel: 'Resource Title',
      titleValid: 'Title entered.',
      titleRequired: 'A title is required to identify this resource.',
      urlLabel: 'URL Format',
      urlValid: 'URL format is safe.',
      urlInvalid: 'Start with http:// or https://.',
      memoLabel: 'Usage Memo',
      memoValid: 'Learning context memo is present.',
      memoOptional: 'Optional. Add where this resource should be used.',
      authRequired: 'Check your login status and try again.',
      parseFailed: 'Could not check URL metadata. Check the URL again.',
      parseFailedKeep: 'Could not check URL metadata. Your current input is preserved.',
      parsePartial: 'The URL was reachable, but some metadata could not be read. It has not been saved yet.',
      parseSuccess: 'URL metadata checked. It has not been saved yet.',
      contentInvalid: 'Check the title and URL. The content has not been saved yet.',
      contentSaveFailed: 'Could not save the content. Try again later.',
      contentSaveFailedKeep: 'Could not save the content. Your current input is preserved.',
      contentResponseMissing: 'Could not verify the content save response.',
      contentCreated: 'Saved as content. It is not attached to the selected region yet.',
      attachNotFound: 'Could not find the selected region or content. Refresh and check again.',
      attachConflict: 'This resource is already attached to the region. The current notebook state is preserved.',
      attachInvalidState: 'Could not verify the attachment state. The current notebook state is preserved.',
      attachFailed: 'Could not attach the resource to the region. The current notebook state is preserved.',
      attachFailedKeep: 'Could not attach the resource to the region. The current notebook state is preserved.',
      attachResponseMissing: 'Could not verify the attach response. The current notebook state is preserved.',
      contentCreatedAndAttached: 'Content saved and attached to the region.',
      attachSuccess: 'Attached the selected resource to the selected region.',
      attachCandidateMessage: 'Directly added content was attached to the selected region as a selected resource.',
      insufficientPoints: 'Not enough points to refresh recommendations.',
      candidateRefreshFailed: 'Could not refresh recommendation candidates. Try again later.',
      candidateRefreshFailedKeep: 'Could not refresh recommendation candidates. The current notebook state is preserved.',
      candidateResponseMissing: 'Could not verify the recommendation response. The current notebook state is preserved.',
      selectNotFound: 'Could not find the resource to select. Refresh and check again.',
      selectCandidateOnly: 'Only candidate resources can be selected.',
      selectFailed: 'Could not confirm this selection. The current notebook state is preserved.',
      selectFailedKeep: 'Could not confirm this selection. The current notebook state is preserved.',
      selectResponseMissing: 'Could not verify the selection response. The current notebook state is preserved.',
      selectSuccess: 'Confirmed one candidate resource as selected.',
      costPreview: (cost) => `Estimated cost ${cost}pt`,
      selectedRegionBasis: 'for the selected region',
      allRegionBasis: 'for all regions',
      candidateRefreshSelected: (basis, count) => `Loaded ${count} candidate resources ${basis}.`,
      candidateRefreshAll: (basis, count) => `Loaded ${count} candidate resources ${basis}.`,
    },
  },
  planning: {
    backgroundAlt: {
      smallPlanning: 'Planning - small device',
      smallJournal: 'Diary - small device',
      expandedPlanning: 'Planning - map expanded',
      expandedJournal: 'Diary - map expanded',
      collapsedPlanning: 'Planning - map collapsed',
      collapsedJournal: 'Diary - map collapsed',
    },
    inactiveBanner: 'Inactive Planet',
    collapseMap: 'Collapse Map',
    expandMap: 'Expand Map',
    pointModal: {
      aria: 'Confirm point learning entry',
      region: (name) => `Region ${name}`,
      subRegion: (name) => `Sub-region ${name}`,
      type: (value) => `Type ${value === 'research' ? 'Research Point' : 'Exploration Point'}`,
      question: 'Open this point learning screen?',
      cancel: 'Cancel',
      open: 'Open First Step',
    },
    transition: {
      aria: 'Confirm view switch',
      toPlanningTitle: 'Move to plan editing?',
      toJournalTitle: 'Start the first learning step?',
      toPlanningMessage: 'Explorer Diary is the learning screen, and Exploration Plan is the structure editor. Move to the editor?',
      toJournalMessage: 'Once the plan looks ready, open the diary and start from the first point. Continue?',
      cancel: 'Cancel',
      confirm: 'Confirm',
    },
    rightPanel: {
      bookmarks: [
        { key: 'planning', label: 'Exploration Plan', shortLabel: 'Plan' },
        { key: 'journal', label: 'Explorer Diary', shortLabel: 'Diary' },
        { key: 'records', label: 'Records', shortLabel: 'Log' },
        { key: 'artifacts', label: 'Artifacts', shortLabel: 'Work' },
        { key: 'community', label: 'Community', shortLabel: 'Comm' },
        { key: 'civilization', label: 'Civilization View', shortLabel: 'Civ' },
      ],
      inactiveTitle: 'Inactive planets can only use the exploration plan.',
      lockedTitle: (label) => `${label} is not ready yet`,
      planned: { ...plannedEn, eyebrow: 'Coming Soon', confirm: 'OK' },
    },
    lumiGuide: {
      eyebrow: 'Lumi Guide',
      inactive: 'This planet is inactive. Use Activate Planet to continue the Explorer Diary.',
      unsaved: 'Before starting your first step, save the updated exploration plan.',
      course: 'Rename or deactivate the planet, and add regions.',
      region: 'Edit the region title or order, and add exploration elements such as sub-regions or points.',
      subregion: 'Edit the sub-region title or order, and add learning points.',
      explorationNode: 'This is a content point for your first learning step. Check the title and order before you begin.',
      researchNode: 'This is a research point where you will learn by writing. Check the title and order before you begin.',
      default: 'Your exploration plan is the map for the first learning step. Check the regions and points, then begin from the diary.',
    },
    map: {
      headerBadge: 'EXPLORER DIARY',
      courseMapTitle: 'Exploration Map',
      regionMapFallback: 'Region Map',
      subregionMapFallback: 'Sub-region Map',
      returnToStarSystemAria: 'Return to Star System',
      returnToStarSystemLine1: '← Star',
      returnToStarSystemLine2: 'System',
      returnToCourseMapAria: 'Return to all regions map',
      returnToCourseMapLine1: '← All',
      returnToCourseMapLine2: 'Regions',
      returnToRegionMapAria: 'Return to region map',
      returnToRegionMapLine1: '← Region',
      returnToRegionMapLine2: 'Map',
      emptyCourseLine1: 'No regions have been created yet.',
      emptyCourseLine2: 'Add a region to show exploration boundaries here.',
      missingSubregion: 'Could not find the selected sub-region.',
      missingRegion: 'Could not find the selected region.',
      exitTitle: 'Return to Star System?',
      exitEyebrow: 'Explorer Exit',
      exitMessage: 'You will leave the exploration plan and return to Star System.\nCheck any unsaved edits before leaving.',
      exitCancel: 'Keep Viewing',
      exitConfirm: 'Go to Star System',
      previousPageAria: 'Previous region page',
      nextPageAria: 'Next region page',
      subregionLegend: 'Sub-regions',
      explorationLegend: 'Exploration Points',
      researchLegend: 'Research Points',
      formatRegionTitle: (index, name) => `${index + 1}. ${name}`,
      formatRegionMeta: (subregions, exploration, research) => `📁 ${subregions}   🎯 ${exploration}   🔬 ${research}`,
      formatRegionAria: (title, subregions, exploration, research) => `${title}, ${subregions} sub-regions, ${exploration} exploration points, ${research} research points`,
      regionEmptyLine1: 'No sub-regions have been created yet.',
      regionEmptyLine2: 'Add sub-regions or points from the tree on the right',
      regionEmptyLine3: 'to show them on this region map.',
      subregionEmptyLine1: 'No points have been created yet.',
      subregionEmptyLine2: 'Add exploration or research points from the tree on the right',
      subregionEmptyLine3: 'to show them on this sub-region map.',
      formatSubregionNodeCounts: (exploration, research) => `🎯 ${exploration}  🔬 ${research}`,
      formatRegionSummary: (subregions, exploration, research) => [`📁 Sub-regions ${subregions}`, `🎯 Exploration ${exploration}`, `🔬 Research ${research}`],
      formatSubregionSummary: (exploration, research) => [`🎯 Exploration ${exploration}`, `🔬 Research ${research}`],
    },
  },
  supportSections: {
    common: {
      currentRegion: 'Current Region',
      currentPoint: 'Current Point',
      currentStage: 'Current Stage',
      regionNotSelected: 'No region selected',
      pointNotSelected: 'No point selected',
      selectRegionFirst: 'Select a region first.',
      reset: 'Restore Original',
      draftPreview: 'Draft Preview',
    },
    records: {
      eyebrow: 'Record Summary',
      title: 'Records',
      statusChanged: 'Local record changes',
      statusConnected: 'Point record connected',
      statusNeedRegion: 'Select a region',
      regionDescriptionFallback: 'Select a region to organize learning data for its points.',
      pointDescriptionFallback: 'When an exploration or research point is selected, records focus on that point growth data.',
      focusMinutes: 'Focus time (min)',
      practiceCount: 'Practice count',
      confidence: 'Confidence',
      applicationNote: 'Application note',
      applicationPlaceholder: 'Write examples you tried, blockers, or next repetitions.',
      applicationDisabledPlaceholder: 'Select a region first.',
      saving: 'Saving record...',
      save: 'Save Record',
      noticeReady: 'Records summarize growth data for the current point. The current storage is region-based, but the product meaning is point-based.',
      noticeNeedRegion: 'Select a region to record growth data for its points.',
    },
    artifacts: {
      eyebrow: 'Artifact Entry',
      title: 'Artifacts',
      statusChanged: 'Local artifact changes',
      statusConnected: 'Point artifact connected',
      statusNeedRegion: 'Select a region',
      regionDescriptionFallback: 'Select a region to organize learning artifacts for its points.',
      pointDescriptionFallback: 'When an exploration or research point is selected, artifacts focus on that learning output.',
      type: 'Artifact type',
      typeOptions: [
        { value: 'note', label: 'Note' },
        { value: 'code', label: 'Code' },
        { value: 'project', label: 'Project' },
        { value: 'writing', label: 'Writing' },
        { value: 'design', label: 'Design' },
        { value: 'link', label: 'External Link' },
      ],
      artifactTitle: 'Artifact title',
      link: 'Artifact link',
      description: 'Artifact description',
      descriptionPlaceholder: 'Describe what this learning output shows.',
      descriptionDisabledPlaceholder: 'Select a region first.',
      saving: 'Saving artifact...',
      save: 'Save Artifact',
      noticeReady: 'Artifacts are output metadata for the current point. The current storage is region-based, but the product meaning is point-based.',
      noticeNeedRegion: 'Select a region to save artifact metadata for its points.',
    },
    community: {
      eyebrow: 'Community Preview',
      title: 'Community',
      status: 'Activates after learning starts',
      pointDescription: 'A knowledge hub for questions, answers, tips, and artifact feedback related to the selected point.',
      stageDescription: 'Question writing and feedback requests open in the learning diary. For now, this preview shows structure and entry guidance.',
      items: [
        { title: 'Planned features', body: 'Point-based Q&A search, creator help requests, and artifact feedback flow.' },
        { title: 'Why this is preview', body: 'Community opens with the learning diary, so writing tools stay closed during draft planning.' },
        { title: 'Connection rule', body: 'After learning diary runtime is connected, question writing and point search results will be added.' },
      ],
    },
    civilization: {
      eyebrow: 'Civilization Preview',
      title: 'Civilization View',
      status: 'Activates after sharing',
      stageDescription: 'Civilization View summarizes spread data after a completed course is shared.',
      metricLabel: 'Planned metrics',
      metricValue: 'Copies · Reactions · Spread stage',
      metricDescription: 'For now, this preview shows entry conditions and metric categories. Real numbers will connect after shared APIs are separated.',
      items: [
        { title: 'Activation condition', body: 'After completed + shared, this view shows civilization spread stages and public reaction metrics.' },
        { title: 'Why this is preview', body: 'Draft plans have no public data to aggregate, so only purpose and structure are shown.' },
        { title: 'Connection rule', body: 'When the sharing flow and shared-state aggregate API are ready, this switches to real data.' },
      ],
    },
    selectedPoint: {
      researchLabel: 'Research Point',
      researchTemplateDescription: (template) => `A learner note point using the ${template} template.`,
      researchTemplateLabel: (templateType) => {
        switch (templateType) {
          case 'concept_summary':
            return 'Concept Summary'
          case 'practice_strategy':
            return 'Practice Strategy'
          case 'problem_solving':
            return 'Problem Solving'
          case 'free_research':
            return 'Free Research'
          default:
            return null
        }
      },
      researchDescription: 'A research point where the learner writes their own notes.',
      explorationLabel: 'Exploration Point',
      selectedExploration: 'Selected exploration point',
      candidateExploration: 'Recommended exploration candidate',
      genericExploration: 'Exploration point',
      formatExplorationDescription: (label, resourceType) => `${label} · ${resourceType}`,
    },
  },
  tree: {
    loading: 'Loading...',
    emptyRegions: 'No exploration regions yet.',
    formatPlanTitle: (title) => {
      const trimmed = title?.trim()
      return trimmed ? `${trimmed} Exploration Plan` : 'Planet Exploration Plan'
    },
    formatPlanetName: (title) => title?.trim() || 'Untitled',
    planetPrefix: 'Planet',
    actionLabel: 'Edit',
    pendingDelete: 'To delete',
    pendingCreate: 'New',
    openLinkTitle: 'Open in new window',
    regionLabel: (index, name) => (/^Region\s*\d+$/i.test(name) ? `Region ${index}` : `Region ${index} · ${name}`),
    regionPlaceholder: (index) => `Region ${index}`,
    collapseRegion: 'Collapse region',
    expandRegion: 'Expand region',
    subRegionLabel: (index, name) => (/^(Sub|Sub-region)\s*\d+$/i.test(name) ? `Sub ${index}` : `Sub ${index} · ${name}`),
    subRegionPlaceholder: (index) => `Sub ${index}`,
    collapseSubRegion: 'Collapse sub-region',
    expandSubRegion: 'Expand sub-region',
    actionBar: {
      save: 'Save',
      cancel: 'Cancel',
      delete: 'Delete',
      undoDelete: 'Undo Delete',
      activatePlanet: 'Activate Planet',
      activate: 'Activate',
      hideInactive: 'Hide Inactive',
      showInactive: 'Show Inactive',
      inactiveActivateTitle: 'Reactivate this inactive planet and move to Explorer Diary.',
      activateTitle: 'Reactivate this course and move to Explorer Diary.',
      undoDeleteTitle: 'Undo the delete marker for the selected item.',
      deleteCourseTitle: 'Deactivate this course and return to Galaxy Dashboard.',
      deleteWithChildrenTitle: 'Child items will also be marked for deletion. Before learning starts, saving removes them fully.',
      inactiveToggleTitle: 'Toggle inactive item visibility.',
      reactivateTitle: 'Reactivate the selected item.',
    },
    editPanel: {
      action: 'Action',
      pendingDeleteNotice: 'This item is marked for deletion. Save to apply it, or undo deletion to restore it.',
      pendingCreateNotice: 'This item is pending. Save to add it to the exploration plan.',
      renamePlanet: 'Rename Planet',
      addRegion: '+ Region',
      renameRegion: 'Rename Region',
      renameSubRegion: 'Rename Sub-region',
      addChild: '+ Element',
      editNode: 'Edit Point',
      moveUp: 'Move up one step',
      moveDown: 'Move down one step',
      modals: {
        close: 'Close',
        cancel: 'Cancel',
        apply: 'Apply',
        edit: 'Update',
        add: 'Add',
        change: 'Change',
        check: 'Check',
        checking: 'Checking',
        valid: '✓ Valid',
        error: '✗ Error',
        openLink: 'Open in new window',
        internalContent: 'Internal content',
        urlMustPass: 'The URL must pass the check before it can be added.',
        parentPath: (value) => `Parent path: ${value}`,
        recommendation: {
          aria: 'Content recommendation',
          title: 'Content Recommendation',
          queryLabel: 'Recommendation query',
          queryPlaceholder: 'Content title, description, or author',
          search: 'Recommend',
          searching: 'Searching',
          recommendAgain: 'Recommend Again',
          empty: 'Press Recommend to load candidate content.',
          pointErrorEmpty: 'Add points, then try again.',
          loading: 'Loading candidates.',
          noResults: 'No matching content. Try another query.',
          loaded: (count) => `${count} candidates loaded.`,
          missingCandidate: 'This candidate has no title or content data. Choose another one.',
          loginRequired: 'Check your login status, then try again.',
          insufficientPoints: (balance, cost) => `Not enough points. Current ${balance}pt (required: ${cost}pt)`,
          loadFailed: 'Could not load content candidates.',
        },
        addExploration: {
          aria: 'Add exploration point',
          title: 'Add Exploration Point',
          recommendationTab: 'Recommendation',
          urlTab: 'Add URL',
          titleLabel: 'Point title',
          titlePlaceholder: 'Enter a point title',
          urlLabel: 'Link URL',
          addButton: 'Add Point',
        },
        addChild: {
          aria: 'Add exploration element',
          title: 'Add Exploration Element',
          tabAria: 'Choose element type',
          subregionTab: 'Sub-lesson',
          explorationTab: 'Exploration Point',
          researchTab: 'Research Point',
          subregionDescription: 'A sub-lesson splits the selected lesson into a smaller learning flow. Use it to group related exploration and research points.',
          subregionLimit: 'Sub-lessons can only be added to lessons, up to 3 per lesson.',
          subregionNameLabel: 'Sub-lesson name',
          subregionNamePlaceholder: 'Enter a sub-lesson name',
          addSubregion: 'Add Sub-lesson',
          explorationDescription: 'An exploration point is external content the learner will open and study. Pick a recommendation or add a checked URL.',
          explorationMethodAria: 'Choose exploration point input method',
          recommendationDescription: 'Recommendations use the current region, goal context, and query to find learning candidates. Review a candidate and apply it to add a point.',
          urlDescription: 'Direct URL add lets you add a resource you found yourself. The link must pass the check before it is added.',
          urlDirectTab: 'Direct URL',
          researchDescription: 'A research point is a place for the learner to investigate or organize their own learning material. Material can be added later during learning.',
          researchTitleLabel: 'Research point title',
          researchTitlePlaceholder: 'Enter a research point title',
          addResearch: 'Add Research Point',
        },
        addUrl: {
          explorationTitle: 'Add Exploration Point (URL)',
          researchTitle: 'Add Research Point',
        },
        nodeEdit: {
          explorationAria: 'Edit exploration point',
          researchAria: 'Edit research point',
          explorationTitle: 'Edit Exploration Point',
          researchTitle: 'Edit Research Point',
          titleLabel: 'Title',
          explorationPlaceholder: 'Exploration point title',
          researchPlaceholder: 'Research point title',
          linkLabel: 'Link',
          recommend: 'Recommend',
        },
        deleteConfirm: {
          deleteAria: 'Confirm deletion',
          activateAria: 'Confirm activation',
          alertAria: 'Cannot delete',
          warningEyebrow: 'Lumi Warning',
          alertTitle: 'Cannot delete',
          alertSuffix: 'Keep at least one lesson and continue refining the structure.',
          deleteTitle: 'Confirm Delete',
          activateTitle: 'Confirm Activation',
          targetLabel: 'Target',
          courseConfirmLabel: 'Confirm course name',
          exactTitleHelp: (action) => `Enter the exact planet name to ${action === 'activate' ? 'activate' : 'delete'} it.`,
          delete: 'Delete',
          activate: 'Activate',
        },
        courseRename: {
          aria: 'Rename planet',
          title: 'Rename Planet',
          label: 'Planet name',
          placeholder: 'Enter a new planet name',
          help: 'The planet name updates the plan title and the map header.',
        },
        regionRename: {
          regionTitle: 'Rename Region',
          subregionTitle: 'Rename Sub-region',
          regionLabel: 'Region name',
          subregionLabel: 'Sub-region name',
          aria: 'Rename region',
          placeholder: (label) => `Enter ${label.toLowerCase()}`,
          help: 'The changed name appears on this screen until you save the exploration plan.',
        },
        regionAdd: {
          aria: 'Add region',
          title: 'Add Region',
          label: 'Region name',
          description: 'A region is a major lesson stage in the course. Use regions to divide the main topics you will explore toward the goal.',
          placeholder: 'New region name',
          help: 'The new region appears on this screen until you save the exploration plan.',
          addButton: 'Add Region',
        },
        messages: {
          urlValid: 'This link is valid.',
          urlInvalid: 'This link is not valid.',
          urlCheckFailed: 'Link check failed.',
          contentIdBased: 'This exploration point is based on an internal content ID.',
          deleteCourse: 'This deactivates the course and returns to Galaxy Dashboard. Enter the course name exactly to continue.',
          activateCourse: 'This reactivates the inactive course and moves to Explorer Diary. Enter the planet name exactly to continue.',
          lastRegionTitle: 'Region',
          lastRegionBlock: 'The last lesson cannot be marked for deletion. Keep at least one lesson to avoid an empty course.',
          deleteRegionWithChildren: 'This region and its child items will be marked for deletion. Before learning starts, saving removes them fully.',
          deleteRegion: 'This region will be marked for deletion.',
          deleteSubregionWithChildren: 'This sub-region and child points will be marked for deletion. Before learning starts, saving removes them fully.',
          deleteSubregion: 'This sub-region will be marked for deletion. Before learning starts, saving removes it fully.',
          deleteNode: 'This point will be marked for deletion. Before learning starts, saving removes it fully.',
          fallbackRegion: 'Region',
          fallbackSubregion: 'Sub-region',
        },
      },
    },
  },
}

export function getDashboardCourseDraftCopy(locale: Locale | string | null | undefined): DashboardCourseDraftCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
