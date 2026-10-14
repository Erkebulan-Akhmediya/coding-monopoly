import { getBaseHttpUrl } from './serverUrls'

export interface LocalizedText {
  en: string
  ru: string
  kz: string
}

export function emptyLocalizedText(): LocalizedText {
  return { en: '', ru: '', kz: '' }
}

export function localizedFromEn(en: string, existing?: LocalizedText | null): LocalizedText {
  return {
    en,
    ru: existing?.ru ?? '',
    kz: existing?.kz ?? '',
  }
}

export interface OptionInput {
  text: LocalizedText
  is_correct: boolean
}

export interface Option {
  id: string
  text: LocalizedText
  is_correct: boolean
}

export interface ProblemInput {
  type: 'mcq' | 'text'
  difficulty: 'easy' | 'medium' | 'hard'
  topic: string
  title: LocalizedText
  prompt: LocalizedText
  is_published: boolean
  options?: OptionInput[]
  accepted_answers?: LocalizedText[]
}

export interface Problem {
  id: string
  type: 'mcq' | 'text'
  difficulty: 'easy' | 'medium' | 'hard'
  topic: string
  title: LocalizedText
  prompt: LocalizedText
  is_published: boolean
  options?: Option[]
  accepted_answers?: LocalizedText[]
  created_at: string
  updated_at: string
}

export interface ProblemFilters {
  type?: string
  difficulty?: string
  topic?: string
  is_published?: string
}

export interface RoomSummary {
  room_id: string
  topic: string
  board_id?: string
  board_name?: string
  initial_xp: number
  player_count: number
  is_started: boolean
  is_paused: boolean
  is_finished: boolean
  active_turn: string
  players: string[]
}

export interface CreateRoomInput {
  room_id: string
  topic: string
  board_id: string
  initial_xp: number
}

export interface CreatedRoom extends CreateRoomInput {
  board_name: string
}

export interface BoardCell {
  cell_index: number
  name: LocalizedText
  explanation: LocalizedText
  type: string
  params: Record<string, number>
}

export interface BoardSummary {
  id: string
  name: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface Board extends BoardSummary {
  cells: BoardCell[]
}

export interface BoardInput {
  name: string
  cells: BoardCell[]
}

export interface CellParamSpec {
  key: string
  default: number
  min: number
  max: number
}

export interface CellTypeSpec {
  type: string
  start_only: boolean
  params: CellParamSpec[]
}

export interface CellTypeCatalog {
  board_size: number
  cell_types: CellTypeSpec[]
}

export interface ValidationError {
  key: string
  params?: Record<string, number | string>
}

export function validateProblemInput(input: ProblemInput): ValidationError[] {
  const errors: ValidationError[] = []

  if (input.type !== 'mcq' && input.type !== 'text') {
    errors.push({ key: 'invalidType' })
  }
  if (!['easy', 'medium', 'hard'].includes(input.difficulty)) {
    errors.push({ key: 'invalidDifficulty' })
  }
  if (!input.topic || input.topic.trim() === '') {
    errors.push({ key: 'topicRequired' })
  }
  if (!input.title?.en || input.title.en.trim() === '') {
    errors.push({ key: 'titleRequired' })
  }
  if (!input.prompt?.en || input.prompt.en.trim() === '') {
    errors.push({ key: 'promptRequired' })
  }

  if (input.type === 'mcq') {
    const opts = input.options || []
    if (opts.length < 2) {
      errors.push({ key: 'mcqMinOptions' })
    }
    let hasCorrect = false
    opts.forEach((opt, idx) => {
      if (!opt.text?.en || opt.text.en.trim() === '') {
        errors.push({ key: 'optionTextRequired', params: { index: idx + 1 } })
      }
      if (opt.is_correct) {
        hasCorrect = true
      }
    })
    if (!hasCorrect) {
      errors.push({ key: 'mcqCorrectRequired' })
    }
  } else if (input.type === 'text') {
    const answers = input.accepted_answers || []
    if (answers.length === 0) {
      errors.push({ key: 'textMinAnswers' })
    }
    answers.forEach((ans, idx) => {
      if (!ans?.en || ans.en.trim() === '') {
        errors.push({ key: 'answerRequired', params: { index: idx + 1 } })
      }
    })
  }

  return errors
}

/** Defaults for every param of a cell type, used when a cell changes type. */
export function defaultCellParams(spec: CellTypeSpec | undefined): Record<string, number> {
  const params: Record<string, number> = {}
  for (const param of spec?.params || []) {
    params[param.key] = param.default
  }
  return params
}

/**
 * Client-side mirror of the server's board rules (room.ValidateBoard), keyed by
 * cell index so the editor can flag the offending cells. Index -1 holds
 * board-level errors.
 */
export function validateBoardInput(
  input: BoardInput,
  catalog: CellTypeCatalog,
): Array<ValidationError & { cell: number }> {
  const errors: Array<ValidationError & { cell: number }> = []
  if (!input.name || input.name.trim() === '') {
    errors.push({ cell: -1, key: 'boardNameRequired' })
  }
  if (input.cells.length !== catalog.board_size) {
    errors.push({ cell: -1, key: 'boardCellCount', params: { count: catalog.board_size } })
  }
  input.cells.forEach((cell) => {
    const index = cell.cell_index
    const spec = catalog.cell_types.find((s) => s.type === cell.type)
    if (!spec) {
      errors.push({ cell: index, key: 'cellTypeUnknown', params: { index } })
      return
    }
    if (index === 0 && !spec.start_only) {
      errors.push({ cell: index, key: 'startCellDeploy' })
    }
    if (spec.start_only && index !== 0) {
      errors.push({ cell: index, key: 'deployOnlyStart', params: { index } })
    }
    if (!cell.name?.en || cell.name.en.trim() === '') {
      errors.push({ cell: index, key: 'cellNameRequired', params: { index } })
    }
    for (const param of spec.params) {
      const value = cell.params?.[param.key]
      if (typeof value !== 'number' || !Number.isInteger(value) || value < param.min || value > param.max) {
        errors.push({
          cell: index,
          key: 'cellParamRange',
          params: { index, min: param.min, max: param.max },
        })
      }
    }
    if (cell.type === 'teleport' && cell.params?.target_position === index) {
      errors.push({ cell: index, key: 'teleportSelf', params: { index } })
    }
  })
  return errors
}

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    let msg = `HTTP ${res.status} ${res.statusText}`
    try {
      const errJson = await res.json()
      if (errJson && errJson.error) {
        msg = errJson.error
      }
    } catch {
      // ignore json parse error
    }
    throw new Error(msg)
  }
  if (res.status === 204) {
    return {} as T
  }
  return res.json()
}

