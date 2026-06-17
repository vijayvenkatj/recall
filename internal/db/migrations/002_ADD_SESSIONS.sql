-- +goose Up
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,

    repo TEXT NOT NULL,

    start_ts INTEGER NOT NULL,
    end_ts INTEGER NOT NULL,

    command_count INTEGER NOT NULL DEFAULT 0,

    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_sessions_repo
ON sessions(repo);

CREATE INDEX idx_sessions_end_ts
ON sessions(end_ts);

-- +goose Down
DROP TABLE sessions;