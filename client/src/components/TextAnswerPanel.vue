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
        placeholder="Type your answer here..."
        :disabled="disabled"
        rows="4"
        @keydown.enter.ctrl.exact="submitAnswer"
      ></textarea>
      <span class="input-hint">Press Ctrl+Enter to submit</span>
    </div>

    <div class="submit-action">
      <button
        class="submit-btn"
        :disabled="disabled || !textValue.trim()"
        @click="submitAnswer"
      >
        {{ disabled ? 'SUBMITTED' : 'SUBMIT ANSWER' }}
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
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 0.75rem;
  color: #f1f5f9;
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
  border-color: #38bdf8;
  background: rgba(15, 23, 42, 0.6);
  box-shadow: 0 0 10px rgba(56, 189, 248, 0.1);
}

.text-input:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.input-hint {
  font-size: 0.65rem;
  color: #64748b;
  text-align: right;
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
