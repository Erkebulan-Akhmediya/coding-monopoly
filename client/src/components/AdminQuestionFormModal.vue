<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import {
  adminApiService,
  validateProblemInput,
  type Problem,
  type ProblemInput,
  type OptionInput,
} from '../services/adminApiService'
import { adminStore } from '../adminStore'

export default defineComponent({
  name: 'AdminQuestionFormModal',

  props: {
    problem: {
      type: Object as PropType<Problem | null>,
      default: null,
    },
    initialType: {
      type: String as PropType<'mcq' | 'text'>,
      default: 'mcq',
    },
  },

  emits: ['saved', 'cancel'],

  data() {
    const isEdit = !!this.problem
    let type: 'mcq' | 'text' = this.initialType
    let difficulty: 'easy' | 'medium' | 'hard' = 'easy'
    let title = ''
    let prompt = ''
    let is_published = false
    let options: OptionInput[] = [
      { text: '', is_correct: true },
      { text: '', is_correct: false },
    ]
    let accepted_answers: string[] = ['']

    if (this.problem) {
      type = this.problem.type
      difficulty = this.problem.difficulty
      title = this.problem.title
      prompt = this.problem.prompt
      is_published = this.problem.is_published
      if (this.problem.type === 'mcq' && this.problem.options && this.problem.options.length > 0) {
        options = this.problem.options.map((opt) => ({
          text: opt.text,
          is_correct: opt.is_correct,
        }))
      }
      if (
        this.problem.type === 'text' &&
        this.problem.accepted_answers &&
        this.problem.accepted_answers.length > 0
      ) {
        accepted_answers = [...this.problem.accepted_answers]
      }
    }

    return {
      isEdit,
      form: {
        type,
        difficulty,
        title,
        prompt,
        is_published,
        options,
        accepted_answers,
      },
      validationErrors: [] as string[],
      serverError: '' as string,
      saving: false as boolean,
    }
  },

  methods: {
    handleTypeChange() {
      if (this.form.type === 'mcq' && this.form.options.length < 2) {
        this.form.options = [
          { text: '', is_correct: true },
          { text: '', is_correct: false },
        ]
      } else if (this.form.type === 'text' && this.form.accepted_answers.length < 1) {
        this.form.accepted_answers = ['']
      }
      this.validationErrors = []
    },

    addOption() {
      this.form.options.push({ text: '', is_correct: false })
    },

    removeOption(index: number) {
      if (this.form.options.length > 2) {
        this.form.options.splice(index, 1)
      }
    },

    addAnswer() {
      this.form.accepted_answers.push('')
    },

    removeAnswer(index: number) {
      if (this.form.accepted_answers.length > 1) {
        this.form.accepted_answers.splice(index, 1)
      }
    },

    async handleSubmit() {
      this.serverError = ''
      this.validationErrors = []

      const payload: ProblemInput = {
        type: this.form.type,
        difficulty: this.form.difficulty,
        title: this.form.title,
        prompt: this.form.prompt,
        is_published: this.form.is_published,
      }

      if (this.form.type === 'mcq') {
        payload.options = this.form.options.map((opt) => ({
          text: opt.text,
          is_correct: opt.is_correct,
        }))
      } else {
        payload.accepted_answers = this.form.accepted_answers.map((ans) => ans)
      }

      // Client-side validation mirroring Phase 4 server-side rules
      const errors = validateProblemInput(payload)
      if (errors.length > 0) {
        this.validationErrors = errors
        return
      }

      this.saving = true
      try {
        let result: Problem
        if (this.isEdit && this.problem) {
          result = await adminApiService.updateProblem(adminStore.token, this.problem.id, payload)
        } else {
          result = await adminApiService.createProblem(adminStore.token, payload)
        }
        this.$emit('saved', result)
      } catch (err: any) {
        this.serverError = err.message || 'Failed to save question.'
      } finally {
        this.saving = false
      }
    },

    handleCancel() {
      this.$emit('cancel')
    },
  },
})
</script>

