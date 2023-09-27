-- name: CreateProductSpec :exec
INSERT INTO product_spec (merchant_id, product_id, `level`, type, `name`, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetProductSpec :many
SELECT id, `level`, type, `name`
FROM product_spec
WHERE product_id = ?
  AND merchant_id = ?
ORDER BY id;

-- name: GetProductSpecByID :one
SELECT *
FROM product_spec
WHERE id = ?
  AND merchant_id = ?;

-- name: UpdateProductSpec :exec
UPDATE product_spec
SET `name`     = ?,
    updated_at = ?
WHERE id = ?
  AND merchant_id = ?;

-- name: DeleteProductSpec :exec
DELETE
FROM product_spec
WHERE id = ?
  AND merchant_id = ?;

-- name: DeleteProductSpecByProductID :exec
DELETE
FROM product_spec
WHERE product_id = ?
  AND merchant_id = ?;

-- name: GetProductSpecTitlesByProductID :many
SELECT name
FROM product_spec
WHERE product_id = ?
  AND merchant_id = ?
  AND type = 1;

-- name: ListProductSecondLevelSpecIDs :many
SELECT id
FROM product_spec
WHERE product_id = ?
  AND merchant_id = ?
  AND `level` = 2
  AND `type` = 2;

-- name: DeleteProductSecondLevelSpec :exec
DELETE
FROM product_spec
WHERE product_id = ?
  AND merchant_id = ?
  AND `level` = 2
  AND `type` = 2;