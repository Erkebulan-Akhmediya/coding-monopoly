<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'TextAnswerPanel',
  props: {
    question: {
      type: Object,
      required: true
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      textValue: ''
    }
  },
  methods: {
    submitAnswer() {
      if (this.disabled || !this.textValue.trim()) return
      this.$emit('submit', this.textValue)
    }
  }
})
</script>

<template>
  <div class="text-panel">
    <div class="input-container">
      <textarea
        v-model="textValue"
        class="text-input"
        :placeholder="$t('problem.textPlaceholder')"
        :disabled="disabled"
        rows="4"
        @keydown.enter.ctrl.exact="submitAnswer"
      ></textarea>
      <span class="input-hint">{{ $t('problem.textHint') }}</span>
    </div>

    <div class="submit-action">
      <button
        class="submit-btn"
        :disabled="disabled || !textValue.trim()"
        @click="submitAnswer"
      >
        {{ disabled ? $t('problem.submitted') : $t('problem.submitAnswer') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.text-panel {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input-container {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.text-input {
  background: color-mix(in srgb, var(--c-slate-900) 40%, transparent);
  border: 1px solid color-mix(in srgb, var(--c-bright) 5%, transparent);
  border-radius: 8px;
  padding: 0.75rem;
  color: var(--c-slate-100);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8rem;
  line-height: 1.5;
  resize: vertical;
  outline: none;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.text-input:focus:not(:disabled) {
  border-color: var(--c-sky-400);
  background: color-mix(in srgb, var(--c-slate-900) 60%, transparent);
  box-shadow: 0 0 10px color-mix(in srgb, var(--c-sky-400) 10%, transparent);
}

.text-input:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.input-hint {
  font-size: 0.65rem;
  color: var(--c-slate-500);
  text-align: right;
}

.submit-action {
  display: flex;
  justify-content: flex-end;
}

.submit-btn {
  background: linear-gradient(90deg, var(--c-blue-500), var(--c-violet-500));
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1.25rem;
  font-weight: 700;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  color: var(--c-white);
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 10px color-mix(in srgb, var(--c-blue-500) 20%, transparent);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 14px color-mix(in srgb, var(--c-blue-500) 30%, transparent);
}

.submit-btn:disabled {
  background: var(--c-slate-600);
  color: var(--c-slate-400);
  cursor: not-allowed;
  box-shadow: none;
}
</style>
