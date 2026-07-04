import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DashboardSettingsAICopy = {
  loading: string
  shell: {
    eyebrow: string
    title: string
    copy: string
    tabs: { profile: string; points: string; ai: string }
  }
  byok: {
    title: string
    subtitle: string
    statusUsing: string
    statusDisabled: string
    statusManaged: string
    summarySaved: (provider: string, enabled: boolean) => string
    summaryManaged: string
    recentlyValidated: (date: string) => string
    usageStateTitle: string
    usageEnabled: string
    usageDisabled: string
    usageEmpty: string
    registrationTitle: string
    savedKeyTitle: string
    savedProvider: (provider: string) => string
    savedValid: string
    validationStatus: (status: string) => string
    providerLabel: string
    endpointLabel: string
    endpointPlaceholder: string
    apiKeyLabel: string
    validatingStatus: string
    validStatus: string
    needsValidationStatus: string
    validate: string
    validating: string
    save: string
    saving: string
    inputHelp: string
    deleteSummary: string
    deleteTitle: string
    deleteText: string
    deleteConfirmLabel: string
    deleteConfirmValue: string
    deleteConfirmAria: string
    deleteAction: string
    deleting: string
  }
  usage: {
    noticeTitle: string
    noticeBody1: string
    noticeBody2: string
    historyTitle: string
    historySubtitle: string
    startDate: string
    endDate: string
    search: string
    searching: string
    invalidRange: string
    rangeResult: (start: string, end: string, total: number) => string
    totalTokens: string
    dailyTokens: string
    totalCost: string
    dailyCost: string
    input: string
    output: string
    success: string
    fail: string
    loading: string
    empty: string
    count: (start: number, end: number, total: number) => string
    zero: string
    previous: string
    next: string
    unknownTokens: string
    unknownCost: string
    estimatedCost: (value: string) => string
    lessThanOneWon: string
    won: (value: string) => string
    features: Record<string, string>
  }
}

