-- name: CreateAdmin :exec
INSERT INTO admin(merchant_id, email, password, real_name, mobile, sex, is_enabled, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: AdminLogin :one
SELECT *
FROM admin
WHERE merchant_id = ?
  AND email = ?
  AND password = ?;