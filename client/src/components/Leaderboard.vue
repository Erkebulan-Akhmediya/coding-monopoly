<script lang="ts">
import { defineComponent } from 'vue'
import { store, type Player } from '../store'
import PlayerToken from './PlayerToken.vue'

export default defineComponent({
  name: 'Leaderboard',
  components: { PlayerToken },
  data() {
    return {
      store: store,
    }
  },
  computed: {
    sortedPlayers(): Player[] {
      return [...store.players].sort((a, b) => (b.xp || 0) - (a.xp || 0))
    },
    targetXP(): number {
      return store.targetXP || 500
    },
  },
  methods: {
    rankFor(index: number): number {
      return index + 1
    },
    isActive(p: Player): boolean {
      return store.currentTurnPlayer === p.name
    },
    progress(p: Player): number {
      if (!this.targetXP) return 0
      return Math.max(0, Math.min(100, Math.round(((p.xp || 0) / this.targetXP) * 100)))
    },
  },
})
</script>

<template>
  <div class="leaderboard">
    <h3>Live Standings</h3>
    <p class="target">First to {{ targetXP }} XP</p>
    <div class="player-list">
      <div
        v-for="(p, idx) in sortedPlayers"
        :key="p.id"
        class="player-card"
        :class="{
          active: isActive(p),
          offline: p.is_connected === false,
        }"
      >
        <span class="rank">{{ rankFor(idx) }}</span>
        <PlayerToken :player="p" size="small" />
        <div class="player-info">
          <span class="player-name">
            {{ p.name }}
            <span v-if="p.name === store.playerName" class="you-tag">(You)</span>
            <span v-if="p.is_connected === false" class="offline-tag">offline</span>
          </span>
          <span class="player-pos">Cell #{{ p.position ?? 0 }}</span>
          <div class="xp-bar" aria-hidden="true">
            <div class="xp-fill" :style="{ width: progress(p) + '%' }" />
          </div>
        </div>
        <div class="player-xp-badge">{{ p.xp ?? 0 }} XP</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard {
  background: #1e293b;
  border-radius: 8px;
  padding: 0.75rem;
  border: 1px solid #334155;
  position: absolute;
  top: 10px;
  left: 10px;
  width: 270px;
  z-index: 30;
}

.leaderboard h3 {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 1px;
  color: #94a3b8;
  text-align: center;
  text-transform: uppercase;
}

.target {
  margin: 0.2rem 0 0.55rem;
  text-align: center;
  font-size: 0.7rem;
  color: #64748b;
}

.player-list {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  max-height: 520px;
  overflow-y: auto;
}

.player-card {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.35rem 0.55rem;
  background: #0f172a;
  border-radius: 6px;
  border: 1px solid #334155;
  font-size: 0.8rem;
}

.player-card.active {
  border-color: #fbbf24;
  background: #172554;
}

.player-card.offline {
  opacity: 0.55;
}

.rank {
  width: 1.1rem;
  font-weight: 800;
  color: #64748b;
  font-size: 0.75rem;
}

.player-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.player-name {
  font-weight: 700;
  color: #f8fafc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.you-tag {
  color: #60a5fa;
  font-size: 0.7rem;
  margin-left: 2px;
}

.offline-tag {
  color: #f87171;
  font-size: 0.65rem;
  margin-left: 4px;
  font-weight: 600;
}

.player-pos {
  font-size: 0.65rem;
  color: #94a3b8;
}

.xp-bar {
  margin-top: 0.2rem;
  height: 4px;
  border-radius: 999px;
  background: #1e293b;
  overflow: hidden;
}

.xp-fill {
  height: 100%;
  background: linear-gradient(90deg, #22c55e, #34d399);
  transition: width 0.35s ease;
}

.player-xp-badge {
  font-weight: 800;
  color: #34d399;
  font-size: 0.85rem;
  white-space: nowrap;
}
</style>
