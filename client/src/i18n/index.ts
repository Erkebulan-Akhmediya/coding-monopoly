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

export const localeCycle: AppLocale[] = LOCALES

document.documentElement.lang = getSavedLocale()
