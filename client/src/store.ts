import { reactive } from 'vue'

export interface Player {
  id: string
  name: string
  token?: string
  position: number
  xp: number
  is_connected?: boolean
  in_code_freeze?: boolean
  skip_next_turn?: boolean
  double_xp?: boolean
  free_passes?: number
}

export interface Question {
  id: string
  type: string
  difficulty: string
  prompt: string
  options?: Array<{ id: string; text: string }>
}

export interface PendingTokenMove {
  playerId: string
  from: number
  to: number
  dieRoll: number
}

export interface EffectToast {
  id: number
  effectType: string
  description: string
  xpDelta: number
  cellIndex: number | null
}

export interface PlayerStanding {
  player_id: string
  name: string
  xp: number
  position: number
  rank: number
  is_connected: boolean
}

export interface GameOverSummary {
  winner_id: string
  winner_name: string
  reason: string
  target_xp: number
  standings: PlayerStanding[]
}

// Simple global reactive store for the client app (Options API friendly)
export const store = reactive({
  // Player state
  playerName: '' as string,
  playerId: '' as string,
  roomId: 'default' as string,
  connected: false as boolean,
  // List of all players (including briefly disconnected slots)
  players: [] as Player[],

  // Board cells state received from server
  boardCells: [] as any[],

  // Turn management
  currentTurnPlayer: '' as string,
  questionActive: false as boolean,
  deadline: 0 as number,
  activeQuestion: null as Question | null,
  diceRolls: [] as number[],
  lastEffect: '' as string,
  lastEffectType: '' as string,

  // Effect feedback
  effectToasts: [] as EffectToast[],
  highlightedCellIndex: null as number | null,

  // Win condition / end screen
  targetXP: 500 as number,
  gameOver: null as GameOverSummary | null,

  // Room / game control state from server
  isStarted: false as boolean,
  isPaused: false as boolean,
  joinError: '' as string,

  // Board token animation (visual cell ≠ logical position while hopping)
  tokenVisualPositions: {} as Record<string, number>,
  pendingTokenMoves: [] as PendingTokenMove[],
  hoppingPlayerId: '' as string,
  /** Magnified destination cell shown in the board center after hops finish. */
  landedCellIndex: null as number | null,
  landedCellPlayerId: '' as string,
})
