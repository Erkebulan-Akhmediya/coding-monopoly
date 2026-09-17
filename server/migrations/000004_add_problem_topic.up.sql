-- Add required topic to problems for room-scoped question banks.
ALTER TABLE problems
    ADD COLUMN topic VARCHAR(128);

UPDATE problems SET topic = 'Go' WHERE topic IS NULL;

ALTER TABLE problems
    ALTER COLUMN topic SET NOT NULL;

DROP INDEX IF EXISTS idx_problems_lookup;
CREATE INDEX idx_problems_lookup ON problems (topic, difficulty, type, is_published);
CREATE INDEX idx_problems_topic ON problems (topic);
