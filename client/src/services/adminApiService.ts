export interface OptionInput {
  text: string
  is_correct: boolean
}

export interface Option {
  id: string
  text: string
  is_correct: boolean
}

export interface ProblemInput {
  type: 'mcq' | 'text'
  difficulty: 'easy' | 'medium' | 'hard'
  title: string
  prompt: string
  is_published: boolean
  options?: OptionInput[]
  accepted_answers?: string[]
}

export interface Problem {
  id: string
  type: 'mcq' | 'text'
  difficulty: 'easy' | 'medium' | 'hard'
  title: string
  prompt: string
  is_published: boolean
  options?: Option[]
  accepted_answers?: string[]
  created_at: string
  updated_at: string
}

export interface ProblemFilters {
  type?: string
  difficulty?: string
  is_published?: string
}

export function getBaseHttpUrl(): string {
  const wsUrl = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080/ws'
  return wsUrl.replace(/^wss:/, 'https:').replace(/^ws:/, 'http:').replace(/\/ws$/, '')
}

export function validateProblemInput(input: ProblemInput): string[] {
  const errors: string[] = []

  if (input.type !== 'mcq' && input.type !== 'text') {
    errors.push('Type must be mcq or text')
  }
  if (!['easy', 'medium', 'hard'].includes(input.difficulty)) {
    errors.push('Difficulty must be easy, medium, or hard')
  }
  if (!input.title || input.title.trim() === '') {
    errors.push('Title is required')
  }
  if (!input.prompt || input.prompt.trim() === '') {
    errors.push('Prompt is required')
  }

  if (input.type === 'mcq') {
    const opts = input.options || []
    if (opts.length < 2) {
      errors.push('MCQ requires at least two options')
    }
    let hasCorrect = false
    opts.forEach((opt, idx) => {
      if (!opt.text || opt.text.trim() === '') {
        errors.push(`Option ${idx + 1} text is required`)
      }
      if (opt.is_correct) {
        hasCorrect = true
      }
    })
    if (!hasCorrect) {
      errors.push('MCQ requires at least one correct option')
    }
  } else if (input.type === 'text') {
    const answers = input.accepted_answers || []
    if (answers.length === 0) {
      errors.push('Text requires at least one accepted answer')
    }
    answers.forEach((ans, idx) => {
      if (!ans || ans.trim() === '') {
        errors.push(`Accepted answer ${idx + 1} is required`)
      }
    })
  }

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
}
