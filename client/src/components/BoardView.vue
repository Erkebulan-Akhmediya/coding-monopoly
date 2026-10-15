<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import type { Player } from '../store'
import { localeText } from '../i18n'
import { getTokenVisualPosition } from '../services/tokenMovement'
import PlayerToken from './PlayerToken.vue'
import DiceOverlay from './DiceOverlay.vue'
import LandedCellPreview from './LandedCellPreview.vue'
import Leaderboard from './Leaderboard.vue'
import GameActionPanel from './GameActionPanel.vue'
import EffectToastStack from './EffectToastStack.vue'
import EndGameSummary from './EndGameSummary.vue'
import PauseOverlay from './PauseOverlay.vue'
import InstructionsOverlay from './InstructionsOverlay.vue'
import CellDetailCard from './CellDetailCard.vue'

export default defineComponent({
  name: 'BoardView',
  components: { 
    PlayerToken, 
    DiceOverlay,
    LandedCellPreview,
    Leaderboard, 
    GameActionPanel,
    EffectToastStack,
    EndGameSummary,
    PauseOverlay,
    InstructionsOverlay,
    CellDetailCard,
  },
  data() {
    return {
      remaining: 0 as number,
      intervalId: null as number | null,
      store: store,
      showInstructions: true as boolean,
      selectedCellIndex: null as number | null,
    }
  },
  computed: {
    cells(): any[] {
      return store.boardCells
    },
    cornerIndexes(): number[] {
      return [0, 8, 16, 24]
    },
    isMyTurn(): boolean {
      return store.currentTurnPlayer === store.playerName && store.playerName !== ''
    },
    turnMessage(): string {
      if (!store.currentTurnPlayer) return this.$t('board.waitingTurn')
      if (this.isMyTurn) {
        return `⚡ ${this.$t('board.yourTurn')}`
      }
      return this.$t('board.waitingFor', { player: store.currentTurnPlayer })
    },
    showCountdown(): boolean {
      return store.questionActive && store.deadline > 0
    },
    countdown(): number {
      const now = Date.now()
      const diff = Math.max(0, store.deadline - now)
      return Math.ceil(diff / 1000)
    },
  },
  watch: {
    'store.deadline'() {
      this.updateRemaining()
    },
  },
  methods: {
    isCorner(idx: number): boolean {
      return this.cornerIndexes.includes(idx)
    },
    getPlayersAtCell(cellIndex: number): Player[] {
      // Use animated visual positions so tokens hop cell-to-cell after dice settle.
      return store.players.filter(
        p => getTokenVisualPosition(p.id, p.position ?? 0) === cellIndex
      )
    },
    isTokenHopping(playerId: string): boolean {
      return store.hoppingPlayerId === playerId
    },
    getCellGridStyle(idx: number) {
      let row = 1
      let col = 1
      if (idx >= 0 && idx <= 8) {
        row = 9
        col = 9 - idx
      } else if (idx >= 8 && idx <= 16) {
        col = 1
        row = 9 - (idx - 8)
      } else if (idx >= 16 && idx <= 24) {
        row = 1
        col = 1 + (idx - 16)
      } else if (idx >= 24 && idx <= 31) {
        col = 9
        row = 1 + (idx - 24)
      }
      return {
        gridRowStart: row,
        gridColumnStart: col,
      }
    },
    getCellIcon(type: string): string {
      switch (type) {
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
    getCellSubtitle(cell: any): string {
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
    cellName(cell: any): string {
      return localeText(cell?.name) || this.$t('board.cellFallback', { index: cell?.cell_index ?? '?' })
    },
    getXpGainClass(cell: any): string {
      if (cell && cell.type === 'xp_gain' && cell.params && typeof cell.params.amount === 'number') {
        const amount = cell.params.amount
        if (amount <= 20) return 'xp-gain-sm'
        if (amount <= 50) return 'xp-gain-md'
        return 'xp-gain-lg'
      }
      return ''
    },
    updateRemaining() {
      if (this.showCountdown) {
        this.remaining = this.countdown
      } else {
        this.remaining = 0
      }
    },
    startTimer() {
      if (this.intervalId !== null) return
      this.intervalId = setInterval(() => {
        this.updateRemaining()
      }, 1000) as unknown as number
    },
    stopTimer() {
      if (this.intervalId !== null) {
        clearInterval(this.intervalId)
        this.intervalId = null
      }
    },
    openCellCard(idx: number) {
      this.selectedCellIndex = idx
    },
    closeCellCard() {
      this.selectedCellIndex = null
    },
  },
  beforeMount() {
    if (!store.playerName)
      store.playerName = sessionStorage.getItem('playerName') || ''
  },
  mounted() {
    this.startTimer()
  },
  beforeUnmount() {
    this.stopTimer()
  },
})
</script>

<template>
  <div class="board-view">
    <!-- Monopoly 9x9 Grid Layout -->
    <div class="board-grid">
      <!-- 32 Perimeter Cells -->
      <div
        v-for="(cell, idx) in cells"
        :key="idx"
        class="board-cell"
        :class="[
          'cell-' + (cell.type || 'generic'),
          getXpGainClass(cell),
          {
            corner: isCorner(idx),
            'effect-highlight': store.highlightedCellIndex === idx,
          }
        ]"
        :style="getCellGridStyle(idx)"
        role="button"
        tabindex="0"
        :title="cellName(cell)"
        @click="openCellCard(idx)"
        @keydown.enter="openCellCard(idx)"
        @keydown.space.prevent="openCellCard(idx)"
      >
        <div class="cell-header">
          <span class="cell-index">#{{ idx }}</span>
          <span class="cell-icon">{{ getCellIcon(cell.type) }}</span>
        </div>

        <div class="cell-body">
          <span class="cell-name">{{ cellName(cell) }}</span>
          <span v-if="getCellSubtitle(cell)" class="cell-subtitle">
            {{ getCellSubtitle(cell) }}
          </span>
        </div>

        <!-- Players Tokens on this cell -->
        <div class="cell-tokens" v-if="getPlayersAtCell(idx).length > 0">
          <PlayerToken
            v-for="p in getPlayersAtCell(idx)"
            :key="p.id"
            :player="p"
            size="small"
            :hopping="isTokenHopping(p.id)"
          />
        </div>
      </div>

      <!-- Center Hub Area (Grid Rows 2-8, Cols 2-8) -->
      <div class="board-center-hub">
        <div class="hub-header">
          <div class="hub-top-row">
            <h1 class="game-logo">⚡ {{ $t('board.gameLogo') }} ⚡</h1>
            <button
              id="open-instructions-btn"
              class="instructions-btn"
              type="button"
              :title="$t('instructions.openBtn')"
              @click="showInstructions = true"
            >
              {{ $t('instructions.openBtn') }}
            </button>
          </div>
          <div class="turn-card" :class="{ 'my-turn': isMyTurn }">
            <span class="turn-text">{{ turnMessage }}</span>
            <div v-if="showCountdown" class="countdown-badge">
              ⏳ {{ $t('board.countdownLeft', { seconds: remaining }) }}
            </div>
          </div>
        </div>

        <!-- Dice Roll & Effect Overlay Component -->
        <DiceOverlay />

        <!-- Magnified destination cell after token finishes hopping -->
        <LandedCellPreview />

        <!-- Cell Detail Card (shown when a cell is clicked) -->
        <CellDetailCard
          v-if="selectedCellIndex !== null && cells[selectedCellIndex]"
          :cell="cells[selectedCellIndex]"
          :cell-index="selectedCellIndex"
          @close="closeCellCard"
        />
      </div>
    </div>

    <!-- Leaderboard / Player Overview Panel -->
    <Leaderboard />

    <!-- Game Action Panel: LevelPicker, ProblemPanel -->
    <GameActionPanel />

    <EffectToastStack />
    <EndGameSummary />
    <PauseOverlay />
    <InstructionsOverlay
      v-if="showInstructions"
      @close="showInstructions = false"
    />
  </div>
</template>

<style scoped>
.board-view {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem;
  min-height: 100vh;
  background: var(--c-slate-900);
  color: var(--c-slate-50);
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  box-sizing: border-box;
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(9, minmax(0, 130px));
  grid-template-rows: repeat(9, minmax(0, 85px));
  gap: 4px;
  background: var(--c-slate-800);
  padding: 8px;
  border-radius: 12px;
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 60%, transparent), 0 0 20px color-mix(in srgb, var(--c-blue-500) 15%, transparent);
  border: 1px solid var(--c-slate-700);
  width: max-content;
}

/* Base Board Cell */
.board-cell {
  position: relative;
  background: var(--c-slate-700);
  border-radius: 6px;
  border: 1px solid var(--c-slate-600);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 4px;
  box-sizing: border-box;
  overflow: hidden;
  transition: all 0.2s ease;
}

.board-cell:hover {
  transform: scale(1.04);
  z-index: 5;
  box-shadow: 0 6px 16px color-mix(in srgb, var(--c-shadow) 40%, transparent);
}

.board-cell.effect-highlight {
  animation: cell-effect-pulse 0.9s ease-in-out 2;
  box-shadow: 0 0 0 2px var(--c-amber-400), 0 0 18px color-mix(in srgb, var(--c-amber-400) 55%, transparent);
  z-index: 6;
}

@keyframes cell-effect-pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.06);
  }
}