<template>
  <div class="modal-backdrop" @click.self="handleCancel">
    <div class="modal-card">
      <header class="modal-header">
        <h2>{{ isEdit ? 'Edit Question' : 'Create New Question' }}</h2>
        <button class="close-btn" @click="handleCancel">✕</button>
      </header>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <!-- Client-side & server validation errors -->
        <div v-if="validationErrors.length > 0" class="error-banner" id="form-validation-errors">
          <strong>Validation Errors:</strong>
          <ul>
            <li v-for="(err, i) in validationErrors" :key="i">{{ err }}</li>
          </ul>
        </div>
        <div v-if="serverError" class="error-banner" id="form-server-error">
          <strong>Server Error:</strong> {{ serverError }}
        </div>

        <!-- Question Type & Difficulty Row -->
        <div class="form-row">
          <div class="form-group flex-1">
            <label for="problem-type">Question Type</label>
            <select
              id="problem-type"
              v-model="form.type"
              class="form-input"
              @change="handleTypeChange"
            >
              <option value="mcq">MCQ (Multiple Choice)</option>
              <option value="text">Text Answer</option>
            </select>
          </div>

          <div class="form-group flex-1">
            <label for="problem-difficulty">Difficulty Level</label>
            <select id="problem-difficulty" v-model="form.difficulty" class="form-input">
              <option value="easy">Easy (30s countdown)</option>
              <option value="medium">Medium (45s countdown)</option>
              <option value="hard">Hard (60s countdown)</option>
            </select>
          </div>
        </div>

        <!-- Title -->
        <div class="form-group">
          <label for="problem-title">Title</label>
          <input
            id="problem-title"
            v-model="form.title"
            type="text"
            placeholder="e.g. Array indexing in C"
            class="form-input"
          />
        </div>

        <!-- Prompt -->
        <div class="form-group">
          <label for="problem-prompt">Prompt / Code Snippet / Question Text</label>
          <textarea
            id="problem-prompt"
            v-model="form.prompt"
            rows="4"
            placeholder="Enter the question prompt or code listing..."
            class="form-input code-font"
          ></textarea>
        </div>

        <!-- Published Toggle -->
        <div class="form-group checkbox-group">
          <label class="toggle-label">
            <input
              id="problem-published"
              v-model="form.is_published"
              type="checkbox"
              class="toggle-checkbox"
            />
            <span class="toggle-text">Publish immediately (visible in game pools)</span>
          </label>
        </div>

        <!-- MCQ OPTIONS SUB-FORM -->
        <div v-if="form.type === 'mcq'" class="subform-section">
          <div class="subform-header">
            <h3>MCQ Options (at least 2 required, check correct answer(s))</h3>
            <button
              id="add-option-btn"
              type="button"
              class="btn-secondary btn-sm"
              @click="addOption"
            >
              + Add Option
            </button>
          </div>

          <div class="options-list">
            <div
              v-for="(opt, idx) in form.options"
              :key="idx"
              class="option-row"
            >
              <span class="opt-num">#{{ idx + 1 }}</span>
              <input
                v-model="opt.text"
                type="text"
                :placeholder="`Option ${idx + 1} text...`"
                class="form-input option-text-input"
              />
              <label class="correct-checkbox-label">
                <input
                  v-model="opt.is_correct"
                  type="checkbox"
                  class="option-correct-checkbox"
                />
                <span :class="{ 'text-correct': opt.is_correct }">Correct</span>
              </label>
              <button
                type="button"
                class="btn-danger btn-xs remove-option-btn"
                :disabled="form.options.length <= 2"
                @click="removeOption(idx)"
                title="Remove option"
              >
                ✕
              </button>
            </div>
          </div>
        </div>

        <!-- TEXT ACCEPTED ANSWERS SUB-FORM -->
        <div v-else-if="form.type === 'text'" class="subform-section">
          <div class="subform-header">
            <h3>Accepted Answers (at least 1 required, case-folded & trimmed)</h3>
            <button
              id="add-answer-btn"
              type="button"
              class="btn-secondary btn-sm"
              @click="addAnswer"
            >
              + Add Accepted Answer
            </button>
          </div>

          <div class="answers-list">
            <div
              v-for="(_, idx) in form.accepted_answers"
              :key="idx"
              class="answer-row"
            >
              <span class="ans-num">#{{ idx + 1 }}</span>
              <input
                v-model="form.accepted_answers[idx]"
                type="text"
                :placeholder="`Accepted answer variant ${idx + 1}...`"
                class="form-input answer-text-input"
              />
              <button
                type="button"
                class="btn-danger btn-xs remove-answer-btn"
                :disabled="form.accepted_answers.length <= 1"
                @click="removeAnswer(idx)"
                title="Remove answer"
              >
                ✕
              </button>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <footer class="modal-footer">
          <button
            id="cancel-problem-btn"
            type="button"
            class="btn-secondary"
            @click="handleCancel"
          >
            Cancel
          </button>
          <button
            id="save-problem-btn"
            type="submit"
            class="btn-primary"
            :disabled="saving"
          >
            <span v-if="saving">Saving...</span>
            <span v-else>{{ isEdit ? 'Update Question' : 'Create Question' }}</span>
          </button>
        </footer>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(4, 9, 20, 0.82);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 1rem;
}

