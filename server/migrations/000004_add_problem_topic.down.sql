DROP INDEX IF EXISTS idx_problems_topic;
DROP INDEX IF EXISTS idx_problems_lookup;
CREATE INDEX idx_problems_lookup ON problems (difficulty, type, is_published);

ALTER TABLE problems DROP COLUMN IF EXISTS topic;
