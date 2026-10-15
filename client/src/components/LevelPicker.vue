<script lang="ts">
import { defineComponent } from 'vue'
import websocketService from '../services/websocketService'

export default defineComponent({
  name: 'LevelPicker',
  methods: {
    selectDifficulty(level: 'easy' | 'medium' | 'hard') {
      websocketService.send({
        type: 'choose_level',
        payload: { difficulty: level }
      })
    }
  }
})
</script>

<template>
  <div class="level-picker">
    <h3 class="picker-title">{{ $t('levelPicker.title') }}</h3>
    <p class="picker-desc">
      {{ $t('levelPicker.description') }}
    </p>
    
    <div class="difficulty-cards">
      <button class="diff-card easy" @click="selectDifficulty('easy')">
        <div class="card-glow"></div>
        <div class="diff-icon">🟢</div>
        <span class="diff-label">{{ $t('levelPicker.easy') }}</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ {{ $t('levelPicker.timeLimit', { seconds: 30 }) }}</span>
          <span class="reward">🎲 {{ $t('levelPicker.rolls', 1) }}</span>
        </div>
      </button>

      <button class="diff-card medium" @click="selectDifficulty('medium')">
        <div class="card-glow"></div>
        <div class="diff-icon">🟡</div>
        <span class="diff-label">{{ $t('levelPicker.medium') }}</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ {{ $t('levelPicker.timeLimit', { seconds: 45 }) }}</span>
          <span class="reward">🎲 {{ $t('levelPicker.rolls', 2) }}</span>
        </div>
      </button>

      <button class="diff-card hard" @click="selectDifficulty('hard')">
        <div class="card-glow"></div>
        <div class="diff-icon">🔴</div>
        <span class="diff-label">{{ $t('levelPicker.hard') }}</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ {{ $t('levelPicker.timeLimit', { seconds: 60 }) }}</span>
          <span class="reward">🎲 {{ $t('levelPicker.rolls', 3) }}</span>
        </div>
      </button>
    </div>
  </div>
</template>

<style scoped>
.level-picker {
  background: color-mix(in srgb, var(--c-slate-800) 75%, transparent);
  backdrop-filter: blur(12px);
  border: 1px solid color-mix(in srgb, var(--c-bright) 10%, transparent);
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  box-shadow: 0 8px 32px 0 color-mix(in srgb, var(--c-shadow) 40%, transparent);
  margin: 1rem 0;
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.picker-title {
  font-size: 1.1rem;
  letter-spacing: 2px;
  background: linear-gradient(90deg, var(--c-emerald-400), var(--c-amber-400), var(--c-red-500));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-weight: 800;
}

.picker-desc {
  font-size: 0.75rem;
  color: var(--c-slate-400);
  margin-bottom: 1.5rem;
  line-height: 1.4;
}

.difficulty-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.diff-card {
  position: relative;
  background: color-mix(in srgb, var(--c-slate-900) 60%, transparent);
  border: 1px solid color-mix(in srgb, var(--c-bright) 5%, transparent);
  border-radius: 10px;
  padding: 1.25rem 0.5rem;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  color: var(--c-slate-50);
}

.diff-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at center, var(--card-accent-alpha) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}

.diff-card:hover {
  transform: translateY(-4px) scale(1.02);
  border-color: var(--card-accent);
  box-shadow: 0 10px 20px -5px var(--card-shadow);
}

.diff-card:hover::before {
  opacity: 0.6;
}

.diff-card:active {
  transform: translateY(-1px);
}

.card-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--card-accent);
}

.diff-card.easy {
  --card-accent: var(--c-emerald-500);
  --card-accent-alpha: color-mix(in srgb, var(--c-emerald-500) 15%, transparent);
  --card-shadow: color-mix(in srgb, var(--c-emerald-500) 30%, transparent);
}

.diff-card.medium {
  --card-accent: var(--c-amber-400);
  --card-accent-alpha: color-mix(in srgb, var(--c-amber-400) 15%, transparent);
  --card-shadow: color-mix(in srgb, var(--c-amber-400) 30%, transparent);
}

.diff-card.hard {
  --card-accent: var(--c-red-500);
  --card-accent-alpha: color-mix(in srgb, var(--c-red-500) 15%, transparent);
  --card-shadow: color-mix(in srgb, var(--c-red-500) 30%, transparent);
}

.diff-icon {
  font-size: 1.5rem;
  filter: drop-shadow(0 2px 4px color-mix(in srgb, var(--c-shadow) 30%, transparent));
}

.diff-label {
  font-weight: 800;
  font-size: 0.85rem;
  letter-spacing: 1px;
}

.diff-stats {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.65rem;
  color: var(--c-slate-400);
}

.time-limit {
  font-weight: 600;
}

.reward {
  color: var(--c-amber-400);
  font-weight: 700;
}

@media (max-width: 600px) {
  .difficulty-cards {
    grid-template-columns: 1fr;
  }
}
</style>
