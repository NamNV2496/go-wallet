-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: UpdateExpiredTime :one
UPDATE sessions
SET
  expires_at = $1
WHERE
  id = $2
RETURNING *;
