<script lang="ts">
import { defineComponent } from 'vue'
import { adminStore } from '../adminStore'
import { adminApiService, type RoomSummary } from '../services/adminApiService'
import adminWS from '../services/adminWebsocketService'
import AdminQuestionList from './AdminQuestionList.vue'
import AdminSpectatorView from './AdminSpectatorView.vue'

export default defineComponent({
  name: 'AdminView',

  components: {
    AdminQuestionList,
    AdminSpectatorView,
  },

  data() {
    return {
      store: adminStore,
      loginPassword: '' as string,
      loginError: '' as string,
      loginLoading: false as boolean,
      activeTab: 'questions' as 'questions' | 'rooms',
      // rooms list
      rooms: [] as RoomSummary[],
      roomsLoading: false as boolean,
      roomsError: '' as string,
      roomsRefreshInterval: null as number | null,
      // selected room for spectating
      selectedRoom: null as string | null,
    }
  },

  computed: {
    isLoggedIn(): boolean {
      return !!adminStore.token
    },
    isWatchingRoom(): boolean {
      return this.activeTab === 'rooms' && this.selectedRoom !== null
    },
  },

  mounted() {
    const savedToken = sessionStorage.getItem('admin_token')
    if (savedToken) {
      adminStore.token = savedToken
    }
  },

  beforeUnmount() {
    this.stopRoomsRefresh()
  },

  methods: {
    async handleLogin() {
      if (!this.loginPassword) return
      this.loginLoading = true
      this.loginError = ''
      try {
        const data = await adminApiService.login(this.loginPassword)
        adminStore.token = data.token
        sessionStorage.setItem('admin_token', data.token)
      } catch (err: any) {
        this.loginError = err.message || 'Invalid password or connection failed'
      } finally {
        this.loginLoading = false
        this.loginPassword = ''
      }
    },

    async handleLogout() {
      adminWS.disconnect()
      adminStore.token = ''
      adminStore.connected = false
      adminStore.roomID = ''
      this.selectedRoom = null
      this.rooms = []
      sessionStorage.removeItem('admin_token')
      sessionStorage.removeItem('admin_room')
    },

    async switchTab(tab: 'questions' | 'rooms') {
      this.activeTab = tab
      if (tab === 'rooms') {
        // Deselect room so we show the list, not a spectator view
        this.selectedRoom = null
        adminWS.disconnect()
        adminStore.connected = false
        await this.loadRooms()
        this.startRoomsRefresh()
      } else {
        this.stopRoomsRefresh()
      }
    },

    async loadRooms() {
      this.roomsLoading = true
      this.roomsError = ''
      try {
        this.rooms = await adminApiService.listRooms(adminStore.token)
      } catch (err: any) {
        this.roomsError = err.message || 'Failed to load rooms'
      } finally {
        this.roomsLoading = false
      }
    },

    startRoomsRefresh() {
      this.stopRoomsRefresh()
      this.roomsRefreshInterval = setInterval(() => {
        if (this.selectedRoom === null) this.loadRooms()
      }, 3000) as unknown as number
    },

    stopRoomsRefresh() {
      if (this.roomsRefreshInterval !== null) {
        clearInterval(this.roomsRefreshInterval)
        this.roomsRefreshInterval = null
      }
    },

    async selectRoom(roomID: string) {
      this.stopRoomsRefresh()
      this.selectedRoom = roomID
      adminStore.roomID = roomID
      sessionStorage.setItem('admin_room', roomID)
      try {
        await adminWS.connect(adminStore.token, roomID)
      } catch (e) {
        console.error('Spectator connection failed:', e)
      }
    },

    backToRoomList() {
      adminWS.disconnect()
      adminStore.connected = false
      this.selectedRoom = null
      this.loadRooms()
      this.startRoomsRefresh()
    },

    roomStatusLabel(room: RoomSummary): string {
      if (room.is_finished) return 'Finished'
      if (room.is_paused) return 'Paused'
      if (room.is_started) return 'In Progress'
      return 'Waiting'
    },

    roomStatusClass(room: RoomSummary): string {
      if (room.is_finished) return 'status-finished'
      if (room.is_paused) return 'status-paused'
      if (room.is_started) return 'status-active'
      return 'status-waiting'
    },
  },
})
</script>

