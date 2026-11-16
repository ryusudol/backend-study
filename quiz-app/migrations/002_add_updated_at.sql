ALTER TABLE notes ADD COLUMN updated_at TIMESTAMPTZ;

UPDATE notes SET updated_at = created_at;

ALTER TABLE notes ALTER COLUMN updated_at SET NOT NULL;
