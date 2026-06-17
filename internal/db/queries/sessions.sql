-- name: CreateSession :one
INSERT INTO sessions (
    id,
    repo,
    start_ts,
    end_ts,
    command_count,
    created_at,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = ?;

-- name: ListSessionsByRepo :many
SELECT *
FROM sessions
WHERE repo = ?
ORDER BY end_ts DESC, updated_at DESC
LIMIT ? OFFSET ?;

-- name: ListRecentSessions :many
SELECT *
FROM sessions
ORDER BY end_ts DESC, updated_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateSession :one
UPDATE sessions
SET
    repo = ?,
    start_ts = ?,
    end_ts = ?,
    command_count = ?,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: TouchSessionForCommand :one
UPDATE sessions
SET
    end_ts = ?,
    command_count = command_count,
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;
