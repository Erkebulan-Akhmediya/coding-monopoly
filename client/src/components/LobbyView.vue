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
      <div class="lobby-logo">🎲 {{ $t('lobby.title') }}</div>
      <p class="lobby-subtitle">{{ reconnecting ? $t('lobby.subtitleRejoin') : $t('lobby.subtitleJoin') }}</p>

      <p v-if="reconnecting" class="hint">{{ $t('lobby.hintRejoin') }}</p>

      <div class="field">
        <label for="lobby-name">{{ $t('lobby.yourName') }}</label>
        <input
          id="lobby-name"
          v-model="nameInput"
          type="text"
          :placeholder="$t('lobby.namePlaceholder')"
          class="lobby-input"
          :disabled="reconnecting || joining"
          @keydown.enter="join"
        />
      </div>

      <div class="field">
        <label for="lobby-room">{{ $t('lobby.roomId') }}</label>
        <input
          id="lobby-room"
          v-model="roomInput"
          type="text"
          :placeholder="$t('lobby.roomPlaceholder')"
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
        {{ reconnecting ? $t('lobby.joinAsNew') : joining ? $t('lobby.joining') : $t('lobby.joinGame') }}
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
  background: radial-gradient(ellipse at top, var(--c-slate-800) 0%, var(--c-slate-900) 60%);
  color: var(--c-slate-50);
  padding: 1rem;
}

.lobby-card {
  background: linear-gradient(145deg, var(--c-navy-800) 0%, var(--c-gray-900) 100%);
  border: 1px solid var(--c-navy-700);
  border-radius: 16px;
  padding: 2.5rem 2.75rem;
  width: 360px;
  box-shadow: 0 24px 60px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 0 1px color-mix(in srgb, var(--c-blue-400) 10%, transparent);
  display: flex;
  flex-direction: column;
  gap: 0;
}

.lobby-logo {
  font-size: 1.6rem;
  font-weight: 900;
  text-align: center;
  background: linear-gradient(90deg, var(--c-blue-400), var(--c-violet-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.3rem;
}

.lobby-subtitle {
  text-align: center;
  color: var(--c-slate-500);
  font-size: 0.85rem;
  margin: 0 0 1.5rem 0;
}

.hint {
  text-align: center;
  color: var(--c-slate-400);
  font-size: 0.8rem;
  margin: -0.75rem 0 0.75rem 0;
  font-style: italic;
}

.join-error {
  text-align: center;
  color: var(--c-red-400);
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
  color: var(--c-slate-400);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.lobby-input {
  background: var(--c-slate-900);
  border: 1px solid var(--c-slate-700);
  border-radius: 8px;
  padding: 0.55rem 0.85rem;
  color: var(--c-slate-100);
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.lobby-input:focus {
  border-color: var(--c-blue-400);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--c-blue-400) 20%, transparent);
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
  background: linear-gradient(90deg, var(--c-solid--blue-600), var(--c-solid--violet-600));
  color: var(--c-white);
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
