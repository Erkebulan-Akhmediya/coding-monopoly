<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import websocketService from '../services/websocketService'

export default defineComponent({
  name: 'LobbyView',
  data() {
    return {
      nameInput: '' as string,
      roomInput: '' as string,
      reconnecting: false as boolean,
      joining: false as boolean,
      store,
    }
  },
  created() {
    const savedName = sessionStorage.getItem('playerName') || ''
    const savedId = sessionStorage.getItem('playerId') || ''
    const savedRoom = sessionStorage.getItem('roomId') || ''
    if (savedName) this.nameInput = savedName
    if (savedRoom) this.roomInput = savedRoom
    // Auto-resume if we already have a prior player id from this browser tab.
    if (savedName && savedId) {
      this.reconnecting = true
      this.resume()
    }
  },
  methods: {
    async resume() {
      store.playerName = this.nameInput.trim() || sessionStorage.getItem('playerName') || ''
      store.roomId = this.roomInput.trim() || sessionStorage.getItem('roomId') || ''
      if (!store.playerName || !store.roomId) {
        this.reconnecting = false
        return
      }
      store.joinError = ''
      this.joining = true
      await websocketService.connect()
      websocketService.sendJoin(store.playerName, store.roomId)
      this.joining = false
    },
    async join() {
      const name = this.nameInput.trim()
      const roomId = this.roomInput.trim()
      if (!name || !roomId) return
      this.reconnecting = false
      store.joinError = ''
      store.roomId = roomId
      this.joining = true
      await websocketService.connect()
      websocketService.sendJoin(name, roomId)
      this.joining = false
    },
  },
})
</script>

<template>
  <div class="lobby">
    <div class="lobby-card">
      <div class="lobby-logo">🎲 Coding Monopoly</div>
      <p class="lobby-subtitle">{{ reconnecting ? 'Rejoining your game…' : 'Enter your details to join' }}</p>

      <p v-if="reconnecting" class="hint">Restoring your seat, XP, and position</p>

      <div class="field">
        <label for="lobby-name">Your Name</label>
        <input
          id="lobby-name"
          v-model="nameInput"
          type="text"
          placeholder="Alice"
          class="lobby-input"
          :disabled="reconnecting || joining"
          @keydown.enter="join"
        />
      </div>

      <div class="field">
        <label for="lobby-room">Room ID</label>
        <input
          id="lobby-room"
          v-model="roomInput"
          type="text"
          placeholder="Ask your instructor for the room ID"
          class="lobby-input"
          :disabled="reconnecting || joining"
          @keydown.enter="join"
        />
      </div>

      <p v-if="store.joinError" class="join-error">{{ store.joinError }}</p>

      <button
        id="lobby-join-btn"
        class="lobby-btn"
        :disabled="(!nameInput.trim() || !roomInput.trim()) && !reconnecting || joining"
        @click="join"
      >
        {{ reconnecting ? 'Join as New Player' : joining ? 'Joining…' : 'Join Game' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.lobby {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: radial-gradient(ellipse at top, #1e293b 0%, #0f172a 60%);
  color: #f8fafc;
  padding: 1rem;
}

.lobby-card {
  background: linear-gradient(145deg, #0f1e3a 0%, #111827 100%);
  border: 1px solid #1e3a5f;
  border-radius: 16px;
  padding: 2.5rem 2.75rem;
  width: 360px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(96, 165, 250, 0.1);
  display: flex;
  flex-direction: column;
  gap: 0;
}

.lobby-logo {
  font-size: 1.6rem;
  font-weight: 900;
  text-align: center;
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.3rem;
}

.lobby-subtitle {
  text-align: center;
  color: #64748b;
  font-size: 0.85rem;
  margin: 0 0 1.5rem 0;
}

.hint {
  text-align: center;
  color: #94a3b8;
  font-size: 0.8rem;
  margin: -0.75rem 0 0.75rem 0;
  font-style: italic;
}

.join-error {
  text-align: center;
  color: #f87171;
  font-size: 0.85rem;
  margin: 0 0 0.75rem 0;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  margin-bottom: 1rem;
}

.field label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.lobby-input {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 0.55rem 0.85rem;
  color: #f1f5f9;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.lobby-input:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}

.lobby-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.lobby-btn {
  width: 100%;
  padding: 0.65rem;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 700;
  font-size: 0.95rem;
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  transition: opacity 0.2s, transform 0.15s;
  margin-top: 0.25rem;
}

.lobby-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.lobby-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
