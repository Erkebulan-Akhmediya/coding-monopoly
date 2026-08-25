/**
 * websocketService – manages a single WebSocket connection to the game server,
 * automatic reconnect with exponential back-off, and routes messages into `store`.
 *
 * On reconnect, the prior player_id (sessionStorage) is sent with join so the
 * server restores position, XP, and turn-order slot — including mid-question
 * resume with the original deadline.
 */
import { store } from '../store'
import type { EffectToast, GameOverSummary, Player } from '../store'
import { t } from '../i18n'
import { getWsBaseUrl } from './serverUrls'
import { playTurnPingSound } from './soundService'
import {
  clearLandedCellPreview,
  checkAnimationFinished,
  onDiceAnimationComplete,
  queueTokenMove,
  syncTokenVisualPositions,
} from './tokenMovement'

interface Message {
  type: string
  payload: any
  error?: string
}

let toastSeq = 1
let highlightClearTimer: ReturnType<typeof setTimeout> | null = null
/** Last `turn_started` active player id seen on this client (avoids ping on reconnect). */
let lastActivePlayerId = ''

function pushEffectFeedback(effect: any, cellIndex: number | null): void {
  const effectType = (effect?.effect_type || 'generic') as string
  const description = (effect?.description || store.lastEffect || t('common.cellEffect')) as string
  const xpDelta = typeof effect?.xp_delta === 'number' ? effect.xp_delta : 0

  store.lastEffect = description
  store.lastEffectType = effectType

  const toast: EffectToast = {
    id: toastSeq++,
    effectType,
    description,
    xpDelta,
    cellIndex,
  }
  store.effectToasts = [...store.effectToasts, toast].slice(-5)

  if (typeof cellIndex === 'number') {
    store.highlightedCellIndex = cellIndex
    if (highlightClearTimer !== null) {
      clearTimeout(highlightClearTimer)
    }
    highlightClearTimer = setTimeout(() => {
      highlightClearTimer = null
      store.highlightedCellIndex = null
    }, 2800)
  }

  setTimeout(() => {
    store.effectToasts = store.effectToasts.filter((t) => t.id !== toast.id)
  }, 4200)
}

class WebSocketService {
  private socket: WebSocket | null = null
  private reconnectAttempts: number = 0
  private _connectPromise?: Promise<void>
  private maxBackoff: number = 30000
  private diceClearTimer: ReturnType<typeof setTimeout> | null = null
  private static readonly FALLBACK_DICE_CLEAR_MS = 20000
  private intentionalClose = false
  private hasJoined = false

  private clearDiceOverlay(): void {
    store.diceRolls = []
    store.lastEffect = ''
    store.lastEffectType = ''
    onDiceAnimationComplete()
    checkAnimationFinished()
  }

  private scheduleDiceOverlayClear(): void {
    if (this.diceClearTimer !== null) {
      clearTimeout(this.diceClearTimer)
    }
    this.diceClearTimer = setTimeout(() => {
      this.diceClearTimer = null
      this.clearDiceOverlay()
    }, WebSocketService.FALLBACK_DICE_CLEAR_MS)
  }

  private persistIdentity(): void {
    if (store.playerId) {
      sessionStorage.setItem('playerId', store.playerId)
    }
    if (store.playerName) {
      sessionStorage.setItem('playerName', store.playerName)
    }
    if (store.roomId) {
      sessionStorage.setItem('roomId', store.roomId)
    }
  }

  private loadIdentity(): void {
    store.playerId = sessionStorage.getItem('playerId') || store.playerId || ''
    store.playerName = sessionStorage.getItem('playerName') || store.playerName || ''
    store.roomId = sessionStorage.getItem('roomId') || store.roomId || ''
    if (store.playerId && store.playerName) {
      this.hasJoined = true
    }
  }

  /** Connect (or reconnect) to the server */
  async connect(): Promise<void> {
    this.intentionalClose = false
    this.loadIdentity()

    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return Promise.resolve()
    }
    if (this._connectPromise) {
      return this._connectPromise
    }

