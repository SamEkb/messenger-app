-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    username TEXT NOT NULL,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens
(
    token      TEXT PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;