<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import type { Player } from '../store'
import { localeText } from '../i18n'
import PlayerToken from './PlayerToken.vue'

export default defineComponent({
  name: 'LandedCellPreview',
  components: { PlayerToken },
  computed: {
    cellIndex(): number | null {
      return store.landedCellIndex
    },
    cell(): any | null {
      if (this.cellIndex === null) return null
      return store.boardCells[this.cellIndex] ?? null
    },
    player(): Player | null {
      if (!store.landedCellPlayerId) return null
      return store.players.find(p => p.id === store.landedCellPlayerId) || null
    },
    visible(): boolean {
      return this.cellIndex !== null && !!this.cell
    },
    isCorner(): boolean {
      return [0, 8, 16, 24].includes(this.cellIndex ?? -1)
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
    cellDisplayName(): string {
      return localeText(this.cell?.name) || this.$t('board.cellFallback', { index: this.cellIndex ?? '?' })
    },
    effectText(): string {
      return store.lastEffect || ''
    },
  },
})
</script>

<template>
  <div v-if="visible" class="landed-preview" aria-live="polite">
    <p class="landed-label">{{ $t('board.landedOnLabel') }}</p>
    <div
      class="landed-cell"
      :class="[cellTypeClass, xpGainClass, { corner: isCorner }]"
    >
      <div class="landed-header">
        <span class="landed-index">#{{ cellIndex }}</span>
        <span class="landed-icon">{{ cellIcon }}</span>
      </div>
      <div class="landed-body">
        <span class="landed-name">{{ cellDisplayName }}</span>
        <span v-if="cellSubtitle" class="landed-subtitle">{{ cellSubtitle }}</span>
      </div>
      <div v-if="player" class="landed-token">
        <PlayerToken :player="player" size="large" />
      </div>
    </div>
    <p v-if="effectText" class="landed-effect">{{ effectText }}</p>
  </div>
</template>

<style scoped>
.landed-preview {
  position: absolute;
  inset: 0;
  z-index: 20;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem;
  background: color-mix(in srgb, var(--c-slate-950) 78%, transparent);
  box-sizing: border-box;
  pointer-events: none;
  animation: preview-fade-in 0.25s ease-out;
  container-type: size;
}

.landed-label {
  margin: 0;
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--c-slate-400);
}

/*
 * Board cells are 130×85 in the main grid. Keep that ratio and size the
 * copy to half the hub’s available height (width follows; never overflows).
 */
.landed-cell {
  position: relative;
  box-sizing: border-box;
  height: min(50cqh, calc(100cqw * 85 / 130));
  width: auto;
  aspect-ratio: 130 / 85;
  max-width: 100%;
  border-radius: 10px;
  border: 2px solid var(--c-slate-600);
  background: var(--c-slate-700);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 0.65rem 0.85rem;
  box-shadow: 0 16px 40px color-mix(in srgb, var(--c-shadow) 45%, transparent);
  animation: preview-pop 0.35s cubic-bezier(0.22, 1, 0.36, 1);
}

.landed-cell.corner {
  background: var(--c-slate-800);
  border-width: 3px;
}

.landed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.landed-index {
  color: var(--c-slate-400);
  font-weight: 700;
  font-size: clamp(0.75rem, 2.2vh, 1.1rem);
}

.landed-icon {
  font-size: clamp(1.4rem, 5vh, 2.4rem);
  line-height: 1;
}

.landed-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 0.25rem;
  flex: 1;
  justify-content: center;
  min-height: 0;
}

.landed-name {
  font-size: clamp(0.95rem, 3.2vh, 1.55rem);
  font-weight: 800;
  color: var(--c-slate-50);
  line-height: 1.15;
}

.landed-subtitle {
  font-size: clamp(0.75rem, 2.4vh, 1.15rem);
  font-weight: 700;
  color: var(--c-sky-400);
}

.landed-token {
  display: flex;
  justify-content: center;
  padding-top: 0.15rem;
}

.landed-effect {
  margin: 0;
  max-width: 28rem;
  font-size: 0.9rem;
  font-style: italic;
  color: var(--c-slate-300);
  text-align: center;
}

/* Match BoardView cell type accents */
.cell-deploy {
  background: linear-gradient(135deg, var(--c-deploy--emerald-800) 0%, var(--c-deploy--emerald-700) 100%);
  border-color: var(--c-deploy--emerald-500);
}

.cell-code_freeze {
  background: linear-gradient(135deg, var(--c-code_freeze--blue-900) 0%, var(--c-code_freeze--blue-700) 100%);
  border-color: var(--c-code_freeze--blue-500);
}

.cell-coffee_break {
  background: linear-gradient(135deg, var(--c-coffee_break--amber-900) 0%, var(--c-coffee_break--amber-700) 100%);
  border-color: var(--c-coffee_break--amber-500);
}

.cell-deadline {
  background: linear-gradient(135deg, var(--c-deadline--red-900) 0%, var(--c-deadline--red-700) 100%);
  border-color: var(--c-deadline--red-500);
}

.cell-xp_gain {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-xp_gain--emerald-500);
}

.cell-xp_gain.xp-gain-sm {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-600) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top-color: var(--c-xp_gain--emerald-600);
}

.cell-xp_gain.xp-gain-md {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-500) 45%, transparent) 0%, var(--c-slate-700) 100%);
  border-top-color: var(--c-xp_gain--emerald-500);
}

.cell-xp_gain.xp-gain-lg {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-400) 65%, transparent) 0%, var(--c-slate-700) 100%);
  border-top-color: var(--c-xp_gain--emerald-400);
}

.cell-xp_loss {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_loss--red-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-xp_loss--red-500);
}

.cell-mystery {
  background: linear-gradient(180deg, var(--c-mystery--indigo-900) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-mystery--indigo-400);
}

.cell-teleport {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-teleport--violet-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-teleport--violet-500);
}

.cell-skip_next {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-skip_next--orange-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-skip_next--orange-500);
}

.cell-double_xp {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-double_xp--yellow-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-double_xp--yellow-500);
}

.cell-free_pass {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-free_pass--cyan-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-free_pass--cyan-500);
}

.cell-special_challenge {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-special_challenge--pink-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 5px solid var(--c-special_challenge--pink-500);
}

@keyframes preview-fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes preview-pop {
  from {
    opacity: 0;
    transform: scale(0.86);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@media (max-width: 820px) {
  .landed-cell {
    padding: 0.4rem 0.5rem;
  }
  .landed-label {
    font-size: 0.7rem;
  }
  .landed-effect {
    font-size: 0.75rem;
  }
}
</style>
