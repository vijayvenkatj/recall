-- +goose Up
CREATE TABLE memories (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL UNIQUE,

    title TEXT,
    summary TEXT NOT NULL,

    created_at INTEGER NOT NULL,

    FOREIGN KEY(session_id)
        REFERENCES sessions(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_memories_session_id
ON memories(session_id);

-- +goose Down
DROP TABLE memories;