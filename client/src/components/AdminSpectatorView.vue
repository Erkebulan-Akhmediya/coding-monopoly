<script lang="ts">
/**
 * AdminSpectatorView – live read-only board + event feed + game controls.
 *
 * Layout: three-column grid.
 *   Left  – mini read-only board (same 9×9 grid as BoardView, no interactive overlays)
 *   Centre – live scrolling game_event feed
 *   Right  – player list + admin control panel
 *
 * Connection flow:
 *   1. On mount, if a token is already stored in sessionStorage connect immediately.
 *   2. Otherwise show the login panel.
 *
 * Security: ALL writes go through adminWebsocketService. This component never
 * calls any player-action methods. Even if a bug somehow sent choose_level or
 * submit_answer on the admin socket, the server rejects them outright.
 */
import { defineComponent } from 'vue'
import { adminStore, type AdminPlayer, type GameEvent } from '../adminStore'
import adminWS from '../services/adminWebsocketService'
import { getBaseHttpUrl } from '../services/serverUrls'

const CELL_ICONS: Record<string, string> = {
  deploy: '🚩',
  code_freeze: '🧊',
  coffee_break: '☕',
  deadline: '🚨',
  xp_gain: '📈',
  xp_loss: '📉',
  mystery: '❓',
  teleport: '🌀',
  skip_next: '⏭️',
  double_xp: '⚡',
  free_pass: '🎟️',
  special_challenge: '🏆',
}

const EVENT_COLORS: Record<string, string> = {
  turn_started: '#60a5fa',
  turn_ended: '#94a3b8',
  question_started: '#f59e0b',
  answer_result: '#34d399',
  roll_resolved: '#a78bfa',
  admin_action: '#f472b6',
  presence: '#38bdf8',
  error: '#ef4444',
  event: '#e2e8f0',
}

const PLAYER_COLORS = [
  '#f472b6', '#60a5fa', '#34d399', '#fbbf24',
  '#a78bfa', '#f87171', '#38bdf8', '#4ade80',
]

