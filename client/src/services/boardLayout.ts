import type { BoardCell } from './adminApiService'

/** Emoji shown for each cell type on the admin board views. */
const CELL_ICONS: Record<string, string> = {
  deploy: '🚩',
  code_freeze: '🧊',
  coffee_break: '☕',
  deadline: '🚨',
  xp_gain: '📈',
  xp_loss: '📉',
  mystery: '❓',
  teleport: '🌀',
  skip_next: '⏭️',
  double_xp: '⚡',
  free_pass: '🎟️',
  special_challenge: '🏆',
}

export function cellIcon(type: string): string {
  return CELL_ICONS[type] || '📍'
}

export const CORNER_INDEXES = [0, 8, 16, 24]

/** 1-based row/column of a perimeter cell on the 9×9 grid; mirrors BoardView. */
export function cellGridPosition(idx: number): { gridRowStart: number; gridColumnStart: number } {
  let row = 1
  let col = 1
  if (idx >= 0 && idx <= 8) {
    row = 9
    col = 9 - idx
  } else if (idx <= 16) {
    col = 1
    row = 9 - (idx - 8)
  } else if (idx <= 24) {
    row = 1
    col = 1 + (idx - 16)
  } else {
    col = 9
    row = 1 + (idx - 24)
  }
  return { gridRowStart: row, gridColumnStart: col }
}

/** Short language-neutral summary of a cell's effect values, e.g. "+25" or "→ #8". */
export function cellParamSummary(cell: BoardCell): string {
  const p = cell.params || {}
  switch (cell.type) {
    case 'xp_gain':
    case 'coffee_break':
      return `+${p.amount}`
    case 'xp_loss':
      return `−${p.amount}`
    case 'special_challenge':
      return `+${p.bonus}`
    case 'deploy':
      return `+${p.amount} / +${p.pass_bonus}`
    case 'deadline':
      return `+${p.gain} / −${p.loss}`
    case 'teleport':
      return `→ #${p.target_position}`
    default:
      return ''
  }
}
