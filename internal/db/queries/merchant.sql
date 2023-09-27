-- name: GetMerchantByHost :one
SELECT *
FROM merchant
WHERE host = ?;

-- name: GetMerchant :one
SELECT *
FROM merchant
WHERE id = ?;