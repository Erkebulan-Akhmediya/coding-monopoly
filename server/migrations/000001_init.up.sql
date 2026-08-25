-- Enable extension for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Define Enums for Problems and Board Cells
CREATE TYPE problem_type AS ENUM ('mcq', 'text');
CREATE TYPE difficulty_level AS ENUM ('easy', 'medium', 'hard');
CREATE TYPE cell_type AS ENUM (
    'deploy',
    'code_freeze',
    'coffee_break',
    'deadline',
    'xp_gain',
    'mystery',
    'xp_loss',
    'double_xp',
    'skip_turn',
    'teleport',
    'free_pass',
    'bonus_challenge'
);

-- 1. games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(50) NOT NULL DEFAULT 'lobby',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. players table
-- (game_id, name, xp, position, status, is_connected)
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    xp INTEGER NOT NULL DEFAULT 0,
    position INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    is_connected BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (game_id, name)
);

-- 3. board_cells table
-- (game_id, cell_index, name_en, name_ru, name_kz, type, params jsonb)
CREATE TABLE board_cells (
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    cell_index INTEGER NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    name_ru VARCHAR(255) NOT NULL,
    name_kz VARCHAR(255) NOT NULL,
    type cell_type NOT NULL,
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (game_id, cell_index)
);

-- 4. problems table
-- (id, type enum[mcq, text], difficulty enum[easy, medium, hard], title_en/ru/kz, prompt_en/ru/kz, is_published bool, created_at, updated_at)
CREATE TABLE problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type problem_type NOT NULL,
    difficulty difficulty_level NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    title_ru VARCHAR(255),
    title_kz VARCHAR(255),
    prompt_en TEXT NOT NULL,
    prompt_ru TEXT,
    prompt_kz TEXT,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, type)
);

-- 5. problem_options table
-- (problem_id, text_en/ru/kz, is_correct) for mcq
CREATE TABLE problem_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL,
    problem_type problem_type NOT NULL DEFAULT 'mcq' CHECK (problem_type = 'mcq'),
    text_en TEXT NOT NULL,
    text_ru TEXT,
    text_kz TEXT,
    is_correct BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (problem_id, problem_type) REFERENCES problems(id, type) ON DELETE CASCADE
);

-- 6. problem_accepted_answers table
-- (problem_id, answer_text) for text
CREATE TABLE problem_accepted_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL,
    problem_type problem_type NOT NULL DEFAULT 'text' CHECK (problem_type = 'text'),
    answer_text TEXT NOT NULL,
    FOREIGN KEY (problem_id, problem_type) REFERENCES problems(id, type) ON DELETE CASCADE
);

-- 7. submissions table
-- (player_id, problem_id, payload jsonb, verdict, created_at)
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    verdict VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 8. game_events table
-- (game_id, type, payload jsonb, created_at)
CREATE TABLE game_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
-- lookup-by-game
CREATE INDEX idx_players_game_id ON players (game_id);
CREATE INDEX idx_board_cells_game_id ON board_cells (game_id);
CREATE INDEX idx_game_events_game_id ON game_events (game_id);

-- lookup-by-difficulty+type+is_published
CREATE INDEX idx_problems_lookup ON problems (difficulty, type, is_published);

-- lookup-by-problem-for-options-and-answers
CREATE INDEX idx_problem_options_problem_id ON problem_options (problem_id);
CREATE INDEX idx_problem_accepted_answers_problem_id ON problem_accepted_answers (problem_id);

-- lookup-for-submissions
CREATE INDEX idx_submissions_player_id ON submissions (player_id);
CREATE INDEX idx_submissions_problem_id ON submissions (problem_id);
