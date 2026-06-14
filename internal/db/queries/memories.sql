-- name: CreateMemory :one
INSERT INTO memories (
    id,
    session_id,
    title,
    summary,
    created_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: GetMemory :one
SELECT *
FROM memories
WHERE id = ?;

-- name: GetMemoryBySession :one
SELECT *
FROM memories
WHERE session_id = ?;

-- name: ListMemories :many
SELECT *
FROM memories
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateMemory :one
UPDATE memories
SET
    title = ?,
    summary = ?
WHERE id = ?
RETURNING *;

-- name: DeleteMemory :exec
DELETE FROM memories
WHERE id = ?;
