SET lock_timeout = '1s';
SET statement_timeout = '5s';
ALTER TABLE snippets ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;
