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
  border-left: 4px solid var(--c-slate-500);
  background: color-mix(in srgb, var(--c-slate-900) 94%, transparent);
  color: var(--c-slate-50);
  box-shadow: 0 8px 24px color-mix(in srgb, var(--c-shadow) 35%, transparent);
  animation: toast-in 0.28s ease-out;
}

.toast-kind {
  font-size: 0.68rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--c-slate-400);
  font-weight: 700;
}

.toast-body {
  font-size: 0.88rem;
  line-height: 1.3;
}

.toast-xp {
  font-size: 0.8rem;
  font-weight: 800;
  color: var(--c-emerald-400);
}

.effect-xp_gain { border-left-color: var(--c-xp_gain--emerald-400); }
.effect-xp_loss { border-left-color: var(--c-xp_loss--red-400); }
.effect-xp_loss .toast-xp { color: var(--c-xp_loss--red-400); }
.effect-teleport { border-left-color: var(--c-teleport--cyan-400); }
.effect-skip_next { border-left-color: var(--c-skip_next--orange-400); }
.effect-double_xp { border-left-color: var(--c-double_xp--amber-400); }
.effect-free_pass { border-left-color: var(--c-free_pass--violet-400); }
.effect-special_challenge { border-left-color: var(--c-special_challenge--pink-400); }
.effect-mystery { border-left-color: var(--c-mystery--purple-400); }
.effect-deploy { border-left-color: var(--c-deploy--green-400); }
.effect-code_freeze { border-left-color: var(--c-code_freeze--cyan-300); }
.effect-coffee_break { border-left-color: var(--c-coffee_break--latte-350); }
.effect-deadline { border-left-color: var(--c-deadline--red-400); }

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
