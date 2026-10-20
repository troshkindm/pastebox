ALTER TABLE snippets ADD COLUMN title text NOT NULL;
ALTER TABLE snippets DROP COLUMN content;
CREATE INDEX idx_snippets_created_at ON snippets (created_at);
