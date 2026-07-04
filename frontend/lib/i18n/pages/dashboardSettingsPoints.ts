import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DashboardSettingsPointsCopy = {
  loading: string
  shell: {
    eyebrow: string
    title: string
    copy: string
    tabs: { profile: string; points: string; ai: string }
  }
  points: {
    title: string
    subtitle: string
    standard: string
    account: (email?: string) => string
    unknown: string
    breakdown: string
    free: string
    paid: string
  }
  usage: {
    title: string
    subtitle: string
    startDate: string
    endDate: string
    search: string
    searching: string
    invalidRange: string
    rangeResult: (startDate: string, endDate: string, total: number) => string
    granted: string
    purchased: string
    used: string
    refunded: string
    netChange: string
    loading: string
    empty: string
    previous: string
    next: string
    count: (start: number, end: number, total: number) => string
    zero: string
    transactionTypes: Record<string, string>
  }
}

const ko: DashboardSettingsPointsCopy = {
  loading: '포인트 정보를 불러오는 중...',
  shell: {
    eyebrow: '학습자 설정',
    title: '포인트 현황',
    copy: '현재 보유 포인트와 무료/유료 포인트 구성을 확인합니다. 학습 흐름에서 사용 가능한 여유분을 여기서 빠르게 점검할 수 있습니다.',
    tabs: { profile: '학습자정보', points: '포인트', ai: 'AI 설정' },
  },
  points: {
    title: '포인트 요약',
    subtitle: '탐험계획 생성과 AI 기능에 사용할 수 있는 현재 포인트 현황입니다.',
    standard: 'Standard',
    account: (email) => `현재 계정 ${email || '미확인'}`,
    unknown: '미확인',
    breakdown: '포인트 구성',
    free: '무료 포인트',
    paid: '유료 포인트',
  },
  usage: {
    title: '포인트 사용 내역',
    subtitle: 'BYOK 사용량처럼 기간별 포인트 지급, 사용, 환불 흐름을 확인합니다.',
    startDate: '시작일',
    endDate: '종료일',
    search: '조회',
    searching: '조회 중...',
    invalidRange: '조회 기간을 확인해 주세요. 최대 3개월까지 조회할 수 있습니다.',
    rangeResult: (startDate, endDate, total) => `${startDate} ~ ${endDate} · ${total.toLocaleString('ko-KR')}건`,
    granted: '무료 지급',
    purchased: '구매/충전',
    used: '사용',
    refunded: '환불/복구',
    netChange: '순변동',
    loading: '포인트 내역을 불러오는 중입니다.',
    empty: '선택한 기간에 포인트 내역이 없습니다.',
    previous: '이전',
    next: '다음',
    count: (start, end, total) => `${total.toLocaleString('ko-KR')}건 중 ${start.toLocaleString('ko-KR')}~${end.toLocaleString('ko-KR')}건`,
    zero: '표시할 포인트 내역이 없습니다.',
    transactionTypes: {
      grant_free: '무료 포인트 지급',
      purchase: '포인트 구매',
      use_course_gen: '탐험계획 생성',
      use_lesson_rec: '추천 검색',
      refund: '포인트 환불/복구',
    },
  },
}

const en: DashboardSettingsPointsCopy = {
  loading: 'Loading point balance...',
  shell: {
    eyebrow: 'Learner Settings',
    title: 'Points',
    copy: 'Review your current point balance and the split between free and paid points available for learning flows.',
    tabs: { profile: 'Profile', points: 'Points', ai: 'AI Settings' },
  },
  points: {
    title: 'Point Summary',
    subtitle: 'Current points available for exploration planning and AI features.',
    standard: 'Standard',
    account: (email) => `Current account ${email || 'Unknown'}`,
    unknown: 'Unknown',
    breakdown: 'Point Breakdown',
    free: 'Free Points',
    paid: 'Paid Points',
  },
  usage: {
    title: 'Point History',
    subtitle: 'Review point grants, spending, and refunds by date range, similar to BYOK usage.',
    startDate: 'Start Date',
    endDate: 'End Date',
    search: 'Search',
    searching: 'Searching...',
    invalidRange: 'Check the date range. You can search up to 3 months.',
    rangeResult: (startDate, endDate, total) => `${startDate} - ${endDate} · ${total.toLocaleString('en-US')} items`,
    granted: 'Free Grants',
    purchased: 'Purchases',
    used: 'Spent',
    refunded: 'Refunded',
    netChange: 'Net Change',
    loading: 'Loading point history.',
    empty: 'No point history in the selected range.',
    previous: 'Previous',
    next: 'Next',
    count: (start, end, total) => `${start.toLocaleString('en-US')}-${end.toLocaleString('en-US')} of ${total.toLocaleString('en-US')}`,
    zero: 'No point history to show.',
    transactionTypes: {
      grant_free: 'Free point grant',
      purchase: 'Point purchase',
      use_course_gen: 'Exploration plan generation',
      use_lesson_rec: 'Recommendation search',
      refund: 'Point refund/restore',
    },
  },
}

export function getDashboardSettingsPointsCopy(locale: Locale | string | null | undefined): DashboardSettingsPointsCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
