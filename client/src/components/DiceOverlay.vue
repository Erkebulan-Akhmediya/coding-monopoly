<script lang="ts">
import { defineComponent } from 'vue'
import rollADie from 'roll-a-die'
import { store } from '../store'

/** CSS animation in roll-a-die is fixed at 3s. */
const ROLL_ANIMATION_MS = 3000
/** Brief settle so multi-roll websocket bursts animate once. */
const ROLL_BURST_MS = 80

export default defineComponent({
  name: 'DiceOverlay',
  data() {
    return {
      rolling: false,
      settled: false,
      burstTimer: null as ReturnType<typeof setTimeout> | null,
      settleTimer: null as ReturnType<typeof setTimeout> | null,
      animatedKey: '' as string,
    }
  },
  computed: {
    diceRolls(): number[] {
      return store.diceRolls || []
    },
    lastEffect(): string {
      return store.lastEffect || ''
    },
    shouldShow(): boolean {
      return this.diceRolls.length > 0 || this.lastEffect !== ''
    },
  },
  watch: {
    diceRolls: {
      handler(rolls: number[]) {
        if (!rolls.length) {
          this.resetAnimationState()
          return
        }
        if (this.burstTimer !== null) {
          clearTimeout(this.burstTimer)
        }
        this.burstTimer = setTimeout(() => {
          this.burstTimer = null
          this.startRollAnimation(rolls.slice())
        }, ROLL_BURST_MS)
      },
      deep: true,
    },
  },
  beforeUnmount() {
    this.clearTimers()
  },
  methods: {
    clearTimers() {
      if (this.burstTimer !== null) {
        clearTimeout(this.burstTimer)
        this.burstTimer = null
      }
      if (this.settleTimer !== null) {
        clearTimeout(this.settleTimer)
        this.settleTimer = null
      }
    },
    resetAnimationState() {
      this.clearTimers()
      this.rolling = false
      this.settled = false
      this.animatedKey = ''
      const mount = this.$refs.diceMount as HTMLElement | undefined
      if (mount) {
        mount.innerHTML = ''
      }
    },
    startRollAnimation(values: number[]) {
      const key = values.join(',')
      if (!values.length || key === this.animatedKey) {
        return
      }
      this.animatedKey = key

      this.$nextTick(() => {
        const mount = this.$refs.diceMount as HTMLElement | undefined
        if (!mount) {
          return
        }

        mount.innerHTML = ''
        this.rolling = true
        this.settled = false
        if (this.settleTimer !== null) {
          clearTimeout(this.settleTimer)
          this.settleTimer = null
        }

        try {
          rollADie({
            element: mount,
            numberOfDice: values.length,
            values,
            // Classroom / shared screens: keep the effect visual-only.
            soundVolume: 0,
            // Keep faces until the store clears the overlay (library default is 3s).
            delay: 10_000,
            callback: () => {
              /* fires when dice DOM is built, not when CSS animation ends */
            },
          })
        } catch (err) {
          console.error('Dice animation failed', err)
          this.rolling = false
          this.settled = true
          return
        }

        this.settleTimer = setTimeout(() => {
          this.settleTimer = null
          this.rolling = false
          this.settled = true
        }, ROLL_ANIMATION_MS)
      })
    },
  },
})
</script>

<template>
  <div v-if="shouldShow" class="dice-overlay" aria-live="polite">
    <div class="dice-stage" :class="{ rolling }">
      <div ref="diceMount" class="dice-mount" />
    </div>

    <div v-if="settled && diceRolls.length" class="dice-result">
      Rolled
      <span v-for="(v, i) in diceRolls" :key="i" class="dice-value">{{ v }}</span>
    </div>

    <div v-if="lastEffect" class="effect" :class="{ visible: settled || !diceRolls.length }">
      {{ lastEffect }}
    </div>
  </div>
</template>

<style scoped>
.dice-overlay {
  position: absolute;
  inset: 0;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(2, 6, 23, 0.72);
  color: #f8fafc;
  text-align: center;
  pointer-events: none;
  box-sizing: border-box;
}

.dice-stage {
  position: relative;
  width: min(100%, 280px);
  height: 230px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  overflow: visible;
}

.dice-mount {
  position: relative;
  min-width: 96px;
  min-height: 32px;
  /* roll-a-die places .dice-outer with large translate; give the fall room */
  padding-left: 1.5rem;
}

.dice-result {
  font-size: 0.95rem;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.dice-value {
  display: inline-block;
  min-width: 1.6rem;
  height: 1.6rem;
  line-height: 1.6rem;
  margin: 0 0.2rem;
  padding: 0 0.25rem;
  background: #2563eb;
  border-radius: 4px;
  font-variant-numeric: tabular-nums;
}

.effect {
  max-width: 28rem;
  font-size: 0.9rem;
  font-style: italic;
  color: #cbd5e1;
  opacity: 0;
  transition: opacity 0.25s ease;
}

.effect.visible {
  opacity: 1;
}
</style>