const ko: DashboardSettingsAICopy = {
  loading: 'AI 설정을 불러오는 중...',
  shell: {
    eyebrow: '학습자 설정',
    title: 'AI 설정',
    copy: '생성형 AI에 사용할 개인 AI 제공자를 선택하고, BYOK 검증과 저장 상태를 관리합니다. 콘텐츠 검색과 추천은 LearnCosmos의 설치형 임베딩 기반으로 별도 처리됩니다.',
    tabs: { profile: '학습자정보', points: '포인트', ai: 'AI 설정' },
  },
  byok: {
    title: 'BYOK 설정',
    subtitle: '목표 채팅, 코스 생성, 학습 코칭에 사용할 개인 AI 제공자와 API 키 상태를 관리합니다.',
    statusUsing: '내 키 사용 중',
    statusDisabled: '내 키 비활성',
    statusManaged: '관리형 기본 사용',
    summarySaved: (provider, enabled) => `${provider} 개인 키 ${enabled ? '활성화됨' : '비활성화됨'}`,
    summaryManaged: '현재는 LearnCosmos 관리형 기본 설정 사용 중',
    recentlyValidated: (date) => `최근 검증 ${date}`,
    usageStateTitle: 'BYOK 사용 상태',
    usageEnabled: '생성형 AI 호출에 개인 키를 우선 사용합니다. 추천 검색용 설치형 임베딩 경로는 BYOK LLM 키와 분리됩니다.',
    usageDisabled: '저장된 키는 보관되지만 AI 호출에는 사용하지 않습니다. AI 기능은 LearnCosmos 시스템 기본 키와 포인트 정책을 따릅니다.',
    usageEmpty: '개인 키를 저장한 뒤 생성형 AI 호출에 BYOK를 사용할지 켜고 끌 수 있습니다.',
    registrationTitle: '개인 키 등록',
    savedKeyTitle: '등록된 키 정보',
    savedProvider: (provider) => `AI 제공자: ${provider}`,
    savedValid: '저장된 BYOK는 유효한 API 키 값입니다.',
    validationStatus: (status) => `검증 상태: ${status || '미확인'}`,
    providerLabel: '1. AI 제공자',
    endpointLabel: 'Endpoint URL',
    endpointPlaceholder: 'Llama 호환 endpoint URL',
    apiKeyLabel: '2. API Key 등록',
    validatingStatus: '현재 입력한 키를 검증하고 있습니다.',
    validStatus: '입력한 API 키는 유효합니다. BYOK 저장을 진행할 수 있습니다.',
    needsValidationStatus: 'API Key를 입력한 뒤 검증을 진행해 주세요.',
    validate: '3. API 키 검증',
    validating: '검증 중...',
    save: '4. BYOK 저장',
    saving: '저장 중...',
    inputHelp: '입력한 API Key는 검증과 저장 요청에만 사용되며, 저장 후에는 화면에 다시 표시하지 않습니다.',
    deleteSummary: '저장된 BYOK 키 삭제',
    deleteTitle: '삭제 전 확인',
    deleteText: '삭제하면 저장된 개인 API 키와 검증 상태가 제거되고, AI 기능은 LearnCosmos 시스템 기본 키와 포인트 정책을 따릅니다.',
    deleteConfirmLabel: '확인 문구',
    deleteConfirmValue: 'BYOK 삭제',
    deleteConfirmAria: 'BYOK 삭제 확인 문구',
    deleteAction: '키 삭제',
    deleting: '삭제 중...',
  },
  usage: {
    noticeTitle: 'BYOK 사용량 기록 안내',
    noticeBody1: 'BYOK를 사용하는 생성형 AI 호출에서는 사용량 확인과 장애 대응을 위해 입력 토큰 수, 출력 토큰 수, 사용 모델, 호출 시각, 기능 구분, 원화 환산 예상비용을 기록합니다.',
    noticeBody2: 'API 키 원문은 사용량 기록, 일반 조회 화면, API 응답, 운영 로그에 포함하지 않습니다. 키는 암호화 저장되며 생성형 AI 호출에 필요한 순간에만 사용됩니다. 콘텐츠 추천과 유사 자료 검색은 설치형 임베딩 기반으로 별도 처리됩니다.',
    historyTitle: '최근 BYOK 사용량',
    historySubtitle: '기본은 오늘 하루 기준입니다. 시작 날짜는 제한 없이 선택할 수 있고, 한 번에 최대 3개월 범위까지 10개씩 페이지로 확인합니다. 예상비용은 1달러=1,400원 정책 환율로 표시합니다.',
    startDate: '시작 날짜',
    endDate: '종료 날짜',
    search: '사용량 검색',
    searching: '검색 중...',
    invalidRange: '종료 날짜는 시작 날짜 이후, 최대 3개월 범위 안에서 선택해 주세요.',
    rangeResult: (start, end, total) => `${start} ~ ${end} 기준 ${total.toLocaleString('ko-KR')}건`,
    totalTokens: '검색기간 총 토큰 사용량',
    dailyTokens: '검색기간 하루 평균 사용량',
    totalCost: '검색기간 총 예상비용',
    dailyCost: '검색기간 하루 평균 예상비용',
    input: '입력',
    output: '출력',
    success: '성공',
    fail: '실패',
    loading: 'BYOK 사용량을 불러오는 중입니다.',
    empty: '선택한 기간에 기록된 BYOK 호출 내역이 없습니다.',
    count: (start, end, total) => `${start.toLocaleString('ko-KR')}~${end.toLocaleString('ko-KR')} / ${total.toLocaleString('ko-KR')}건`,
    zero: '0건',
    previous: '이전',
    next: '다음',
    unknownTokens: '미집계',
    unknownCost: '예상 비용 미집계',
    estimatedCost: (value) => `예상 ${value}`,
    lessThanOneWon: '1원 미만',
    won: (value) => `${value}원`,
    features: {
      course_draft_create: '코스 초안 생성',
      lesson_recommendation_search: '리슨 추천 검색',
      point_question_generate: '질문 생성',
      point_question_feedback: '답변 피드백',
      point_self_evaluation_draft: '자기평가 초안',
      point_ai_summary: 'AI 요약',
    },
  },
}

