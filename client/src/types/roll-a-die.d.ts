declare module 'roll-a-die' {
  export interface RollADieOptions {
    element: HTMLElement
    numberOfDice: number
    callback: (values: number[]) => void
    soundVolume?: number
    delay?: number
    values?: number[]
  }

  export default function rollADie(options: RollADieOptions): void
}
