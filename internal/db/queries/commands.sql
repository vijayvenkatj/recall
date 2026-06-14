-- name: CreateCommand :one
INSERT INTO commands (
    id,
    session_id,
    timestamp,
    command,
    cwd,
    repo,
    branch,
    exit_code,
    created_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetCommand :one
SELECT *
FROM commands
WHERE id = ?;

-- name: ListCommandsBySession :many
SELECT *
FROM commands
WHERE session_id = ?
ORDER BY timestamp DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: ListCommandsByRepo :many
SELECT *
FROM commands
WHERE repo = ?
ORDER BY timestamp DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: ListRecentCommands :many
SELECT *
FROM commands
ORDER BY timestamp DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: ListCommandsInTimeRange :many
SELECT *
FROM commands
WHERE timestamp >= ?
  AND timestamp <= ?
ORDER BY timestamp DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: DeleteCommand :exec
DELETE FROM commands
WHERE id = ?;
