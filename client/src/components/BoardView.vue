<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import type { Player } from '../store'
import PlayerToken from './PlayerToken.vue'
import DiceOverlay from './DiceOverlay.vue'

export default defineComponent({
  name: 'BoardView',
  components: { PlayerToken, DiceOverlay },
  data() {
    return {
      remaining: 0 as number,
      intervalId: null as number | null,
      store: store,
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
      if (!store.currentTurnPlayer) return 'Waiting for game turn...'
      if (this.isMyTurn) {
        return '⚡ Your turn!'
      }
      return `Waiting for ${store.currentTurnPlayer}`
    },
    showCountdown(): boolean {
      return store.questionActive && store.deadline > 0
    },
    countdown(): number {
      const now = Date.now()
      const diff = Math.max(0, store.deadline - now)
      return Math.ceil(diff / 1000)
    },
    sortedPlayers(): Player[] {
      return [...store.players].sort((a, b) => (b.xp || 0) - (a.xp || 0))
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
      return store.players.filter(p => (p.position ?? 0) === cellIndex)
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
        if (cell.type === 'xp_gain') return `+${cell.params.amount} XP`
        if (cell.type === 'xp_loss') return `-${cell.params.amount} XP`
      }
      if (cell.type === 'deploy') return 'START (+100)'
      if (cell.type === 'double_xp') return '2x XP'
      if (cell.type === 'skip_next') return 'Skip Turn'
      if (cell.type === 'free_pass') return 'CI Shield'
      if (cell.type === 'teleport') return 'Fast-Track'
      if (cell.type === 'special_challenge') return cell.params?.bonus ? `+${cell.params.bonus} XP` : 'Bonus'
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
          { corner: isCorner(idx) }
        ]"
        :style="getCellGridStyle(idx)"
      >
        <div class="cell-header">
          <span class="cell-index">#{{ idx }}</span>
          <span class="cell-icon">{{ getCellIcon(cell.type) }}</span>
        </div>

        <div class="cell-body">
          <span class="cell-name">{{ cell.name || ('Cell ' + idx) }}</span>
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
          />
        </div>
      </div>

      <!-- Center Hub Area (Grid Rows 2-8, Cols 2-8) -->
      <div class="board-center-hub">
        <div class="hub-header">
          <h1 class="game-logo">⚡ CODING MONOPOLY ⚡</h1>
          <div class="turn-card" :class="{ 'my-turn': isMyTurn }">
            <span class="turn-text">{{ turnMessage }}</span>
            <div v-if="showCountdown" class="countdown-badge">
              ⏳ {{ remaining }}s left
            </div>
          </div>
        </div>

        <!-- Dice Roll & Effect Overlay Component -->
        <DiceOverlay />

        <!-- Leaderboard / Player Overview Panel -->
        <div class="leaderboard">
          <h3>🎮 PLAYERS ({{ sortedPlayers.length }})</h3>
          <div class="player-list">
            <div
              v-for="p in sortedPlayers"
              :key="p.id"
              class="player-card"
              :class="{ active: store.currentTurnPlayer === p.name }"
            >
              <PlayerToken :player="p" size="small" />
              <div class="player-info">
                <span class="player-name">
                  {{ p.name }}
                  <span v-if="p.name === store.playerName" class="you-tag">(You)</span>
                </span>
                <span class="player-pos">Cell #{{ p.position ?? 0 }}</span>
              </div>
              <div class="player-xp-badge">{{ p.xp ?? 0 }} XP</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.board-view {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem;
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  box-sizing: border-box;
}

.board-grid {
  display: grid;
  grid-template-columns: repeat(9, minmax(0, 85px));
  grid-template-rows: repeat(9, minmax(0, 85px));
  gap: 4px;
  background: #1e293b;
  padding: 8px;
  border-radius: 12px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6), 0 0 20px rgba(59, 130, 246, 0.15);
  border: 1px solid #334155;
  width: max-content;
}

/* Base Board Cell */
.board-cell {
  position: relative;
  background: #334155;
  border-radius: 6px;
  border: 1px solid #475569;
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
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.4);
}

