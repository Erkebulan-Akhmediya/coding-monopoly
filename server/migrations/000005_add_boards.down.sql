-- Only boards whose id matches a game (the migrated originals) can map back
-- onto games; custom boards are dropped.
DELETE FROM board_cells WHERE board_id NOT IN (SELECT id FROM games);

UPDATE board_cells SET params = '{"lap_bonus": 100}'::jsonb WHERE type = 'deploy';

ALTER TABLE board_cells DROP CONSTRAINT IF EXISTS board_cells_cell_index_check;
ALTER TABLE board_cells DROP CONSTRAINT IF EXISTS board_cells_board_id_fkey;
ALTER TABLE board_cells DROP CONSTRAINT IF EXISTS board_cells_pkey;
ALTER TABLE board_cells RENAME COLUMN board_id TO game_id;
ALTER TABLE board_cells ADD CONSTRAINT board_cells_pkey PRIMARY KEY (game_id, cell_index);
ALTER TABLE board_cells ADD CONSTRAINT board_cells_game_id_fkey
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
CREATE INDEX idx_board_cells_game_id ON board_cells (game_id);

DROP INDEX IF EXISTS idx_boards_single_default;
DROP TABLE IF EXISTS boards;
