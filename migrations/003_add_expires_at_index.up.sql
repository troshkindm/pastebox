-- squawk-ignore require-lock-timeout, require-statement-timeout
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_snippets_expires_at ON snippets (expires_at) WHERE expires_at IS NOT NULL;
