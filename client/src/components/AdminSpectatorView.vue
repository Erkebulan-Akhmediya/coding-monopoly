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
  turn_started: 'var(--c-blue-400)',
  turn_ended: 'var(--c-slate-400)',
  question_started: 'var(--c-amber-500)',
  answer_result: 'var(--c-emerald-400)',
  roll_resolved: 'var(--c-violet-400)',
  admin_action: 'var(--c-pink-400)',
  presence: 'var(--c-sky-400)',
  error: 'var(--c-red-500)',
  event: 'var(--c-slate-200)',
}

const PLAYER_COLORS = [
  'var(--c-pink-400)', 'var(--c-blue-400)', 'var(--c-emerald-400)', 'var(--c-amber-400)',
  'var(--c-violet-400)', 'var(--c-red-400)', 'var(--c-sky-400)', 'var(--c-green-400)',
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
      <div class="login-logo">🛡️ {{ $t('admin.spectator') }}</div>
      <p class="login-subtitle">{{ $t('admin.spectatorLoginSubtitle') }}</p>

      <div class="login-field">
        <label for="admin-password">{{ $t('admin.passwordShort') }}</label>
        <input id="admin-password" v-model="loginPassword" type="password" placeholder="••••••••"
               class="login-input" autofocus @keydown.enter="handleLogin" />
      </div>

      <p v-if="loginError" class="login-error">{{ loginError }}</p>

      <button id="admin-login-btn" class="login-btn" :disabled="loginLoading || !loginPassword" @click="handleLogin">
        <span v-if="loginLoading">{{ $t('admin.connecting') }}</span>
        <span v-else>{{ $t('admin.connectAsAdmin') }}</span>
      </button>
    </div>
  </div>

  <!-- MAIN ADMIN SPECTATOR VIEW -->
  <div v-else class="admin-view">

    <!-- Top bar -->
    <header class="admin-topbar">
      <div class="topbar-left">
        <span class="topbar-logo">🛡️ {{ $t('admin.spectator') }}</span>
        <span class="topbar-room">{{ $t('admin.roomLabel', { id: store.roomID }) }}</span>
        <span class="topbar-conn" :class="{ connected: store.connected }">
          {{ store.connected ? $t('admin.live') : $t('admin.disconnected') }}
        </span>
      </div>
      <div class="topbar-right">
        <span v-if="showCountdown" class="topbar-countdown">⏳ {{ remaining }}s</span>
        <button id="admin-logout-btn" class="btn-secondary btn-sm" @click="handleLogout">{{ $t('admin.logout') }}</button>
      </div>
    </header>

    <!-- Three-column layout -->
    <div class="admin-columns">

      <!-- LEFT: read-only board -->
      <section class="col-board">
        <h2 class="section-title">📋 {{ $t('admin.board') }}</h2>
        <div v-if="cells.length === 0" class="empty-hint">{{ $t('admin.waitingFirstPlayer') }}</div>
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
              <template v-else>{{ $t('admin.waiting') }}</template>
            </div>
          </div>
        </div>
      </section>

      <!-- CENTRE: event feed -->
      <section class="col-feed">
        <h2 class="section-title">📡 {{ $t('admin.liveEventFeed') }}</h2>
        <div class="event-feed" ref="eventFeed">
          <div v-for="evt in reversedEvents" :key="evt.id" class="event-entry">
            <span class="evt-time">{{ formatTime(evt.timestamp) }}</span>
            <span class="evt-dot" :style="{ background: eventColor(evt.kind) }"></span>
            <span class="evt-kind" :style="{ color: eventColor(evt.kind) }">{{ evt.kind }}</span>
            <span class="evt-msg">{{ evt.message }}</span>
          </div>
          <div v-if="store.events.length === 0" class="empty-hint">{{ $t('admin.noEvents') }}</div>
        </div>
      </section>

      <!-- RIGHT: players + controls -->
      <section class="col-controls">

        <h2 class="section-title">🎮 {{ $t('admin.players', { count: sortedPlayers.length }) }}</h2>
        <div class="player-list">
          <div
            v-for="p in sortedPlayers" :key="p.id"
            class="player-row" :class="{ 'is-active': p.name === store.currentTurnPlayer }"
          >
            <span class="p-token" :style="{ background: playerColor(p) }">{{ p.name[0] }}</span>
            <div class="p-info">
              <span class="p-name">{{ p.name }}</span>
              <span class="p-meta">
                Cell #{{ p.position ?? 0 }} · {{ p.xp ?? 0 }} {{ $t('common.xp') }}
                <span v-if="!p.is_connected" class="p-offline">{{ $t('common.offline') }}</span>
                <span v-if="p.in_code_freeze" class="p-badge freeze">❄️ CF</span>
                <span v-if="p.skip_next_turn" class="p-badge skip">⏭</span>
                <span v-if="p.double_xp" class="p-badge double">⚡2x</span>
              </span>
            </div>
            <div class="p-actions">
              <button class="btn-danger btn-xs" :id="`kick-${p.id}`" :title="$t('admin.kickPlayer')"
                      @click="requestKick(p.id)">⛔</button>
              <button class="btn-warning btn-xs" :id="`skip-${p.id}`" :title="$t('admin.skipPlayerTurn')"
                      :disabled="p.name !== store.currentTurnPlayer"
                      @click="skipTurn(p.id)">⏭</button>
            </div>
          </div>
          <div v-if="sortedPlayers.length === 0" class="empty-hint">{{ $t('admin.noPlayers') }}</div>
        </div>

        <!-- Kick confirm -->
        <div v-if="confirmKick" class="confirm-overlay">
          <div class="confirm-card">
            <p>{{ $t('admin.kickConfirm', { name: sortedPlayers.find(p => p.id === confirmKick)?.name ?? confirmKick }) }}</p>
            <div class="confirm-btns">
              <button class="btn-danger ctrl-btn" id="confirm-kick-yes" @click="confirmKickAction">{{ $t('admin.yesKick') }}</button>
              <button class="btn-secondary ctrl-btn" id="confirm-kick-no" @click="cancelKick">{{ $t('common.cancel') }}</button>
            </div>
          </div>
        </div>

        <h2 class="section-title mt">🕹️ {{ $t('admin.gameControls') }}</h2>
        <div class="controls-grid">
          <button
            v-if="!store.isStarted"
            id="admin-start-btn"
            class="ctrl-btn btn-success"
            @click="startGame"
          >{{ $t('admin.startGame') }}</button>
          <button
            v-else
            id="admin-pause-btn"
            class="ctrl-btn btn-warning"
            @click="togglePause"
          >{{ store.isPaused ? $t('admin.resume') : $t('admin.pause') }}</button>
          <button id="admin-skip-turn-btn" class="ctrl-btn btn-info"
                  :disabled="!store.currentTurnPlayer" @click="skipTurn()">
            {{ $t('admin.skipActiveTurn') }}
          </button>
          <button id="admin-end-btn" class="ctrl-btn btn-danger" @click="endGame">{{ $t('admin.endGame') }}</button>
        </div>

        <div class="turn-strip" v-if="store.currentTurnPlayer">
          <span class="turn-label">{{ $t('admin.active') }}</span>
          <span class="turn-name">{{ store.currentTurnPlayer }}</span>
          <span v-if="store.questionActive" class="turn-qa">{{ $t('admin.questionActive') }}</span>
          <span v-else class="turn-nq">{{ $t('admin.awaitingLevelPick') }}</span>
        </div>
      </section>

    </div>
  </div>
</template>

<style scoped>
.admin-view {
  display: flex; flex-direction: column; min-height: 100vh;
  background: var(--c-navy-950); color: var(--c-slate-200);
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

/* ---- Login ---- */
.login-backdrop {
  min-height: 100vh; background: var(--c-navy-950);
  display: flex; align-items: center; justify-content: center;
}
.login-card {
  background: linear-gradient(145deg, var(--c-navy-800) 0%, var(--c-gray-900) 100%);
  border: 1px solid var(--c-navy-700); border-radius: 16px;
  padding: 2.5rem 2.75rem; width: 360px;
  box-shadow: 0 24px 60px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 0 1px color-mix(in srgb, var(--c-blue-400) 10%, transparent);
}
.login-logo {
  font-size: 1.6rem; font-weight: 900; text-align: center;
  background: linear-gradient(90deg,var(--c-blue-400),var(--c-violet-400));
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  margin-bottom: .4rem;
}
.login-subtitle { text-align: center; color: var(--c-slate-500); font-size: .85rem; margin-bottom: 1.5rem; }
.login-field { display: flex; flex-direction: column; gap: .3rem; margin-bottom: 1rem; }
.login-field label { font-size: .75rem; font-weight: 600; color: var(--c-slate-400); text-transform: uppercase; letter-spacing: .5px; }
.login-input {
  background: var(--c-slate-900); border: 1px solid var(--c-slate-700); border-radius: 8px;
  padding: .55rem .85rem; color: var(--c-slate-100); font-size: .9rem; outline: none;
  transition: border-color .2s;
}
.login-input:focus { border-color: var(--c-blue-400); box-shadow: 0 0 0 2px color-mix(in srgb, var(--c-blue-400) 20%, transparent); }
.login-error { color: var(--c-red-400); font-size: .8rem; text-align: center; margin-bottom: .75rem; }
.login-btn {
  width: 100%; padding: .65rem; border-radius: 8px; border: none; cursor: pointer;
  font-weight: 700; font-size: .95rem;
  background: linear-gradient(90deg,var(--c-solid--blue-600),var(--c-solid--violet-600)); color: var(--c-white);
  transition: opacity .2s, transform .15s;
}
.login-btn:hover:not(:disabled) { opacity: .9; transform: translateY(-1px); }
.login-btn:disabled { opacity: .45; cursor: not-allowed; }

/* ---- Top bar ---- */
.admin-topbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: .65rem 1.25rem; background: var(--c-navy-900); border-bottom: 1px solid var(--c-slate-800);
  position: sticky; top: 0; z-index: 50;
}
.topbar-left  { display: flex; align-items: center; gap: 1rem; }
.topbar-right { display: flex; align-items: center; gap: .75rem; }
.topbar-logo {
  font-weight: 900; font-size: 1rem;
  background: linear-gradient(90deg,var(--c-blue-400),var(--c-violet-400));
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}
.topbar-room { font-size: .8rem; color: var(--c-slate-500); font-family: monospace; }
.topbar-conn { font-size: .75rem; font-weight: 700; color: var(--c-red-500); }
.topbar-conn.connected { color: var(--c-emerald-400); }
.topbar-countdown { font-weight: 700; font-size: .85rem; color: var(--c-amber-400); animation: pulse-amber 1s infinite; }
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
.col-board { border-right: 1px solid var(--c-slate-800); }
.col-feed  { border-right: 1px solid var(--c-slate-800); }

.section-title {
  font-size: .7rem; font-weight: 800; letter-spacing: 1.5px; text-transform: uppercase;
  color: var(--c-slate-500); margin: 0 0 .75rem 0;
}
.section-title.mt { margin-top: 1.25rem; }

/* ---- Mini board ---- */
.mini-board-grid {
  display: grid;
  grid-template-columns: repeat(9, minmax(0,1fr));
  grid-template-rows: repeat(9, minmax(0,1fr));
  gap: 2px; aspect-ratio: 1;
  background: var(--c-slate-900); border-radius: 8px;
  border: 1px solid var(--c-slate-800); padding: 4px;
}
.mini-cell {
  position: relative; background: var(--c-slate-800); border-radius: 3px; border: 1px solid var(--c-slate-700);
  display: flex; flex-direction: column; justify-content: space-between;
  padding: 1px 2px; overflow: hidden; box-sizing: border-box;
}
.mini-cell-top { display: flex; justify-content: space-between; align-items: center; }
.mini-idx  { font-size: .38rem; color: var(--c-slate-600); font-weight: 700; line-height: 1; }
.mini-icon { font-size: .55rem; line-height: 1; }
.mini-tokens { display: flex; flex-wrap: wrap; gap: 1px; justify-content: center; }
.mini-token {
  display: inline-flex; align-items: center; justify-content: center;
  width: 10px; height: 10px; border-radius: 50%;
  font-size: .35rem; font-weight: 800; color: var(--c-white);
  text-shadow: 0 0 2px color-mix(in srgb, var(--c-shadow) 80%, transparent); line-height: 1;
}
.mini-cell.corner            { background: var(--c-slate-900); border-width: 2px; }
.mini-cell.cell-deploy       { background: var(--c-deploy--emerald-800); border-color: var(--c-deploy--emerald-500); }
.mini-cell.cell-code_freeze  { background: var(--c-code_freeze--blue-900); border-color: var(--c-code_freeze--blue-500); }
.mini-cell.cell-coffee_break { background: var(--c-coffee_break--amber-900); border-color: var(--c-coffee_break--amber-500); }
.mini-cell.cell-deadline     { background: var(--c-deadline--red-900); border-color: var(--c-deadline--red-500); }
.mini-cell.cell-xp_gain          { border-top: 2px solid var(--c-xp_gain--emerald-500); }
.mini-cell.cell-xp_loss          { border-top: 2px solid var(--c-xp_loss--red-500); }
.mini-cell.cell-mystery           { border-top: 2px solid var(--c-mystery--indigo-400); }
.mini-cell.cell-teleport          { border-top: 2px solid var(--c-teleport--violet-500); }
.mini-cell.cell-skip_next         { border-top: 2px solid var(--c-skip_next--orange-500); }
.mini-cell.cell-double_xp         { border-top: 2px solid var(--c-double_xp--yellow-500); }
.mini-cell.cell-free_pass         { border-top: 2px solid var(--c-free_pass--cyan-500); }
.mini-cell.cell-special_challenge { border-top: 2px solid var(--c-special_challenge--pink-500); }
.mini-board-center {
  grid-row: 2/span 7; grid-column: 2/span 7;
  display: flex; align-items: center; justify-content: center;
  background: color-mix(in srgb, var(--c-navy-950) 75%, transparent); border-radius: 4px;
}
.mini-turn-info { text-align: center; font-size: .75rem; font-weight: 700; color: var(--c-slate-500); padding: .5rem; }
.mini-turn-info.active { color: var(--c-blue-400); }
.mini-timer {
  display: inline-block; margin-left: .4rem;
  background: var(--c-red-500); color: var(--c-white);
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
  background: var(--c-navy-900); animation: fadeSlide .25s ease; transition: background .2s;
}
.event-entry:hover { background: var(--c-navy-850); }
@keyframes fadeSlide { from{opacity:0;transform:translateY(-4px)} to{opacity:1;transform:none} }
.evt-time  { color: var(--c-slate-600); font-size: .65rem; font-family: monospace; white-space: nowrap; padding-top: 1px; }
.evt-dot   { width: 8px; height: 8px; border-radius: 50%; margin-top: 4px; flex-shrink: 0; }
.evt-kind  { font-size: .62rem; font-weight: 700; text-transform: uppercase; letter-spacing: .4px; white-space: nowrap; padding-top: 1px; }
.evt-msg   { color: var(--c-slate-300); word-break: break-word; }

/* ---- Player list ---- */
.player-list { display: flex; flex-direction: column; gap: .4rem; margin-bottom: .75rem; }
.player-row {
  display: flex; align-items: center; gap: .5rem;
  padding: .4rem .6rem; background: var(--c-navy-900); border-radius: 8px;
  border: 1px solid var(--c-slate-800); transition: border-color .2s, background .2s;
}
.player-row.is-active { border-color: var(--c-blue-400); background: var(--c-navy-800); box-shadow: 0 0 8px color-mix(in srgb, var(--c-blue-400) 20%, transparent); }
.p-token {
  width: 28px; height: 28px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  font-weight: 800; font-size: .75rem; color: var(--c-white);
  text-shadow: 0 1px 3px color-mix(in srgb, var(--c-shadow) 50%, transparent); flex-shrink: 0;
}
.p-info { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.p-name { font-weight: 700; font-size: .85rem; color: var(--c-slate-100); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.p-meta { font-size: .68rem; color: var(--c-slate-500); display: flex; flex-wrap: wrap; gap: .3rem; align-items: center; }
.p-offline { color: var(--c-red-500); font-weight: 700; }
.p-badge   { padding: 1px 4px; border-radius: 4px; font-size: .6rem; font-weight: 700; }
.p-badge.freeze { background: var(--c-blue-900); color: var(--c-blue-300); }
.p-badge.skip   { background: var(--c-orange-950); color: var(--c-orange-400); }
.p-badge.double { background: var(--c-yellow-900); color: var(--c-yellow-300); }
.p-actions { display: flex; gap: .25rem; flex-shrink: 0; }

/* ---- Kick confirm ---- */
.confirm-overlay {
  position: fixed; inset: 0; background: color-mix(in srgb, var(--c-shadow) 65%, transparent);
  display: flex; align-items: center; justify-content: center; z-index: 200;
}
.confirm-card {
  background: var(--c-slate-900); border: 1px solid var(--c-slate-700); border-radius: 12px;
  padding: 1.5rem 2rem; text-align: center;
  box-shadow: 0 20px 50px color-mix(in srgb, var(--c-shadow) 70%, transparent);
}
.confirm-card p { margin: 0 0 1rem 0; font-size: .95rem; color: var(--c-slate-200); }
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
  padding: .45rem .7rem; background: var(--c-navy-900);
  border: 1px solid var(--c-slate-800); border-radius: 8px; font-size: .8rem;
}
.turn-label { color: var(--c-slate-500); }
.turn-name  { font-weight: 700; color: var(--c-blue-400); }
.turn-qa    { color: var(--c-amber-400); font-size: .72rem; font-weight: 700; }
.turn-nq    { color: var(--c-slate-600); font-size: .72rem; }

/* ---- Button colours ---- */
.btn-success  { background: linear-gradient(90deg,var(--c-solid--emerald-800),var(--c-emerald-500)); color: var(--c-white); }
.btn-warning  { background: linear-gradient(90deg,var(--c-solid--amber-800),var(--c-amber-500)); color: var(--c-white); }
.btn-info     { background: linear-gradient(90deg,var(--c-solid--blue-900),var(--c-blue-500)); color: var(--c-white); }
.btn-danger   { background: linear-gradient(90deg,var(--c-solid--red-900),var(--c-red-500)); color: var(--c-white); }
.btn-secondary { background: var(--c-slate-800); border: 1px solid var(--c-slate-700); color: var(--c-slate-200); cursor: pointer; }
.btn-sm  { padding: .3rem .65rem; border-radius: 6px; font-size: .78rem; font-weight: 700; border: none; cursor: pointer; transition: opacity .2s; }
.btn-xs  { padding: .2rem .4rem; border-radius: 5px; font-size: .7rem; font-weight: 700; border: none; cursor: pointer; transition: opacity .2s; line-height: 1.2; }
.btn-secondary:hover { opacity: .85; }

.empty-hint { color: var(--c-slate-600); font-size: .8rem; text-align: center; padding: 1rem 0; font-style: italic; }

@media (max-width: 1100px) {
  .admin-columns { grid-template-columns: 1fr 1fr; }
  .col-board { grid-column: 1/-1; height: auto; border-right: none; border-bottom: 1px solid var(--c-slate-800); }
  .col-feed, .col-controls { height: calc(50vh - 49px); }
}
@media (max-width: 680px) {
  .admin-columns { grid-template-columns: 1fr; }
  .col-feed, .col-controls { height: auto; max-height: 50vh; }
}
</style>
