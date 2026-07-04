import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DashboardGateCopy = {
  checkingTitle: string
  checkingBody: string
}

const copies: Record<Locale, DashboardGateCopy> = {
  ko: {
    checkingTitle: '알파 테스트 접근권을 확인하는 중입니다',
    checkingBody: '잠시만 기다려 주세요.',
  },
  en: {
    checkingTitle: 'Checking alpha access',
    checkingBody: 'Please wait a moment.',
  },
}

export function getDashboardGateCopy(locale: Locale | string | null | undefined): DashboardGateCopy {
  return copies[normalizeLocale(locale)]
}
