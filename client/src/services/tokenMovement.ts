/**
 * Queues and plays cell-to-cell token hops after the dice overlay settles.
 * Visual positions are separate from logical player.position so placement
 * rules (grid of tokens in a cell) stay unchanged.
 */
import { store } from '../store'

export interface TokenMove {
  playerId: string
  from: number
  to: number
  dieRoll: number
}

const BOARD_SIZE = 32
/** Duration of the hop leap itself. */
const HOP_MS = 260
/** Brief pause on a cell before the next hop. */
const STOP_MS = 110
/** How long the magnified landing cell stays in the hub. */
const LANDED_PREVIEW_MS = 2800

let running = false
let stepTimer: ReturnType<typeof setTimeout> | null = null
let previewTimer: ReturnType<typeof setTimeout> | null = null
/** Dice are mid-animation; hold the move queue until settle/clear. */
let waitingForDice = false

function clearStepTimer() {
  if (stepTimer !== null) {
    clearTimeout(stepTimer)
    stepTimer = null
  }
}

function clearPreviewTimer() {
  if (previewTimer !== null) {
    clearTimeout(previewTimer)
    previewTimer = null
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    stepTimer = setTimeout(() => {
      stepTimer = null
      resolve()
    }, ms)
  })
}

export function clearLandedCellPreview(): void {
  clearPreviewTimer()
  store.landedCellIndex = null
  store.landedCellPlayerId = ''
}

function showLandedCellPreview(cellIndex: number, playerId: string): void {
  clearPreviewTimer()
  store.landedCellIndex = cellIndex
  store.landedCellPlayerId = playerId
  previewTimer = setTimeout(() => {
    previewTimer = null
    store.landedCellIndex = null
    store.landedCellPlayerId = ''
  }, LANDED_PREVIEW_MS)
}

/** Ensure every known player has a visual cell (used after state sync). */
export function syncTokenVisualPositions(): void {
  const next: Record<string, number> = {}
  for (const p of store.players) {
    next[p.id] = p.position ?? 0
  }
  store.tokenVisualPositions = next
  store.hoppingPlayerId = ''
  store.pendingTokenMoves = []
  waitingForDice = false
  running = false
  clearStepTimer()
  clearLandedCellPreview()
}

/** Drop a player from visual maps (presence left). */
export function removeTokenVisual(playerId: string): void {
  delete store.tokenVisualPositions[playerId]
  store.pendingTokenMoves = store.pendingTokenMoves.filter((m) => m.playerId !== playerId)
  if (store.hoppingPlayerId === playerId) {
    store.hoppingPlayerId = ''
  }
  if (store.landedCellPlayerId === playerId) {
    clearLandedCellPreview()
  }
}

/**
 * Record a dice move for later animation. Keeps the token on `from` until
 * hops run. Call while dice overlay is (or will be) visible.
 */
export function queueTokenMove(move: TokenMove): void {
  if (!move.playerId) return
  clearLandedCellPreview()
  if (store.tokenVisualPositions[move.playerId] === undefined) {
    store.tokenVisualPositions[move.playerId] = move.from
  } else if (!running && store.pendingTokenMoves.length === 0) {
    // Snap to the move start if we were idle (covers reconnect / missed sync).
    store.tokenVisualPositions[move.playerId] = move.from
  }
  store.pendingTokenMoves.push(move)
  waitingForDice = true
}

/** Dice CSS animation finished (or overlay was cleared) — start hops. */
export function onDiceAnimationComplete(): void {
  waitingForDice = false
  void pumpMoves()
}

function buildPath(from: number, to: number, dieRoll: number): number[] {
  if (from === to) return []

  // dieRoll === 0: effect jump (teleport) — land on target in one hop.
  if (dieRoll <= 0) {
    return [to]
  }

  const path: number[] = []
  for (let i = 1; i <= dieRoll; i++) {
    path.push((from + i) % BOARD_SIZE)
  }
  // Prefer ending on reported `to` if path math drifted (should not).
  if (path.length && path[path.length - 1] !== to) {
    path[path.length - 1] = to
  }
  return path
}

async function animateMove(move: TokenMove): Promise<void> {
  const path = buildPath(move.from, move.to, move.dieRoll)
  if (!path.length) {
    store.tokenVisualPositions[move.playerId] = move.to
    return
  }

  for (const cell of path) {
    store.hoppingPlayerId = move.playerId
    store.tokenVisualPositions[move.playerId] = cell
    await sleep(HOP_MS)
    store.hoppingPlayerId = ''
    await sleep(STOP_MS)
  }
}

async function pumpMoves(): Promise<void> {
  if (running || waitingForDice) return
  if (!store.pendingTokenMoves.length) return

  running = true
  let lastMove: TokenMove | null = null
  while (store.pendingTokenMoves.length && !waitingForDice) {
    const move = store.pendingTokenMoves.shift()!
    lastMove = move
    await animateMove(move)
  }
  running = false

  // Another roll may have arrived while we finished.
  if (store.pendingTokenMoves.length && !waitingForDice) {
    void pumpMoves()
    return
  }

  // Final resting cell after this hop wave — show a large center copy.
  if (lastMove && !waitingForDice) {
    showLandedCellPreview(lastMove.to, lastMove.playerId)
  }
}

export function getTokenVisualPosition(playerId: string, fallback: number): number {
  const visual = store.tokenVisualPositions[playerId]
  return visual === undefined ? fallback : visual
}