<template>
  <!-- LOGIN BACKDROP -->
  <div v-if="!isLoggedIn" class="login-backdrop">
    <div class="login-card">
      <div class="login-logo">🛡️ Admin Portal</div>
      <p class="login-subtitle">Enter admin password to manage content &amp; rooms</p>

      <div class="login-field">
        <label for="admin-password">Admin Password</label>
        <input
          id="admin-password"
          v-model="loginPassword"
          type="password"
          placeholder="••••••••"
          class="login-input"
          autofocus
          @keydown.enter="handleLogin"
        />
      </div>

      <p v-if="loginError" class="login-error">{{ loginError }}</p>

      <button
        id="admin-login-btn"
        class="login-btn"
        :disabled="loginLoading || !loginPassword"
        @click="handleLogin"
      >
        <span v-if="loginLoading">Authenticating...</span>
        <span v-else>Log In to Admin</span>
      </button>
    </div>
  </div>

  <!-- AUTHENTICATED ADMIN AREA -->
  <div v-else class="admin-app-layout">
    <!-- Top Navigation Bar -->
    <header class="admin-topbar">
      <div class="topbar-left">
        <span class="topbar-logo">🛡️ Admin Panel</span>
        <nav class="admin-tabs">
          <button
            id="tab-questions"
            class="tab-btn"
            :class="{ active: activeTab === 'questions' }"
            @click="switchTab('questions')"
          >
            📚 Question Bank
          </button>
          <button
            id="tab-rooms"
            class="tab-btn"
            :class="{ active: activeTab === 'rooms' }"
            @click="switchTab('rooms')"
          >
            📡 Live Rooms
          </button>
        </nav>
      </div>

      <div class="topbar-right">
        <span v-if="isWatchingRoom" class="room-indicator">Room: {{ store.roomID }}</span>
        <button id="admin-logout-btn" class="btn-secondary btn-sm" @click="handleLogout">
          Log Out
        </button>
      </div>
    </header>

    <!-- Main View Content -->
    <main class="admin-main-content">
      <!-- QUESTION BANK TAB -->
      <AdminQuestionList v-if="activeTab === 'questions'" />

      <!-- ROOMS TAB: list or spectator -->
      <template v-else-if="activeTab === 'rooms'">

        <!-- SPECTATOR VIEW (a room is selected) -->
        <div v-if="isWatchingRoom" class="spectator-wrapper">
          <div class="spectator-back-bar">
            <button id="back-to-rooms-btn" class="back-btn" @click="backToRoomList">
              ← All Rooms
            </button>
            <span class="spectator-room-label">Watching: <strong>{{ selectedRoom }}</strong></span>
            <span class="conn-pill" :class="{ connected: store.connected }">
              {{ store.connected ? '● Live' : '○ Disconnected' }}
            </span>
          </div>
          <AdminSpectatorView :isEmbedded="true" />
        </div>

        <!-- ROOM LIST -->
        <div v-else class="rooms-tab">
          <div class="rooms-header">
            <h2 class="rooms-title">🏠 Active Rooms</h2>
            <button id="refresh-rooms-btn" class="btn-secondary btn-sm" :disabled="roomsLoading" @click="loadRooms">
              {{ roomsLoading ? 'Refreshing…' : '↺ Refresh' }}
            </button>
          </div>

          <p v-if="roomsError" class="rooms-error">{{ roomsError }}</p>

          <div v-if="rooms.length === 0 && !roomsLoading" class="rooms-empty">
            <div class="empty-icon">🌐</div>
            <p>No rooms found. Rooms appear here once the first player connects.</p>
          </div>

          <div v-else class="rooms-grid">
            <div
              v-for="room in rooms"
              :key="room.room_id"
              class="room-card"
              :class="roomStatusClass(room)"
              :id="`room-card-${room.room_id}`"
              @click="selectRoom(room.room_id)"
            >
              <div class="room-card-top">
                <span class="room-id">{{ room.room_id }}</span>
                <span class="room-status" :class="roomStatusClass(room)">
                  {{ roomStatusLabel(room) }}
                </span>
              </div>

              <div class="room-card-stats">
                <span class="stat">
                  <span class="stat-icon">👥</span>
                  {{ room.player_count }} player{{ room.player_count !== 1 ? 's' : '' }}
                </span>
                <span v-if="room.active_turn" class="stat">
                  <span class="stat-icon">🎯</span>
                  {{ room.active_turn }}'s turn
                </span>
              </div>

              <div v-if="room.players && room.players.length > 0" class="room-players">
                <span
                  v-for="(name, idx) in room.players.slice(0, 6)"
                  :key="idx"
                  class="player-chip"
                >{{ name }}</span>
                <span v-if="room.players.length > 6" class="player-chip more">
                  +{{ room.players.length - 6 }}
                </span>
              </div>

              <div class="room-card-footer">
                <span class="view-room-link">View Room →</span>
              </div>
            </div>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>

<style scoped>
.admin-app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #080f1e;
  color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

/* Login Backdrop */
.login-backdrop {
  min-height: 100vh;
  background: #080f1e;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-card {
  background: linear-gradient(145deg, #0f1e3a 0%, #111827 100%);
  border: 1px solid #1e3a5f;
  border-radius: 16px;
  padding: 2.5rem 2.75rem;
  width: 360px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(96, 165, 250, 0.1);
}
.login-logo {
  font-size: 1.6rem;
  font-weight: 900;
  text-align: center;
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.4rem;
}
.login-subtitle {
  text-align: center;
  color: #64748b;
  font-size: 0.85rem;
  margin-bottom: 1.5rem;
}
.login-field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  margin-bottom: 1rem;
}
.login-field label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.login-input {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 0.55rem 0.85rem;
  color: #f1f5f9;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s;
}
.login-input:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}
.login-error {
  color: #f87171;
  font-size: 0.8rem;
  text-align: center;
  margin-bottom: 0.75rem;
}
.login-btn {
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
}
.login-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}
.login-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