const en: DashboardSettingsAICopy = {
  loading: 'Loading AI settings...',
  shell: {
    eyebrow: 'Learner Settings',
    title: 'AI Settings',
    copy: 'Choose the personal AI provider used for generative AI and manage BYOK validation, saved key status, and recent usage. Content search and recommendations use LearnCosmos installed embeddings separately.',
    tabs: { profile: 'Profile', points: 'Points', ai: 'AI Settings' },
  },
  byok: {
    title: 'BYOK Settings',
    subtitle: 'Manage the personal AI provider and API key used for goal chat, course creation, and learning coaching.',
    statusUsing: 'Using My Key',
    statusDisabled: 'My Key Off',
    statusManaged: 'Managed Default',
    summarySaved: (provider, enabled) => `${provider} personal key ${enabled ? 'enabled' : 'disabled'}`,
    summaryManaged: 'Currently using the LearnCosmos managed default.',
    recentlyValidated: (date) => `Last validated ${date}`,
    usageStateTitle: 'BYOK Status',
    usageEnabled: 'Your personal key is used first for generative AI calls. The installed embedding path for recommendation search is separate from your BYOK LLM key.',
    usageDisabled: 'Your saved key is kept, but AI calls use the LearnCosmos system key and point policy.',
    usageEmpty: 'Save a personal key, then turn BYOK usage for generative AI calls on or off here.',
    registrationTitle: 'Register Personal Key',
    savedKeyTitle: 'Saved Key',
    savedProvider: (provider) => `AI Provider: ${provider}`,
    savedValid: 'The saved BYOK key is valid.',
    validationStatus: (status) => `Validation: ${status || 'Not checked'}`,
    providerLabel: '1. AI Provider',
    endpointLabel: 'Endpoint URL',
    endpointPlaceholder: 'Llama-compatible endpoint URL',
    apiKeyLabel: '2. API Key',
    validatingStatus: 'Validating the key you entered.',
    validStatus: 'This API key is valid. You can save BYOK now.',
    needsValidationStatus: 'Enter an API key, then validate it.',
    validate: '3. Validate Key',
    validating: 'Validating...',
    save: '4. Save BYOK',
    saving: 'Saving...',
    inputHelp: 'The API key you enter is used only for validation and saving. It is not shown again after storage.',
    deleteSummary: 'Delete Saved BYOK Key',
    deleteTitle: 'Confirm Deletion',
    deleteText: 'Deleting removes the saved personal API key and validation status. AI features will use the LearnCosmos system key and point policy.',
    deleteConfirmLabel: 'Confirmation Text',
    deleteConfirmValue: 'DELETE BYOK',
    deleteConfirmAria: 'BYOK delete confirmation text',
    deleteAction: 'Delete Key',
    deleting: 'Deleting...',
  },
  usage: {
    noticeTitle: 'BYOK Usage Logging',
    noticeBody1: 'For generative AI calls that use BYOK, LearnCosmos records input tokens, output tokens, model, call time, feature, and estimated KRW cost for usage review and troubleshooting.',
    noticeBody2: 'Raw API keys are not included in usage logs, general views, API responses, or operational logs. Keys are encrypted and used only when needed for generative AI calls. Content recommendations and similar-resource search use a separate installed embedding path.',
    historyTitle: 'Recent BYOK Usage',
    historySubtitle: 'The default range is today. You can choose any start date and view up to three months at a time, 10 records per page. Estimated cost uses the policy rate of 1 USD = 1,400 KRW.',
    startDate: 'Start Date',
    endDate: 'End Date',
    search: 'Search Usage',
    searching: 'Searching...',
    invalidRange: 'Choose an end date after the start date and within a three-month range.',
    rangeResult: (start, end, total) => `${total.toLocaleString('en-US')} records for ${start} to ${end}`,
    totalTokens: 'Total Tokens',
    dailyTokens: 'Daily Avg. Tokens',
    totalCost: 'Total Estimated Cost',
    dailyCost: 'Daily Avg. Cost',
    input: 'Input',
    output: 'Output',
    success: 'Success',
    fail: 'Failed',
    loading: 'Loading BYOK usage...',
    empty: 'No BYOK calls were recorded in the selected range.',
    count: (start, end, total) => `${start.toLocaleString('en-US')}-${end.toLocaleString('en-US')} / ${total.toLocaleString('en-US')} records`,
    zero: '0 records',
    previous: 'Previous',
    next: 'Next',
    unknownTokens: 'Not tracked',
    unknownCost: 'Cost not tracked',
    estimatedCost: (value) => `Est. ${value}`,
    lessThanOneWon: 'under KRW 1',
    won: (value) => `KRW ${value}`,
    features: {
      course_draft_create: 'Course Draft',
      lesson_recommendation_search: 'Lesson Search',
      point_question_generate: 'Question Generation',
      point_question_feedback: 'Answer Feedback',
      point_self_evaluation_draft: 'Review Draft',
      point_ai_summary: 'AI Summary',
    },
  },
}

export function getDashboardSettingsAICopy(locale: Locale | string | null | undefined): DashboardSettingsAICopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
