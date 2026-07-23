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
  private url: string = import.meta.env.VITE_WS_BASE_URL
  private socket: WebSocket | null = null
  private reconnectAttempts: number = 0
  private _connectPromise?: Promise<void>
  private maxBackoff: number = 30000 // 30 s

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
        store.players = payload.players
        store.boardCells = payload.board_cells
        break
      case 'presence':
        if (payload.event === 'joined')
          store.players.push(payload.player)
        else if (payload.event === 'left')
          store.players = store.players.filter(player => payload.player.id !== player.id)
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
        // clear any lingering dice/effect UI
        store.diceRolls = []
        store.lastEffect = ''
        break
      default:
        console.warn('Unhandled message type', type)
        break
    }
  }
}

// Export a singleton instance
export default new WebSocketService()