.modal-card {
  background: linear-gradient(145deg, #0f1e3a 0%, #111827 100%);
  border: 1px solid #1e3a5f;
  border-radius: 16px;
  width: 100%;
  max-width: 680px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.8), 0 0 0 1px rgba(96, 165, 250, 0.15);
  color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #1e293b;
  background: #0c1527;
}

.modal-header h2 {
  font-size: 1.15rem;
  font-weight: 800;
  margin: 0;
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.close-btn {
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}
.close-btn:hover {
  color: #f1f5f9;
  background: #1e293b;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}

.error-banner {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid #ef4444;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  color: #f87171;
  font-size: 0.85rem;
}
.error-banner ul {
  margin: 0.3rem 0 0 1.2rem;
  padding: 0;
}

.form-row {
  display: flex;
  gap: 1rem;
}
.flex-1 {
  flex: 1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.form-group label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.form-input {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 0.6rem 0.85rem;
  color: #f1f5f9;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s;
}
.form-input:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
}
.code-font {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
}
.toggle-label {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  cursor: pointer;
  font-size: 0.88rem;
  color: #cbd5e1;
}
.toggle-checkbox {
  width: 18px;
  height: 18px;
  accent-color: #3b82f6;
  cursor: pointer;
}

.subform-section {
  background: #0c1527;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.subform-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.subform-header h3 {
  font-size: 0.82rem;
  font-weight: 700;
  color: #94a3b8;
  margin: 0;
}

.options-list,
.answers-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.option-row,
.answer-row {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.opt-num,
.ans-num {
  font-size: 0.75rem;
  font-weight: 700;
  color: #64748b;
  width: 24px;
}

.correct-checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.8rem;
  font-weight: 600;
  color: #94a3b8;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
}
.option-correct-checkbox {
  accent-color: #10b981;
  width: 16px;
  height: 16px;
  cursor: pointer;
}
.text-correct {
  color: #34d399;
  font-weight: 700;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid #1e293b;
}

.btn-primary {
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1.25rem;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.15s;
}
.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-weight: 600;
  font-size: 0.88rem;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-secondary:hover {
  background: #334155;
}

.btn-danger {
  background: #7f1d1d;
  color: #fca5a5;
  border: 1px solid #991b1b;
  border-radius: 6px;
  cursor: pointer;
}
.btn-danger:hover:not(:disabled) {
  background: #991b1b;
}
.btn-danger:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.btn-sm {
  padding: 0.3rem 0.65rem;
  font-size: 0.78rem;
  border-radius: 6px;
}
.btn-xs {
  padding: 0.25rem 0.45rem;
  font-size: 0.75rem;
}
</style>