export const adminApiService = {
  async login(password: string): Promise<{ token: string; expires_at: string }> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    })
    return handleResponse<{ token: string; expires_at: string }>(res)
  },

  async listProblems(token: string, filters?: ProblemFilters): Promise<Problem[]> {
    const baseUrl = getBaseHttpUrl()
    const params = new URLSearchParams()
    if (filters) {
      if (filters.type) params.append('type', filters.type)
      if (filters.difficulty) params.append('difficulty', filters.difficulty)
      if (filters.topic) params.append('topic', filters.topic)
      if (filters.is_published !== undefined && filters.is_published !== '') {
        params.append('is_published', filters.is_published)
      }
    }
    const queryStr = params.toString() ? `?${params.toString()}` : ''
    const res = await fetch(`${baseUrl}/admin/problems${queryStr}`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    const data = await handleResponse<{ problems: Problem[] }>(res)
    return data.problems || []
  },

  async getProblem(token: string, id: string): Promise<Problem> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/problems/${id}`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    return handleResponse<Problem>(res)
  },

  async createProblem(token: string, input: ProblemInput): Promise<Problem> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/problems`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(input),
    })
    return handleResponse<Problem>(res)
  },

  async updateProblem(token: string, id: string, input: ProblemInput): Promise<Problem> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/problems/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(input),
    })
    return handleResponse<Problem>(res)
  },

  async deleteProblem(token: string, id: string): Promise<void> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/problems/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    await handleResponse<void>(res)
  },

  async publishProblem(token: string, id: string): Promise<Problem> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/problems/${id}/publish`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    return handleResponse<Problem>(res)
  },

  async listRooms(token: string): Promise<RoomSummary[]> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/rooms`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleResponse<{ rooms: RoomSummary[] }>(res)
    return data.rooms || []
  },

  async createRoom(token: string, input: CreateRoomInput): Promise<CreatedRoom> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/rooms`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ ...input, room_id: input.room_id.trim(), topic: input.topic.trim() }),
    })
    return handleResponse<CreatedRoom>(res)
  },

  async listTopics(token: string): Promise<string[]> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/topics`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleResponse<{ topics: string[] }>(res)
    return data.topics || []
  },

  async listCellTypes(token: string): Promise<CellTypeCatalog> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/cell-types`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    return handleResponse<CellTypeCatalog>(res)
  },

  async listBoards(token: string): Promise<BoardSummary[]> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleResponse<{ boards: BoardSummary[] }>(res)
    return data.boards || []
  },

  async getBoard(token: string, id: string): Promise<Board> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards/${id}`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    return handleResponse<Board>(res)
  },

  async createBoard(token: string, input: BoardInput): Promise<Board> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(input),
    })
    return handleResponse<Board>(res)
  },

  async updateBoard(token: string, id: string, input: BoardInput): Promise<Board> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(input),
    })
    return handleResponse<Board>(res)
  },

  async deleteBoard(token: string, id: string): Promise<void> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    })
    await handleResponse<void>(res)
  },

  async setDefaultBoard(token: string, id: string): Promise<BoardSummary[]> {
    const baseUrl = getBaseHttpUrl()
    const res = await fetch(`${baseUrl}/admin/boards/${id}/default`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await handleResponse<{ boards: BoardSummary[] }>(res)
    return data.boards || []
  },
}
