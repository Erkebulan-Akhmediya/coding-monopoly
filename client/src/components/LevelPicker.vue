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
    <h3 class="picker-title">CHOOSE YOUR CHALLENGE</h3>
    <p class="picker-desc">
      Select a difficulty level to receive your coding task. Higher difficulty grants more dice rolls on a correct answer!
    </p>
    
    <div class="difficulty-cards">
      <!-- Easy Card -->
      <button class="diff-card easy" @click="selectDifficulty('easy')">
        <div class="card-glow"></div>
        <div class="diff-icon">🟢</div>
        <span class="diff-label">EASY</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ 30s limit</span>
          <span class="reward">🎲 1 Roll</span>
        </div>
      </button>

      <!-- Medium Card -->
      <button class="diff-card medium" @click="selectDifficulty('medium')">
        <div class="card-glow"></div>
        <div class="diff-icon">🟡</div>
        <span class="diff-label">MEDIUM</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ 45s limit</span>
          <span class="reward">🎲 2 Rolls</span>
        </div>
      </button>

      <!-- Hard Card -->
      <button class="diff-card hard" @click="selectDifficulty('hard')">
        <div class="card-glow"></div>
        <div class="diff-icon">🔴</div>
        <span class="diff-label">HARD</span>
        <div class="diff-stats">
          <span class="time-limit">⏱️ 60s limit</span>
          <span class="reward">🎲 3 Rolls</span>
        </div>
      </button>
    </div>
  </div>
</template>

<style scoped>
.level-picker {
  background: rgba(30, 41, 59, 0.75);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.4);
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
  background: linear-gradient(90deg, #34d399, #fbbf24, #ef4444);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-weight: 800;
}

.picker-desc {
  font-size: 0.75rem;
  color: #94a3b8;
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
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  padding: 1.25rem 0.5rem;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  color: #f8fafc;
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

/* Card Glow Effect */
.card-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--card-accent);
}

/* Easy Color Theme */
.diff-card.easy {
  --card-accent: #10b981;
  --card-accent-alpha: rgba(16, 185, 129, 0.15);
  --card-shadow: rgba(16, 185, 129, 0.3);
}

/* Medium Color Theme */
.diff-card.medium {
  --card-accent: #fbbf24;
  --card-accent-alpha: rgba(251, 191, 36, 0.15);
  --card-shadow: rgba(251, 191, 36, 0.3);
}

/* Hard Color Theme */
.diff-card.hard {
  --card-accent: #ef4444;
  --card-accent-alpha: rgba(239, 68, 68, 0.15);
  --card-shadow: rgba(239, 68, 68, 0.3);
}

.diff-icon {
  font-size: 1.5rem;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
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
  color: #94a3b8;
}

.time-limit {
  font-weight: 600;
}

.reward {
  color: #fbbf24;
  font-weight: 700;
}

@media (max-width: 600px) {
  .difficulty-cards {
    grid-template-columns: 1fr;
  }
}
</style>
