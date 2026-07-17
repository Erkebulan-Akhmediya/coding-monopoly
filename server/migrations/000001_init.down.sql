-- Drop indexes
DROP INDEX IF EXISTS idx_problems_lookup;
DROP INDEX IF EXISTS idx_game_events_game_id;
DROP INDEX IF EXISTS idx_board_cells_game_id;
DROP INDEX IF EXISTS idx_players_game_id;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS game_events;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS problem_accepted_answers;
DROP TABLE IF EXISTS problem_options;
DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS board_cells;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS games;

-- Drop types and extensions
DROP TYPE IF EXISTS difficulty_level;
DROP TYPE IF EXISTS problem_type;
DROP EXTENSION IF EXISTS "uuid-ossp";
