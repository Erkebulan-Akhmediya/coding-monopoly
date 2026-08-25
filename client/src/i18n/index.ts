import { createI18n } from 'vue-i18n'
import en from './locales/en'
import ru from './locales/ru'
import kk from './locales/kk'

export type AppLocale = 'en' | 'ru' | 'kk'

const STORAGE_KEY = 'ui-locale'

const LOCALES: AppLocale[] = ['en', 'ru', 'kk']

function getSavedLocale(): AppLocale {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved && LOCALES.includes(saved as AppLocale)) {
    return saved as AppLocale
  }
  return 'en'
}

export const i18n = createI18n({
  legacy: true,
  locale: getSavedLocale(),
  fallbackLocale: 'en',
  messages: { en, ru, kk },
})

export function setLocale(locale: AppLocale): void {
  i18n.global.locale = locale
  localStorage.setItem(STORAGE_KEY, locale)
  document.documentElement.lang = locale
}

export function t(key: string, params?: Record<string, unknown> | number): string {
  return (i18n.global.t as (key: string, arg?: Record<string, unknown> | number) => string)(
    key,
    params,
  )
}

/** Server-side locale.Text uses `kz` for Kazakh; the frontend i18n uses `kk`. */
type ServerLocale = 'en' | 'ru' | 'kz'

const LOCALE_MAP: Record<AppLocale, ServerLocale> = {
  en: 'en',
  ru: 'ru',
  kk: 'kz',
}

/**
 * Return the string for the current UI locale from a server locale.Text object.
 * Falls back to `.en` if the current locale key is missing or empty.
 */
export function localeText(
  text: { en?: string; ru?: string; kz?: string } | string | null | undefined,
): string {
  if (typeof text === 'string') return text
  if (!text) return ''
  const key = LOCALE_MAP[(i18n.global.locale as AppLocale) || 'en'] || 'en'
  return text[key] || text.en || ''
}

export const localeCycle: AppLocale[] = LOCALES

document.documentElement.lang = getSavedLocale()
