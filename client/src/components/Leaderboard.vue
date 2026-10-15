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
    <h3>{{ $t('leaderboard.title') }}</h3>
    <p class="target">{{ $t('leaderboard.target', { xp: targetXP }) }}</p>
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
            <span v-if="p.name === store.playerName" class="you-tag">{{ $t('common.you') }}</span>
            <span v-if="p.is_connected === false" class="offline-tag">{{ $t('common.offline') }}</span>
          </span>
          <span class="player-pos">{{ $t('leaderboard.cellPosition', { index: p.position ?? 0 }) }}</span>
          <div class="xp-bar" aria-hidden="true">
            <div class="xp-fill" :style="{ width: progress(p) + '%' }" />
          </div>
        </div>
        <div class="player-xp-badge">{{ p.xp ?? 0 }} {{ $t('common.xp') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard {
  background: var(--c-slate-800);
  border-radius: 8px;
  padding: 0.75rem;
  border: 1px solid var(--c-slate-700);
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
  color: var(--c-slate-400);
  text-align: center;
  text-transform: uppercase;
}

.target {
  margin: 0.2rem 0 0.55rem;
  text-align: center;
  font-size: 0.7rem;
  color: var(--c-slate-500);
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
  background: var(--c-slate-900);
  border-radius: 6px;
  border: 1px solid var(--c-slate-700);
  font-size: 0.8rem;
}

.player-card.active {
  border-color: var(--c-amber-400);
  background: var(--c-blue-950);
}

.player-card.offline {
  opacity: 0.55;
}

.rank {
  width: 1.1rem;
  font-weight: 800;
  color: var(--c-slate-500);
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
  color: var(--c-slate-50);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.you-tag {
  color: var(--c-blue-400);
  font-size: 0.7rem;
  margin-left: 2px;
}

.offline-tag {
  color: var(--c-red-400);
  font-size: 0.65rem;
  margin-left: 4px;
  font-weight: 600;
}

.player-pos {
  font-size: 0.65rem;
  color: var(--c-slate-400);
}

.xp-bar {
  margin-top: 0.2rem;
  height: 4px;
  border-radius: 999px;
  background: var(--c-slate-800);
  overflow: hidden;
}

.xp-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--c-green-500), var(--c-emerald-400));
  transition: width 0.35s ease;
}

.player-xp-badge {
  font-weight: 800;
  color: var(--c-emerald-400);
  font-size: 0.85rem;
  white-space: nowrap;
}
</style>
