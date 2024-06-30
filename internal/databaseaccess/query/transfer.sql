-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount,
  currency,
  status,
  message
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE 
    from_account_id = $1
    AND
    ($2::bigint IS NULL OR $2::bigint = 0 OR to_account_id = $2::bigint)
ORDER BY id asc
LIMIT $3
OFFSET $4;

-- name: UpdateTransferStatus :one
UPDATE transfers
SET
  status = COALESCE($1, status),
  message = COALESCE($2, message)
WHERE id = $3
RETURNING *;