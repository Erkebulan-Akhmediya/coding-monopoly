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

// Simple global reactive store for the client app (Options API friendly)
export const store = reactive({
  // Player state
  playerName: '' as string,
  connected: false as boolean,
  // List of all connected players (including this client)
  players: [] as Player[],

  // Board cells state received from server
  boardCells: [] as any[], // each cell object as defined by the server

  // Turn management
  currentTurnPlayer: '' as string, // name of player whose turn it is
  questionActive: false as boolean,
  // deadline timestamp in ms (UTC) for the active question countdown
  deadline: 0 as number,
  diceRolls: [] as number[], // recent dice roll values
  lastEffect: '' as string, // description of latest effect
})

