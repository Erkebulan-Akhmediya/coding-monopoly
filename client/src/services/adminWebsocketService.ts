/**
 * adminWebsocketService – singleton WebSocket connection for the admin
 * spectator view. Connects to /ws/admin?token=...&room_id=... and routes
 * incoming messages into the reactive adminStore.
 *
 * This service is entirely separate from websocketService.ts so that an
 * admin tab and a player tab can coexist without interfering with each other.
 */
import { adminStore } from '../adminStore'
import { getWsBaseUrl } from './serverUrls'

interface Message {
  type: string
  payload: any
}

class AdminWebSocketService {
  private socket: WebSocket | null = null
  private reconnectAttempts: number = 0
  private maxBackoff: number = 30000
  private _token: string = ''
  private _roomID: string = 'default'
  private _connectPromise?: Promise<void>

  /** Connect to the admin WS endpoint. Token and roomID must be set first. */
  async connect(token: string, roomID: string = 'default'): Promise<void> {
    this._token = token
    this._roomID = roomID

    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return Promise.resolve()
    }
    if (this._connectPromise) return this._connectPromise

    // Build the admin WS URL — pass token + room_id as query params so the
    // server can validate before upgrading the connection.
    const adminUrl = this.buildUrl()

    this._connectPromise = new Promise<void>((resolve, reject) => {
      this.socket = new WebSocket(adminUrl)

      this.socket.onopen = () => {
        console.log('[Admin WS] connected')
        this.reconnectAttempts = 0
        adminStore.connected = true
        // Request a full state snapshot immediately after connecting.
        this.send({ type: 'state_request', payload: {} })
        resolve()
        this._connectPromise = undefined
      }

      this.socket.onclose = () => {
        console.warn('[Admin WS] closed – scheduling reconnect')
        adminStore.connected = false
        this.scheduleReconnect()
        if (this._connectPromise) {
          reject(new Error('Admin WS closed before opening'))
          this._connectPromise = undefined
        }
      }

      this.socket.onerror = (err) => {
        console.error('[Admin WS] error', err)
      }

      this.socket.onmessage = (ev) => this.handleMessage(ev.data)
    })

