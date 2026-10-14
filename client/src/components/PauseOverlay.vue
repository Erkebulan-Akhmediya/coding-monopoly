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
      <h2 class="pause-title">{{ $t('pause.title') }}</h2>
      <p class="pause-message">{{ $t('pause.message') }}</p>
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
  background: color-mix(in srgb, var(--c-navy-950) 82%, transparent);
  backdrop-filter: blur(4px);
  pointer-events: all;
}

.pause-card {
  background: linear-gradient(145deg, var(--c-navy-800) 0%, var(--c-gray-900) 100%);
  border: 1px solid var(--c-amber-500);
  border-radius: 16px;
  padding: 2rem 2.5rem;
  max-width: 420px;
  text-align: center;
  box-shadow: 0 24px 60px color-mix(in srgb, var(--c-shadow) 70%, transparent), 0 0 40px color-mix(in srgb, var(--c-amber-500) 15%, transparent);
}

.pause-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.pause-title {
  margin: 0 0 0.75rem 0;
  font-size: 1.4rem;
  font-weight: 900;
  color: var(--c-amber-400);
}

.pause-message {
  margin: 0;
  color: var(--c-slate-400);
  font-size: 0.95rem;
  line-height: 1.5;
}
</style>
