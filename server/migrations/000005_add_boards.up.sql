-- Boards are reusable 32-cell layouts an admin picks when creating a room.
CREATE TABLE boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL UNIQUE,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- At most one board is the default (preselected for new rooms, never deletable).
CREATE UNIQUE INDEX idx_boards_single_default ON boards (is_default) WHERE is_default;

-- Each game's existing cell set becomes a board with the same id; the seeded
-- default game becomes the default board.
INSERT INTO boards (id, name, is_default)
SELECT game_id,
       CASE WHEN game_id = '00000000-0000-0000-0000-000000000000' THEN 'Default'
            ELSE 'Board ' || left(game_id::text, 8) END,
       game_id = '00000000-0000-0000-0000-000000000000'
FROM (SELECT DISTINCT game_id FROM board_cells) AS existing;

-- Re-key board_cells from games to boards.
ALTER TABLE board_cells DROP CONSTRAINT board_cells_game_id_fkey;
ALTER TABLE board_cells DROP CONSTRAINT board_cells_pkey;
DROP INDEX IF EXISTS idx_board_cells_game_id;
ALTER TABLE board_cells RENAME COLUMN game_id TO board_id;
ALTER TABLE board_cells ADD CONSTRAINT board_cells_pkey PRIMARY KEY (board_id, cell_index);
ALTER TABLE board_cells ADD CONSTRAINT board_cells_board_id_fkey
    FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE;
ALTER TABLE board_cells ADD CONSTRAINT board_cells_cell_index_check
    CHECK (cell_index >= 0 AND cell_index < 32);

-- Deploy's landing bonus was stored as an unused lap_bonus; store the keys the
-- engine reads (landing bonus and pass-GO bonus) with their current values.
UPDATE board_cells SET params = '{"amount": 100, "pass_bonus": 50}'::jsonb WHERE type = 'deploy';
