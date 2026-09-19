-- +goose Up
CREATE TABLE IF NOT EXISTS users (id BIGSERIAL PRIMARY KEY, email TEXT NOT NULL, name TEXT NOT NULL DEFAULT '', provider TEXT NOT NULL, provider_subject TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_subject ON users(provider, provider_subject);
CREATE TABLE IF NOT EXISTS sessions (id BIGSERIAL PRIMARY KEY, token TEXT NOT NULL UNIQUE, user_id BIGINT REFERENCES users(id), admin_id BIGINT, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
-- +goose Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
