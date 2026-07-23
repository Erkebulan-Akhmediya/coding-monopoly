<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import websocketService from '../services/websocketService'

export default defineComponent({
  name: 'LobbyView',
  data() {
    return {
      nameInput: '' as string,
    }
  },
  methods: {
    async join() {
      if (!this.nameInput.trim()) return
      store.playerName = this.nameInput.trim()
      await websocketService.connect()
      // after connection, server will broadcast presence and set store.connected
      console.log('connected:', store.connected)
    },
  },
})
</script>

<template>
  <div class="lobby">
    <h2>Enter Lobby</h2>
    <input v-model="nameInput" placeholder="Your name" />
    <button @click="join">Join Game</button>
  </div>
</template>

<style scoped>
.lobby {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}
</style>