/* Let the hop jump paint above the cell edge */
.board-cell:has(.player-token.hopping) {
  overflow: visible;
  z-index: 8;
}

.board-cell:hover .cell-tokens {
  opacity: 0.2;
  pointer-events: none;
}

.cell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.65rem;
}

.cell-index {
  color: var(--c-slate-400);
  font-weight: 700;
  font-size: 0.6rem;
}

.cell-icon {
  font-size: 0.85rem;
}

.cell-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  margin: 2px 0;
}

.cell-name {
  font-size: 0.65rem;
  font-weight: 700;
  color: var(--c-slate-100);
  line-height: 1.1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.cell-subtitle {
  font-size: 0.55rem;
  font-weight: 600;
  color: var(--c-sky-400);
  margin-top: 2px;
}

/* Tokens Container inside Cell */
.cell-tokens {
  display: grid;
  grid-template-columns: repeat(3, auto);
  grid-template-rows: repeat(3, auto);
  place-items: center;

  position: absolute;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;

  min-height: 1.25rem;
  padding-top: 2px;
  transition: opacity 0.2s ease;
}

/* Corner Cells Special Styling */
.board-cell.corner {
  background: var(--c-slate-800);
  border-width: 2px;
}

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

/* Cell Types Accent Colors (Using gradients matching their borders for better visibility) */
.cell-xp_gain {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-xp_gain--emerald-500);
}

