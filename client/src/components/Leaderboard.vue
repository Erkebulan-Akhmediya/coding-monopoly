<script lang="ts">
import { defineComponent } from 'vue'
import { store, type Player } from '../store';
export default defineComponent({
  name: 'Leaderboard',
  data() {
    return {
      store: store,
    }
  },
  computed: {
    sortedPlayers(): Player[] {
      return [...store.players].sort((a, b) => (b.xp || 0) - (a.xp || 0))
    },
  }
});
</script>

<template>
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
</template>

<style scoped>
/* Leaderboard */
.leaderboard {
  background: #1e293b;
  border-radius: 8px;
  padding: 0.75rem;
  border: 1px solid #334155;

  position: absolute;
  top: 10px;
  left: 10px;

  width: 300px;
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
  max-height: 500px;
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
</style>