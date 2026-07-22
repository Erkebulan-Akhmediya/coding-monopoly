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

interface Message {
  type: string
  payload: any
}

class WebSocketService {
  private url: string = `${window.location.protocol.replace('http', 'ws')}//${window.location.host}/ws`
  private socket: WebSocket | null = null
  private reconnectAttempts: number = 0
  private maxBackoff: number = 30000 // 30 s

  /** Connect (or reconnect) to the server */
  async connect(): Promise<void> {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) return
    this.socket = new WebSocket(this.url)
    this.socket.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
      store.connected = true
    }
    this.socket.onclose = () => {
      console.warn('WebSocket closed – attempting reconnection')
      store.connected = false
      this.scheduleReconnect()
    }
    this.socket.onerror = (err) => {
      console.error('WebSocket error', err)
      // Let onclose handle reconnection
    }
    this.socket.onmessage = (ev) => this.handleMessage(ev.data)
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
      case 'presence':
        // payload: { players: [{id, name, token}], selfId }
        store.players = payload.players
        break
      case 'board_state':
        // payload: { cells: [...] }
        store.boardCells = payload.cells
        break
      case 'turn':
        // payload: { currentPlayer: string }
        store.currentTurnPlayer = payload.currentPlayer
        store.questionActive = false
        break
      case 'question_start':
        // payload: { deadline: number (ms since epoch) }
        store.questionActive = true
        store.deadline = payload.deadline
        break
      case 'question_end':
        store.questionActive = false
        store.deadline = 0
        break
      default:
        console.warn('Unhandled message type', type)
    }
  }
}

// Export a singleton instance
export default new WebSocketService()