    this._connectPromise = new Promise<void>((resolve, reject) => {
      this.socket = new WebSocket(getWsBaseUrl())
      this.socket.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
        store.connected = true
        resolve()
        this._connectPromise = undefined

        // Auto-rejoin after drop so the server restores the prior slot.
        if (this.hasJoined && store.playerName) {
          this.sendJoin()
        } else {
          this.send({ type: 'state_request', payload: {} })
        }
      }
      this.socket.onclose = () => {
        console.warn('WebSocket closed – attempting reconnection')
        store.connected = false
        if (!this.intentionalClose) {
          this.scheduleReconnect()
        }
        if (this._connectPromise) {
          reject(new Error('WebSocket connection closed before opening'))
          this._connectPromise = undefined
        }
      }
      this.socket.onerror = (err) => {
        console.error('WebSocket error', err)
      }
      this.socket.onmessage = (ev) => this.handleMessage(ev.data)
    })
    return this._connectPromise
  }

  /** Join (or rejoin) the room with optional prior player_id for resume. */
  sendJoin(name?: string, roomId?: string): void {
    if (name) {
      store.playerName = name.trim()
    }
    if (roomId) {
      store.roomId = roomId
    }
    this.loadIdentity()

    const payload: Record<string, string> = {
      name: store.playerName,
      room_id: store.roomId,
    }
    if (store.playerId) {
      payload.player_id = store.playerId
    }
    this.send({ type: 'join', payload })
  }

  send(msg: Message): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(msg))
    } else {
      console.warn('WebSocket not open – message dropped', msg)
    }
  }

  /** Intentionally close the connection (e.g. return to lobby after game over). */
  disconnect(): void {
    this.intentionalClose = true
    this.hasJoined = false
    this.socket?.close()
    this.socket = null
    store.connected = false
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++
    const backoff = Math.min(1000 * 2 ** (this.reconnectAttempts - 1), this.maxBackoff)
    setTimeout(() => this.connect(), backoff)
  }

  private applyPlayers(list: any[]): void {
    store.players = (list || []).map((p: any) => ({
      ...p,
      position: p.position ?? 0,
      xp: p.xp ?? 0,
      is_connected: p.is_connected !== false,
    })) as Player[]
    syncTokenVisualPositions()
  }

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
      case 'joined':
        store.joinError = ''
        if (payload?.player_id) {
          store.playerId = payload.player_id
        }
        if (payload?.name) {
          store.playerName = payload.name
        }
        if (payload?.room_id) {
          store.roomId = payload.room_id
        }
        this.hasJoined = true
        this.persistIdentity()
        break
      case 'state_sync':
        this.applyPlayers(payload.players || [])
        store.boardCells = payload.board_cells || []
        syncTokenVisualPositions()
        if (payload.current_turn_player !== undefined) {
          const activeP = store.players.find((p) => p.id === payload.current_turn_player)
          store.currentTurnPlayer = activeP ? activeP.name : payload.current_turn_player || ''
        }
        if (payload.question_active !== undefined) {
          store.questionActive = !!payload.question_active
        }
        if (payload.deadline !== undefined) {
          store.deadline =
            typeof payload.deadline === 'number'
              ? payload.deadline
              : payload.deadline
                ? new Date(payload.deadline).getTime()
                : 0
        }
        if (typeof payload.target_xp === 'number' && payload.target_xp > 0) {
          store.targetXP = payload.target_xp
        }
        if (payload.game_over) {
          store.gameOver = payload.game_over as GameOverSummary
        }
        if (payload.is_started !== undefined) {
          store.isStarted = !!payload.is_started
        }
        if (payload.is_paused !== undefined) {
          store.isPaused = !!payload.is_paused
        }
        break
      case 'error': {
        const errMsg = msg.error || 'Unknown error'
        if (!store.boardCells.length) {
          store.joinError = errMsg
        }
        break
      }
      case 'presence':
        if (payload.event === 'joined') {
          const newPlayer: Player = {
            ...payload.player,
            position: payload.player.position ?? 0,
            xp: payload.player.xp ?? 0,
            is_connected: true,
          }
          const existingIdx = store.players.findIndex((p) => p.id === newPlayer.id)
          if (existingIdx >= 0) {
            store.players[existingIdx] = {
              ...store.players[existingIdx],
              ...newPlayer,
              is_connected: true,
            }
          } else {
            store.players.push(newPlayer)
          }
          store.tokenVisualPositions[newPlayer.id] = newPlayer.position
        } else if (payload.event === 'left') {
          const leftId = payload.player.id
          const idx = store.players.findIndex((p) => p.id === leftId)
          if (idx >= 0) {
            // Keep slot for reconnect/resume; mark offline instead of removing.
            store.players[idx] = {
              ...store.players[idx],
              is_connected: false,
            }
          }
        }
        break
      case 'turn':
        store.currentTurnPlayer = payload.currentPlayer || ''
        store.questionActive = false
        store.activeQuestion = null
        break
      case 'turn_started': {
        const activeId = payload.active_player_id || ''
        const activeP = store.players.find((p) => p.id === activeId)
        store.currentTurnPlayer = activeP ? activeP.name : activeId
        store.questionActive = false
        store.activeQuestion = null
        if (activeId && activeId === store.playerId) {
          if (!store.isMoveAnimating) {
            store.pendingTurnNotification = false
            if (lastActivePlayerId !== activeId) {
              playTurnPingSound()
            }
          } else {
            store.pendingTurnNotification = true
          }
        } else {
          store.pendingTurnNotification = false
        }
        lastActivePlayerId = activeId
        break
      }
      case 'turn_ended':
        store.questionActive = false
        store.activeQuestion = null
        break
      case 'roll_resolved': {
        const player = store.players.find((p) => p.id === payload.player_id)
        const oldPos =
          typeof payload.old_position === 'number' ? payload.old_position : (player?.position ?? 0)
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
        const effectPos = payload.effect?.new_position
        const hasTeleport = typeof effectPos === 'number' && effectPos !== newPos

        queueTokenMove({
          playerId: payload.player_id,
          from: oldPos,
          to: newPos,
          dieRoll: typeof payload.die_roll === 'number' ? payload.die_roll : 0,
          effect: hasTeleport ? undefined : payload.effect,
        })
        let feedbackCell = newPos
        if (hasTeleport) {
          if (player) {
            player.position = effectPos
          }
          queueTokenMove({
            playerId: payload.player_id,
            from: newPos,
            to: effectPos,
            dieRoll: 0,
            effect: payload.effect,
          })
          feedbackCell = effectPos
        }
        if (payload.effect) {
          pushEffectFeedback(payload.effect, feedbackCell)
        } else if (payload.landed_cell?.name) {
          const cellName =
            typeof payload.landed_cell.name === 'string'
              ? payload.landed_cell.name
              : payload.landed_cell.name.en || ''
          pushEffectFeedback(
            { effect_type: payload.landed_cell.type || 'generic', description: t('board.landedOn', { name: cellName }) },
            feedbackCell,
          )
        }
        this.scheduleDiceOverlayClear()
        break
      }
      case 'answer_result': {
        const player = store.players.find((p) => p.id === payload.player_id)
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
        clearLandedCellPreview()
        store.pendingTurnNotification = false
        store.isMoveAnimating = false
        store.questionActive = true
        // Genuine remaining time: use server deadline as-is (never reset locally).
        store.deadline =
          typeof payload.deadline === 'number'
            ? payload.deadline
            : payload.deadline
              ? new Date(payload.deadline).getTime()
              : 0
        store.activeQuestion = {
          id: payload.problem_id || '',
          type: payload.type || '',
          difficulty: payload.difficulty || '',
          prompt: payload.prompt || { en: '', ru: '', kz: '' },
          options: payload.options || [],
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
        clearLandedCellPreview()
        break
      case 'game_over':
        store.gameOver = payload as GameOverSummary
        store.questionActive = false
        store.activeQuestion = null
        store.currentTurnPlayer = ''
        store.pendingTurnNotification = false
        store.isMoveAnimating = false
        break
      default:
        console.warn('Unhandled message type', type)
        break
    }
  }
}

export default new WebSocketService()
