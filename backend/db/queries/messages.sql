-- name: CreateMessage :one
INSERT INTO messages (
    session_id, role, content, tool_type, tool_payload, image_path
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMessagesBySession :many
SELECT * FROM messages WHERE session_id = $1 ORDER BY created_at ASC;