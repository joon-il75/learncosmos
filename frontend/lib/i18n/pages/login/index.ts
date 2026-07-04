import type { Locale } from '../../locales'
import { loginEn } from './en'
import { loginKo } from './ko'
import type { LoginPageCopy } from './types'

export function getLoginPageCopy(locale: Locale): LoginPageCopy {
  return locale === 'en' ? loginEn : loginKo
}
