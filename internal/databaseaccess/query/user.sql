-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  phone,
  role,
  is_email_verified
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUsersByUsernameOrPhone :many
SELECT accounts.id as account_id, users.username, users.full_name, users.phone
FROM users JOIN accounts ON users.id = accounts.user_id
WHERE users.username like $1
OR users.phone like $2
ORDER BY account_id
LIMIT $3;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  email = COALESCE(sqlc.narg(email), email),
  phone = COALESCE(sqlc.narg(phone), phone)
WHERE
  username = sqlc.arg(username)
RETURNING *;

-- name: VerifyEmail :one
UPDATE users
SET
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
  username = sqlc.arg(username)
RETURNING *;
