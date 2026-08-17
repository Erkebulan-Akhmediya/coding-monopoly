<script lang="ts">
import { defineComponent } from 'vue'
import { store, type EffectToast } from '../store'

export default defineComponent({
  name: 'EffectToastStack',
  data() {
    return { store }
  },
  computed: {
    toasts(): EffectToast[] {
      return store.effectToasts
    },
  },
  methods: {
    labelFor(type: string): string {
      const key = `effects.${type}`
      return this.$te(key) ? this.$t(key) : this.$t('common.effect')
    },
    xpText(toast: EffectToast): string {
      if (!toast.xpDelta) return ''
      return toast.xpDelta > 0 ? `+${toast.xpDelta} ${this.$t('common.xp')}` : `${toast.xpDelta} ${this.$t('common.xp')}`
    },
  },
})
</script>

<template>
  <div class="toast-stack" aria-live="polite">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="toast"
      :class="'effect-' + toast.effectType"
    >
      <span class="toast-kind">{{ labelFor(toast.effectType) }}</span>
      <span class="toast-body">{{ toast.description }}</span>
      <span v-if="xpText(toast)" class="toast-xp">{{ xpText(toast) }}</span>
    </div>
  </div>
</template>

<style scoped>
.toast-stack {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: 80;
  display: flex;
  flex-direction: column-reverse;
  gap: 0.5rem;
  max-width: min(360px, calc(100vw - 32px));
  pointer-events: none;
}

.toast {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  padding: 0.7rem 0.85rem;
  border-radius: 8px;
  border-left: 4px solid #64748b;
  background: rgba(15, 23, 42, 0.94);
  color: #f8fafc;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
  animation: toast-in 0.28s ease-out;
}

.toast-kind {
  font-size: 0.68rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #94a3b8;
  font-weight: 700;
}

.toast-body {
  font-size: 0.88rem;
  line-height: 1.3;
}

.toast-xp {
  font-size: 0.8rem;
  font-weight: 800;
  color: #34d399;
}

.effect-xp_gain { border-left-color: #34d399; }
.effect-xp_loss { border-left-color: #f87171; }
.effect-xp_loss .toast-xp { color: #f87171; }
.effect-teleport { border-left-color: #22d3ee; }
.effect-skip_next { border-left-color: #fb923c; }
.effect-double_xp { border-left-color: #fbbf24; }
.effect-free_pass { border-left-color: #a78bfa; }
.effect-special_challenge { border-left-color: #f472b6; }
.effect-mystery { border-left-color: #c084fc; }
.effect-deploy { border-left-color: #4ade80; }
.effect-code_freeze { border-left-color: #67e8f9; }
.effect-coffee_break { border-left-color: #d6b28a; }
.effect-deadline { border-left-color: #f87171; }

@keyframes toast-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
