export const supportedLocales = ['ko', 'en'] as const

export type Locale = (typeof supportedLocales)[number]

export const defaultLocale: Locale = 'ko'

export function normalizeLocale(value: string | null | undefined): Locale {
  return value === 'en' ? 'en' : defaultLocale
}
