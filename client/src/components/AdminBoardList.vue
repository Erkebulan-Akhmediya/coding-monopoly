<script lang="ts">
import { defineComponent } from 'vue'
import {
  adminApiService,
  type Board,
  type BoardSummary,
  type CellTypeCatalog,
} from '../services/adminApiService'
import { adminStore } from '../adminStore'
import AdminBoardEditor from './AdminBoardEditor.vue'

export default defineComponent({
  name: 'AdminBoardList',

  components: {
    AdminBoardEditor,
  },

  data() {
    return {
      boards: [] as BoardSummary[],
      catalog: null as CellTypeCatalog | null,
      loading: false as boolean,
      error: '' as string,
      // Board open in the editor; a draft without an id is a new board.
      editing: null as Pick<Board, 'id' | 'name' | 'cells'> | null,
      busyBoardId: '' as string,
      confirmDeleteBoard: null as BoardSummary | null,
      deleting: false as boolean,
    }
  },

  async mounted() {
    await this.loadAll()
  },

  methods: {
    async loadAll() {
      this.loading = true
      this.error = ''
      try {
        const [boards, catalog] = await Promise.all([
          adminApiService.listBoards(adminStore.token),
          this.catalog ? Promise.resolve(this.catalog) : adminApiService.listCellTypes(adminStore.token),
        ])
        this.boards = boards
        this.catalog = catalog
      } catch (err: any) {
        this.error = err.message || 'Failed to load boards'
      } finally {
        this.loading = false
      }
    },

    async openEditor(board: BoardSummary) {
      await this.withBoard(board, (full) => {
        this.editing = full
      })
    },

    async openCopy(board: BoardSummary, name: string) {
      await this.withBoard(board, (full) => {
        this.editing = { id: '', name, cells: full.cells }
      })
    },

    async createNew() {
      const source = this.boards.find((b) => b.is_default) || this.boards[0]
      if (source) await this.openCopy(source, '')
    },

    async withBoard(board: BoardSummary, open: (full: Board) => void) {
      this.busyBoardId = board.id
      this.error = ''
      try {
        open(await adminApiService.getBoard(adminStore.token, board.id))
      } catch (err: any) {
        this.error = err.message || 'Failed to load board'
      } finally {
        this.busyBoardId = ''
      }
    },

    async makeDefault(board: BoardSummary) {
      this.busyBoardId = board.id
      this.error = ''
      try {
        this.boards = await adminApiService.setDefaultBoard(adminStore.token, board.id)
      } catch (err: any) {
        this.error = err.message || 'Failed to set default board'
      } finally {
        this.busyBoardId = ''
      }
    },

    async confirmDelete() {
      if (!this.confirmDeleteBoard) return
      this.deleting = true
      try {
        await adminApiService.deleteBoard(adminStore.token, this.confirmDeleteBoard.id)
        this.confirmDeleteBoard = null
        await this.loadAll()
      } catch (err: any) {
        this.error = err.message || 'Failed to delete board'
        this.confirmDeleteBoard = null
      } finally {
        this.deleting = false
      }
    },

    async onSaved() {
      this.editing = null
      await this.loadAll()
    },

    formatDate(value: string): string {
      return new Date(value).toLocaleString()
    },
  },
})
</script>

