import { reactive } from 'vue'

export type ThemeId = 'midnight' | 'ivory' | 'terminal' | 'mono' | 'solarpunk'

export interface ThemeOption {
  id: ThemeId
  /** Background, surface and accent colors for the picker preview (theme-independent). */
  swatch: [string, string, string]
}

/** Available themes in picker order; colors themselves live in themes.css. */
export const THEMES: ThemeOption[] = [
  { id: 'midnight', swatch: ['#0f172a', '#334155', '#3b82f6'] },
  { id: 'ivory', swatch: ['#f5f4ed', '#e8e6dc', '#c96442'] },
  { id: 'terminal', swatch: ['#1a1b26', '#2f3549', '#7aa2f7'] },
  { id: 'mono', swatch: ['#ffffff', '#c8c8c8', '#111111'] },
  { id: 'solarpunk', swatch: ['#f6f2df', '#ddd6b6', '#2a8c82'] },
]

const STORAGE_KEY = 'ui-theme'
const DEFAULT_THEME: ThemeId = 'midnight'

function isTheme(value: unknown): value is ThemeId {
  return THEMES.some((theme) => theme.id === value)
}

function savedTheme(): ThemeId {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (isTheme(saved)) return saved
  } catch {
    // Storage can be unavailable (private mode, blocked site data).
  }
  return DEFAULT_THEME
}

export const themeState = reactive({
  current: savedTheme(),
})

export function setTheme(id: ThemeId): void {
  themeState.current = id
  document.documentElement.dataset.theme = id
  try {
    localStorage.setItem(STORAGE_KEY, id)
  } catch {
    // The choice still applies for this page load.
  }
}

// Apply before the app mounts so the first paint already uses the saved theme.
document.documentElement.dataset.theme = themeState.current
