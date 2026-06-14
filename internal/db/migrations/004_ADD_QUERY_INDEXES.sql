-- +goose Up
CREATE INDEX idx_commands_session_timestamp_created
ON commands(session_id, timestamp DESC, created_at DESC);

CREATE INDEX idx_commands_repo_timestamp_created
ON commands(repo, timestamp DESC, created_at DESC);

CREATE INDEX idx_commands_timestamp_created
ON commands(timestamp DESC, created_at DESC);

CREATE INDEX idx_sessions_repo_end_updated
ON sessions(repo, end_ts DESC, updated_at DESC);

CREATE INDEX idx_sessions_end_updated
ON sessions(end_ts DESC, updated_at DESC);

CREATE INDEX idx_memories_created_at
ON memories(created_at DESC);

-- +goose Down
DROP INDEX idx_memories_created_at;
DROP INDEX idx_sessions_end_updated;
DROP INDEX idx_sessions_repo_end_updated;
DROP INDEX idx_commands_timestamp_created;
DROP INDEX idx_commands_repo_timestamp_created;
DROP INDEX idx_commands_session_timestamp_created;
