import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'
import type { GoalFlowStage } from '@/components/goal-interview/GoalFlowLoadingOverlay'

export type DashboardGoalCopy = {
  loading: {
    stages: Record<GoalFlowStage, { title: string; description: string }>
    creatingPlanetWait: {
      normal: string
      long: string
      busyHint: string
      busy: string
    }
    waitEstimate: {
      duration: (minSeconds: number, maxSeconds: number) => string
      queued: (position: number, label: string) => string
      running: (label: string) => string
      activeOnly: (activeJobs: number, label: string) => string
    }
    atlasCaption: {
      readyNoSelection: string
      selected: (name: string) => string
      redirecting: string
      creating: string
    }
    picker: {
      selected: string
      ready: string
      idle: string
      empty: string
      ariaLabel: (name: string, selected: boolean) => string
      chooseTexture: string
      startWithTexture: string
      preparing: string
    }
    complete: string
  }
  page: {
    headerEyebrow: string
    headerTitle: string
    loadErrorTitle: string
    backToDashboard: string
    confirmedAction: string
    confirmedActionBusy: string
    exitAria: string
    exitEyebrow: string
    exitTitle: string
    exitMessage: string
    continue: string
    dashboardShort: string
    notices: {
      busyTitle: string
      busyBody: string
      retryTitle: string
      retryBody: string
    }
    errors: {
      resetFailed: string
      missingConfirmedGoal: string
      createFailed: string
    }
  }
  interview: {
    subtitle: string
    retryMessage: string
    inputPlaceholder: string
    send: string
    rebuildTitle: string
    rebuildDescription: string
    rebuildOptions: {
      rebuildAll: { label: string; desc: string }
      keepStructure: { label: string; desc: string }
    }
    cancel: string
    confirmedTitle: string
    confirmedActionBusy: string
    typing: string
  }
  proposal: {
    badge: string
    confirming: string
    confirm: string
    retry: string
  }
  points: {
    byokTitle: string
    byokDetail: string
    checkingTitle: string
    checkingDetail: (cost: number) => string
    insufficientDetail: (cost: number) => string
    readyDetail: (cost: number, remaining: number) => string
    points: (value: number) => string
  }
  hookErrors: {
    loadFailed: string
    loadFallback: string
    startFailed: string
    startFallback: string
    sendFailed: string
    sendFallback: string
    confirmFailed: string
    confirmFallback: string
    reviseFailed: string
    reviseFallback: string
    rebuildFailed: string
    rebuildFallback: string
    cancelFailed: string
    cancelFallback: string
    attachFailed: string
    resetFailed: string
  }
}