/* Top bar */
.admin-topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1.25rem;
  background: #0c1527;
  border-bottom: 1px solid #1e293b;
  position: sticky;
  top: 0;
  z-index: 50;
}
.topbar-left {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.topbar-logo {
  font-weight: 900;
  font-size: 1.05rem;
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.admin-tabs {
  display: flex;
  gap: 0.5rem;
}
.tab-btn {
  background: transparent;
  border: 1px solid transparent;
  color: #94a3b8;
  padding: 0.4rem 0.85rem;
  border-radius: 8px;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}
.tab-btn:hover {
  color: #f1f5f9;
  background: #1e293b;
}
.tab-btn.active {
  color: #60a5fa;
  background: #0f1e3a;
  border-color: #1e3a5f;
  box-shadow: 0 0 10px rgba(96, 165, 250, 0.15);
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.room-indicator {
  font-size: 0.8rem;
  color: #64748b;
  font-family: monospace;
}
.admin-main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* Buttons */
.btn-secondary {
  background: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  border-radius: 6px;
  cursor: pointer;
  transition: opacity 0.2s;
}
.btn-secondary:hover:not(:disabled) { opacity: 0.8; }
.btn-secondary:disabled { opacity: 0.45; cursor: not-allowed; }
.btn-sm {
  padding: 0.3rem 0.65rem;
  font-size: 0.78rem;
  font-weight: 700;
}

/* Spectator wrapper */
.spectator-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}
.spectator-back-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 1.25rem;
  background: #0c1527;
  border-bottom: 1px solid #1e293b;
  flex-shrink: 0;
}
.back-btn {
  background: transparent;
  border: 1px solid #334155;
  color: #94a3b8;
  padding: 0.3rem 0.7rem;
  border-radius: 6px;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}
.back-btn:hover {
  color: #f1f5f9;
  border-color: #60a5fa;
  background: #0f1e3a;
}
.spectator-room-label {
  font-size: 0.85rem;
  color: #64748b;
}
.spectator-room-label strong {
  color: #f1f5f9;
  font-family: monospace;
}
.conn-pill {
  font-size: 0.72rem;
  font-weight: 700;
  color: #ef4444;
  margin-left: auto;
}
.conn-pill.connected { color: #34d399; }

/* Rooms tab */
.rooms-tab {
  padding: 1.5rem;
  flex: 1;
}
.rooms-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}
.rooms-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 800;
  color: #f1f5f9;
}
.rooms-error {
  color: #f87171;
  font-size: 0.85rem;
  margin-bottom: 1rem;
}
.rooms-empty {
  text-align: center;
  padding: 3rem 1rem;
  color: #475569;
}
.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.75rem;
}
.rooms-empty p {
  font-size: 0.9rem;
  max-width: 360px;
  margin: 0 auto;
}

/* Room cards grid */
.rooms-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1rem;
}

.room-card {
  background: #0c1527;
  border: 1px solid #1e293b;
  border-radius: 12px;
  padding: 1rem 1.1rem;
  cursor: pointer;
  transition: border-color 0.2s, transform 0.15s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}
.room-card:hover {
  border-color: #3b82f6;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}
.room-card.status-active  { border-left: 3px solid #34d399; }
.room-card.status-paused  { border-left: 3px solid #f59e0b; }
.room-card.status-finished { border-left: 3px solid #64748b; }
.room-card.status-waiting  { border-left: 3px solid #60a5fa; }

.room-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.room-id {
  font-weight: 800;
  font-size: 0.95rem;
  color: #f1f5f9;
  font-family: monospace;
}
.room-status {
  font-size: 0.65rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
}
.room-status.status-active   { background: #064e3b; color: #34d399; }
.room-status.status-paused   { background: #78350f; color: #fbbf24; }
.room-status.status-finished { background: #1e293b; color: #64748b; }
.room-status.status-waiting  { background: #1e3a8a; color: #93c5fd; }

.room-card-stats {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}
.stat {
  font-size: 0.78rem;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}
.stat-icon { font-size: 0.85rem; }

.room-players {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
}
.player-chip {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 99px;
  padding: 0.1rem 0.5rem;
  font-size: 0.72rem;
  color: #cbd5e1;
  white-space: nowrap;
}
.player-chip.more {
  background: #0f172a;
  color: #64748b;
}

.room-card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: auto;
  padding-top: 0.25rem;
  border-top: 1px solid #1e293b;
}
.view-room-link {
  font-size: 0.75rem;
  font-weight: 700;
  color: #3b82f6;
  transition: color 0.2s;
}
.room-card:hover .view-room-link {
  color: #60a5fa;
}
</style>
