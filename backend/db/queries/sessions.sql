-- name: CreateSession :one
INSERT INTO sessions (device_id) VALUES ($1) RETURNING *;

-- name: ListSessionsByDevice :many
SELECT s.*, COUNT(m.id) as message_count
FROM sessions s LEFT JOIN messages m ON s.id = m.session_id
WHERE s.device_id = $1
GROUP BY s.id
ORDER BY s.updated_at DESC;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1 AND device_id = $2;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1 AND device_id = $2;

-- name: UpdateSessionTitle :one
UPDATE sessions SET title = $1, updated_at = NOW()
WHERE id = $2 AND device_id = $3 RETURNING *;
