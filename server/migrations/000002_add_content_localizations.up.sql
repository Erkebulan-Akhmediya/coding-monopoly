-- Localize user-facing content columns: rename original to *_en and add *_ru / *_kz.

-- board_cells.name
ALTER TABLE board_cells RENAME COLUMN name TO name_en;
ALTER TABLE board_cells ADD COLUMN name_ru VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE board_cells ADD COLUMN name_kz VARCHAR(255) NOT NULL DEFAULT '';

-- problems.title / problems.prompt
ALTER TABLE problems RENAME COLUMN title TO title_en;
ALTER TABLE problems ADD COLUMN title_ru VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN title_kz VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE problems RENAME COLUMN prompt TO prompt_en;
ALTER TABLE problems ADD COLUMN prompt_ru TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN prompt_kz TEXT NOT NULL DEFAULT '';

-- problem_options.text
ALTER TABLE problem_options RENAME COLUMN text TO text_en;
ALTER TABLE problem_options ADD COLUMN text_ru TEXT NOT NULL DEFAULT '';
ALTER TABLE problem_options ADD COLUMN text_kz TEXT NOT NULL DEFAULT '';

-- problem_accepted_answers.answer_text
ALTER TABLE problem_accepted_answers RENAME COLUMN answer_text TO answer_text_en;
ALTER TABLE problem_accepted_answers ADD COLUMN answer_text_ru TEXT NOT NULL DEFAULT '';
ALTER TABLE problem_accepted_answers ADD COLUMN answer_text_kz TEXT NOT NULL DEFAULT '';
