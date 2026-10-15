<script lang="ts">
import { defineComponent } from 'vue'
import { adminApiService, type BoardSummary, type CreateRoomInput } from '../services/adminApiService'
import { adminStore } from '../adminStore'

// Mirrors the server's limits (admin.createRoom / room.DefaultTargetXP).
const MAX_ROOM_ID_LENGTH = 64
const TARGET_XP = 500
const MAX_INITIAL_XP = TARGET_XP - 1

export default defineComponent({
  name: 'AdminCreateRoomModal',

  emits: ['created', 'cancel'],

  data() {
    return {
      form: {
        roomID: '' as string,
        topic: '' as string,
        boardID: '' as string,
        initialXP: 0 as number | string,
      },
      topics: [] as string[],
      boards: [] as BoardSummary[],
      loadingOptions: true as boolean,
      validationErrors: [] as string[],
      serverError: '' as string,
      saving: false as boolean,
      maxInitialXP: MAX_INITIAL_XP,
      targetXP: TARGET_XP,
      maxRoomIDLength: MAX_ROOM_ID_LENGTH,
    }
  },

  async mounted() {
    window.addEventListener('keydown', this.handleKeyDown)
    try {
      const [topics, boards] = await Promise.all([
        adminApiService.listTopics(adminStore.token),
        adminApiService.listBoards(adminStore.token),
      ])
      this.topics = topics
      this.boards = boards
      if (topics.length === 1) this.form.topic = topics[0]
      const defaultBoard = boards.find((board) => board.is_default) || boards[0]
      if (defaultBoard) this.form.boardID = defaultBoard.id
    } catch (err: any) {
      this.serverError = err.message || 'Failed to load room options'
    } finally {
      this.loadingOptions = false
    }
  },

  beforeUnmount() {
    window.removeEventListener('keydown', this.handleKeyDown)
  },

  methods: {
    validate(): CreateRoomInput | null {
      const errors: string[] = []
      const roomID = this.form.roomID.trim()
      const initialXP = Number(this.form.initialXP)
      if (!roomID) {
        errors.push(this.$t('admin.createRoomDialog.validation.roomIdRequired'))
      } else if (roomID.length > MAX_ROOM_ID_LENGTH) {
        errors.push(this.$t('admin.createRoomDialog.validation.roomIdTooLong', { max: MAX_ROOM_ID_LENGTH }))
      }
      if (!this.form.topic) {
        errors.push(this.$t('admin.createRoomDialog.validation.topicRequired'))
      }
      if (this.form.initialXP === '' || !Number.isInteger(initialXP) || initialXP < 0 || initialXP > MAX_INITIAL_XP) {
        errors.push(this.$t('admin.createRoomDialog.validation.initialXpRange', { max: MAX_INITIAL_XP }))
      }
      this.validationErrors = errors
      if (errors.length > 0) return null
      return { room_id: roomID, topic: this.form.topic, board_id: this.form.boardID, initial_xp: initialXP }
    },

    async handleSubmit() {
      this.serverError = ''
      const input = this.validate()
      if (!input) return
      this.saving = true
      try {
        const created = await adminApiService.createRoom(adminStore.token, input)
        this.$emit('created', created)
      } catch (err: any) {
        this.serverError = err.message || 'Failed to create room'
      } finally {
        this.saving = false
      }
    },

    handleCancel() {
      if (!this.saving) this.$emit('cancel')
    },

    handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') this.handleCancel()
    },
  },
})
</script>

