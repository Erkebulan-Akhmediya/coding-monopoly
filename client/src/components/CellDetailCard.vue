<script lang="ts">
import { defineComponent } from 'vue'
import { localeText } from '../i18n'

export default defineComponent({
  name: 'CellDetailCard',
  props: {
    cell: {
      type: Object,
      required: true,
    },
    cellIndex: {
      type: Number,
      required: true,
    },
  },
  emits: ['close'],
  computed: {
    isCorner(): boolean {
      return [0, 8, 16, 24].includes(this.cellIndex)
    },
    cellTypeClass(): string {
      return 'cell-' + (this.cell?.type || 'generic')
    },
    xpGainClass(): string {
      const cell = this.cell
      if (cell && cell.type === 'xp_gain' && cell.params && typeof cell.params.amount === 'number') {
        const amount = cell.params.amount
        if (amount <= 20) return 'xp-gain-sm'
        if (amount <= 50) return 'xp-gain-md'
        return 'xp-gain-lg'
      }
      return ''
    },
    cellIcon(): string {
      switch (this.cell?.type) {
        case 'deploy': return '🚩'
        case 'code_freeze': return '🧊'
        case 'coffee_break': return '☕'
        case 'deadline': return '🚨'
        case 'xp_gain': return '📈'
        case 'xp_loss': return '📉'
        case 'mystery': return '❓'
        case 'teleport': return '🌀'
        case 'skip_next': return '⏭️'
        case 'double_xp': return '⚡'
        case 'free_pass': return '🎟️'
        case 'special_challenge': return '🏆'
        default: return '📍'
      }
    },
    cellDisplayName(): string {
      return localeText(this.cell?.name) || this.$t('board.cellFallback', { index: this.cellIndex })
    },
    cellExplanation(): string {
      const expl = localeText(this.cell?.explanation)
      return expl || this.$t('cellCard.noExplanation')
    },
    cellSubtitle(): string {
      const cell = this.cell
      if (!cell) return ''
      if (cell.params && typeof cell.params.amount === 'number') {
        if (cell.type === 'xp_gain') return `+${cell.params.amount} ${this.$t('common.xp')}`
        if (cell.type === 'xp_loss') return `-${cell.params.amount} ${this.$t('common.xp')}`
      }
      if (cell.type === 'special_challenge') {
        return cell.params?.bonus ? `+${cell.params.bonus} ${this.$t('common.xp')}` : this.$t('board.cellSubtitle.bonus')
      }
      const key = `board.cellSubtitle.${cell.type}`
      if (this.$te(key)) return this.$t(key, cell.params || {})
      return ''
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
  },
})
</script>

<template>
  <div class="cell-detail-overlay" role="dialog" aria-modal="true" @click.self="close">
    <div class="cell-card" :class="[cellTypeClass, xpGainClass, { corner: isCorner }]">
      <!-- Card Header -->
      <div class="card-header">
        <div class="cell-tag">
          <span class="cell-idx">#{{ cellIndex }}</span>
          <span class="cell-emoji">{{ cellIcon }}</span>
        </div>
        <button
          class="card-close-btn"
          type="button"
          aria-label="Close"
          @click="close"
        >
          ✕
        </button>
      </div>

      <!-- Card Title & Subtitle Badge -->
      <div class="card-main">
        <h3 class="cell-title">{{ cellDisplayName }}</h3>
        <div v-if="cellSubtitle" class="cell-subtitle-badge">
          {{ cellSubtitle }}
        </div>
      </div>

      <!-- Card Explanation Text -->
      <div class="card-explanation">
        <p class="explanation-text">{{ cellExplanation }}</p>
      </div>

      <!-- Card Footer -->
      <div class="card-footer">
        <button
          id="cell-card-close-btn"
          class="dismiss-btn"
          type="button"
          @click="close"
        >
          {{ $t('cellCard.close') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cell-detail-overlay {
  position: absolute;
  inset: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--c-slate-950) 75%, transparent);
  backdrop-filter: blur(4px);
  padding: 1rem;
  box-sizing: border-box;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.cell-card {
  position: relative;
  background: linear-gradient(145deg, var(--c-navy-860) 0%, var(--c-slate-900) 100%);
  border: 2px solid var(--c-slate-600);
  border-radius: 16px;
  width: 100%;
  max-width: 440px;
  padding: 1.25rem 1.5rem;
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-blue-500) 20%, transparent);
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  color: var(--c-slate-100);
  box-sizing: border-box;
  animation: popIn 0.25s cubic-bezier(0.2, 0.9, 0.3, 1.1);
}

