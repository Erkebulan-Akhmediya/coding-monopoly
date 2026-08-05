import { reactive } from 'vue'

export interface AdminPlayer {
  id: string
  name: string
  position: number
  xp: number
  is_connected?: boolean
  in_code_freeze?: boolean
  skip_next_turn?: boolean
  double_xp?: boolean
  free_passes?: number
}

export interface GameEvent {
  id: number            // monotonically increasing, used as :key
  kind: string          // 'turn_started' | 'turn_ended' | 'answer_result' | 'admin_action' | 'presence' | …
  message: string
  timestamp: number     // Date.now() at the moment we received the event
}

const MAX_EVENTS = 200

let eventCounter = 0

export const adminStore = reactive({
  // Connection state
  connected: false as boolean,
  token: '' as string,
  roomID: 'default' as string,

  // Mirrored game state
  players: [] as AdminPlayer[],
  boardCells: [] as any[],
  currentTurnPlayer: '' as string,
  questionActive: false as boolean,
  deadline: 0 as number,

  // Live event feed
  events: [] as GameEvent[],

  /** Append an entry to the event feed and trim to MAX_EVENTS. */
  appendEvent(kind: string, message: string) {
    this.events.push({ id: eventCounter++, kind, message, timestamp: Date.now() })
    if (this.events.length > MAX_EVENTS) {
      this.events.splice(0, this.events.length - MAX_EVENTS)
    }
  },
})
