<script lang="ts">
import { defineComponent } from 'vue'
import { store } from '../store'

export default defineComponent({
  name: 'PauseOverlay',
  data() {
    return {
      store,
    }
  },
  computed: {
    visible(): boolean {
      return store.isPaused && !store.gameOver
    },
  },
})
</script>

<template>
  <div v-if="visible" class="pause-overlay" aria-live="polite">
    <div class="pause-card">
      <div class="pause-icon">⏸</div>
      <h2 class="pause-title">Game Paused</h2>
      <p class="pause-message">An admin has paused this room. Please wait until the game resumes.</p>
    </div>
  </div>
</template>

<style scoped>
.pause-overlay {
  position: fixed;
  inset: 0;
  z-index: 9000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(8, 15, 30, 0.82);
  backdrop-filter: blur(4px);
  pointer-events: all;
}

.pause-card {
  background: linear-gradient(145deg, #0f1e3a 0%, #111827 100%);
  border: 1px solid #f59e0b;
  border-radius: 16px;
  padding: 2rem 2.5rem;
  max-width: 420px;
  text-align: center;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.7), 0 0 40px rgba(245, 158, 11, 0.15);
}

.pause-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.pause-title {
  margin: 0 0 0.75rem 0;
  font-size: 1.4rem;
  font-weight: 900;
  color: #fbbf24;
}

.pause-message {
  margin: 0;
  color: #94a3b8;
  font-size: 0.95rem;
  line-height: 1.5;
}
</style>