.cell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.65rem;
}

.cell-index {
  color: #94a3b8;
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
  color: #f1f5f9;
  line-height: 1.1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.cell-subtitle {
  font-size: 0.55rem;
  font-weight: 600;
  color: #38bdf8;
  margin-top: 2px;
}

/* Tokens Container inside Cell */
.cell-tokens {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  justify-content: center;
  align-items: center;
  min-height: 1.25rem;
  padding-top: 2px;
}

/* Corner Cells Special Styling */
.board-cell.corner {
  background: #1e293b;
  border-width: 2px;
}

.cell-deploy {
  background: linear-gradient(135deg, #065f46 0%, #047857 100%);
  border-color: #10b981;
}

.cell-code_freeze {
  background: linear-gradient(135deg, #1e3a8a 0%, #1d4ed8 100%);
  border-color: #3b82f6;
}

.cell-coffee_break {
  background: linear-gradient(135deg, #78350f 0%, #b45309 100%);
  border-color: #f59e0b;
}

.cell-deadline {
  background: linear-gradient(135deg, #7f1d1d 0%, #b91c1c 100%);
  border-color: #ef4444;
}

/* Cell Types Accent Colors */
.cell-xp_gain {
  border-top: 3px solid #10b981;
}

.cell-xp_loss {
  border-top: 3px solid #ef4444;
}

.cell-mystery {
  background: linear-gradient(180deg, #312e81 0%, #334155 100%);
  border-top: 3px solid #818cf8;
}

.cell-teleport {
  border-top: 3px solid #8b5cf6;
}

.cell-skip_next {
  border-top: 3px solid #f97316;
}

.cell-double_xp {
  border-top: 3px solid #eab308;
}

.cell-free_pass {
  border-top: 3px solid #06b6d4;
}

.cell-special_challenge {
  border-top: 3px solid #ec4899;
}

/* Center Hub Area (Grid Rows 2..8, Cols 2..8) */
.board-center-hub {
  grid-row: 2 / span 7;
  grid-column: 2 / span 7;
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(8px);
  border-radius: 8px;
  border: 1px dashed #475569;
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
}

.game-logo {
  font-size: 1.4rem;
  font-weight: 900;
  margin: 0;
  letter-spacing: 2px;
  background: linear-gradient(90deg, #60a5fa, #a78bfa, #f472b6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 2px 10px rgba(96, 165, 250, 0.3);
}

.turn-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.4rem 1rem;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 20px;
  font-weight: 700;
  font-size: 0.9rem;
}

.turn-card.my-turn {
  background: linear-gradient(90deg, #065f46, #047857);
  border-color: #10b981;
  color: #a7f3d0;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.4);
}

.countdown-badge {
  background: #ef4444;
  color: white;
  padding: 0.15rem 0.5rem;
  border-radius: 12px;
  font-size: 0.8rem;
  animation: pulse-red 1s infinite;
}

@keyframes pulse-red {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* Leaderboard */
.leaderboard {
  background: #1e293b;
  border-radius: 8px;
  padding: 0.75rem;
  border: 1px solid #334155;
}

.leaderboard h3 {
  margin: 0 0 0.5rem 0;
  font-size: 0.75rem;
  letter-spacing: 1px;
  color: #94a3b8;
  text-align: center;
}

.player-list {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  max-height: 160px;
  overflow-y: auto;
}

.player-card {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.3rem 0.6rem;
  background: #0f172a;
  border-radius: 6px;
  border: 1px solid #334155;
  font-size: 0.8rem;
}

.player-card.active {
  border-color: #fbbf24;
  background: #172554;
}

.player-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.player-name {
  font-weight: 700;
  color: #f8fafc;
}

.you-tag {
  color: #60a5fa;
  font-size: 0.7rem;
  margin-left: 2px;
}

.player-pos {
  font-size: 0.65rem;
  color: #94a3b8;
}

.player-xp-badge {
  font-weight: 800;
  color: #34d399;
  font-size: 0.85rem;
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