@keyframes popIn {
  from {
    transform: scale(0.9);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.cell-card.corner {
  background: linear-gradient(145deg, var(--c-navy-840) 0%, var(--c-slate-900) 100%);
  border-width: 3px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cell-tag {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.cell-idx {
  background: var(--c-slate-800);
  border: 1px solid var(--c-slate-700);
  color: var(--c-slate-400);
  font-weight: 800;
  font-size: 0.85rem;
  padding: 0.2rem 0.55rem;
  border-radius: 8px;
}

.cell-emoji {
  font-size: 1.6rem;
}

.card-close-btn {
  background: color-mix(in srgb, var(--c-bright) 6%, transparent);
  border: 1px solid var(--c-slate-700);
  border-radius: 8px;
  color: var(--c-slate-400);
  font-size: 1rem;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.card-close-btn:hover {
  background: color-mix(in srgb, var(--c-red-500) 20%, transparent);
  border-color: var(--c-red-500);
  color: var(--c-red-400);
}

.card-main {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.cell-title {
  margin: 0;
  font-size: 1.3rem;
  font-weight: 800;
  color: var(--c-bright);
  line-height: 1.25;
}

.cell-subtitle-badge {
  display: inline-block;
  align-self: flex-start;
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--c-sky-400);
  background: color-mix(in srgb, var(--c-sky-400) 15%, transparent);
  border: 1px solid color-mix(in srgb, var(--c-sky-400) 30%, transparent);
  padding: 0.2rem 0.6rem;
  border-radius: 6px;
}

.card-explanation {
  background: color-mix(in srgb, var(--c-slate-900) 70%, transparent);
  border: 1px solid var(--c-navy-700);
  border-radius: 10px;
  padding: 0.9rem 1.1rem;
}

.explanation-text {
  margin: 0;
  font-size: 0.95rem;
  line-height: 1.55;
  color: var(--c-slate-300);
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.25rem;
}

.dismiss-btn {
  background: var(--c-slate-800);
  border: 1px solid var(--c-blue-500);
  color: var(--c-blue-300);
  font-weight: 700;
  font-size: 0.88rem;
  padding: 0.45rem 1.1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dismiss-btn:hover {
  background: var(--c-solid--blue-600);
  color: var(--c-white);
  box-shadow: 0 0 12px color-mix(in srgb, var(--c-solid--blue-600) 40%, transparent);
}

/* Cell Type Accent Borders & Glows */
.cell-deploy {
  border-color: var(--c-deploy--emerald-500);
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-deploy--emerald-500) 25%, transparent);
}

.cell-code_freeze {
  border-color: var(--c-code_freeze--blue-500);
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-code_freeze--blue-500) 25%, transparent);
}

.cell-coffee_break {
  border-color: var(--c-coffee_break--amber-500);
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-coffee_break--amber-500) 25%, transparent);
}

.cell-deadline {
  border-color: var(--c-deadline--red-500);
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-deadline--red-500) 25%, transparent);
}

.cell-xp_gain {
  border-color: var(--c-xp_gain--emerald-500);
}

.cell-xp_gain.xp-gain-lg {
  border-color: var(--c-xp_gain--emerald-400);
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 25px color-mix(in srgb, var(--c-xp_gain--emerald-400) 30%, transparent);
}

.cell-xp_loss {
  border-color: var(--c-xp_loss--red-500);
}

.cell-mystery {
  border-color: var(--c-mystery--indigo-400);
}

.cell-teleport {
  border-color: var(--c-teleport--violet-500);
}

.cell-skip_next {
  border-color: var(--c-skip_next--orange-500);
}

.cell-double_xp {
  border-color: var(--c-double_xp--yellow-500);
}

.cell-free_pass {
  border-color: var(--c-free_pass--cyan-500);
}

.cell-special_challenge {
  border-color: var(--c-special_challenge--pink-500);
}
</style>
