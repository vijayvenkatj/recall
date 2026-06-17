-- name: UpsertMetadata :one
INSERT INTO metadata (
    key,
    value
)
VALUES (
    ?,
    ?
)
ON CONFLICT(key) DO UPDATE SET value = excluded.value
RETURNING *;

-- name: GetMetadata :one
SELECT * FROM metadata WHERE key = ?;