/* Distinct Intensity for XP Gain Sizes */
.cell-xp_gain.xp-gain-sm {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-600) 15%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-xp_gain--emerald-600);
}

.cell-xp_gain.xp-gain-md {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-500) 35%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-xp_gain--emerald-500);
}

.cell-xp_gain.xp-gain-lg {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_gain--emerald-400) 60%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-xp_gain--emerald-400);
  box-shadow: inset 0 2px 10px color-mix(in srgb, var(--c-xp_gain--emerald-400) 15%, transparent);
}

.cell-xp_loss {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-xp_loss--red-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-xp_loss--red-500);
}

.cell-mystery {
  background: linear-gradient(180deg, var(--c-mystery--indigo-900) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-mystery--indigo-400);
}

.cell-teleport {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-teleport--violet-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-teleport--violet-500);
}

.cell-skip_next {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-skip_next--orange-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-skip_next--orange-500);
}

.cell-double_xp {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-double_xp--yellow-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-double_xp--yellow-500);
}

.cell-free_pass {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-free_pass--cyan-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-free_pass--cyan-500);
}

.cell-special_challenge {
  background: linear-gradient(180deg, color-mix(in srgb, var(--c-special_challenge--pink-500) 25%, transparent) 0%, var(--c-slate-700) 100%);
  border-top: 3px solid var(--c-special_challenge--pink-500);
}

/* Center Hub Area (Grid Rows 2..8, Cols 2..8) */
.board-center-hub {
  grid-row: 2 / span 7;
  grid-column: 2 / span 7;
  background: color-mix(in srgb, var(--c-slate-900) 85%, transparent);
  backdrop-filter: blur(8px);
  border-radius: 8px;
  border: 1px dashed var(--c-slate-600);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1.25rem;
  box-sizing: border-box;
  overflow: hidden;
  position: relative;
}

.hub-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
}

.hub-top-row {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: 100%;
}

.instructions-btn {
  position: absolute;
  right: 0;
  background: color-mix(in srgb, var(--c-slate-800) 85%, transparent);
  border: 1px solid var(--c-slate-700);
  color: var(--c-blue-300);
  border-radius: 20px;
  padding: 0.3rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.35rem;
  transition: all 0.2s ease;
}

.instructions-btn:hover {
  background: var(--c-navy-700);
  border-color: var(--c-blue-500);
  color: var(--c-bright);
  transform: translateY(-1px);
}

.game-logo {
  font-size: 1.4rem;
  font-weight: 900;
  margin: 0;
  letter-spacing: 2px;
  background: linear-gradient(90deg, var(--c-blue-400), var(--c-violet-400), var(--c-pink-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 2px 10px color-mix(in srgb, var(--c-blue-400) 30%, transparent);
}

.turn-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.4rem 1rem;
  background: var(--c-slate-800);
  border: 1px solid var(--c-slate-700);
  border-radius: 20px;
  font-weight: 700;
  font-size: 0.9rem;
}

.turn-card.my-turn {
  background: linear-gradient(90deg, var(--c-emerald-800), var(--c-emerald-700));
  border-color: var(--c-emerald-500);
  color: var(--c-emerald-200);
  box-shadow: 0 0 12px color-mix(in srgb, var(--c-emerald-500) 40%, transparent);
}

.countdown-badge {
  background: var(--c-red-500);
  color: var(--c-white);
  padding: 0.15rem 0.5rem;
  border-radius: 12px;
  font-size: 0.8rem;
  animation: pulse-red 1s infinite;
}

@keyframes pulse-red {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

@media (max-width: 820px) {
  .board-grid {
    grid-template-columns: repeat(9, minmax(0, 36px));
    grid-template-rows: repeat(9, minmax(0, 36px));
    gap: 2px;
  }
  .cell-name, .cell-subtitle {
    display: none;
  }
  .board-center-hub {
    padding: 0.5rem;
  }
  .game-logo {
    font-size: 0.9rem;
  }
}
</style>