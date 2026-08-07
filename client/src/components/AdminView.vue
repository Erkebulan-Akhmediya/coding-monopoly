<script lang="ts">
import { defineComponent } from 'vue'
import { adminStore } from '../adminStore'
import { adminApiService } from '../services/adminApiService'
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
      loginRoomID: 'default' as string,
      loginError: '' as string,
      loginLoading: false as boolean,
      activeTab: 'questions' as 'questions' | 'spectator',
    }
  },

  computed: {
    isLoggedIn(): boolean {
      return !!adminStore.token
    },
  },

  mounted() {
    const savedToken = sessionStorage.getItem('admin_token')
    const savedRoom = sessionStorage.getItem('admin_room') || 'default'
    if (savedToken) {
      adminStore.token = savedToken
      adminStore.roomID = savedRoom
      if (this.activeTab === 'spectator') {
        this.connectSpectator()
      }
    }
  },

  methods: {
    async handleLogin() {
      if (!this.loginPassword) return
      this.loginLoading = true
      this.loginError = ''
      try {
        const data = await adminApiService.login(this.loginPassword)
        const token = data.token
        const roomID = this.loginRoomID || 'default'

        adminStore.token = token
        adminStore.roomID = roomID
        sessionStorage.setItem('admin_token', token)
        sessionStorage.setItem('admin_room', roomID)

        if (this.activeTab === 'spectator') {
          await this.connectSpectator()
        }
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
      sessionStorage.removeItem('admin_token')
    },

    async switchTab(tab: 'questions' | 'spectator') {
      this.activeTab = tab
      if (tab === 'spectator' && adminStore.token && !adminStore.connected) {
        await this.connectSpectator()
      }
    },

    async connectSpectator() {
      if (!adminStore.token) return
      try {
        await adminWS.connect(adminStore.token, adminStore.roomID || 'default')
      } catch (e) {
        console.error('Spectator connection failed:', e)
      }
    },
  },
})
</script>

<template>
  <!-- LOGIN BACKDROP -->
  <div v-if="!isLoggedIn" class="login-backdrop">
    <div class="login-card">
      <div class="login-logo">🛡️ Admin Portal</div>
      <p class="login-subtitle">Enter admin password to manage content & room</p>

      <div class="login-field">
        <label for="admin-room-id">Room ID</label>
        <input
          id="admin-room-id"
          v-model="loginRoomID"
          type="text"
          placeholder="default"
          class="login-input"
          @keydown.enter="handleLogin"
        />
      </div>

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
            id="tab-spectator"
            class="tab-btn"
            :class="{ active: activeTab === 'spectator' }"
            @click="switchTab('spectator')"
          >
            📡 Live Spectator & Control
          </button>
        </nav>
      </div>

      <div class="topbar-right">
        <span class="room-indicator">Room: {{ store.roomID }}</span>
        <button id="admin-logout-btn" class="btn-secondary btn-sm" @click="handleLogout">
          Log Out
        </button>
      </div>
    </header>

    <!-- Main View Content -->
    <main class="admin-main-content">
      <AdminQuestionList v-if="activeTab === 'questions'" />
      <AdminSpectatorView v-else-if="activeTab === 'spectator'" :isEmbedded="true" />
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

.btn-secondary {
  background: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  border-radius: 6px;
  cursor: pointer;
}
.btn-sm {
  padding: 0.3rem 0.65rem;
  font-size: 0.78rem;
  font-weight: 700;
}
</style>