export default defineComponent({
  name: 'AdminSpectatorView',

  props: {
    isEmbedded: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      store: adminStore,
      adminWS,
      loginPassword: '' as string,
      loginError: '' as string,
      loginLoading: false as boolean,
      countdownInterval: null as number | null,
      remaining: 0 as number,
      confirmKick: null as string | null,
    }
  },

  computed: {
    isLoggedIn(): boolean {
      return (adminStore.connected || this.isEmbedded) && !!adminStore.token
    },

    sortedPlayers(): AdminPlayer[] {
      return [...adminStore.players].sort((a, b) => (b.xp ?? 0) - (a.xp ?? 0))
    },

    reversedEvents(): GameEvent[] {
      return [...adminStore.events].reverse()
    },

    showCountdown(): boolean {
      return adminStore.questionActive && adminStore.deadline > 0
    },

    cells(): any[] {
      return adminStore.boardCells
    },

    cornerIndexes(): number[] {
      return [0, 8, 16, 24]
    },
  },

  watch: {
    'store.deadline'() {
      this.updateRemaining()
    },
  },

  mounted() {
    const saved = sessionStorage.getItem('admin_token')
    const savedRoom = sessionStorage.getItem('admin_room') || 'default'
    if (saved) {
      adminStore.token = saved
      adminStore.roomID = savedRoom
      adminWS.connect(saved, savedRoom).catch(() => {
        adminStore.token = ''
      })
    }
    this.countdownInterval = setInterval(() => this.updateRemaining(), 500) as unknown as number
  },

  beforeUnmount() {
    if (this.countdownInterval !== null) clearInterval(this.countdownInterval)
  },

  methods: {
    async handleLogin() {
      if (!this.loginPassword) return
      this.loginLoading = true
      this.loginError = ''
      try {
        const res = await fetch(`${getBaseHttpUrl()}/admin/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ password: this.loginPassword }),
        })
        if (!res.ok) {
          this.loginError = 'Invalid password'
          return
        }
        const data = await res.json()
        const token: string = data.token
        const roomID = adminStore.roomID || 'default'
        adminStore.token = token
        sessionStorage.setItem('admin_token', token)
        sessionStorage.setItem('admin_room', roomID)
        await adminWS.connect(token, roomID)
      } catch {
        this.loginError = 'Connection failed'
      } finally {
        this.loginLoading = false
        this.loginPassword = ''
      }
    },

    handleLogout() {
      adminWS.disconnect()
      adminStore.token = ''
      sessionStorage.removeItem('admin_token')
    },

    startGame() { adminWS.startGame() },
    togglePause() { adminWS.togglePause() },

    requestKick(playerID: string) { this.confirmKick = playerID },
    confirmKickAction() {
      if (this.confirmKick) { adminWS.kickPlayer(this.confirmKick); this.confirmKick = null }
    },
    cancelKick() { this.confirmKick = null },

    skipTurn(playerID?: string) { adminWS.skipTurn(playerID) },
    endGame() { adminWS.endGame() },

    isCorner(idx: number): boolean { return this.cornerIndexes.includes(idx) },

    getCellGridStyle(idx: number) {
      let row = 1, col = 1
      if (idx >= 0 && idx <= 8) { row = 9; col = 9 - idx }
      else if (idx >= 8 && idx <= 16) { col = 1; row = 9 - (idx - 8) }
      else if (idx >= 16 && idx <= 24) { row = 1; col = 1 + (idx - 16) }
      else if (idx >= 24 && idx <= 31) { col = 9; row = 1 + (idx - 24) }
      return { gridRowStart: row, gridColumnStart: col }
    },

    getCellIcon(type: string): string { return CELL_ICONS[type] ?? '📍' },

    getPlayersAtCell(cellIndex: number): AdminPlayer[] {
      return adminStore.players.filter((p) => (p.position ?? 0) === cellIndex)
    },

    playerColor(player: AdminPlayer): string {
      const idx = adminStore.players.findIndex((p) => p.id === player.id)
      return PLAYER_COLORS[idx % PLAYER_COLORS.length]
    },

    eventColor(kind: string): string { return EVENT_COLORS[kind] ?? EVENT_COLORS.event },

    formatTime(ts: number): string {
      return new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    },

    updateRemaining() {
      if (this.showCountdown) {
        this.remaining = Math.ceil(Math.max(0, adminStore.deadline - Date.now()) / 1000)
      } else {
        this.remaining = 0
      }
    },
  },
})
</script>

<template>
  <!-- LOGIN PANEL -->
  <div v-if="!isLoggedIn" class="login-backdrop">
    <div class="login-card">
      <div class="login-logo">🛡️ Admin Spectator</div>
      <p class="login-subtitle">Enter the admin password to connect</p>

      <div class="login-field">
        <label for="admin-password">Password</label>
        <input id="admin-password" v-model="loginPassword" type="password" placeholder="••••••••"
               class="login-input" autofocus @keydown.enter="handleLogin" />
      </div>

      <p v-if="loginError" class="login-error">{{ loginError }}</p>

      <button id="admin-login-btn" class="login-btn" :disabled="loginLoading || !loginPassword" @click="handleLogin">
        <span v-if="loginLoading">Connecting…</span>
        <span v-else>Connect as Admin</span>
      </button>
    </div>
  </div>

  <!-- MAIN ADMIN SPECTATOR VIEW -->
  <div v-else class="admin-view">

    <!-- Top bar -->
    <header class="admin-topbar">
      <div class="topbar-left">
        <span class="topbar-logo">🛡️ Admin Spectator</span>
        <span class="topbar-room">Room: {{ store.roomID }}</span>
        <span class="topbar-conn" :class="{ connected: store.connected }">
          {{ store.connected ? '● Live' : '○ Disconnected' }}
        </span>
      </div>
      <div class="topbar-right">
        <span v-if="showCountdown" class="topbar-countdown">⏳ {{ remaining }}s</span>
        <button id="admin-logout-btn" class="btn-secondary btn-sm" @click="handleLogout">Logout</button>
      </div>
    </header>

    <!-- Three-column layout -->
    <div class="admin-columns">

      <!-- LEFT: read-only board -->
      <section class="col-board">
        <h2 class="section-title">📋 Board</h2>
        <div v-if="cells.length === 0" class="empty-hint">Waiting for the first player to join…</div>
        <div v-else class="mini-board-grid">
          <div
            v-for="(cell, idx) in cells" :key="idx"
            class="mini-cell"
            :class="['cell-' + (cell.type || 'generic'), { corner: isCorner(idx) }]"
            :style="getCellGridStyle(idx)"
          >
            <div class="mini-cell-top">
              <span class="mini-idx">#{{ idx }}</span>
              <span class="mini-icon">{{ getCellIcon(cell.type) }}</span>
            </div>
            <div class="mini-tokens">
              <span
                v-for="p in getPlayersAtCell(idx)" :key="p.id"
                class="mini-token" :style="{ background: playerColor(p) }" :title="p.name"
              >{{ p.name[0] }}</span>
            </div>
          </div>
          <!-- Centre label -->
          <div class="mini-board-center">
            <div class="mini-turn-info" :class="{ active: !!store.currentTurnPlayer }">
              <template v-if="store.currentTurnPlayer">
                🎯 {{ store.currentTurnPlayer }}
                <span v-if="showCountdown" class="mini-timer">{{ remaining }}s</span>
              </template>
              <template v-else>Waiting…</template>
            </div>
          </div>
        </div>
      </section>

      <!-- CENTRE: event feed -->
      <section class="col-feed">
        <h2 class="section-title">📡 Live Event Feed</h2>
        <div class="event-feed" ref="eventFeed">
          <div v-for="evt in reversedEvents" :key="evt.id" class="event-entry">
            <span class="evt-time">{{ formatTime(evt.timestamp) }}</span>
            <span class="evt-dot" :style="{ background: eventColor(evt.kind) }"></span>
            <span class="evt-kind" :style="{ color: eventColor(evt.kind) }">{{ evt.kind }}</span>
            <span class="evt-msg">{{ evt.message }}</span>
          </div>
          <div v-if="store.events.length === 0" class="empty-hint">No events yet…</div>
        </div>
      </section>

      <!-- RIGHT: players + controls -->
      <section class="col-controls">

        <h2 class="section-title">🎮 Players ({{ sortedPlayers.length }})</h2>
        <div class="player-list">
          <div
            v-for="p in sortedPlayers" :key="p.id"
            class="player-row" :class="{ 'is-active': p.name === store.currentTurnPlayer }"
          >
            <span class="p-token" :style="{ background: playerColor(p) }">{{ p.name[0] }}</span>
            <div class="p-info">
              <span class="p-name">{{ p.name }}</span>
              <span class="p-meta">
                Cell #{{ p.position ?? 0 }} · {{ p.xp ?? 0 }} XP
                <span v-if="!p.is_connected" class="p-offline">offline</span>
                <span v-if="p.in_code_freeze" class="p-badge freeze">❄️ CF</span>
                <span v-if="p.skip_next_turn" class="p-badge skip">⏭</span>
                <span v-if="p.double_xp" class="p-badge double">⚡2x</span>
              </span>
            </div>
            <div class="p-actions">
              <button class="btn-danger btn-xs" :id="`kick-${p.id}`" title="Kick player"
                      @click="requestKick(p.id)">⛔</button>
              <button class="btn-warning btn-xs" :id="`skip-${p.id}`" title="Skip this player's turn"
                      :disabled="p.name !== store.currentTurnPlayer"
                      @click="skipTurn(p.id)">⏭</button>
            </div>
          </div>
          <div v-if="sortedPlayers.length === 0" class="empty-hint">No players yet.</div>
        </div>

        <!-- Kick confirm -->
        <div v-if="confirmKick" class="confirm-overlay">
          <div class="confirm-card">
            <p>Kick <strong>{{ sortedPlayers.find(p => p.id === confirmKick)?.name ?? confirmKick }}</strong>?</p>
            <div class="confirm-btns">
              <button class="btn-danger ctrl-btn" id="confirm-kick-yes" @click="confirmKickAction">Yes, kick</button>
              <button class="btn-secondary ctrl-btn" id="confirm-kick-no" @click="cancelKick">Cancel</button>
            </div>
          </div>
        </div>

        <h2 class="section-title mt">🕹️ Game Controls</h2>
        <div class="controls-grid">
          <button
            v-if="!store.isStarted"
            id="admin-start-btn"
            class="ctrl-btn btn-success"
            @click="startGame"
          >▶ Start Game</button>
          <button
            v-else
            id="admin-pause-btn"
            class="ctrl-btn btn-warning"
            @click="togglePause"
          >{{ store.isPaused ? '▶ Resume' : '⏸ Pause' }}</button>
          <button id="admin-skip-turn-btn" class="ctrl-btn btn-info"
                  :disabled="!store.currentTurnPlayer" @click="skipTurn()">
            ⏭ Skip Active Turn
          </button>
          <button id="admin-end-btn" class="ctrl-btn btn-danger" @click="endGame">🏁 End Game</button>
        </div>

        <div class="turn-strip" v-if="store.currentTurnPlayer">
          <span class="turn-label">Active:</span>
          <span class="turn-name">{{ store.currentTurnPlayer }}</span>
          <span v-if="store.questionActive" class="turn-qa">⏱ Question active</span>
          <span v-else class="turn-nq">Awaiting level pick</span>
        </div>
      </section>

    </div>
  </div>
</template>

<style scoped>
.admin-view {
  display: flex; flex-direction: column; min-height: 100vh;
  background: #080f1e; color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

/* ---- Login ---- */
.login-backdrop {
  min-height: 100vh; background: #080f1e;
  display: flex; align-items: center; justify-content: center;
}
.login-card {
  background: linear-gradient(145deg, #0f1e3a 0%, #111827 100%);
  border: 1px solid #1e3a5f; border-radius: 16px;
  padding: 2.5rem 2.75rem; width: 360px;
  box-shadow: 0 24px 60px rgba(0,0,0,.7), 0 0 0 1px rgba(96,165,250,.1);
}
.login-logo {
  font-size: 1.6rem; font-weight: 900; text-align: center;
  background: linear-gradient(90deg,#60a5fa,#a78bfa);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  margin-bottom: .4rem;
}
.login-subtitle { text-align: center; color: #64748b; font-size: .85rem; margin-bottom: 1.5rem; }
.login-field { display: flex; flex-direction: column; gap: .3rem; margin-bottom: 1rem; }
.login-field label { font-size: .75rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .5px; }
.login-input {
  background: #0f172a; border: 1px solid #334155; border-radius: 8px;
  padding: .55rem .85rem; color: #f1f5f9; font-size: .9rem; outline: none;
  transition: border-color .2s;
}
.login-input:focus { border-color: #60a5fa; box-shadow: 0 0 0 2px rgba(96,165,250,.2); }
.login-error { color: #f87171; font-size: .8rem; text-align: center; margin-bottom: .75rem; }
.login-btn {
  width: 100%; padding: .65rem; border-radius: 8px; border: none; cursor: pointer;
  font-weight: 700; font-size: .95rem;
  background: linear-gradient(90deg,#2563eb,#7c3aed); color: white;
  transition: opacity .2s, transform .15s;
}
.login-btn:hover:not(:disabled) { opacity: .9; transform: translateY(-1px); }
.login-btn:disabled { opacity: .45; cursor: not-allowed; }

/* ---- Top bar ---- */
.admin-topbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: .65rem 1.25rem; background: #0c1527; border-bottom: 1px solid #1e293b;
  position: sticky; top: 0; z-index: 50;
}
.topbar-left  { display: flex; align-items: center; gap: 1rem; }
.topbar-right { display: flex; align-items: center; gap: .75rem; }
.topbar-logo {
  font-weight: 900; font-size: 1rem;
  background: linear-gradient(90deg,#60a5fa,#a78bfa);
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}
.topbar-room { font-size: .8rem; color: #64748b; font-family: monospace; }
.topbar-conn { font-size: .75rem; font-weight: 700; color: #ef4444; }
.topbar-conn.connected { color: #34d399; }
.topbar-countdown { font-weight: 700; font-size: .85rem; color: #fbbf24; animation: pulse-amber 1s infinite; }
@keyframes pulse-amber { 0%,100%{opacity:1} 50%{opacity:.6} }

/* ---- Columns ---- */
.admin-columns {
  display: grid;
  grid-template-columns: minmax(0,3fr) minmax(0,2.2fr) minmax(0,2fr);
  flex: 1; overflow: hidden;
}
.col-board, .col-feed, .col-controls {
  padding: 1rem 1.1rem; overflow-y: auto;
  height: calc(100vh - 49px); box-sizing: border-box;
}
.col-board { border-right: 1px solid #1e293b; }
.col-feed  { border-right: 1px solid #1e293b; }

.section-title {
  font-size: .7rem; font-weight: 800; letter-spacing: 1.5px; text-transform: uppercase;
  color: #64748b; margin: 0 0 .75rem 0;
}
.section-title.mt { margin-top: 1.25rem; }

/* ---- Mini board ---- */
.mini-board-grid {
  display: grid;
  grid-template-columns: repeat(9, minmax(0,1fr));
  grid-template-rows: repeat(9, minmax(0,1fr));
  gap: 2px; aspect-ratio: 1;
  background: #0f172a; border-radius: 8px;
  border: 1px solid #1e293b; padding: 4px;
}
.mini-cell {
  position: relative; background: #1e293b; border-radius: 3px; border: 1px solid #334155;
  display: flex; flex-direction: column; justify-content: space-between;
  padding: 1px 2px; overflow: hidden; box-sizing: border-box;
}
.mini-cell-top { display: flex; justify-content: space-between; align-items: center; }
.mini-idx  { font-size: .38rem; color: #475569; font-weight: 700; line-height: 1; }
.mini-icon { font-size: .55rem; line-height: 1; }
.mini-tokens { display: flex; flex-wrap: wrap; gap: 1px; justify-content: center; }
.mini-token {
  display: inline-flex; align-items: center; justify-content: center;
  width: 10px; height: 10px; border-radius: 50%;
  font-size: .35rem; font-weight: 800; color: #fff;
  text-shadow: 0 0 2px rgba(0,0,0,.8); line-height: 1;
}
.mini-cell.corner            { background: #0f172a; border-width: 2px; }
.mini-cell.cell-deploy       { background: #065f46; border-color: #10b981; }
.mini-cell.cell-code_freeze  { background: #1e3a8a; border-color: #3b82f6; }
.mini-cell.cell-coffee_break { background: #78350f; border-color: #f59e0b; }
.mini-cell.cell-deadline     { background: #7f1d1d; border-color: #ef4444; }
.mini-cell.cell-xp_gain          { border-top: 2px solid #10b981; }
.mini-cell.cell-xp_loss          { border-top: 2px solid #ef4444; }
.mini-cell.cell-mystery           { border-top: 2px solid #818cf8; }
.mini-cell.cell-teleport          { border-top: 2px solid #8b5cf6; }
.mini-cell.cell-skip_next         { border-top: 2px solid #f97316; }
.mini-cell.cell-double_xp         { border-top: 2px solid #eab308; }
.mini-cell.cell-free_pass         { border-top: 2px solid #06b6d4; }
.mini-cell.cell-special_challenge { border-top: 2px solid #ec4899; }
.mini-board-center {
  grid-row: 2/span 7; grid-column: 2/span 7;
  display: flex; align-items: center; justify-content: center;
  background: rgba(8,15,30,.75); border-radius: 4px;
}
.mini-turn-info { text-align: center; font-size: .75rem; font-weight: 700; color: #64748b; padding: .5rem; }
.mini-turn-info.active { color: #60a5fa; }
.mini-timer {
  display: inline-block; margin-left: .4rem;
  background: #ef4444; color: white;
  border-radius: 8px; padding: 0 .35rem; font-size: .68rem;
}

/* ---- Event feed ---- */
.event-feed {
  display: flex; flex-direction: column; gap: .3rem;
  overflow-y: auto; max-height: calc(100vh - 120px); padding-right: .25rem;
}
.event-entry {
  display: grid; grid-template-columns: 72px 8px 90px 1fr;
  gap: .35rem; align-items: start; font-size: .78rem; line-height: 1.4;
  padding: .3rem .45rem; border-radius: 6px;
  background: #0c1527; animation: fadeSlide .25s ease; transition: background .2s;
}
.event-entry:hover { background: #111e35; }
@keyframes fadeSlide { from{opacity:0;transform:translateY(-4px)} to{opacity:1;transform:none} }
.evt-time  { color: #475569; font-size: .65rem; font-family: monospace; white-space: nowrap; padding-top: 1px; }
.evt-dot   { width: 8px; height: 8px; border-radius: 50%; margin-top: 4px; flex-shrink: 0; }
.evt-kind  { font-size: .62rem; font-weight: 700; text-transform: uppercase; letter-spacing: .4px; white-space: nowrap; padding-top: 1px; }
.evt-msg   { color: #cbd5e1; word-break: break-word; }

/* ---- Player list ---- */
.player-list { display: flex; flex-direction: column; gap: .4rem; margin-bottom: .75rem; }
.player-row {
  display: flex; align-items: center; gap: .5rem;
  padding: .4rem .6rem; background: #0c1527; border-radius: 8px;
  border: 1px solid #1e293b; transition: border-color .2s, background .2s;
}
.player-row.is-active { border-color: #60a5fa; background: #0f1e3a; box-shadow: 0 0 8px rgba(96,165,250,.2); }
.p-token {
  width: 28px; height: 28px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  font-weight: 800; font-size: .75rem; color: #fff;
  text-shadow: 0 1px 3px rgba(0,0,0,.5); flex-shrink: 0;
}
.p-info { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.p-name { font-weight: 700; font-size: .85rem; color: #f1f5f9; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.p-meta { font-size: .68rem; color: #64748b; display: flex; flex-wrap: wrap; gap: .3rem; align-items: center; }
.p-offline { color: #ef4444; font-weight: 700; }
.p-badge   { padding: 1px 4px; border-radius: 4px; font-size: .6rem; font-weight: 700; }
.p-badge.freeze { background: #1e3a8a; color: #93c5fd; }
.p-badge.skip   { background: #431407; color: #fb923c; }
.p-badge.double { background: #713f12; color: #fde047; }
.p-actions { display: flex; gap: .25rem; flex-shrink: 0; }

/* ---- Kick confirm ---- */
.confirm-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,.65);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.confirm-card {
  background: #0f172a; border: 1px solid #334155; border-radius: 12px;
  padding: 1.5rem 2rem; text-align: center;
  box-shadow: 0 20px 50px rgba(0,0,0,.7);
}
.confirm-card p { margin: 0 0 1rem 0; font-size: .95rem; color: #e2e8f0; }
.confirm-btns { display: flex; gap: .75rem; justify-content: center; }

/* ---- Game controls ---- */
.controls-grid {
  display: grid; grid-template-columns: 1fr 1fr; gap: .5rem; margin-bottom: .75rem;
}
.ctrl-btn {
  padding: .55rem .4rem; border-radius: 8px; border: none; cursor: pointer;
  font-weight: 700; font-size: .8rem; transition: opacity .2s, transform .15s;
}
.ctrl-btn:hover:not(:disabled) { opacity: .85; transform: translateY(-1px); }
.ctrl-btn:disabled { opacity: .4; cursor: not-allowed; }

.turn-strip {
  display: flex; align-items: center; gap: .5rem; flex-wrap: wrap;
  padding: .45rem .7rem; background: #0c1527;
  border: 1px solid #1e293b; border-radius: 8px; font-size: .8rem;
}
.turn-label { color: #64748b; }
.turn-name  { font-weight: 700; color: #60a5fa; }
.turn-qa    { color: #fbbf24; font-size: .72rem; font-weight: 700; }
.turn-nq    { color: #475569; font-size: .72rem; }

/* ---- Button colours ---- */
.btn-success  { background: linear-gradient(90deg,#065f46,#10b981); color: white; }
.btn-warning  { background: linear-gradient(90deg,#92400e,#f59e0b); color: white; }
.btn-info     { background: linear-gradient(90deg,#1e3a8a,#3b82f6); color: white; }
.btn-danger   { background: linear-gradient(90deg,#7f1d1d,#ef4444); color: white; }
.btn-secondary { background: #1e293b; border: 1px solid #334155; color: #e2e8f0; cursor: pointer; }
.btn-sm  { padding: .3rem .65rem; border-radius: 6px; font-size: .78rem; font-weight: 700; border: none; cursor: pointer; transition: opacity .2s; }
.btn-xs  { padding: .2rem .4rem; border-radius: 5px; font-size: .7rem; font-weight: 700; border: none; cursor: pointer; transition: opacity .2s; line-height: 1.2; }
.btn-secondary:hover { opacity: .85; }

.empty-hint { color: #475569; font-size: .8rem; text-align: center; padding: 1rem 0; font-style: italic; }

@media (max-width: 1100px) {
  .admin-columns { grid-template-columns: 1fr 1fr; }
  .col-board { grid-column: 1/-1; height: auto; border-right: none; border-bottom: 1px solid #1e293b; }
  .col-feed, .col-controls { height: calc(50vh - 49px); }
}
@media (max-width: 680px) {
  .admin-columns { grid-template-columns: 1fr; }
  .col-feed, .col-controls { height: auto; max-height: 50vh; }
}
</style>
