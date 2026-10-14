<script lang="ts">
import { defineComponent } from 'vue'
import { localeCycle, setLocale, type AppLocale } from '../i18n'

const LOCALE_LABELS: Record<AppLocale, string> = {
  en: 'EN',
  ru: 'RU',
  kk: 'KK',
}

export default defineComponent({
  name: 'LanguageSwitcher',
  computed: {
    currentLocale(): AppLocale {
      return this.$i18n.locale as AppLocale
    },
    currentLabel(): string {
      return LOCALE_LABELS[this.currentLocale] || 'EN'
    },
    nextLocale(): AppLocale {
      const idx = localeCycle.indexOf(this.currentLocale)
      return localeCycle[(idx + 1) % localeCycle.length]
    },
    nextLanguageName(): string {
      return this.$t(`language.${this.nextLocale}`)
    },
  },
  methods: {
    cycleLanguage() {
      setLocale(this.nextLocale)
    },
  },
})
</script>

<template>
  <button
    class="lang-switcher"
    type="button"
    :title="$t('language.switchTo', { language: nextLanguageName })"
    :aria-label="$t('language.switch')"
    @click="cycleLanguage"
  >
    🌐 {{ currentLabel }}
  </button>
</template>

<style scoped>
.lang-switcher {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.4rem 0.75rem;
  border-radius: 8px;
  border: 1px solid color-mix(in srgb, var(--c-bright) 15%, transparent);
  background: color-mix(in srgb, var(--c-slate-900) 92%, transparent);
  color: var(--c-slate-200);
  font-size: 0.78rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  cursor: pointer;
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 16px color-mix(in srgb, var(--c-shadow) 35%, transparent);
  transition: border-color 0.2s, transform 0.15s, background 0.2s;
}

.lang-switcher:hover {
  border-color: var(--c-blue-400);
  background: color-mix(in srgb, var(--c-slate-800) 95%, transparent);
  transform: translateY(-1px);
}
</style>