    return this._connectPromise
  }

  /** Close the connection (e.g. on logout). */
  disconnect(): void {
    this.socket?.close()
    this.socket = null
    adminStore.connected = false
  }

  send(msg: Message): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(msg))
    } else {
      console.warn('[Admin WS] not open – message dropped', msg)
    }
  }

  // ---- Admin commands -------------------------------------------------------

  startGame(): void {
    this.send({ type: 'admin_start', payload: {} })
  }

  togglePause(): void {
    this.send({ type: 'admin_pause', payload: {} })
  }

  kickPlayer(playerID: string): void {
    this.send({ type: 'admin_kick', payload: { player_id: playerID } })
  }

  skipTurn(playerID?: string): void {
    this.send({ type: 'admin_skip_turn', payload: { player_id: playerID ?? '' } })
  }

  endGame(winnerID?: string): void {
    this.send({ type: 'admin_end_game', payload: { winner_id: winnerID ?? '' } })
  }

  // ---- Internals ------------------------------------------------------------

  private buildUrl(): string {
    // Replace /ws with /ws/admin
    const base = getWsBaseUrl().replace(/\/ws$/, '') + '/ws/admin'
    return `${base}?token=${encodeURIComponent(this._token)}&room_id=${encodeURIComponent(this._roomID)}`
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++
    const backoff = Math.min(1000 * 2 ** (this.reconnectAttempts - 1), this.maxBackoff)
    setTimeout(() => {
      if (this._token) this.connect(this._token, this._roomID)
    }, backoff)
  }

  private handleMessage(raw: string): void {
    let msg: Message
    try {
      msg = JSON.parse(raw)
    } catch {
      console.error('[Admin WS] invalid JSON', raw)
      return
    }

    const { type, payload } = msg

    switch (type) {
      case 'state_sync':
        adminStore.players = (payload.players || []).map((p: any) => ({
          ...p,
          position: p.position ?? 0,
          xp: p.xp ?? 0,
        }))
        adminStore.boardCells = payload.board_cells || []
        if (payload.current_turn_player !== undefined) {
          const ap = adminStore.players.find((p) => p.id === payload.current_turn_player)
          adminStore.currentTurnPlayer = ap ? ap.name : payload.current_turn_player || ''
        }
        adminStore.questionActive = !!payload.question_active
        if (payload.deadline !== undefined) {
          adminStore.deadline =
            typeof payload.deadline === 'number'
              ? payload.deadline
              : payload.deadline
                ? new Date(payload.deadline).getTime()
                : 0
        }
        break

      case 'presence':
        if (payload.event === 'joined') {
          const np = {
            ...payload.player,
            position: payload.player.position ?? 0,
            xp: payload.player.xp ?? 0,
            is_connected: true,
          }
          const idx = adminStore.players.findIndex((p) => p.id === np.id)
          if (idx >= 0) adminStore.players[idx] = { ...adminStore.players[idx], ...np }
          else adminStore.players.push(np)
          adminStore.appendEvent('presence', `${np.name} joined`)
        } else if (payload.event === 'left') {
          const left = adminStore.players.find((p) => p.id === payload.player.id)
          if (left) {
            left.is_connected = false
          }
          adminStore.appendEvent('presence', `${left?.name ?? payload.player.id} disconnected`)
        }
        break

      case 'turn_started': {
        const ap = adminStore.players.find((p) => p.id === payload.active_player_id)
        adminStore.currentTurnPlayer = ap ? ap.name : payload.active_player_id || ''
        adminStore.questionActive = false
        adminStore.deadline = 0
        adminStore.appendEvent('turn_started', `Turn started: ${adminStore.currentTurnPlayer}`)
        break
      }

      case 'turn_ended':
        adminStore.questionActive = false
        adminStore.deadline = 0
        adminStore.appendEvent('turn_ended', `Turn ended for player id: ${payload.player_id}`)
        break

      case 'question_started':
        adminStore.questionActive = true
        adminStore.deadline =
          typeof payload.deadline === 'number'
            ? payload.deadline
            : payload.deadline
              ? new Date(payload.deadline).getTime()
              : 0
        adminStore.appendEvent(
          'question_started',
          `Question started (${payload.difficulty}) – ${adminStore.currentTurnPlayer}`
        )
        break

      case 'answer_result': {
        const pl = adminStore.players.find((p) => p.id === payload.player_id)
        const outcome = payload.timed_out ? 'timed out' : payload.correct ? '✅ correct' : '❌ incorrect'
        adminStore.appendEvent('answer_result', `${pl?.name ?? payload.player_id}: ${outcome}`)
        if (payload.rolls && Array.isArray(payload.rolls)) {
          payload.rolls.forEach((r: any) => {
            if (pl) {
              if (typeof r.new_position === 'number') pl.position = r.new_position
              if (typeof r.player_xp === 'number') pl.xp = r.player_xp
            }
          })
        }
        break
      }

      case 'roll_resolved': {
        const rp = adminStore.players.find((p) => p.id === payload.player_id)
        if (rp) {
          if (typeof payload.new_position === 'number') rp.position = payload.new_position
          if (typeof payload.player_xp === 'number') rp.xp = payload.player_xp
        }
        const cellName = payload.landed_cell?.name ?? `cell ${payload.new_position}`
        adminStore.appendEvent('roll_resolved', `${rp?.name ?? payload.player_id} rolled ${payload.die_roll} → ${cellName}`)
        break
      }

      case 'game_event':
        // Push admin action and other server-emitted events into the feed.
        adminStore.appendEvent(payload.kind ?? 'event', payload.message ?? JSON.stringify(payload))
        break

      case 'game_over':
        adminStore.questionActive = false
        adminStore.currentTurnPlayer = ''
        adminStore.appendEvent(
          'admin_action',
          `Game over — winner: ${payload?.winner_name ?? payload?.winner_id ?? 'unknown'}`,
        )
        break

      case 'error':
        adminStore.appendEvent('error', `Server error: ${msg.payload ?? (msg as any).error ?? raw}`)
        break

      case 'pong':
        break

      default:
        console.warn('[Admin WS] unhandled message', type, payload)
    }
  }
}

export default new AdminWebSocketService()
