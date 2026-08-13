<script lang="ts">
import { defineComponent } from 'vue'
import { store, type GameOverSummary, type PlayerStanding } from '../store'

export default defineComponent({
  name: 'EndGameSummary',
  data() {
    return { store }
  },
  computed: {
    summary(): GameOverSummary | null {
      return store.gameOver
    },
    standings(): PlayerStanding[] {
      return this.summary?.standings || []
    },
    reasonText(): string {
      if (!this.summary) return ''
      if (this.summary.reason === 'admin') return 'Match ended by admin'
      return `First to ${this.summary.target_xp} XP`
    },
  },
})
</script>

<template>
  <div v-if="summary" class="end-overlay" role="dialog" aria-modal="true">
    <div class="end-card">
      <p class="eyebrow">Match complete</p>
      <h2 class="winner">{{ summary.winner_name || 'Winner' }} wins</h2>
      <p class="reason">{{ reasonText }}</p>

      <ol class="standings">
        <li
          v-for="row in standings"
          :key="row.player_id"
          :class="{ champion: row.rank === 1 }"
        >
          <span class="rank">#{{ row.rank }}</span>
          <span class="name">
            {{ row.name }}
            <span v-if="row.name === store.playerName" class="you">(You)</span>
          </span>
          <span class="xp">{{ row.xp }} XP</span>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped>
.end-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  background: rgba(2, 6, 23, 0.82);
  backdrop-filter: blur(4px);
}

.end-card {
  width: min(420px, 100%);
  padding: 1.75rem 1.5rem;
  border-radius: 12px;
  background: linear-gradient(160deg, #1e293b, #0f172a);
  border: 1px solid #334155;
  color: #f8fafc;
  text-align: center;
  animation: rise 0.35s ease-out;
}

.eyebrow {
  margin: 0;
  font-size: 0.72rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #94a3b8;
}

.winner {
  margin: 0.4rem 0 0.25rem;
  font-size: 1.75rem;
  color: #fbbf24;
}

.reason {
  margin: 0 0 1.25rem;
  color: #cbd5e1;
  font-size: 0.9rem;
}

.standings {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  text-align: left;
}

.standings li {
  display: grid;
  grid-template-columns: 2.5rem 1fr auto;
  gap: 0.5rem;
  align-items: center;
  padding: 0.55rem 0.7rem;
  border-radius: 8px;
  background: #0f172a;
  border: 1px solid #334155;
}

.standings li.champion {
  border-color: #fbbf24;
  background: #172554;
}

.rank {
  font-weight: 800;
  color: #94a3b8;
}

.name {
  font-weight: 700;
}

.you {
  color: #60a5fa;
  font-size: 0.75rem;
  font-weight: 600;
}

.xp {
  font-weight: 800;
  color: #34d399;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