const ko: DashboardGoalCopy = {
  loading: {
    stages: {
      bootstrapping_lumi: {
        title: '루미를 불러오는 중...',
        description: '대화를 시작할 준비를 하고 있어요.',
      },
      starting_interview: {
        title: '루미가 목표를 정리하는 중...',
        description: '당신의 배움 의도를 이해하고 첫 질문을 준비하고 있어요.',
      },
      creating_planet: {
        title: '탐험계획을 만드는 중...',
        description: '코스 초안을 만들고 있어요.',
      },
      redirecting: {
        title: '행성으로 이동하는 중...',
        description: 'Explorer Diary를 곧 열어드릴게요.',
      },
    },
    creatingPlanetWait: {
      normal: '코스 초안을 만들고 있어요.',
      long: '생각보다 조금 더 걸리고 있어요. 학습 목표에 맞춰 리슨을 정리하는 중입니다.',
      busyHint: '현재 생성 요청이 많아 조금 지연되고 있어요.',
      busy: '현재 코스 생성 요청이 많습니다. 잠시 후 다시 시도해 주세요.',
    },
    waitEstimate: {
      duration: (minSeconds, maxSeconds) => minSeconds === maxSeconds ? `${maxSeconds}초` : `${minSeconds}~${maxSeconds}초`,
      queued: (position, label) => `대기 순위 ${position}번째 · 예상 ${label}`,
      running: (label) => `생성 작업이 시작됐어요 · 예상 ${label}`,
      activeOnly: (activeJobs, label) => `현재 처리 중 ${activeJobs}건 · 예상 ${label}`,
    },
    atlasCaption: {
      readyNoSelection: '탐험계획은 준비됐어요. 행성 텍스처를 골라주세요',
      selected: (name) => `${name} 텍스처를 적용할게요`,
      redirecting: '행성 표면을 정돈하고 있어요',
      creating: '새 행성의 표면이 만들어지고 있어요',
    },
    picker: {
      selected: '선택한 텍스처를 확인해 주세요',
      ready: '행성 텍스처를 선택한 뒤 확인을 눌러 이동합니다',
      idle: '행성 텍스처 선택',
      empty: '새 텍스처맵 목록을 불러오는 중...',
      ariaLabel: (name, selected) => `행성 텍스처맵: ${name}${selected ? ' (선택됨)' : ''}`,
      chooseTexture: '텍스처를 선택해 주세요',
      startWithTexture: '이 텍스처로 시작하기',
      preparing: '탐험계획을 준비하는 중...',
    },
    complete: '완료!',
  },
  page: {
    headerEyebrow: 'Goal Flow',
    headerTitle: '목표 찾기',
    loadErrorTitle: '루미를 불러오지 못했어요',
    backToDashboard: '대시보드로 돌아가기',
    confirmedAction: '탐험 시작하기',
    confirmedActionBusy: '탐험계획 준비 중...',
    exitAria: '목표 찾기 나가기 확인',
    exitEyebrow: 'Lumi Confirm',
    exitTitle: '목표 찾기를 나갈까요?',
    exitMessage: '대시보드에서 목표 찾기를 다시 시작할 수 있어요. 지금은 목표를 확정하지 않고 나갑니다.',
    continue: '계속하기',
    dashboardShort: '대시보드로',
    notices: {
      busyTitle: '지금은 탐험계획 생성 요청이 많아요',
      busyBody: '현재 개발서버에서 동시에 처리 가능한 생성 수를 넘으면 일부 요청은 잠시 멈추지 않고 빠르게 안내하고 있어요. 잠시 후 다시 누르면 생성될 가능성이 높습니다.',
      retryTitle: '이번 생성 응답이 불안정했어요',
      retryBody: '리슨 초안을 만드는 첫 응답이 정상 형식으로 돌아오지 않아 이번 생성은 저장하지 않았어요. 같은 목표로 다시 누르면 새로 생성 요청을 보냅니다.',
    },
    errors: {
      resetFailed: '목표 인터뷰를 초기화하지 못했습니다.',
      missingConfirmedGoal: '확정된 목표를 확인하지 못했습니다.',
      createFailed: '행성탐험계획 생성에 실패했습니다.',
    },
  },
  interview: {
    subtitle: '목표 탐험 안내자',
    retryMessage: '다시 정해주세요.',
    inputPlaceholder: '루미에게 답장하기... (Enter로 전송)',
    send: '전송',
    rebuildTitle: '학습 구조를 어떻게 할까요?',
    rebuildDescription: '목표가 바뀌었어요. 기존 탐험계획을 어떻게 처리할지 선택해주세요.',
    rebuildOptions: {
      rebuildAll: { label: '모두 재구성', desc: '현재 탐험계획 전체를 새 목표 기준으로 다시 구성합니다' },
      keepStructure: { label: '학습자 편집', desc: '현재 구조를 유지한 채 학습자가 직접 리슨과 자료를 조정합니다' },
    },
    cancel: '취소',
    confirmedTitle: '목표 확정!',
    confirmedActionBusy: '탐험계획 준비 중...',
    typing: '···',
  },
  proposal: {
    badge: '🎯 목표 제안',
    confirming: '확정 중...',
    confirm: '이 목표로 시작하기',
    retry: '다시 정하기',
  },
  points: {
    byokTitle: 'BYOK LLM 사용 중',
    byokDetail: 'BYOK LLM · 목표 대화/행성 생성에 개인 키 사용',
    checkingTitle: '포인트 확인 중',
    checkingDetail: (cost) => `목표 대화 0pt · 행성 생성 ${cost}pt`,
    insufficientDetail: (cost) => `행성 생성에 ${cost}pt가 필요해요`,
    readyDetail: (cost, remaining) => `목표 대화 0pt · 행성 생성 ${cost}pt · 생성 후 ${remaining}pt`,
    points: (value) => `보유 ${value}pt`,
  },
  hookErrors: {
    loadFailed: '목표 인터뷰를 불러오지 못했습니다.',
    loadFallback: '로드 실패',
    startFailed: '인터뷰 시작에 실패했습니다.',
    startFallback: '인터뷰 시작 실패',
    sendFailed: '메시지 전송에 실패했습니다.',
    sendFallback: '메시지 전송 실패',
    confirmFailed: '목표 확정에 실패했습니다.',
    confirmFallback: '목표 확정 실패',
    reviseFailed: '목표 수정에 실패했습니다.',
    reviseFallback: '목표 수정 실패',
    rebuildFailed: '구조 재구성 선택에 실패했습니다.',
    rebuildFallback: '구조 재구성 선택 실패',
    cancelFailed: '목표 변경 취소에 실패했습니다.',
    cancelFallback: '목표 변경 취소 실패',
    attachFailed: '목표를 행성탐험계획에 연결하지 못했습니다.',
    resetFailed: '이전 목표 인터뷰를 초기화하지 못했습니다.',
  },
}

