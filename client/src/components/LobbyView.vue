<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import websocketService from '../services/websocketService'

export default defineComponent({
  name: 'LobbyView',
  data() {
    return {
      nameInput: '' as string,
      reconnecting: false as boolean,
    }
  },
  created() {
    const savedName = sessionStorage.getItem('playerName') || ''
    const savedId = sessionStorage.getItem('playerId') || ''
    if (savedName) {
      this.nameInput = savedName
    }
    // Auto-resume if we already have a prior player id from this browser tab.
    if (savedName && savedId) {
      this.reconnecting = true
      this.resume()
    }
  },
  methods: {
    async resume() {
      store.playerName = this.nameInput.trim() || sessionStorage.getItem('playerName') || ''
      if (!store.playerName) {
        this.reconnecting = false
        return
      }
      await websocketService.connect()
      websocketService.sendJoin(store.playerName)
    },
    async join() {
      if (!this.nameInput.trim()) return
      this.reconnecting = false
      await websocketService.connect()
      websocketService.sendJoin(this.nameInput.trim())
    },
  },
})
</script>

<template>
  <div class="lobby">
    <h2>{{ reconnecting ? 'Rejoining…' : 'Enter Lobby' }}</h2>
    <p v-if="reconnecting" class="hint">Restoring your seat, XP, and position</p>
    <input v-model="nameInput" placeholder="Your name" :disabled="reconnecting" />
    <button @click="join" :disabled="!nameInput.trim()">
      {{ reconnecting ? 'Join as new name' : 'Join Game' }}
    </button>
  </div>
</template>

<style scoped>
.lobby {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  min-height: 100vh;
  justify-content: center;
  background: radial-gradient(ellipse at top, #1e293b, #0f172a);
  color: #f8fafc;
}

.hint {
  margin: 0;
  font-size: 0.85rem;
  color: #94a3b8;
}

input {
  padding: 0.55rem 0.75rem;
  border-radius: 6px;
  border: 1px solid #334155;
  background: #0f172a;
  color: #f8fafc;
  min-width: 220px;
}

button {
  padding: 0.55rem 1rem;
  border-radius: 6px;
  border: none;
  background: #2563eb;
  color: white;
  font-weight: 600;
  cursor: pointer;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
