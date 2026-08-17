<script lang="ts">
import { defineComponent } from 'vue'
import {
  adminApiService,
  type Problem,
  type ProblemFilters,
} from '../services/adminApiService'
import { adminStore } from '../adminStore'
import AdminQuestionFormModal from './AdminQuestionFormModal.vue'

export default defineComponent({
  name: 'AdminQuestionList',

  components: {
    AdminQuestionFormModal,
  },

  data() {
    return {
      problems: [] as Problem[],
      loading: false as boolean,
      error: '' as string,

      // Filters
      filters: {
        type: '' as string,
        difficulty: '' as string,
        is_published: '' as string,
      },
      searchQuery: '' as string,

      // Modal state
      showFormModal: false as boolean,
      selectedProblem: null as Problem | null,
      initialFormType: 'mcq' as 'mcq' | 'text',

      // Confirm Delete Modal
      confirmDeleteProblem: null as Problem | null,
      deleting: false as boolean,

      // Publishing state tracking per id
      togglingPublishId: null as string | null,
    }
  },

  computed: {
    filteredProblems(): Problem[] {
      let result = this.problems
      if (this.searchQuery.trim()) {
        const q = this.searchQuery.toLowerCase().trim()
        result = result.filter(
          (p) =>
            p.title.toLowerCase().includes(q) ||
            p.prompt.toLowerCase().includes(q) ||
            p.id.toLowerCase().includes(q)
        )
      }
      return result
    },
  },

  mounted() {
    this.fetchProblems()
  },

  methods: {
    async fetchProblems() {
      if (!adminStore.token) return
      this.loading = true
      this.error = ''
      try {
        const reqFilters: ProblemFilters = {
          type: this.filters.type,
          difficulty: this.filters.difficulty,
          is_published: this.filters.is_published,
        }
        const data = await adminApiService.listProblems(adminStore.token, reqFilters)
        this.problems = data
      } catch (err: any) {
        this.error = err.message || 'Failed to fetch questions.'
      } finally {
        this.loading = false
      }
    },

    onFilterChange() {
      this.fetchProblems()
    },

    openCreateModal(type: 'mcq' | 'text') {
      this.selectedProblem = null
      this.initialFormType = type
      this.showFormModal = true
    },

    openEditModal(problem: Problem) {
      this.selectedProblem = problem
      this.showFormModal = true
    },

    closeFormModal() {
      this.showFormModal = false
      this.selectedProblem = null
    },

    onProblemSaved() {
      this.closeFormModal()
      this.fetchProblems()
    },

    requestDelete(problem: Problem) {
      this.confirmDeleteProblem = problem
    },

    cancelDelete() {
      this.confirmDeleteProblem = null
    },

    async confirmDelete() {
      if (!this.confirmDeleteProblem || !adminStore.token) return
      this.deleting = true
      try {
        await adminApiService.deleteProblem(adminStore.token, this.confirmDeleteProblem.id)
        this.confirmDeleteProblem = null
        await this.fetchProblems()
      } catch (err: any) {
        alert(`Delete failed: ${err.message}`)
      } finally {
        this.deleting = false
      }
    },

    async togglePublishStatus(problem: Problem) {
      if (!adminStore.token || this.togglingPublishId) return
      this.togglingPublishId = problem.id
      try {
        if (!problem.is_published) {
          // Publish endpoint validates problem before publishing
          await adminApiService.publishProblem(adminStore.token, problem.id)
        } else {
          // Unpublish endpoint: update is_published to false
          const payload = {
            type: problem.type,
            difficulty: problem.difficulty,
            title: problem.title,
            prompt: problem.prompt,
            is_published: false,
            options: problem.options
              ? problem.options.map((o) => ({ text: o.text, is_correct: o.is_correct }))
              : [],
            accepted_answers: problem.accepted_answers || [],
          }
          await adminApiService.updateProblem(adminStore.token, problem.id, payload)
        }
        await this.fetchProblems()
      } catch (err: any) {
        alert(`Failed to update publish status: ${err.message}`)
      } finally {
        this.togglingPublishId = null
      }
    },

    difficultyBadgeClass(diff: string): string {
      switch (diff) {
        case 'easy':
          return 'badge-easy'
        case 'medium':
          return 'badge-medium'
        case 'hard':
          return 'badge-hard'
        default:
          return ''
      }
    },

    formatDate(isoStr: string): string {
      if (!isoStr) return ''
      return new Date(isoStr).toLocaleString()
    },
  },
})
</script>

