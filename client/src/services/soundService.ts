/**
 * soundService - Audio playback service for game sound effects.
 *
 * Handles:
 * - /move.mp3: Token movement / landing on a new square.
 * - /good.mp3: Positive / rewarding effects (XP gain, double XP, free pass, bonus, deploy).
 * - /bad.mp3: Negative / penalty effects (XP loss, code freeze, skip turn).
 * - Neutral squares (coffee break, teleport, regular landing) have no sound effect.
 */

const SOUND_PATHS = {
  move: '/move.mp3',
  good: '/good.mp3',
  bad: '/bad.mp3',
} as const

export type SoundEffect = keyof typeof SOUND_PATHS
export type EffectSoundCategory = 'good' | 'bad' | 'neutral'

// Preload audio files so they are cached in the browser
if (typeof window !== 'undefined') {
  try {
    Object.values(SOUND_PATHS).forEach((src) => {
      const audio = new Audio(src)
      audio.preload = 'auto'
    })
  } catch (e) {
    console.debug('Failed to preload sound effects:', e)
  }
}

/**
 * Play one of the sound effects (/move.mp3, /good.mp3, /bad.mp3).
 */
export function playSound(effect: SoundEffect): void {
  try {
    const audio = new Audio(SOUND_PATHS[effect])
    audio.currentTime = 0
    audio.play().catch((err) => {
      // Browser autoplay policy might block playback until user interaction
      console.debug(`Could not play ${effect} sound:`, err)
    })
  } catch (err) {
    console.debug(`Error initializing audio for ${effect}:`, err)
  }
}

/**
 * Helper to play the token move sound.
 */
export function playMoveSound(): void {
  playSound('move')
}

/**
 * Classifies a cell effect or event into 'good', 'bad', or 'neutral'.
 */
export function classifyEffect(
  effect:
    | {
        effect_type?: string
        xp_delta?: number
        description?: string
      }
    | null
    | undefined,
): EffectSoundCategory {
  if (!effect) return 'neutral'

  const type = (effect.effect_type || '').toLowerCase().trim()
  const xpDelta = typeof effect.xp_delta === 'number' ? effect.xp_delta : 0
  const desc = (effect.description || '').toLowerCase()

  // 1. Explicit Neutral squares / effects (no sound)
  if (type === 'coffee_break' || type === 'teleport' || type === 'generic' || !type) {
    return 'neutral'
  }

  // 2. Explicit Bad effects (penalties)
  if (type === 'xp_loss' || type === 'code_freeze' || type === 'skip_next') {
    return 'bad'
  }

  // 3. Explicit Good effects (buffs & bonuses)
  if (
    type === 'xp_gain' ||
    type === 'double_xp' ||
    type === 'free_pass' ||
    type === 'special_challenge' ||
    type === 'deploy'
  ) {
    return 'good'
  }

  // 4. Delta-based evaluation (e.g. deadline or mystery outcomes)
  if (xpDelta > 0) {
    return 'good'
  }
  if (xpDelta < 0) {
    return 'bad'
  }

  // 5. Mystery box sub-effects where xp_delta is 0
  if (type === 'mystery') {
    if (
      desc.includes('double xp') ||
      desc.includes('free pass') ||
      desc.includes('bonus') ||
      desc.includes('gain')
    ) {
      return 'good'
    }
    if (
      desc.includes('loss') ||
      desc.includes('setback') ||
      desc.includes('lost') ||
      desc.includes('penalty')
    ) {
      return 'bad'
    }
    return 'neutral'
  }

  return 'neutral'
}

/**
 * Evaluates the effect and plays either 'good' or 'bad' sound effect.
 * Neutral effects produce no sound.
 */
export function playEffectSound(effect: any): void {
  const category = classifyEffect(effect)
  if (category === 'good') {
    playSound('good')
  } else if (category === 'bad') {
    playSound('bad')
  }
}
