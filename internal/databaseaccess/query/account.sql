-- name: CreateAccount :one
INSERT INTO accounts (
  user_id,
  balance,
  currency
) VALUES (
  $1, 0, $2
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByUserId :many
SELECT * FROM accounts
WHERE accounts.user_id = $1 LIMIT 1;

-- name: GetBalance :one
SELECT accounts.balance FROM accounts
WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE accounts.user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;