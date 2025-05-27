-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    nickname    TEXT NOT NULL UNIQUE,
    description TEXT,
    avatar_url  TEXT
);

-- +goose Down
DROP TABLE IF EXISTS users; 