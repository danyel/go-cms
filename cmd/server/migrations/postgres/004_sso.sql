-- +goose Up
-- The CMS no longer stores identities: the upstream SSO proxy owns them.
ALTER TABLE content DROP COLUMN IF EXISTS created_by;
ALTER TABLE content DROP COLUMN IF EXISTS updated_by;
ALTER TABLE content_history DROP COLUMN IF EXISTS actor_id;
ALTER TABLE content_history DROP COLUMN IF EXISTS actor_admin_id;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- +goose Down
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    provider TEXT NOT NULL,
    provider_subject TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT '',
    editor BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_subject ON users(provider, provider_subject);
CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    user_id BIGINT REFERENCES users(id),
    admin_id BIGINT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
ALTER TABLE content ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE content ADD COLUMN IF NOT EXISTS updated_by BIGINT REFERENCES users(id);
ALTER TABLE content_history ADD COLUMN IF NOT EXISTS actor_id BIGINT REFERENCES users(id);
ALTER TABLE content_history ADD COLUMN IF NOT EXISTS actor_admin_id BIGINT;
