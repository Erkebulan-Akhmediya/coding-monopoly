<script lang="ts">
import { defineComponent } from 'vue'
import { THEMES, setTheme, themeState, type ThemeId, type ThemeOption } from '../theme'

export default defineComponent({
  name: 'ThemeSwitcher',

  data() {
    return {
      open: false as boolean,
      themes: THEMES as ThemeOption[],
    }
  },

  computed: {
    current(): ThemeId {
      return themeState.current
    },
  },

  mounted() {
    document.addEventListener('pointerdown', this.handleOutside)
    window.addEventListener('keydown', this.handleKeyDown)
  },

  beforeUnmount() {
    document.removeEventListener('pointerdown', this.handleOutside)
    window.removeEventListener('keydown', this.handleKeyDown)
  },

  methods: {
    toggle() {
      this.open = !this.open
    },

    choose(id: ThemeId) {
      setTheme(id)
      this.open = false
    },

    handleOutside(e: PointerEvent) {
      if (this.open && !this.$el.contains(e.target as Node)) this.open = false
    },

    handleKeyDown(e: KeyboardEvent) {
      if (this.open && e.key === 'Escape') this.open = false
    },
  },
})
</script>

<template>
  <div class="theme-switcher">
    <button
      id="theme-switcher-btn"
      class="theme-btn"
      type="button"
      aria-haspopup="listbox"
      :aria-expanded="open"
      :title="$t('theme.switch')"
      @click="toggle"
    >
      🎨 {{ $t(`theme.names.${current}`) }}
    </button>

    <ul v-if="open" class="theme-menu" role="listbox" :aria-label="$t('theme.title')">
      <li
        v-for="theme in themes"
        :key="theme.id"
        role="option"
        :aria-selected="theme.id === current"
      >
        <button
          :id="`theme-option-${theme.id}`"
          type="button"
          class="theme-option"
          :class="{ active: theme.id === current }"
          @click="choose(theme.id)"
        >
          <span class="swatch" aria-hidden="true">
            <span v-for="(color, i) in theme.swatch" :key="i" :style="{ background: color }"></span>
          </span>
          <span class="theme-text">
            <span class="theme-name">{{ $t(`theme.names.${theme.id}`) }}</span>
            <span class="theme-desc">{{ $t(`theme.descriptions.${theme.id}`) }}</span>
          </span>
          <span v-if="theme.id === current" class="check" aria-hidden="true">✓</span>
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
/* The menu is positioned against the fixed controls cluster in App.vue, so it
   right-aligns with the screen edge instead of overflowing on narrow screens. */
.theme-switcher {
  position: static;
}

.theme-btn {
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

.theme-btn:hover {
  border-color: var(--c-blue-400);
  background: color-mix(in srgb, var(--c-slate-800) 95%, transparent);
  transform: translateY(-1px);
}

.theme-menu {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  width: 300px;
  max-width: calc(100vw - 24px);
  margin: 0;
  padding: 0.35rem;
  list-style: none;
  background: var(--c-slate-900);
  border: 1px solid var(--c-slate-700);
  border-radius: 10px;
  box-shadow: 0 12px 32px color-mix(in srgb, var(--c-shadow) 45%, transparent);
  text-align: left;
}

.theme-option {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  width: 100%;
  padding: 0.5rem 0.55rem;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: var(--c-slate-200);
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.theme-option:hover {
  background: var(--c-slate-800);
}

.theme-option.active {
  border-color: var(--c-blue-400);
  background: var(--c-navy-800);
}

.swatch {
  display: inline-flex;
  flex-shrink: 0;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--c-slate-600);
}

.swatch span {
  width: 12px;
  height: 26px;
}

.theme-text {
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
  min-width: 0;
  flex: 1;
  line-height: 1.3;
}

.theme-name {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--c-slate-100);
}

.theme-desc {
  font-size: 0.7rem;
  color: var(--c-slate-400);
}

.check {
  color: var(--c-blue-400);
  font-weight: 800;
}
</style>
