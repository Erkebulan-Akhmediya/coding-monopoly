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
  max-width: 480px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.8), 0 0 0 1px rgba(96, 165, 250, 0.15);
  color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  text-align: left;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.1rem 1.5rem;
  border-bottom: 1px solid #1e293b;
  background: #0c1527;
}
.modal-header h2 {
  font-size: 1.15rem;
  font-weight: 800;
  margin: 0;
  color: #f1f5f9;
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
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}
.form-input:focus {
  border-color: #60a5fa;
  box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.2);
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
  color: #64748b;
  line-height: 1.4;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 0.25rem;
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
.btn-secondary:hover:not(:disabled) {
  background: #334155;
}
.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
