<script lang="ts">
import { defineComponent } from 'vue'

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
        <span class="option-text">{{ opt.text }}</span>
      </button>
    </div>

    <div class="submit-action">
      <button
        class="submit-btn"
        :disabled="disabled || selectedIDs.length === 0"
        @click="submitAnswer"
      >
        {{ disabled ? 'SUBMITTED' : 'SUBMIT ANSWER' }}
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
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  cursor: pointer;
  color: #f1f5f9;
  text-align: left;
  transition: all 0.2s ease;
  width: 100%;
}

.option-card:not(.disabled):hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(56, 189, 248, 0.4);
  transform: translateX(2px);
}

.option-card.selected {
  background: rgba(56, 189, 248, 0.1);
  border-color: #38bdf8;
  box-shadow: 0 0 10px rgba(56, 189, 248, 0.15);
}

.option-card.disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.checkbox-indicator {
  width: 1.1rem;
  height: 1.1rem;
  border-radius: 4px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 900;
  color: #38bdf8;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.option-card.selected .checkbox-indicator {
  border-color: #38bdf8;
  background: rgba(56, 189, 248, 0.1);
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
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1.25rem;
  font-weight: 700;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 10px rgba(59, 130, 246, 0.2);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(59, 130, 246, 0.3);
}

.submit-btn:disabled {
  background: #475569;
  color: #94a3b8;
  cursor: not-allowed;
  box-shadow: none;
}
</style>
