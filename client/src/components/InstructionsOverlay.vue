<script lang="ts">
import { defineComponent } from 'vue'
import { setLocale, type AppLocale } from '../i18n'
import enMarkdown from '../instructions/en.md?raw'
import ruMarkdown from '../instructions/ru.md?raw'
import kkMarkdown from '../instructions/kk.md?raw'

export default defineComponent({
  name: 'InstructionsOverlay',
  emits: ['close'],
  computed: {
    currentLocale(): AppLocale {
      return (this.$i18n.locale as AppLocale) || 'en'
    },
    activeMarkdown(): string {
      const loc = this.currentLocale
      if (loc === 'ru') return ruMarkdown
      if (loc === 'kk') return kkMarkdown
      return enMarkdown
    },
    renderedHtml(): string {
      return this.renderMarkdown(this.activeMarkdown)
    },
  },
  mounted() {
    window.addEventListener('keydown', this.handleKeyDown)
  },
  beforeUnmount() {
    window.removeEventListener('keydown', this.handleKeyDown)
  },
  methods: {
    close() {
      this.$emit('close')
    },
    handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        this.close()
      }
    },
    changeLanguage(lang: AppLocale) {
      setLocale(lang)
    },
    escapeHtml(text: string): string {
      return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;')
    },
    formatInline(text: string): string {
      // Bold **text**
      let formatted = text.replace(/\*\*(.+?)\*\*/g, '<strong class="inst-bold">$1</strong>')
      // Inline `code`
      formatted = formatted.replace(/`([^`]+)`/g, '<code class="inst-code">$1</code>')
      return formatted
    },
    renderMarkdown(md: string): string {
      if (!md) return ''
      const lines = md.split(/\r?\n/)
      const htmlParts: string[] = []
      let listStack: string[] = [] // 'ul' or 'ol'

      const closeAllLists = () => {
        while (listStack.length > 0) {
          const tag = listStack.pop()
          htmlParts.push(`</${tag}>`)
        }
      }

      for (let i = 0; i < lines.length; i++) {
        const rawLine = lines[i]
        const trimmed = rawLine.trim()

        if (!trimmed) {
          closeAllLists()
          continue
        }

        // Horizontal rule
        if (trimmed === '---' || trimmed === '***' || trimmed === '___') {
          closeAllLists()
          htmlParts.push('<hr class="inst-divider" />')
          continue
        }

        // Heading 1
        if (trimmed.startsWith('# ')) {
          closeAllLists()
          const content = this.formatInline(this.escapeHtml(trimmed.slice(2)))
          htmlParts.push(`<h1 class="inst-h1">${content}</h1>`)
          continue
        }

        // Heading 2
        if (trimmed.startsWith('## ')) {
          closeAllLists()
          const content = this.formatInline(this.escapeHtml(trimmed.slice(3)))
          htmlParts.push(`<h2 class="inst-h2">${content}</h2>`)
          continue
        }

        // Heading 3
        if (trimmed.startsWith('### ')) {
          closeAllLists()
          const content = this.formatInline(this.escapeHtml(trimmed.slice(4)))
          htmlParts.push(`<h3 class="inst-h3">${content}</h3>`)
          continue
        }

        // Check indentation level for sub-lists
        const indentMatch = rawLine.match(/^(\s*)/)
        const indentLevel = indentMatch ? Math.floor(indentMatch[1].length / 2) : 0

        // Unordered list item (- or *)
        const ulMatch = trimmed.match(/^[-*]\s+(.*)$/)
        if (ulMatch) {
          if (listStack.length === 0 || listStack[listStack.length - 1] !== 'ul') {
            closeAllLists()
            htmlParts.push('<ul class="inst-ul">')
            listStack.push('ul')
          }
          const itemText = this.formatInline(this.escapeHtml(ulMatch[1]))
          const subClass = indentLevel > 0 ? ' inst-sub-li' : ''
          htmlParts.push(`<li class="inst-li${subClass}">${itemText}</li>`)
          continue
        }

        // Ordered list item (1. 2. etc)
        const olMatch = trimmed.match(/^(\d+)\.\s+(.*)$/)
        if (olMatch) {
          if (listStack.length === 0 || listStack[listStack.length - 1] !== 'ol') {
            closeAllLists()
            htmlParts.push('<ol class="inst-ol">')
            listStack.push('ol')
          }
          const itemText = this.formatInline(this.escapeHtml(olMatch[2]))
          htmlParts.push(`<li class="inst-oli">${itemText}</li>`)
          continue
        }

        // Standard paragraph
        closeAllLists()
        const pContent = this.formatInline(this.escapeHtml(trimmed))
        htmlParts.push(`<p class="inst-p">${pContent}</p>`)
      }

      closeAllLists()
      return htmlParts.join('\n')
    },
  },
})
</script>

<template>
  <div class="instructions-backdrop" role="dialog" aria-modal="true" @click.self="close">
    <div class="instructions-card">
      <!-- Modal Header -->
      <header class="instructions-header">
        <div class="header-left">
          <span class="header-icon">📖</span>
          <h2 class="instructions-title">{{ $t('instructions.modalTitle') }}</h2>
        </div>

        <div class="header-right">
          <!-- Language Tabs -->
          <div class="language-tabs" role="group" aria-label="Language selection">
            <button
              class="lang-pill"
              :class="{ active: currentLocale === 'en' }"
              type="button"
              @click="changeLanguage('en')"
            >
              EN
            </button>
            <button
              class="lang-pill"
              :class="{ active: currentLocale === 'ru' }"
              type="button"
              @click="changeLanguage('ru')"
            >
              RU
            </button>
            <button
              class="lang-pill"
              :class="{ active: currentLocale === 'kk' }"
              type="button"
              @click="changeLanguage('kk')"
            >
              KK
            </button>
          </div>

          <!-- Close Icon Button -->
          <button
            class="close-icon-btn"
            type="button"
            aria-label="Close instructions"
            @click="close"
          >
            ✕
          </button>
        </div>
      </header>

      <!-- Modal Content Body -->
      <main class="instructions-body">
        <div class="markdown-content" v-html="renderedHtml"></div>
      </main>

      <!-- Modal Footer Action -->
      <footer class="instructions-footer">
        <button
          id="instructions-close-btn"
          class="confirm-btn"
          type="button"
          @click="close"
        >
          🚀 {{ $t('instructions.closeBtn') }}
        </button>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.instructions-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9500;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--c-navy-960) 85%, transparent);
  backdrop-filter: blur(8px);
  padding: 1.5rem;
  box-sizing: border-box;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.97);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.instructions-card {
  background: linear-gradient(145deg, var(--c-navy-880) 0%, var(--c-slate-900) 100%);
  border: 1px solid var(--c-blue-500);
  border-radius: 16px;
  width: 100%;
  max-width: 720px;
  max-height: 86vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 60px color-mix(in srgb, var(--c-shadow) 80%, transparent), 0 0 35px color-mix(in srgb, var(--c-blue-500) 25%, transparent);
  overflow: hidden;
  color: var(--c-slate-100);
  text-align: left;
}

.instructions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.1rem 1.5rem;
  background: color-mix(in srgb, var(--c-slate-900) 90%, transparent);
  border-bottom: 1px solid var(--c-navy-700);
  gap: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.header-icon {
  font-size: 1.35rem;
}

.instructions-title {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 800;
  background: linear-gradient(90deg, var(--c-blue-400), var(--c-violet-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.language-tabs {
  display: flex;
  background: var(--c-slate-800);
  border-radius: 20px;
  padding: 2px;
  border: 1px solid var(--c-slate-700);
}

.lang-pill {
  background: transparent;
  border: none;
  color: var(--c-slate-400);
  padding: 0.25rem 0.65rem;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.lang-pill:hover {
  color: var(--c-slate-50);
}

.lang-pill.active {
  background: var(--c-solid--blue-600);
  color: var(--c-white);
  box-shadow: 0 0 8px color-mix(in srgb, var(--c-solid--blue-600) 50%, transparent);
}

.close-icon-btn {
  background: color-mix(in srgb, var(--c-bright) 5%, transparent);
  border: 1px solid var(--c-slate-700);
  border-radius: 8px;
  color: var(--c-slate-400);
  font-size: 1rem;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.close-icon-btn:hover {
  background: color-mix(in srgb, var(--c-red-500) 20%, transparent);
  border-color: var(--c-red-500);
  color: var(--c-red-400);
}

.instructions-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 2rem;
  line-height: 1.65;
  font-size: 0.95rem;
  color: var(--c-slate-300);
  text-align: left;
}

.instructions-body::-webkit-scrollbar {
  width: 8px;
}

.instructions-body::-webkit-scrollbar-track {
  background: var(--c-slate-900);
}

.instructions-body::-webkit-scrollbar-thumb {
  background: var(--c-slate-700);
  border-radius: 4px;
}

.instructions-body::-webkit-scrollbar-thumb:hover {
  background: var(--c-slate-600);
}

.instructions-footer {
  padding: 1rem 1.5rem;
  background: color-mix(in srgb, var(--c-slate-900) 90%, transparent);
  border-top: 1px solid var(--c-navy-700);
  display: flex;
  justify-content: center;
}

.confirm-btn {
  width: 100%;
  max-width: 360px;
  padding: 0.8rem 1.5rem;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-weight: 700;
  font-size: 1rem;
  background: linear-gradient(90deg, var(--c-solid--blue-600), var(--c-solid--violet-600));
  color: var(--c-white);
  box-shadow: 0 4px 15px color-mix(in srgb, var(--c-solid--blue-600) 35%, transparent);
  transition: all 0.2s ease;
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px color-mix(in srgb, var(--c-blue-600) 55%, transparent);
}

.confirm-btn:active {
  transform: translateY(0);
}

/* Markdown Rendered Elements Styling */
.markdown-content {
  text-align: left;
}

:deep(.inst-h1) {
  font-size: 1.4rem;
  font-weight: 900;
  color: var(--c-blue-400);
  margin: 0 0 0.85rem 0;
  text-align: left;
}

:deep(.inst-h2) {
  font-size: 1.2rem;
  font-weight: 800;
  color: var(--c-blue-300);
  margin: 1.3rem 0 0.5rem 0;
  text-align: left;
}

:deep(.inst-h3) {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--c-sky-400);
  margin: 1.4rem 0 0.6rem 0;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  text-align: left;
  border-bottom: 1px solid color-mix(in srgb, var(--c-sky-400) 20%, transparent);
  padding-bottom: 0.3rem;
}

:deep(.inst-p) {
  margin: 0.6rem 0;
  text-align: left;
  color: var(--c-slate-300);
}

:deep(.inst-divider) {
  border: none;
  border-top: 1px dashed var(--c-slate-700);
  margin: 1.3rem 0;
}

:deep(.inst-bold) {
  color: var(--c-slate-50);
  font-weight: 700;
}

:deep(.inst-code) {
  background: var(--c-slate-800);
  color: var(--c-sky-400);
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-size: 0.85em;
  border: 1px solid var(--c-slate-700);
}

:deep(.inst-ul) {
  margin: 0.5rem 0 0.9rem 1.25rem;
  padding: 0 0 0 1rem;
  list-style-type: disc;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  text-align: left;
}

:deep(.inst-ol) {
  margin: 0.5rem 0 0.9rem 1.25rem;
  padding: 0 0 0 1rem;
  list-style-type: decimal;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  text-align: left;
}

:deep(.inst-li),
:deep(.inst-oli) {
  color: var(--c-slate-300);
  text-align: left;
  line-height: 1.5;
}

:deep(.inst-sub-li) {
  margin-left: 1.2rem;
  list-style-type: circle;
  color: var(--c-slate-400);
}

:deep(.inst-li strong),
:deep(.inst-oli strong) {
  color: var(--c-slate-100);
}
</style>
