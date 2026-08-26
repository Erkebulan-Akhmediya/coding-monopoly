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
  background: rgba(5, 10, 24, 0.85);
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
  background: linear-gradient(145deg, #0d1b32 0%, #0f172a 100%);
  border: 1px solid #3b82f6;
  border-radius: 16px;
  width: 100%;
  max-width: 720px;
  max-height: 86vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.8), 0 0 35px rgba(59, 130, 246, 0.25);
  overflow: hidden;
  color: #f1f5f9;
  text-align: left;
}

.instructions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.1rem 1.5rem;
  background: rgba(15, 23, 42, 0.9);
  border-bottom: 1px solid #1e3a5f;
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
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
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
  background: #1e293b;
  border-radius: 20px;
  padding: 2px;
  border: 1px solid #334155;
}

.lang-pill {
  background: transparent;
  border: none;
  color: #94a3b8;
  padding: 0.25rem 0.65rem;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.lang-pill:hover {
  color: #f8fafc;
}

.lang-pill.active {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 0 8px rgba(37, 99, 235, 0.5);
}

.close-icon-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid #334155;
  border-radius: 8px;
  color: #94a3b8;
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
  background: rgba(239, 68, 68, 0.2);
  border-color: #ef4444;
  color: #f87171;
}

.instructions-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 2rem;
  line-height: 1.65;
  font-size: 0.95rem;
  color: #cbd5e1;
  text-align: left;
}

.instructions-body::-webkit-scrollbar {
  width: 8px;
}

.instructions-body::-webkit-scrollbar-track {
  background: #0f172a;
}

.instructions-body::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 4px;
}

.instructions-body::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.instructions-footer {
  padding: 1rem 1.5rem;
  background: rgba(15, 23, 42, 0.9);
  border-top: 1px solid #1e3a5f;
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
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  box-shadow: 0 4px 15px rgba(37, 99, 235, 0.35);
  transition: all 0.2s ease;
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.55);
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
  color: #60a5fa;
  margin: 0 0 0.85rem 0;
  text-align: left;
}

:deep(.inst-h2) {
  font-size: 1.2rem;
  font-weight: 800;
  color: #93c5fd;
  margin: 1.3rem 0 0.5rem 0;
  text-align: left;
}

:deep(.inst-h3) {
  font-size: 1.05rem;
  font-weight: 700;
  color: #38bdf8;
  margin: 1.4rem 0 0.6rem 0;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  text-align: left;
  border-bottom: 1px solid rgba(56, 189, 248, 0.2);
  padding-bottom: 0.3rem;
}

:deep(.inst-p) {
  margin: 0.6rem 0;
  text-align: left;
  color: #cbd5e1;
}

:deep(.inst-divider) {
  border: none;
  border-top: 1px dashed #334155;
  margin: 1.3rem 0;
}

:deep(.inst-bold) {
  color: #f8fafc;
  font-weight: 700;
}

:deep(.inst-code) {
  background: #1e293b;
  color: #38bdf8;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-size: 0.85em;
  border: 1px solid #334155;
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
  color: #cbd5e1;
  text-align: left;
  line-height: 1.5;
}

:deep(.inst-sub-li) {
  margin-left: 1.2rem;
  list-style-type: circle;
  color: #94a3b8;
}

:deep(.inst-li strong),
:deep(.inst-oli strong) {
  color: #f1f5f9;
}
</style>
