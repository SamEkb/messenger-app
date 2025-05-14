-- +goose Up
CREATE TABLE IF NOT EXISTS friendships
(
    id           UUID PRIMARY KEY,
    requestor_id TEXT                     NOT NULL,
    recipient_id TEXT                     NOT NULL,
    status       TEXT                     NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE (requestor_id, recipient_id)
);

CREATE INDEX IF NOT EXISTS idx_friendships_requestor_id ON friendships (requestor_id);
CREATE INDEX IF NOT EXISTS idx_friendships_recipient_id ON friendships (recipient_id);

-- +goose Down
DROP TABLE IF EXISTS friendships; 