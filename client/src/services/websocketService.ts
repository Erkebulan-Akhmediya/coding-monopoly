/**
 * websocketService – a tiny singleton that manages a single WebSocket connection
 * to the game server, handles automatic reconnect with exponential back‑off,
 * and routes incoming messages into the global reactive `store`.
 *
 * The client uses the Options API, so we expose simple async methods and the
 * `socket` instance. Components import the default export and call
 * `connect()` (or `send()`). The service updates the store directly – this is a
 * lightweight alternative to Vuex that works well with the `store` defined in
 * `src/store.ts`.
 */
import { store } from '../store'
import {
  onDiceAnimationComplete,
  queueTokenMove,
  removeTokenVisual,
  syncTokenVisualPositions,
} from './tokenMovement'

interface Message {
  type: string
  payload: any
}

class WebSocketService {
  private url: string = import.meta.env.VITE_WS_BASE_URL
  private socket: WebSocket | null = null
  private reconnectAttempts: number = 0
  private _connectPromise?: Promise<void>
  private maxBackoff: number = 30000 // 30 s
  /** Keeps dice overlay visible after roll_resolved; turn_started arrives in the same burst. */
  private diceClearTimer: ReturnType<typeof setTimeout> | null = null
  /** Cover roll-a-die's 3s CSS animation plus a short result hold. */
  private static readonly DICE_DISPLAY_MS = 5500

  private clearDiceOverlay(): void {
    store.diceRolls = []
    store.lastEffect = ''
    // Overlay gone — if hops were waiting on dice, let them run.
    onDiceAnimationComplete()
  }

  private scheduleDiceOverlayClear(): void {
    if (this.diceClearTimer !== null) {
      clearTimeout(this.diceClearTimer)
    }
    this.diceClearTimer = setTimeout(() => {
      this.diceClearTimer = null
      this.clearDiceOverlay()
    }, WebSocketService.DICE_DISPLAY_MS)
  }

