<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import {
  adminApiService,
  defaultCellParams,
  validateBoardInput,
  type Board,
  type BoardCell,
  type BoardInput,
  type CellTypeCatalog,
  type CellTypeSpec,
  type LocalizedText,
} from '../services/adminApiService'
import { CORNER_INDEXES, cellGridPosition, cellIcon, cellParamSummary } from '../services/boardLayout'
import { adminStore } from '../adminStore'
import { localeText } from '../i18n'

type TextLang = keyof LocalizedText

function copyCells(cells: BoardCell[]): BoardCell[] {
  return cells
    .map((cell) => ({
      cell_index: cell.cell_index,
      type: cell.type,
      name: { en: cell.name?.en || '', ru: cell.name?.ru || '', kz: cell.name?.kz || '' },
      explanation: {
        en: cell.explanation?.en || '',
        ru: cell.explanation?.ru || '',
        kz: cell.explanation?.kz || '',
      },
      params: { ...(cell.params || {}) },
    }))
    .sort((a, b) => a.cell_index - b.cell_index)
}

export default defineComponent({
  name: 'AdminBoardEditor',

  props: {
    // A stored board, or a draft without an id (new board / duplicate).
    board: {
      type: Object as PropType<Pick<Board, 'id' | 'name' | 'cells'>>,
      required: true,
    },
    catalog: {
      type: Object as PropType<CellTypeCatalog>,
      required: true,
    },
  },

  emits: ['saved', 'cancel'],

  data() {
    const cells = copyCells(this.board.cells)
    return {
      form: {
        name: this.board.name,
        cells,
      },
      initialSnapshot: JSON.stringify({ name: this.board.name, cells }),
      selectedIndex: 0 as number,
      textLang: 'en' as TextLang,
      textLangs: ['en', 'ru', 'kz'] as TextLang[],
      validationErrors: [] as Array<{ cell: number; message: string }>,
      serverError: '' as string,
      saving: false as boolean,
    }
  },

  computed: {
    isNew(): boolean {
      return !this.board.id
    },
    selectedCell(): BoardCell {
      return this.form.cells[this.selectedIndex]
    },
    selectedSpec(): CellTypeSpec | undefined {
      return this.specFor(this.selectedCell.type)
    },
    isStartCell(): boolean {
      return this.selectedIndex === 0
    },
    typeOptions(): CellTypeSpec[] {
      // The start cell is always Deploy, and Deploy can't go anywhere else.
      return this.catalog.cell_types.filter((spec) => spec.start_only === this.isStartCell)
    },
    teleportTarget(): number | null {
      return this.selectedCell.type === 'teleport' ? this.selectedCell.params.target_position : null
    },
    invalidCells(): Set<number> {
      return new Set(this.validationErrors.map((err) => err.cell).filter((cell) => cell >= 0))
    },
    typeCounts(): Array<{ type: string; count: number }> {
      const counts: Record<string, number> = {}
      for (const cell of this.form.cells) {
        counts[cell.type] = (counts[cell.type] || 0) + 1
      }
      return this.catalog.cell_types
        .filter((spec) => counts[spec.type])
        .map((spec) => ({ type: spec.type, count: counts[spec.type] }))
    },
    isDirty(): boolean {
      return JSON.stringify(this.form) !== this.initialSnapshot
    },
  },

  watch: {
    // Once errors are shown, keep them in sync with edits so fixed ones clear.
    form: {
      deep: true,
      handler() {
        if (this.validationErrors.length > 0) this.validate()
      },
    },
  },

  methods: {
    cellIcon,
    cellParamSummary,

    specFor(type: string): CellTypeSpec | undefined {
      return this.catalog.cell_types.find((spec) => spec.type === type)
    },

    isCorner(idx: number): boolean {
      return CORNER_INDEXES.includes(idx)
    },

    gridStyle(idx: number) {
      return cellGridPosition(idx)
    },

    selectCell(idx: number) {
      this.selectedIndex = idx
    },

    onTypeChange() {
      const cell = this.selectedCell
      cell.params = defaultCellParams(this.specFor(cell.type))
      if (cell.type === 'teleport' && cell.params.target_position === cell.cell_index) {
        cell.params.target_position = 0
      }
    },

    paramLabel(key: string): string {
      return this.$t(`admin.boards.param.${this.selectedCell.type}.${key}`)
    },

    teleportOptionLabel(cell: BoardCell): string {
      return `#${cell.cell_index} — ${this.cellTitle(cell)}`
    },

    cellTitle(cell: BoardCell): string {
      return localeText(cell.name) || this.$t(`effects.${cell.type}`)
    },

    payload(): BoardInput {
      return { name: this.form.name.trim(), cells: this.form.cells }
    },

    validate(): boolean {
      this.validationErrors = validateBoardInput(this.payload(), this.catalog).map((err) => ({
        cell: err.cell,
        message: this.$t(`admin.boards.validation.${err.key}`, err.params || {}),
      }))
      return this.validationErrors.length === 0
    },

    async handleSave() {
      this.serverError = ''
      if (!this.validate()) {
        const firstCell = this.validationErrors.find((err) => err.cell >= 0)
        if (firstCell) this.selectedIndex = firstCell.cell
        return
      }

      const payload = this.payload()
      this.saving = true
      try {
        const saved = this.isNew
          ? await adminApiService.createBoard(adminStore.token, payload)
          : await adminApiService.updateBoard(adminStore.token, this.board.id, payload)
        this.initialSnapshot = JSON.stringify(this.form)
        this.$emit('saved', saved)
      } catch (err: any) {
        this.serverError = err.message || 'Failed to save board.'
      } finally {
        this.saving = false
      }
    },

    handleCancel() {
      if (this.isDirty && !window.confirm(this.$t('admin.boards.discardConfirm'))) return
      this.$emit('cancel')
    },
  },
})
</script>

