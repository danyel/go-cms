-- +goose Up
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);
ALTER TABLE content ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES categories(id);
CREATE INDEX IF NOT EXISTS idx_content_category_id ON content(category_id);

CREATE TABLE IF NOT EXISTS badges (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS content_badges (
    content_id BIGINT NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    badge_id BIGINT NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, badge_id)
);
CREATE INDEX IF NOT EXISTS idx_content_badges_badge_id ON content_badges(badge_id);

-- +goose Down
DROP TABLE IF EXISTS content_badges;
DROP TABLE IF EXISTS badges;
DROP INDEX IF EXISTS idx_content_category_id;
ALTER TABLE content DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS categories;