<template>
  <div class="modal-backdrop" @click.self="handleCancel">
    <div class="modal-card" role="dialog" aria-modal="true" aria-labelledby="create-room-title">
      <header class="modal-header">
        <h2 id="create-room-title">🏠 {{ $t('admin.createRoomDialog.title') }}</h2>
        <button class="close-btn" type="button" :aria-label="$t('common.cancel')" @click="handleCancel">✕</button>
      </header>

      <form class="modal-body" novalidate @submit.prevent="handleSubmit">
        <div v-if="validationErrors.length > 0" class="error-banner" id="create-room-validation-errors">
          <strong>{{ $t('admin.validationErrors') }}</strong>
          <ul>
            <li v-for="(err, i) in validationErrors" :key="i">{{ err }}</li>
          </ul>
        </div>
        <div v-if="serverError" class="error-banner" id="create-room-server-error">
          <strong>{{ $t('admin.serverError') }}</strong> {{ serverError }}
        </div>

        <div class="form-group">
          <label for="new-room-id">{{ $t('admin.createRoomDialog.roomId') }}</label>
          <input
            id="new-room-id"
            v-model="form.roomID"
            type="text"
            class="form-input code-font"
            :maxlength="maxRoomIDLength"
            :placeholder="$t('admin.newRoomPlaceholder')"
            autofocus
          />
          <p class="field-hint">{{ $t('admin.createRoomDialog.roomIdHint') }}</p>
        </div>

        <div class="form-group">
          <label for="new-room-topic">{{ $t('admin.createRoomDialog.topic') }}</label>
          <select
            id="new-room-topic"
            v-model="form.topic"
            class="form-input"
            :disabled="loadingOptions || topics.length === 0"
          >
            <option value="" disabled>
              {{ loadingOptions ? $t('admin.loadingTopics') : $t('admin.selectTopic') }}
            </option>
            <option v-for="topic in topics" :key="topic" :value="topic">{{ topic }}</option>
          </select>
          <p class="field-hint">
            {{ !loadingOptions && topics.length === 0 ? $t('admin.createRoomDialog.noTopics') : $t('admin.createRoomDialog.topicHint') }}
          </p>
        </div>

        <div class="form-group">
          <label for="new-room-board">{{ $t('admin.createRoomDialog.board') }}</label>
          <select
            id="new-room-board"
            v-model="form.boardID"
            class="form-input"
            :disabled="loadingOptions || boards.length === 0"
          >
            <option value="" disabled>
              {{ loadingOptions ? $t('admin.loadingBoards') : $t('admin.selectBoard') }}
            </option>
            <option v-for="board in boards" :key="board.id" :value="board.id">
              {{ board.is_default ? $t('admin.boardOptionDefault', { name: board.name }) : board.name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="new-room-initial-xp">{{ $t('admin.createRoomDialog.initialXp') }}</label>
          <input
            id="new-room-initial-xp"
            v-model.number="form.initialXP"
            type="number"
            step="1"
            min="0"
            :max="maxInitialXP"
            class="form-input"
          />
          <p class="field-hint">
            {{ $t('admin.createRoomDialog.initialXpHint', { max: maxInitialXP, target: targetXP }) }}
          </p>
        </div>

        <footer class="modal-footer">
          <button id="cancel-create-room-btn" type="button" class="btn-secondary" :disabled="saving" @click="handleCancel">
            {{ $t('common.cancel') }}
          </button>
          <button id="create-room-btn" type="submit" class="btn-primary" :disabled="saving || loadingOptions">
            {{ saving ? $t('admin.creating') : $t('admin.createRoomDialog.submit') }}
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
  background: color-mix(in srgb, var(--c-navy-975) 82%, transparent);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 1rem;
}

.modal-card {
  background: linear-gradient(145deg, var(--c-navy-800) 0%, var(--c-gray-900) 100%);
  border: 1px solid var(--c-navy-700);
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 60px color-mix(in srgb, var(--c-shadow) 80%, transparent), 0 0 0 1px color-mix(in srgb, var(--c-blue-400) 15%, transparent);
  color: var(--c-slate-200);
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  text-align: left;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.1rem 1.5rem;
  border-bottom: 1px solid var(--c-slate-800);
  background: var(--c-navy-900);
}
.modal-header h2 {
  font-size: 1.15rem;
  font-weight: 800;
  margin: 0;
  color: var(--c-slate-100);
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--c-slate-500);
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}
.close-btn:hover {
  color: var(--c-slate-100);
  background: var(--c-slate-800);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}

.error-banner {
  background: color-mix(in srgb, var(--c-red-500) 15%, transparent);
  border: 1px solid var(--c-red-500);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  color: var(--c-red-400);
  font-size: 0.85rem;
}
.error-banner ul {
  margin: 0.3rem 0 0 1.2rem;
  padding: 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.form-group label {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--c-slate-400);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.form-input {
  background: var(--c-slate-900);
  border: 1px solid var(--c-slate-700);
  border-radius: 8px;
  padding: 0.6rem 0.85rem;
  color: var(--c-slate-100);
  font-size: 0.9rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}
.form-input:focus {
  border-color: var(--c-blue-400);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--c-blue-400) 20%, transparent);
}
.form-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.code-font {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
.field-hint {
  margin: 0;
  font-size: 0.76rem;
  color: var(--c-slate-500);
  line-height: 1.4;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.25rem;
  padding-top: 1rem;
  border-top: 1px solid var(--c-slate-800);
}

.btn-primary {
  background: linear-gradient(90deg, var(--c-solid--blue-600), var(--c-solid--violet-600));
  color: var(--c-white);
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
  background: var(--c-slate-800);
  border: 1px solid var(--c-slate-700);
  color: var(--c-slate-200);
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-weight: 600;
  font-size: 0.88rem;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-secondary:hover:not(:disabled) {
  background: var(--c-slate-700);
}
.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
