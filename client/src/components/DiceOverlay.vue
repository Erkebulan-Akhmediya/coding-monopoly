<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'

export default defineComponent({
  name: 'DiceOverlay',
  computed: {
    // Array of recent dice roll numbers (e.g., [4, 2])
    diceRolls(): number[] {
      // @ts-ignore – store property added dynamically
      return (store as any).diceRolls || []
    },
    // Description of the most recent effect applied after the roll
    lastEffect(): string {
      // @ts-ignore – store property added dynamically
      return (store as any).lastEffect || ''
    },
    shouldShow(): boolean {
      return this.diceRolls.length > 0 || this.lastEffect !== ''
    },
  },
})
</script>

<template>
  <div v-if="shouldShow" class="dice-overlay">
    <div class="dice-list" v-if="diceRolls.length">
      <span v-for="(v, i) in diceRolls" :key="i" class="dice-value">{{ v }}</span>
    </div>
    <div class="effect" v-if="lastEffect">
      {{ lastEffect }}
    </div>
  </div>
</template>

<style scoped>
.dice-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  margin: auto;
  padding: 0.5rem 1rem;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  border-radius: 4px;
  text-align: center;
  font-size: 0.9rem;
  z-index: 10;
}
.dice-list {
  margin-bottom: 0.3rem;
}
.dice-value {
  display: inline-block;
  width: 1.5rem;
  height: 1.5rem;
  line-height: 1.5rem;
  background: #2196f3;
  border-radius: 3px;
  margin: 0 0.2rem;
}
.effect {
  font-style: italic;
}
</style>
