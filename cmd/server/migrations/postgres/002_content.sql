-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS editor BOOLEAN NOT NULL DEFAULT FALSE;
CREATE TABLE IF NOT EXISTS content (
 id BIGSERIAL PRIMARY KEY,
 slug TEXT NOT NULL UNIQUE,
 title TEXT NOT NULL,
 summary TEXT NOT NULL DEFAULT '',
 body TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'draft',
 published BOOLEAN NOT NULL DEFAULT FALSE,
 created_at TIMESTAMPTZ NOT NULL,
 updated_at TIMESTAMPTZ NOT NULL,
 created_by BIGINT REFERENCES users(id),
 updated_by BIGINT REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_content_updated_at ON content(updated_at DESC, created_at DESC, id DESC);
CREATE TABLE IF NOT EXISTS content_history (
 id BIGSERIAL PRIMARY KEY,
 content_id BIGINT NOT NULL REFERENCES content(id) ON DELETE CASCADE,
 operation TEXT NOT NULL,
 actor_id BIGINT REFERENCES users(id),
 actor_admin_id BIGINT,
 created_at TIMESTAMPTZ NOT NULL,
 snapshot JSONB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_content_history_content_id ON content_history(content_id, created_at DESC);
-- +goose Down
DROP TABLE IF EXISTS content_history;
DROP TABLE IF EXISTS content;
ALTER TABLE users DROP COLUMN IF EXISTS editor;
ALTER TABLE users DROP COLUMN IF EXISTS role;
