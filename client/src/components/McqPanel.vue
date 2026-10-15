<script lang="ts">
import { defineComponent } from 'vue'
import { localeText } from '../i18n'

export default defineComponent({
  name: 'McqPanel',
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
      selectedIDs: [] as string[]
    }
  },
  methods: {
    optionText(text: any): string {
      return localeText(text)
    },
    toggleOption(optionID: string) {
      if (this.disabled) return
      const idx = this.selectedIDs.indexOf(optionID)
      if (idx >= 0) {
        this.selectedIDs.splice(idx, 1)
      } else {
        this.selectedIDs.push(optionID)
      }
    },
    isSelected(optionID: string): boolean {
      return this.selectedIDs.includes(optionID)
    },
    submitAnswer() {
      if (this.disabled || this.selectedIDs.length === 0) return
      this.$emit('submit', this.selectedIDs)
    }
  }
})
</script>

<template>
  <div class="mcq-panel">
    <div class="options-list">
      <button
        v-for="opt in question.options"
        :key="opt.id"
        class="option-card"
        :class="{ selected: isSelected(opt.id), disabled: disabled }"
        @click="toggleOption(opt.id)"
        :disabled="disabled"
      >
        <div class="checkbox-indicator">
          <span v-if="isSelected(opt.id)">✓</span>
        </div>
        <span class="option-text">{{ optionText(opt.text) }}</span>
      </button>
    </div>

    <div class="submit-action">
      <button
        class="submit-btn"
        :disabled="disabled || selectedIDs.length === 0"
        @click="submitAnswer"
      >
        {{ disabled ? $t('problem.submitted') : $t('problem.submitAnswer') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.mcq-panel {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.option-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: color-mix(in srgb, var(--c-slate-900) 40%, transparent);
  border: 1px solid color-mix(in srgb, var(--c-bright) 5%, transparent);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  cursor: pointer;
  color: var(--c-slate-100);
  text-align: left;
  transition: all 0.2s ease;
  width: 100%;
}

.option-card:not(.disabled):hover {
  background: color-mix(in srgb, var(--c-bright) 5%, transparent);
  border-color: color-mix(in srgb, var(--c-sky-400) 40%, transparent);
  transform: translateX(2px);
}

.option-card.selected {
  background: color-mix(in srgb, var(--c-sky-400) 10%, transparent);
  border-color: var(--c-sky-400);
  box-shadow: 0 0 10px color-mix(in srgb, var(--c-sky-400) 15%, transparent);
}

.option-card.disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.checkbox-indicator {
  width: 1.1rem;
  height: 1.1rem;
  border-radius: 4px;
  border: 2px solid color-mix(in srgb, var(--c-bright) 20%, transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 900;
  color: var(--c-sky-400);
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.option-card.selected .checkbox-indicator {
  border-color: var(--c-sky-400);
  background: color-mix(in srgb, var(--c-sky-400) 10%, transparent);
}

.option-text {
  font-size: 0.8rem;
  line-height: 1.4;
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