const en: DashboardGoalCopy = {
  loading: {
    stages: {
      bootstrapping_lumi: {
        title: 'Calling Lumi...',
        description: 'Preparing to start the conversation.',
      },
      starting_interview: {
        title: 'Lumi is shaping your goal...',
        description: 'Understanding your learning intent and preparing the first question.',
      },
      creating_planet: {
        title: 'Creating your learning plan...',
        description: 'We are creating your course draft.',
      },
      redirecting: {
        title: 'Moving to your planet...',
        description: 'Explorer Diary will open shortly.',
      },
    },
    creatingPlanetWait: {
      normal: 'We are creating your course draft.',
      long: 'This is taking a little longer than usual. We are organizing lessons around your learning goal.',
      busyHint: 'Course generation is a little delayed because there are many requests right now.',
      busy: 'There are many course generation requests right now. Please try again shortly.',
    },
    waitEstimate: {
      duration: (minSeconds, maxSeconds) => minSeconds === maxSeconds ? `${maxSeconds}s` : `${minSeconds}-${maxSeconds}s`,
      queued: (position, label) => `Queue position ${position} · about ${label}`,
      running: (label) => `Generation has started · about ${label}`,
      activeOnly: (activeJobs, label) => `${activeJobs} active jobs · about ${label}`,
    },
    atlasCaption: {
      readyNoSelection: 'Your plan is ready. Choose a planet texture.',
      selected: (name) => `Applying the ${name} texture.`,
      redirecting: 'Preparing the planet surface.',
      creating: 'Building the new planet surface.',
    },
    picker: {
      selected: 'Review the selected texture',
      ready: 'Choose a planet texture, then confirm to continue.',
      idle: 'Choose Planet Texture',
      empty: 'Loading texture maps...',
      ariaLabel: (name, selected) => `Planet texture map: ${name}${selected ? ' (selected)' : ''}`,
      chooseTexture: 'Choose a texture',
      startWithTexture: 'Start with this texture',
      preparing: 'Preparing the learning plan...',
    },
    complete: 'Done',
  },
  page: {
    headerEyebrow: 'Goal Flow',
    headerTitle: 'Find Your Goal',
    loadErrorTitle: 'Could not load Lumi',
    backToDashboard: 'Back to Dashboard',
    confirmedAction: 'Start Exploration',
    confirmedActionBusy: 'Preparing plan...',
    exitAria: 'Confirm leaving goal flow',
    exitEyebrow: 'Lumi Confirm',
    exitTitle: 'Leave goal flow?',
    exitMessage: 'You can start goal flow again from the dashboard. For now, you will leave without confirming a goal.',
    continue: 'Continue',
    dashboardShort: 'Dashboard',
    notices: {
      busyTitle: 'Learning plan requests are busy right now',
      busyBody: 'The development server is at its current generation limit, so some requests return quickly instead of waiting. Try again shortly.',
      retryTitle: 'This generation response was unstable',
      retryBody: 'The first lesson draft response did not return in a valid format, so this generation was not saved. Try again with the same goal to start a new request.',
    },
    errors: {
      resetFailed: 'Could not reset the goal interview.',
      missingConfirmedGoal: 'Could not find a confirmed goal.',
      createFailed: 'Failed to create the planet learning plan.',
    },
  },
  interview: {
    subtitle: 'Goal Exploration Guide',
    retryMessage: 'Please suggest another goal.',
    inputPlaceholder: 'Reply to Lumi... (Enter to send)',
    send: 'Send',
    rebuildTitle: 'How should we handle the learning structure?',
    rebuildDescription: 'Your goal changed. Choose what to do with the current exploration plan.',
    rebuildOptions: {
      rebuildAll: { label: 'Rebuild All', desc: 'Recreate the full exploration plan around the new goal.' },
      keepStructure: { label: 'Edit Myself', desc: 'Keep the current structure and adjust lessons or resources yourself.' },
    },
    cancel: 'Cancel',
    confirmedTitle: 'Goal Confirmed',
    confirmedActionBusy: 'Preparing plan...',
    typing: '···',
  },
  proposal: {
    badge: 'Goal Suggestion',
    confirming: 'Confirming...',
    confirm: 'Start with this goal',
    retry: 'Try again',
  },
  points: {
    byokTitle: 'Using BYOK LLM',
    byokDetail: 'BYOK LLM · personal key for goal chat / planet creation',
    checkingTitle: 'Checking points',
    checkingDetail: (cost) => `Goal chat 0 pt · planet creation ${cost} pt`,
    insufficientDetail: (cost) => `${cost} pt is required to create a planet.`,
    readyDetail: (cost, remaining) => `Goal chat 0 pt · planet creation ${cost} pt · ${remaining} pt after creation`,
    points: (value) => `${value} pt available`,
  },
  hookErrors: {
    loadFailed: 'Could not load the goal interview.',
    loadFallback: 'Load failed',
    startFailed: 'Could not start the interview.',
    startFallback: 'Interview start failed',
    sendFailed: 'Could not send the message.',
    sendFallback: 'Message send failed',
    confirmFailed: 'Could not confirm the goal.',
    confirmFallback: 'Goal confirmation failed',
    reviseFailed: 'Could not revise the goal.',
    reviseFallback: 'Goal revision failed',
    rebuildFailed: 'Could not save the rebuild choice.',
    rebuildFallback: 'Rebuild choice failed',
    cancelFailed: 'Could not cancel the goal change.',
    cancelFallback: 'Goal change cancellation failed',
    attachFailed: 'Could not attach the goal to the planet learning plan.',
    resetFailed: 'Could not reset the previous goal interview.',
  },
}

export function getDashboardGoalCopy(locale: Locale | string | null | undefined): DashboardGoalCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
