<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'

export default defineComponent({
  name: 'PlayerToken',
  props: {
    playerId: {
      type: String,
      required: true,
    },
  },
  computed: {
    player() {
      return store.players.find(p => p.id === this.playerId) || { name: '???', token: '' }
    },
  },
})
</script>

<template>
  <div class="player-token" :title="player.name">
    <!-- Simple visual token, could be a colored circle -->
    <span class="token">{{ player.token || '•' }}</span>
  </div>
</template>

<style scoped>
.player-token {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  background: #4caf50;
  color: white;
  font-weight: bold;
  font-size: 0.9rem;
}
</style>