<template>
  <div class="board-editor">
    <!-- Header: back, title, board name, actions -->
    <div class="editor-header">
      <button id="board-editor-back-btn" class="back-btn" type="button" @click="handleCancel">
        ← {{ $t('admin.boards.back') }}
      </button>
      <h2 class="section-heading">
        🎲 {{ isNew ? $t('admin.boards.newTitle') : $t('admin.boards.editTitle') }}
      </h2>
      <div class="name-field">
        <label for="board-name">{{ $t('admin.boards.boardName') }}</label>
        <input
          id="board-name"
          v-model="form.name"
          type="text"
          maxlength="128"
          class="form-input"
          :placeholder="$t('admin.boards.boardNamePlaceholder')"
        />
      </div>
      <div class="header-actions">
        <button class="btn-secondary btn-sm" type="button" @click="handleCancel">
          {{ $t('common.cancel') }}
        </button>
        <button
          id="save-board-btn"
          class="btn-primary btn-sm"
          type="button"
          :disabled="saving"
          @click="handleSave"
        >
          <span v-if="saving">{{ $t('admin.boards.saving') }}</span>
          <span v-else>{{ isNew ? $t('admin.boards.create') : $t('admin.boards.save') }}</span>
        </button>
      </div>
    </div>

    <div v-if="validationErrors.length > 0" class="error-banner" id="board-validation-errors">
      <strong>{{ $t('admin.validationErrors') }}</strong>
      <ul>
        <li
          v-for="(err, i) in validationErrors"
          :key="i"
          :class="{ clickable: err.cell >= 0 }"
          @click="err.cell >= 0 && selectCell(err.cell)"
        >
          {{ err.message }}
        </li>
      </ul>
    </div>
    <div v-if="serverError" class="error-banner" id="board-server-error">
      <strong>{{ $t('admin.serverError') }}</strong> {{ serverError }}
    </div>

    <div class="editor-body">
      <!-- Board preview: click a cell to edit it -->
      <div class="editor-grid">
        <button
          v-for="cell in form.cells"
          :key="cell.cell_index"
          type="button"
          class="editor-cell"
          :class="[
            'cell-' + cell.type,
            {
              corner: isCorner(cell.cell_index),
              selected: cell.cell_index === selectedIndex,
              invalid: invalidCells.has(cell.cell_index),
              'teleport-target': cell.cell_index === teleportTarget,
            },
          ]"
          :style="gridStyle(cell.cell_index)"
          :title="cellTitle(cell)"
          :id="`board-editor-cell-${cell.cell_index}`"
          @click="selectCell(cell.cell_index)"
        >
          <span class="editor-cell-top">
            <span class="editor-cell-idx">#{{ cell.cell_index }}</span>
            <span class="editor-cell-icon">{{ cellIcon(cell.type) }}</span>
          </span>
          <span class="editor-cell-name">{{ cellTitle(cell) }}</span>
          <span v-if="cellParamSummary(cell)" class="editor-cell-param">{{ cellParamSummary(cell) }}</span>
        </button>

        <div class="editor-grid-center">
          <p class="center-hint">{{ $t('admin.boards.clickCellHint') }}</p>
          <p class="center-note">{{ $t('admin.boards.fixedLayout') }}</p>
          <div class="type-counts">
            <span v-for="item in typeCounts" :key="item.type" class="type-chip" :class="'chip-' + item.type">
              {{ cellIcon(item.type) }} {{ $t(`effects.${item.type}`) }} × {{ item.count }}
            </span>
          </div>
        </div>
      </div>

      <!-- Selected cell editor -->
      <aside class="cell-panel" :class="'cell-' + selectedCell.type">
        <div class="cell-panel-header">
          <span class="cell-panel-icon">{{ cellIcon(selectedCell.type) }}</span>
          <h3>{{ $t('admin.boards.cellTitle', { index: selectedIndex }) }}</h3>
        </div>

        <div class="form-group">
          <label for="cell-type">{{ $t('admin.boards.cellType') }}</label>
          <select
            id="cell-type"
            v-model="selectedCell.type"
            class="form-input"
            :disabled="isStartCell"
            @change="onTypeChange"
          >
            <option v-for="spec in typeOptions" :key="spec.type" :value="spec.type">
              {{ cellIcon(spec.type) }} {{ $t(`effects.${spec.type}`) }}
            </option>
          </select>
          <p class="field-hint">
            {{ isStartCell ? $t('admin.boards.startCellLocked') : $t(`admin.boards.typeHelp.${selectedCell.type}`) }}
          </p>
        </div>

        <div class="subform-section">
          <h4>{{ $t('admin.boards.effectSettings') }}</h4>
          <p v-if="!selectedSpec || selectedSpec.params.length === 0" class="field-hint">
            {{ $t('admin.boards.noSettings') }}
          </p>
          <template v-else>
            <div v-for="param in selectedSpec.params" :key="param.key" class="form-group">
              <template v-if="param.key === 'target_position'">
                <label :for="`cell-param-${param.key}`">{{ $t('admin.boards.teleportTarget') }}</label>
                <select
                  :id="`cell-param-${param.key}`"
                  v-model.number="selectedCell.params[param.key]"
                  class="form-input"
                >
                  <option
                    v-for="target in form.cells"
                    :key="target.cell_index"
                    :value="target.cell_index"
                    :disabled="target.cell_index === selectedIndex"
                  >
                    {{ teleportOptionLabel(target) }}
                  </option>
                </select>
              </template>
              <template v-else>
                <label :for="`cell-param-${param.key}`">
                  {{ paramLabel(param.key) }}
                  <span class="range-hint">{{ $t('admin.boards.range', { min: param.min, max: param.max }) }}</span>
                </label>
                <input
                  :id="`cell-param-${param.key}`"
                  v-model.number="selectedCell.params[param.key]"
                  type="number"
                  step="1"
                  :min="param.min"
                  :max="param.max"
                  class="form-input"
                />
              </template>
            </div>
          </template>
        </div>

        <div class="subform-section">
          <div class="text-header">
            <h4>{{ $t('admin.boards.cellText') }}</h4>
            <div class="lang-tabs">
              <button
                v-for="lang in textLangs"
                :key="lang"
                type="button"
                class="lang-tab"
                :class="{ active: textLang === lang }"
                @click="textLang = lang"
              >
                {{ $t(`admin.boards.langTab.${lang}`) }}
              </button>
            </div>
          </div>
          <div class="form-group">
            <label for="cell-name">{{ $t('admin.boards.name') }}</label>
            <input
              id="cell-name"
              v-model="selectedCell.name[textLang]"
              type="text"
              maxlength="255"
              class="form-input"
              :placeholder="textLang === 'en' ? '' : selectedCell.name.en"
            />
          </div>
          <div class="form-group">
            <label for="cell-explanation">{{ $t('admin.boards.explanation') }}</label>
            <textarea
              id="cell-explanation"
              v-model="selectedCell.explanation[textLang]"
              rows="5"
              maxlength="2000"
              class="form-input"
              :placeholder="textLang === 'en' ? '' : selectedCell.explanation.en"
            ></textarea>
          </div>
          <p class="field-hint">{{ $t('admin.boards.textHint') }}</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.board-editor {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--c-navy-950);
  color: var(--c-slate-200);
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  min-height: calc(100vh - 49px);
  box-sizing: border-box;
  text-align: left;
}