  /** Connect (or reconnect) to the server */
  async connect(): Promise<void> {
    // If a connection is already open, resolve immediately
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return Promise.resolve();
    }
    // If a connection attempt is already in progress, return its promise
    if (this._connectPromise) {
      return this._connectPromise;
    }
    // Create a new promise that resolves on successful open, rejects on error/close
    this._connectPromise = new Promise<void>((resolve, reject) => {
      this.socket = new WebSocket(this.url);
      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        store.connected = true;
        // ask server for a full state snapshot on (re)connect
        this.send({ type: 'state_request', payload: {} });
        resolve();
        this._connectPromise = undefined;
      };
      this.socket.onclose = () => {
        console.warn('WebSocket closed – attempting reconnection');
        store.connected = false;
        this.scheduleReconnect();
        // If the connection was never opened, reject the pending promise
        if (this._connectPromise) {
          reject(new Error('WebSocket connection closed before opening'));
          this._connectPromise = undefined;
        }
      };
      this.socket.onerror = (err) => {
        console.error('WebSocket error', err);
        // Let onclose handle reconnection and rejection
      };
      this.socket.onmessage = (ev) => this.handleMessage(ev.data);
    });
    return this._connectPromise;
  }

  /** Send a JSON‑serialisable payload */
  send(msg: Message): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(msg))
    } else {
      console.warn('WebSocket not open – message dropped', msg)
    }
  }

  /** Exponential back‑off reconnection */
  private scheduleReconnect(): void {
    this.reconnectAttempts++
    const backoff = Math.min(1000 * 2 ** (this.reconnectAttempts - 1), this.maxBackoff)
    setTimeout(() => this.connect(), backoff)
  }

  /** Dispatch incoming messages into the store */
  private handleMessage(raw: string): void {
    let msg: Message
    try {
      msg = JSON.parse(raw)
    } catch (e) {
      console.error('Invalid JSON from server', raw)
      return
    }
    const { type, payload } = msg
    switch (type) {
      case 'state_sync':
        store.players = (payload.players || []).map((p: any) => ({
          ...p,
          position: p.position ?? 0,
          xp: p.xp ?? 0,
        }))
        store.boardCells = payload.board_cells || []
        syncTokenVisualPositions()
        if (payload.current_turn_player !== undefined) {
          const activeP = store.players.find(p => p.id === payload.current_turn_player)
          store.currentTurnPlayer = activeP ? activeP.name : (payload.current_turn_player || '')
        }
        if (payload.question_active !== undefined) {
          store.questionActive = !!payload.question_active
        }
        if (payload.deadline !== undefined) {
          store.deadline = typeof payload.deadline === 'number'
            ? payload.deadline
            : (payload.deadline ? new Date(payload.deadline).getTime() : 0)
        }
        break
      case 'presence':
        if (payload.event === 'joined') {
          const newPlayer = {
            ...payload.player,
            position: payload.player.position ?? 0,
            xp: payload.player.xp ?? 0,
          }
          const existingIdx = store.players.findIndex(p => p.id === newPlayer.id)
          if (existingIdx >= 0) {
            store.players[existingIdx] = newPlayer
          } else {
            store.players.push(newPlayer)
          }
          store.tokenVisualPositions[newPlayer.id] = newPlayer.position
        } else if (payload.event === 'left') {
          const leftId = payload.player.id
          store.players = store.players.filter(player => leftId !== player.id)
          removeTokenVisual(leftId)
        }
        break
      case 'turn':
        store.currentTurnPlayer = payload.currentPlayer || ''
        store.questionActive = false
        store.activeQuestion = null
        break
      case 'turn_started':
        const activeP = store.players.find(p => p.id === payload.active_player_id)
        store.currentTurnPlayer = activeP ? activeP.name : (payload.active_player_id || '')
        store.questionActive = false
        store.activeQuestion = null
        // Do not clear dice here: server sends turn_started immediately after
        // roll_resolved in the same burst, which would hide the overlay before paint.
        break
      case 'turn_ended':
        store.questionActive = false
        store.activeQuestion = null
        break
      case 'roll_resolved': {
        const player = store.players.find(p => p.id === payload.player_id)
        const oldPos =
          typeof payload.old_position === 'number'
            ? payload.old_position
            : (player?.position ?? 0)
        const newPos =
          typeof payload.new_position === 'number'
            ? payload.new_position
            : (player?.position ?? oldPos)
        if (player) {
          player.position = newPos
          if (typeof payload.player_xp === 'number') {
            player.xp = payload.player_xp
          }
        }
        if (typeof payload.die_roll === 'number') {
          if (payload.roll_index === 1) {
            store.diceRolls = [payload.die_roll]
          } else {
            store.diceRolls.push(payload.die_roll)
          }
        }
        queueTokenMove({
          playerId: payload.player_id,
          from: oldPos,
          to: newPos,
          dieRoll: typeof payload.die_roll === 'number' ? payload.die_roll : 0,
        })
        // Teleport (and similar) may leave logical position beyond dice landing.
        const effectPos = payload.effect?.new_position
        if (typeof effectPos === 'number' && effectPos !== newPos) {
          if (player) {
            player.position = effectPos
          }
          queueTokenMove({
            playerId: payload.player_id,
            from: newPos,
            to: effectPos,
            dieRoll: 0, // single jump to effect target
          })
        }
        if (payload.effect && payload.effect.description) {
          store.lastEffect = payload.effect.description
        } else if (payload.landed_cell && payload.landed_cell.name) {
          store.lastEffect = `Landed on ${payload.landed_cell.name}`
        }
        this.scheduleDiceOverlayClear()
        break
      }
      case 'answer_result': {
        // Positions already advanced via roll_resolved (+ hop queue). Only sync XP.
        const player = store.players.find(p => p.id === payload.player_id)
        if (player && payload.rolls && Array.isArray(payload.rolls)) {
          payload.rolls.forEach((r: any) => {
            if (typeof r.player_xp === 'number') {
              player.xp = r.player_xp
            }
          })
        }
        break
      }
      case 'question_start':
      case 'question_started':
        if (this.diceClearTimer !== null) {
          clearTimeout(this.diceClearTimer)
          this.diceClearTimer = null
        }
        this.clearDiceOverlay()
        store.questionActive = true
        store.deadline = typeof payload.deadline === 'number'
          ? payload.deadline
          : (payload.deadline ? new Date(payload.deadline).getTime() : 0)
        store.activeQuestion = {
          id: payload.problem_id || '',
          type: payload.type || '',
          difficulty: payload.difficulty || '',
          prompt: payload.prompt || '',
          options: payload.options || []
        }
        break
      case 'question_end':
        store.questionActive = false
        store.activeQuestion = null
        store.deadline = 0
        if (this.diceClearTimer !== null) {
          clearTimeout(this.diceClearTimer)
          this.diceClearTimer = null
        }
        this.clearDiceOverlay()
        break
      default:
        console.warn('Unhandled message type', type)
        break
    }
  }
}

// Export a singleton instance
export default new WebSocketService()
