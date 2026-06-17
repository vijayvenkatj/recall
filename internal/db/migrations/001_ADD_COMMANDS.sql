-- +goose Up
CREATE TABLE commands (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,

    timestamp INTEGER NOT NULL,
    command TEXT NOT NULL,

    cwd TEXT,
    repo TEXT,

    exit_code INTEGER,
    created_at INTEGER NOT NULL
);

CREATE INDEX idx_commands_session_id
ON commands(session_id);

CREATE INDEX idx_commands_repo
ON commands(repo);

CREATE INDEX idx_commands_timestamp
ON commands(timestamp);

-- +goose Down
DROP TABLE commands;