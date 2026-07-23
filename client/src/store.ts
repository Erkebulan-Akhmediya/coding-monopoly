import { reactive } from 'vue'

// Simple global reactive store for the client app (Options API friendly)
export const store = reactive({
  // Player state
  playerName: '' as string,
  connected: false as boolean,
  // List of all connected players (including this client)
  // { id, name, token }
  players: [] as Array<{ id: string; name: string; token: string }>,

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