/* Header */
.editor-header {
  display: flex;
  align-items: flex-end;
  gap: 1rem;
  flex-wrap: wrap;
}
.section-heading {
  font-size: 1.25rem;
  font-weight: 800;
  margin: 0 0 0.35rem;
  background: linear-gradient(90deg, var(--c-blue-400), var(--c-violet-400));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.name-field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  flex: 1;
  min-width: 220px;
  max-width: 360px;
}
.name-field label,
.form-group label {
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--c-slate-400);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.header-actions {
  display: flex;
  gap: 0.6rem;
  margin-left: auto;
}
.back-btn {
  background: transparent;
  border: 1px solid var(--c-slate-700);
  color: var(--c-slate-400);
  padding: 0.35rem 0.75rem;
  border-radius: 6px;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  margin-bottom: 0.25rem;
  transition: all 0.2s;
}
.back-btn:hover {
  color: var(--c-slate-100);
  border-color: var(--c-blue-400);
  background: var(--c-navy-800);
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
.error-banner li.clickable {
  cursor: pointer;
  text-decoration: underline dotted;
}

/* Body: grid + side panel */
.editor-body {
  display: flex;
  gap: 1.25rem;
  align-items: flex-start;
  flex-wrap: wrap;
}

.editor-grid {
  display: grid;
  grid-template-columns: repeat(9, minmax(0, 76px));
  grid-template-rows: repeat(9, minmax(0, 72px));
  gap: 3px;
  background: var(--c-slate-800);
  padding: 6px;
  border-radius: 12px;
  border: 1px solid var(--c-slate-700);
  width: max-content;
  max-width: 100%;
}

.editor-cell {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 1px;
  padding: 3px 4px;
  background: var(--c-slate-700);
  border: 1px solid var(--c-slate-600);
  border-radius: 5px;
  color: var(--c-slate-100);
  font: inherit;
  font-size: 12px;
  line-height: 1.2;
  text-align: center;
  cursor: pointer;
  overflow: hidden;
  transition: transform 0.15s, box-shadow 0.15s;
}
.editor-cell:hover {
  transform: scale(1.05);
  z-index: 2;
}
.editor-cell.selected {
  box-shadow: 0 0 0 2px var(--c-amber-400), 0 0 14px color-mix(in srgb, var(--c-amber-400) 50%, transparent);
  z-index: 3;
}
.editor-cell.teleport-target {
  outline: 2px dashed var(--c-teleport--violet-400);
  outline-offset: -3px;
}
.editor-cell.invalid {
  box-shadow: 0 0 0 2px var(--c-red-500);
}
.editor-cell.selected.invalid {
  box-shadow: 0 0 0 2px var(--c-amber-400), 0 0 0 4px var(--c-red-500);
}
.editor-cell-top {
  flex-shrink: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.editor-cell-idx {
  font-size: 0.55rem;
  font-weight: 700;
  color: var(--c-slate-400);
}
.editor-cell-icon {
  font-size: 0.75rem;
}
.editor-cell-name {
  min-height: 0;
  font-size: 0.58rem;
  font-weight: 700;
  line-height: 1.1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.editor-cell-param {
  flex-shrink: 0;
  font-size: 0.55rem;
  font-weight: 700;
  color: var(--c-sky-400);
}

/* Cell type colours (match the player board) */
.editor-cell.corner { border-width: 2px; }
.editor-cell.cell-deploy       { background: linear-gradient(135deg, var(--c-deploy--emerald-800), var(--c-deploy--emerald-700)); border-color: var(--c-deploy--emerald-500); }
.editor-cell.cell-code_freeze  { background: linear-gradient(135deg, var(--c-code_freeze--blue-900), var(--c-code_freeze--blue-700)); border-color: var(--c-code_freeze--blue-500); }
.editor-cell.cell-coffee_break { background: linear-gradient(135deg, var(--c-coffee_break--amber-900), var(--c-coffee_break--amber-700)); border-color: var(--c-coffee_break--amber-500); }
.editor-cell.cell-deadline     { background: linear-gradient(135deg, var(--c-deadline--red-900), var(--c-deadline--red-700)); border-color: var(--c-deadline--red-500); }
.editor-cell.cell-xp_gain           { border-top: 3px solid var(--c-xp_gain--emerald-500); }
.editor-cell.cell-xp_loss           { border-top: 3px solid var(--c-xp_loss--red-500); }
.editor-cell.cell-mystery           { border-top: 3px solid var(--c-mystery--indigo-400); }
.editor-cell.cell-teleport          { border-top: 3px solid var(--c-teleport--violet-500); }
.editor-cell.cell-skip_next         { border-top: 3px solid var(--c-skip_next--orange-500); }
.editor-cell.cell-double_xp         { border-top: 3px solid var(--c-double_xp--yellow-500); }
.editor-cell.cell-free_pass         { border-top: 3px solid var(--c-free_pass--cyan-500); }
.editor-cell.cell-special_challenge { border-top: 3px solid var(--c-special_challenge--pink-500); }

.editor-grid-center {
  grid-row: 2 / span 7;
  grid-column: 2 / span 7;
  background: color-mix(in srgb, var(--c-slate-900) 85%, transparent);
  border: 1px dashed var(--c-slate-600);
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
  text-align: center;
  overflow: auto;
}
.center-hint {
  margin: 0;
  font-weight: 800;
  color: var(--c-blue-300);
}
.center-note {
  margin: 0;
  font-size: 0.78rem;
  color: var(--c-slate-500);
  max-width: 340px;
}
.type-counts {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.3rem;
  max-width: 420px;
}
.type-chip {
  background: var(--c-slate-800);
  border: 1px solid var(--c-slate-700);
  border-radius: 99px;
  padding: 0.1rem 0.5rem;
  font-size: 0.7rem;
  color: var(--c-slate-300);
  white-space: nowrap;
}

/* Side panel */
.cell-panel {
  flex: 1;
  min-width: 300px;
  max-width: 460px;
  background: var(--c-navy-900);
  border: 1px solid var(--c-slate-800);
  border-top: 3px solid var(--c-slate-700);
  border-radius: 12px;
  padding: 1rem 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}
.cell-panel.cell-deploy            { border-top-color: var(--c-deploy--emerald-500); }
.cell-panel.cell-code_freeze       { border-top-color: var(--c-code_freeze--blue-500); }
.cell-panel.cell-coffee_break      { border-top-color: var(--c-coffee_break--amber-500); }
.cell-panel.cell-deadline          { border-top-color: var(--c-deadline--red-500); }
.cell-panel.cell-xp_gain           { border-top-color: var(--c-xp_gain--emerald-500); }
.cell-panel.cell-xp_loss           { border-top-color: var(--c-xp_loss--red-500); }
.cell-panel.cell-mystery           { border-top-color: var(--c-mystery--indigo-400); }
.cell-panel.cell-teleport          { border-top-color: var(--c-teleport--violet-500); }
.cell-panel.cell-skip_next         { border-top-color: var(--c-skip_next--orange-500); }
.cell-panel.cell-double_xp         { border-top-color: var(--c-double_xp--yellow-500); }
.cell-panel.cell-free_pass         { border-top-color: var(--c-free_pass--cyan-500); }
.cell-panel.cell-special_challenge { border-top-color: var(--c-special_challenge--pink-500); }

.cell-panel-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.cell-panel-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 800;
  color: var(--c-slate-100);
}
.cell-panel-icon {
  font-size: 1.3rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.form-input {
  background: var(--c-slate-900);
  border: 1px solid var(--c-slate-700);
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  color: var(--c-slate-100);
  font-size: 0.88rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}
.form-input:focus {
  border-color: var(--c-blue-400);
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--c-blue-400) 20%, transparent);
}
.form-input:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
textarea.form-input {
  resize: vertical;
}
.field-hint {
  margin: 0;
  font-size: 0.76rem;
  color: var(--c-slate-500);
  line-height: 1.4;
}
.range-hint {
  font-weight: 600;
  color: var(--c-slate-600);
  text-transform: none;
  letter-spacing: 0;
  margin-left: 0.35rem;
}

.subform-section {
  background: var(--c-navy-925);
  border: 1px solid var(--c-slate-800);
  border-radius: 10px;
  padding: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}
.subform-section h4 {
  margin: 0;
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--c-slate-400);
}
.text-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.lang-tabs {
  display: flex;
  gap: 0.25rem;
}
.lang-tab {
  background: transparent;
  border: 1px solid var(--c-slate-700);
  color: var(--c-slate-400);
  border-radius: 6px;
  padding: 0.15rem 0.5rem;
  font-size: 0.72rem;
  font-weight: 700;
  cursor: pointer;
}
.lang-tab.active {
  color: var(--c-blue-400);
  border-color: var(--c-blue-500);
  background: var(--c-navy-800);
}

.btn-primary {
  background: linear-gradient(90deg, var(--c-solid--blue-600), var(--c-solid--violet-600));
  color: var(--c-white);
  border: none;
  border-radius: 8px;
  font-weight: 700;
  cursor: pointer;
  transition: opacity 0.2s;
}
.btn-primary:hover:not(:disabled) { opacity: 0.9; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-secondary {
  background: var(--c-slate-800);
  border: 1px solid var(--c-slate-700);
  color: var(--c-slate-200);
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}
.btn-secondary:hover { background: var(--c-slate-700); }
.btn-sm {
  padding: 0.45rem 0.9rem;
  font-size: 0.82rem;
}

@media (max-width: 820px) {
  .editor-grid {
    grid-template-columns: repeat(9, minmax(0, 36px));
    grid-template-rows: repeat(9, minmax(0, 36px));
    gap: 2px;
  }
  .editor-cell-name,
  .editor-cell-param,
  .center-note,
  .type-counts {
    display: none;
  }
  .cell-panel {
    max-width: none;
  }
}
</style>
