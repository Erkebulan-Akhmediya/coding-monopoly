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
      if (this.$te(key)) return this.$t(key)
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
  background: rgba(2, 6, 23, 0.75);
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
  background: linear-gradient(145deg, #111e38 0%, #0f172a 100%);
  border: 2px solid #475569;
  border-radius: 16px;
  width: 100%;
  max-width: 440px;
  padding: 1.25rem 1.5rem;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(59, 130, 246, 0.2);
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  color: #f1f5f9;
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
  background: linear-gradient(145deg, #162444 0%, #0f172a 100%);
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
  background: #1e293b;
  border: 1px solid #334155;
  color: #94a3b8;
  font-weight: 800;
  font-size: 0.85rem;
  padding: 0.2rem 0.55rem;
  border-radius: 8px;
}

.cell-emoji {
  font-size: 1.6rem;
}

.card-close-btn {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid #334155;
  border-radius: 8px;
  color: #94a3b8;
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
  background: rgba(239, 68, 68, 0.2);
  border-color: #ef4444;
  color: #f87171;
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
  color: #ffffff;
  line-height: 1.25;
}

.cell-subtitle-badge {
  display: inline-block;
  align-self: flex-start;
  font-size: 0.8rem;
  font-weight: 700;
  color: #38bdf8;
  background: rgba(56, 189, 248, 0.15);
  border: 1px solid rgba(56, 189, 248, 0.3);
  padding: 0.2rem 0.6rem;
  border-radius: 6px;
}

.card-explanation {
  background: rgba(15, 23, 42, 0.7);
  border: 1px solid #1e3a5f;
  border-radius: 10px;
  padding: 0.9rem 1.1rem;
}

.explanation-text {
  margin: 0;
  font-size: 0.95rem;
  line-height: 1.55;
  color: #cbd5e1;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.25rem;
}

.dismiss-btn {
  background: #1e293b;
  border: 1px solid #3b82f6;
  color: #93c5fd;
  font-weight: 700;
  font-size: 0.88rem;
  padding: 0.45rem 1.1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dismiss-btn:hover {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 0 12px rgba(37, 99, 235, 0.4);
}

/* Cell Type Accent Borders & Glows */
.cell-deploy {
  border-color: #10b981;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(16, 185, 129, 0.25);
}

.cell-code_freeze {
  border-color: #3b82f6;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(59, 130, 246, 0.25);
}

.cell-coffee_break {
  border-color: #f59e0b;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(245, 158, 11, 0.25);
}

.cell-deadline {
  border-color: #ef4444;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(239, 68, 68, 0.25);
}

.cell-xp_gain {
  border-color: #10b981;
}

.cell-xp_gain.xp-gain-lg {
  border-color: #34d399;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.7), 0 0 25px rgba(52, 211, 153, 0.3);
}

.cell-xp_loss {
  border-color: #ef4444;
}

.cell-mystery {
  border-color: #818cf8;
}

.cell-teleport {
  border-color: #8b5cf6;
}

.cell-skip_next {
  border-color: #f97316;
}

.cell-double_xp {
  border-color: #eab308;
}

.cell-free_pass {
  border-color: #06b6d4;
}

.cell-special_challenge {
  border-color: #ec4899;
}
</style>
