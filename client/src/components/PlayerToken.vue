<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { store } from '../store'
import type { Player } from '../store'

const COLOR_PALETTE = [
  '#3b82f6', // blue
  '#10b981', // emerald
  '#ec4899', // pink
  '#f59e0b', // amber
  '#8b5cf6', // purple
  '#06b6d4', // cyan
  '#f97316', // orange
  '#84cc16', // lime
]

export default defineComponent({
  name: 'PlayerToken',
  props: {
    playerId: {
      type: String,
      default: '',
    },
    player: {
      type: Object as PropType<Player>,
      default: null,
    },
    size: {
      type: String,
      default: 'medium', // 'small' | 'medium' | 'large'
    },
  },
  computed: {
    resolvedPlayer(): Player {
      if (this.player) return this.player
      return store.players.find(p => p.id === this.playerId) || {
        id: this.playerId,
        name: 'Player',
        position: 0,
        xp: 0,
      }
    },
    initials(): string {
      const name = this.resolvedPlayer.name || 'P'
      return name.charAt(0).toUpperCase()
    },
    playerColor(): string {
      let hash = 0
      const str = this.resolvedPlayer.id || this.resolvedPlayer.name || 'default'
      for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash)
      }
      const index = Math.abs(hash) % COLOR_PALETTE.length
      return COLOR_PALETTE[index]
    },
    isCurrentTurn(): boolean {
      return store.currentTurnPlayer === this.resolvedPlayer.name
    },
    tooltipText(): string {
      const p = this.resolvedPlayer
      return `${p.name} (${p.xp || 0} XP) - Cell ${p.position ?? 0}`
    },
  },
})
</script>

<template>
  <div
    class="player-token"
    :class="[size, { 'active-turn': isCurrentTurn }]"
    :style="{ backgroundColor: playerColor }"
    :title="tooltipText"
  >
    <span class="token-symbol">{{ initials }}</span>
    <div v-if="isCurrentTurn" class="turn-halo"></div>
  </div>
</template>

<style scoped>
.player-token {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
  border: 2px solid #ffffff;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  user-select: none;
  flex-shrink: 0;
}

.player-token.small {
  width: 1.25rem;
  height: 1.25rem;
  font-size: 0.65rem;
  border-width: 1.5px;
}

.player-token.medium {
  width: 1.65rem;
  height: 1.65rem;
  font-size: 0.8rem;
}

.player-token.large {
  width: 2.25rem;
  height: 2.25rem;
  font-size: 1rem;
}

.player-token:hover {
  transform: scale(1.15);
  z-index: 10;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.token-symbol {
  line-height: 1;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.player-token.active-turn {
  animation: pulse-border 1.5s infinite ease-in-out;
}

.turn-halo {
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  border: 2px solid #fbbf24;
  animation: halo-ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
  pointer-events: none;
}

@keyframes pulse-border {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(251, 191, 36, 0.7);
  }
  50% {
    box-shadow: 0 0 0 6px rgba(251, 191, 36, 0);
  }
}

@keyframes halo-ping {
  75%, 100% {
    transform: scale(1.3);
    opacity: 0;
  }
}
</style>

