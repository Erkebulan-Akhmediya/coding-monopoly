<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import PlayerToken from './PlayerToken.vue'

export default defineComponent({
  name: 'BoardView',
  components: { PlayerToken },
  data() {
    return {
      remaining: 0 as number,
      intervalId: null as number | null,
      store: store,
    }
  },
  computed: {
    // Full list of cells from server
    cells(): any[] {
      return store.boardCells
    },
    // Indexes of the four corners (assuming 0‑based array of length 32)
    cornerIndexes(): number[] {
      return [0, 7, 24, 31]
    },
    isMyTurn(): boolean {
      return store.currentTurnPlayer === store.playerName && store.playerName !== ''
    },
    turnMessage(): string {
      if (!store.currentTurnPlayer) return ''
      if (this.isMyTurn) {
        return 'Your turn'
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
  },
  watch: {
    // Update remaining whenever deadline changes
    'store.deadline'() {
      this.updateRemaining()
    },
  },
  methods: {
    isCorner(idx: number): boolean {
      return this.cornerIndexes.includes(idx)
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
    <div class="turn-indicator">
      <span class="turn-msg">{{ turnMessage }}</span>
      <span v-if="showCountdown" class="countdown">{{ remaining }}s left</span>
    </div>
    <div class="players">
      <PlayerToken v-for="p in store.players" :key="p.id" :player-id="p.id" />
    </div>
    <div class="board-grid">
      <div
        v-for="(, idx) in cells"
        :key="idx"
        class="board-cell"
        :class="{ corner: isCorner(idx) }"
      >
        <span class="cell-index">{{ idx }}</span>
        <!-- placeholder for cell-specific info, e.g., type/name -->
      </div>
    </div>
  </div>
</template>

<style scoped>
.board-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}
.turn-indicator {
  font-weight: bold;
}
.players {
  display: flex;
  gap: 0.5rem;
}
.board-grid {
  display: grid;
  grid-template-columns: repeat(8, 40px);
  gap: 2px;
}
.board-cell {
  width: 40px;
  height: 40px;
  background: #eee;
  border: 1px solid #999;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
}
.board-cell.corner {
  background: #ffd54f;
}
</style>