<template>
  <AdminBoardEditor
    v-if="editing && catalog"
    :board="editing"
    :catalog="catalog"
    @saved="onSaved"
    @cancel="editing = null"
  />

  <div v-else class="board-list-container">
    <div class="list-header">
      <div class="header-left">
        <h2 class="section-heading">🎲 {{ $t('admin.boards.heading') }}</h2>
        <span class="count-badge">{{ $t('admin.boards.count', { count: boards.length }) }}</span>
      </div>
      <button
        id="create-board-btn"
        class="btn-primary btn-sm"
        :disabled="loading || boards.length === 0 || !catalog"
        @click="createNew"
      >
        {{ $t('admin.boards.newBoard') }}
      </button>
    </div>

    <p class="list-subtitle">{{ $t('admin.boards.subtitle') }}</p>

    <div v-if="error" class="error-box">
      <span>{{ error }}</span>
      <button class="btn-secondary btn-xs" @click="loadAll">{{ $t('admin.retry') }}</button>
    </div>

    <div v-if="loading && boards.length === 0" class="loading-state">
      {{ $t('admin.boards.loading') }}
    </div>

    <div v-else-if="boards.length === 0" class="empty-state">
      <p>{{ $t('admin.boards.empty') }}</p>
    </div>

    <div v-else class="boards-grid">
      <div
        v-for="board in boards"
        :key="board.id"
        class="board-card"
        :class="{ 'is-default': board.is_default }"
        :id="`board-card-${board.id}`"
      >
        <div class="board-card-top">
          <span class="board-name">{{ board.name }}</span>
          <span v-if="board.is_default" class="default-badge">{{ $t('admin.boards.defaultBadge') }}</span>
        </div>
        <span class="board-updated">{{ $t('admin.boards.updated', { date: formatDate(board.updated_at) }) }}</span>
        <div class="board-actions">
          <button
            class="btn-primary btn-xs"
            :disabled="busyBoardId === board.id || !catalog"
            @click="openEditor(board)"
          >
            {{ $t('admin.boards.edit') }}
          </button>
          <button
            class="btn-secondary btn-xs"
            :disabled="busyBoardId === board.id || !catalog"
            @click="openCopy(board, $t('admin.boards.copyName', { name: board.name }))"
          >
            {{ $t('admin.boards.duplicate') }}
          </button>
          <button
            v-if="!board.is_default"
            class="btn-secondary btn-xs"
            :disabled="busyBoardId === board.id"
            @click="makeDefault(board)"
          >
            {{ $t('admin.boards.makeDefault') }}
          </button>
          <button
            v-if="!board.is_default"
            class="btn-danger btn-xs"
            :disabled="busyBoardId === board.id"
            @click="confirmDeleteBoard = board"
          >
            {{ $t('admin.boards.delete') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="confirmDeleteBoard" class="confirm-modal-backdrop">
      <div class="confirm-modal-card">
        <h3>{{ $t('admin.boards.deleteConfirmTitle') }}</h3>
        <p>{{ $t('admin.boards.deleteConfirmBody', { name: confirmDeleteBoard.name }) }}</p>
        <div class="confirm-modal-actions">
          <button class="btn-secondary" :disabled="deleting" @click="confirmDeleteBoard = null">
            {{ $t('common.cancel') }}
          </button>
          <button id="confirm-delete-board-btn" class="btn-danger" :disabled="deleting" @click="confirmDelete">
            <span v-if="deleting">{{ $t('admin.deleting') }}</span>
            <span v-else>{{ $t('admin.yesDelete') }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.board-list-container {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
  padding: 1.5rem;
  background: #080f1e;
  color: #e2e8f0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  min-height: calc(100vh - 49px);
  box-sizing: border-box;
  text-align: left;
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
.list-subtitle {
  margin: 0;
  color: #64748b;
  font-size: 0.85rem;
}

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
.loading-state,
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #64748b;
  background: #0c1527;
  border: 1px solid #1e293b;
  border-radius: 12px;
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}
.board-card {
  background: #0c1527;
  border: 1px solid #1e293b;
  border-left: 3px solid #60a5fa;
  border-radius: 12px;
  padding: 1rem 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}
.board-card.is-default {
  border-left-color: #34d399;
}
.board-card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
}
.board-name {
  font-weight: 800;
  font-size: 0.95rem;
  color: #f1f5f9;
  overflow-wrap: anywhere;
}
.default-badge {
  font-size: 0.65rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
  background: #064e3b;
  color: #34d399;
  white-space: nowrap;
}
.board-updated {
  font-size: 0.75rem;
  color: #64748b;
}
.board-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  padding-top: 0.5rem;
  border-top: 1px solid #1e293b;
}

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
.btn-primary:disabled,
.btn-secondary:disabled,
.btn-danger:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.btn-sm {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}
.btn-xs {
  padding: 0.25rem 0.55rem;
  font-size: 0.75rem;
  border-radius: 5px;
}
</style>
