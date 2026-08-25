-- Revert content localization columns back to single-language originals.

ALTER TABLE problem_accepted_answers DROP COLUMN IF EXISTS answer_text_kz;
ALTER TABLE problem_accepted_answers DROP COLUMN IF EXISTS answer_text_ru;
ALTER TABLE problem_accepted_answers RENAME COLUMN answer_text_en TO answer_text;

ALTER TABLE problem_options DROP COLUMN IF EXISTS text_kz;
ALTER TABLE problem_options DROP COLUMN IF EXISTS text_ru;
ALTER TABLE problem_options RENAME COLUMN text_en TO text;

ALTER TABLE problems DROP COLUMN IF EXISTS prompt_kz;
ALTER TABLE problems DROP COLUMN IF EXISTS prompt_ru;
ALTER TABLE problems RENAME COLUMN prompt_en TO prompt;

ALTER TABLE problems DROP COLUMN IF EXISTS title_kz;
ALTER TABLE problems DROP COLUMN IF EXISTS title_ru;
ALTER TABLE problems RENAME COLUMN title_en TO title;

ALTER TABLE board_cells DROP COLUMN IF EXISTS name_kz;
ALTER TABLE board_cells DROP COLUMN IF EXISTS name_ru;
ALTER TABLE board_cells RENAME COLUMN name_en TO name;
