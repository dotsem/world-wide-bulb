ALTER TABLE toggles ADD COLUMN uuid TEXT NOT NULL DEFAULT '';
UPDATE toggles SET uuid = lower(hex(randomblob(16))) WHERE uuid = '' OR uuid IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_toggles_uuid ON toggles(uuid);