<template>
  <div class="question-list-container">
    <!-- Header & Create Actions -->
    <div class="list-header">
      <div class="header-left">
        <h2 class="section-heading">📚 {{ $t('admin.questionBankManagement') }}</h2>
        <span class="count-badge">{{ $t('admin.questionCount', { count: filteredProblems.length }) }}</span>
      </div>
      <div class="header-right">
        <button
          id="create-mcq-btn"
          class="btn-primary btn-sm"
          @click="openCreateModal('mcq')"
        >
          {{ $t('admin.newMcq') }}
        </button>
        <button
          id="create-text-btn"
          class="btn-secondary btn-sm"
          @click="openCreateModal('text')"
        >
          {{ $t('admin.newText') }}
        </button>
      </div>
    </div>

    <!-- Filters Bar -->
    <div class="filters-bar">
      <div class="filter-group">
        <label for="filter-type">{{ $t('admin.filterType') }}</label>
        <select
          id="filter-type"
          v-model="filters.type"
          class="filter-select"
          @change="onFilterChange"
        >
          <option value="">{{ $t('admin.allTypes') }}</option>
          <option value="mcq">MCQ</option>
          <option value="text">Text</option>
        </select>
      </div>

      <div class="filter-group">
        <label for="filter-difficulty">{{ $t('admin.filterDifficulty') }}</label>
        <select
          id="filter-difficulty"
          v-model="filters.difficulty"
          class="filter-select"
          @change="onFilterChange"
        >
          <option value="">{{ $t('admin.allDifficulties') }}</option>
          <option value="easy">{{ $t('admin.easyTimed') }}</option>
          <option value="medium">{{ $t('admin.mediumTimed') }}</option>
          <option value="hard">{{ $t('admin.hardTimed') }}</option>
        </select>
      </div>

      <div class="filter-group">
        <label for="filter-published">{{ $t('admin.filterStatus') }}</label>
        <select
          id="filter-published"
          v-model="filters.is_published"
          class="filter-select"
          @change="onFilterChange"
        >
          <option value="">{{ $t('admin.allStatuses') }}</option>
          <option value="true">{{ $t('admin.published') }}</option>
          <option value="false">{{ $t('admin.draft') }}</option>
        </select>
      </div>

      <div class="filter-group flex-search">
        <label for="search-input">{{ $t('admin.filterSearch') }}</label>
        <input
          id="search-input"
          v-model="searchQuery"
          type="text"
          :placeholder="$t('admin.searchPlaceholder')"
          class="search-input"
        />
      </div>
    </div>

    <!-- Error Banner -->
    <div v-if="error" class="error-box">
      <span>⚠️ {{ error }}</span>
      <button class="btn-secondary btn-xs" @click="fetchProblems">{{ $t('admin.retry') }}</button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      {{ $t('admin.loadingQuestions') }}
    </div>

    <!-- Empty state -->
    <div v-else-if="filteredProblems.length === 0" class="empty-state">
      <p>{{ $t('admin.noQuestions') }}</p>
      <button class="btn-primary btn-sm" @click="openCreateModal('mcq')">{{ $t('admin.createFirstQuestion') }}</button>
    </div>

    <!-- Questions Table / Grid -->
    <div v-else class="table-wrapper">
      <table class="questions-table">
        <thead>
          <tr>
            <th>{{ $t('admin.colType') }}</th>
            <th>{{ $t('admin.colDifficulty') }}</th>
            <th>{{ $t('admin.colTitlePrompt') }}</th>
            <th>{{ $t('admin.colAnswers') }}</th>
            <th>{{ $t('admin.colStatus') }}</th>
            <th>{{ $t('admin.colActions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filteredProblems" :key="p.id" class="question-row">
            <!-- Type -->
            <td class="col-type">
              <span :class="['type-badge', p.type]">
                {{ p.type.toUpperCase() }}
              </span>
            </td>

            <!-- Difficulty -->
            <td class="col-difficulty">
              <span :class="['diff-badge', difficultyBadgeClass(p.difficulty)]">
                {{ p.difficulty }}
              </span>
            </td>

            <!-- Title & Prompt -->
            <td class="col-title">
              <div class="q-title">{{ p.title }}</div>
              <div class="q-prompt-preview">{{ p.prompt }}</div>
            </td>

            <!-- Content preview -->
            <td class="col-content">
              <div v-if="p.type === 'mcq'" class="mcq-preview">
                <span class="meta-label">{{ $t('admin.optionsCount', { count: p.options ? p.options.length : 0 }) }}</span>
                <ul class="preview-list">
                  <li
                    v-for="(opt, idx) in (p.options || []).slice(0, 3)"
                    :key="idx"
                    :class="{ 'is-correct-preview': opt.is_correct }"
                  >
                    {{ opt.is_correct ? '✓ ' : '• ' }}{{ opt.text }}
                  </li>
                  <li v-if="(p.options || []).length > 3" class="more-hint">
                    {{ $t('admin.moreOptions', { count: p.options!.length - 3 }) }}
                  </li>
                </ul>
              </div>
              <div v-else-if="p.type === 'text'" class="text-preview">
                <span class="meta-label">{{ $t('admin.acceptedAnswers') }}</span>
                <div class="answers-tags">
                  <span
                    v-for="(ans, idx) in (p.accepted_answers || [])"
                    :key="idx"
                    class="ans-tag"
                  >
                    "{{ ans }}"
                  </span>
                </div>
              </div>
            </td>

            <!-- Published Status -->
            <td class="col-status">
              <button
                :id="`toggle-publish-${p.id}`"
                class="toggle-publish-btn"
                :class="p.is_published ? 'status-published' : 'status-draft'"
                :disabled="togglingPublishId === p.id"
                @click="togglePublishStatus(p)"
                :title="p.is_published ? $t('admin.clickUnpublish') : $t('admin.clickPublish')"
              >
                <span v-if="togglingPublishId === p.id">...</span>
                <span v-else>{{ p.is_published ? `🟢 ${$t('admin.publishedStatus')}` : `🟠 ${$t('admin.draftStatus')}` }}</span>
              </button>
            </td>

            <!-- Actions -->
            <td class="col-actions">
              <button
                :id="`edit-btn-${p.id}`"
                class="btn-secondary btn-xs edit-problem-btn"
                @click="openEditModal(p)"
                :title="$t('admin.editQuestion')"
              >
                ✏️ {{ $t('admin.edit') }}
              </button>
              <button
                :id="`delete-btn-${p.id}`"
                class="btn-danger btn-xs delete-problem-btn"
                @click="requestDelete(p)"
                :title="$t('admin.deleteQuestion')"
              >
                🗑️ {{ $t('admin.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Form Modal -->
    <AdminQuestionFormModal
      v-if="showFormModal"
      :problem="selectedProblem"
      :initialType="initialFormType"
      @saved="onProblemSaved"
      @cancel="closeFormModal"
    />

    <!-- Confirm Delete Modal -->
    <div v-if="confirmDeleteProblem" class="confirm-modal-backdrop">
      <div class="confirm-modal-card">
        <h3>{{ $t('admin.deleteConfirmTitle') }}</h3>
        <p>
          {{ $t('admin.deleteConfirmBody', { title: confirmDeleteProblem.title }) }}
        </p>
        <div class="confirm-modal-actions">
          <button class="btn-secondary" @click="cancelDelete" :disabled="deleting">
            {{ $t('common.cancel') }}
          </button>
          <button
            id="confirm-delete-btn"
            class="btn-danger"
            @click="confirmDelete"
            :disabled="deleting"
          >
            <span v-if="deleting">{{ $t('admin.deleting') }}</span>
            <span v-else>{{ $t('admin.yesDelete') }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.question-list-container {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  padding: 1.5rem;
  background: #080f1e;
  color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  min-height: calc(100vh - 49px);
  box-sizing: border-box;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.section-heading {
  font-size: 1.25rem;
  font-weight: 800;
  margin: 0;
  background: linear-gradient(90deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.count-badge {
  background: #1e293b;
  border: 1px solid #334155;
  color: #94a3b8;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
}

.header-right {
  display: flex;
  gap: 0.75rem;
}

/* Filters bar */
.filters-bar {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  background: #0c1527;
  border: 1px solid #1e293b;
  border-radius: 12px;
  padding: 0.85rem 1.1rem;
  align-items: center;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.filter-group label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-select,
.search-input {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 0.45rem 0.75rem;
  color: #f1f5f9;
  font-size: 0.85rem;
  outline: none;
  transition: border-color 0.2s;
}
.filter-select:focus,
.search-input:focus {
  border-color: #60a5fa;
}

.flex-search {
  flex: 1;
  min-width: 200px;
}
.search-input {
  width: 100%;
}

/* Error box */
.error-box {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid #ef4444;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #f87171;
}

/* Loading & Empty states */
.loading-state,
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #64748b;
  background: #0c1527;
  border: 1px solid #1e293b;
  border-radius: 12px;
}
.empty-state p {
  margin-bottom: 1rem;
}

/* Table layout */
.table-wrapper {
  overflow-x: auto;
  border: 1px solid #1e293b;
  border-radius: 12px;
  background: #0c1527;
}

.questions-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.88rem;
}

.questions-table th {
  background: #0f172a;
  color: #64748b;
  font-size: 0.7rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  padding: 0.85rem 1rem;
  border-bottom: 1px solid #1e293b;
}

.questions-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid #1e293b;
  vertical-align: top;
}

.question-row:hover {
  background: #111e35;
}

/* Badges */
.type-badge {
  display: inline-block;
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 800;
  letter-spacing: 0.5px;
}
.type-badge.mcq {
  background: #1e3a8a;
  color: #93c5fd;
  border: 1px solid #3b82f6;
}
.type-badge.text {
  background: #581c87;
  color: #e9d5ff;
  border: 1px solid #a855f7;
}

.diff-badge {
  display: inline-block;
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: capitalize;
}
.badge-easy {
  background: #065f46;
  color: #a7f3d0;
}
.badge-medium {
  background: #78350f;
  color: #fde68a;
}
.badge-hard {
  background: #7f1d1d;
  color: #fca5a5;
}

/* Title & Prompt */
.q-title {
  font-weight: 700;
  color: #f1f5f9;
  margin-bottom: 0.25rem;
}
.q-prompt-preview {
  color: #94a3b8;
  font-size: 0.8rem;
  max-width: 320px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Content previews */
.col-content {
  max-width: 260px;
}
.meta-label {
  font-size: 0.7rem;
  color: #64748b;
  display: block;
  margin-bottom: 0.2rem;
}
.preview-list {
  list-style: none;
  margin: 0;
  padding: 0;
  font-size: 0.78rem;
  color: #cbd5e1;
}
.preview-list li {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.is-correct-preview {
  color: #34d399;
  font-weight: 700;
}
.more-hint {
  color: #64748b;
  font-style: italic;
  font-size: 0.72rem;
}

.answers-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
}
.ans-tag {
  background: #1e293b;
  border: 1px solid #334155;
  color: #38bdf8;
  font-size: 0.72rem;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-family: monospace;
}

/* Status Button */
.toggle-publish-btn {
  border: none;
  border-radius: 6px;
  padding: 0.35rem 0.65rem;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.15s;
}
.toggle-publish-btn:hover:not(:disabled) {
  opacity: 0.85;
  transform: translateY(-1px);
}
.toggle-publish-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.status-published {
  background: #064e3b;
  color: #34d399;
  border: 1px solid #10b981;
}
.status-draft {
  background: #451a03;
  color: #fb923c;
  border: 1px solid #f97316;
}

/* Action buttons */
.col-actions {
  white-space: nowrap;
}
.col-actions .btn-xs {
  margin-right: 0.35rem;
}

/* Confirm modal */
.confirm-modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(4, 9, 20, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 250;
}
.confirm-modal-card {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 1.5rem 2rem;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.8);
}
.confirm-modal-card h3 {
  margin-top: 0;
  color: #f87171;
}
.confirm-modal-card p {
  color: #cbd5e1;
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
}
.confirm-modal-actions {
  display: flex;
  justify-content: center;
  gap: 0.75rem;
}

.btn-primary {
  background: linear-gradient(90deg, #2563eb, #7c3aed);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
}
.btn-secondary {
  background: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
}
.btn-danger {
  background: linear-gradient(90deg, #7f1d1d, #ef4444);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
}
.btn-sm {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}
.btn-xs {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  border-radius: 5px;
}
</style>
