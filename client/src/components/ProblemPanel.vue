<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'
import websocketService from '../services/websocketService'
import McqPanel from './McqPanel.vue'
import TextAnswerPanel from './TextAnswerPanel.vue'

export default defineComponent({
  name: 'ProblemPanel',
  components: {
    McqPanel,
    TextAnswerPanel
  },
  data() {
    return {
      store: store,
      remainingTime: 0,
      timerId: null as number | null,
      submitted: false,
    }
  },
  computed: {
    isMyTurn(): boolean {
      return store.currentTurnPlayer === store.playerName && store.playerName !== ''
    },
    currentPanelComponent(): string {
      const qType = store.activeQuestion?.type
      if (qType === 'mcq') {
        return 'McqPanel'
      } else if (qType === 'text') {
        return 'TextAnswerPanel'
      }
      return ''
    },
    isTimedOut(): boolean {
      return this.remainingTime <= 0
    },
    difficultyClass(): string {
      const difficulty = store.activeQuestion?.difficulty || ''
      return `diff-${difficulty.toLowerCase()}`
    },
    difficultyLabel(): string {
      const difficulty = (store.activeQuestion?.difficulty || 'unknown').toLowerCase()
      const key = `difficulty.${difficulty}`
      return this.$te(key) ? this.$t(key) : difficulty.toUpperCase()
    }
  },
  watch: {
    'store.deadline': {
      immediate: true,
      handler() {
        this.startTimer()
        this.submitted = false
      }
    }
  },
  methods: {
    updateTimer() {
      const now = Date.now()
      const diff = store.deadline - now
      this.remainingTime = Math.max(0, Math.ceil(diff / 1000))
      if (this.remainingTime <= 0) {
        this.stopTimer()
      }
    },
    startTimer() {
      this.stopTimer()
      this.updateTimer()
      this.timerId = setInterval(() => {
        this.updateTimer()
      }, 200) as unknown as number
    },
    stopTimer() {
      if (this.timerId) {
        clearInterval(this.timerId)
        this.timerId = null
      }
    },
    handleSubmit(answer: any) {
      if (this.submitted || this.isTimedOut) return
      this.submitted = true
      websocketService.send({
        type: 'submit_answer',
        payload: {
          problem_id: store.activeQuestion?.id,
          answer: answer
        }
      })
    }
  },
  mounted() {
    this.startTimer()
  },
  beforeUnmount() {
    this.stopTimer()
  }
})
</script>

<template>
  <div class="problem-container">
    <!-- Header: shows countdown timer and difficulty level -->
    <div class="problem-header" :class="difficultyClass">
      <div class="header-left">
        <span class="header-badge">{{ $t('problem.challengeActive') }}</span>
        <span class="difficulty-tag">{{ difficultyLabel }}</span>
      </div>
      <div class="header-timer" :class="{ 'warning-time': remainingTime <= 10 }">
        <span class="timer-icon">⏳</span>
        <span class="timer-val">{{ remainingTime }}s</span>
      </div>
    </div>

    <!-- Active Player: render the exact panel according to question type -->
    <div v-if="isMyTurn" class="problem-active-panel">
      <div v-if="!store.activeQuestion" class="waiting-question">
        <div class="spinner"></div>
        <p>{{ $t('problem.loading') }}</p>
      </div>
      <div v-else class="question-content">
        <!-- Render prompt -->
        <div class="prompt-box">
          <p class="prompt-text">{{ store.activeQuestion.prompt }}</p>
        </div>

        <div v-if="isTimedOut && !submitted" class="timeout-overlay">
          <p class="timeout-msg">⏳ {{ $t('problem.timeUp') }}</p>
        </div>
        <div v-else-if="submitted" class="submitting-overlay">
          <p class="submitting-msg">🚀 {{ $t('problem.answerSubmitted') }}</p>
        </div>

        <!-- Computed dynamic component -->
        <component
          v-if="currentPanelComponent"
          :is="currentPanelComponent"
          :question="store.activeQuestion"
          :disabled="submitted || isTimedOut"
          @submit="handleSubmit"
        />
      </div>
    </div>

    <!-- Redacted Spectator View: shows only metadata, no content -->
    <div v-else class="problem-spectator-panel">
      <div class="spectator-visual">
        <div class="redacted-lock">🔒</div>
        <h4>{{ $t('problem.spectating') }}</h4>
        <p class="spectator-hint">
          {{ $t('problem.spectatorHint', { player: store.currentTurnPlayer }) }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.problem-container {
  background: rgba(30, 41, 59, 0.75);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.4);
  margin: 1rem 0;
  animation: slideIn 0.4s ease-out;
  text-align: left;
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(15px); }
  to { opacity: 1; transform: translateY(0); }
}

.problem-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1.25rem;
  background: rgba(15, 23, 42, 0.6);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-badge {
  font-size: 0.6rem;
  font-weight: 800;
  color: #38bdf8;
  letter-spacing: 1px;
  background: rgba(56, 189, 248, 0.15);
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
}

.difficulty-tag {
  font-size: 0.65rem;
  font-weight: 800;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  letter-spacing: 0.5px;
  color: #fff;
}

.diff-easy .difficulty-tag { background: #10b981; }
.diff-medium .difficulty-tag { background: #fbbf24; color: #1e293b; }
.diff-hard .difficulty-tag { background: #ef4444; }

.header-timer {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
  font-weight: 800;
  background: rgba(255, 255, 255, 0.05);
  padding: 0.25rem 0.6rem;
  border-radius: 20px;
}

.header-timer.warning-time {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.03); }
}

.timer-val {
  font-family: monospace;
}

.prompt-box {
  background: rgba(15, 23, 42, 0.4);
  padding: 1rem 1.25rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.prompt-text {
  margin: 0;
  font-size: 0.85rem;
  line-height: 1.5;
  color: #f1f5f9;
  white-space: pre-wrap;
}

/* Overlays */
.timeout-overlay, .submitting-overlay {
  padding: 0.75rem 1.25rem;
  font-weight: 700;
  font-size: 0.8rem;
  text-align: center;
  background: rgba(15, 23, 42, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.timeout-msg {
  margin: 0;
  color: #ef4444;
}

.submitting-msg {
  margin: 0;
  color: #34d399;
}

/* Spectator Panel */
.problem-spectator-panel {
  padding: 2rem 1.5rem;
  text-align: center;
}

.spectator-visual {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.redacted-lock {
  font-size: 2rem;
  filter: drop-shadow(0 2px 8px rgba(0,0,0,0.5));
  animation: hoverLock 3s ease-in-out infinite;
}

@keyframes hoverLock {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.spectator-visual h4 {
  margin: 0.5rem 0 0 0;
  font-size: 0.9rem;
  letter-spacing: 1px;
  color: #f8fafc;
}

.spectator-hint {
  margin: 0;
  font-size: 0.75rem;
  color: #94a3b8;
  max-width: 320px;
  line-height: 1.4;
}

.waiting-question {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  gap: 1rem;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid rgba(255,255,255,0.1);
  border-radius: 50%;
  border-top-color: #38bdf8;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
