-- name: CreateProductCategory :execresult
INSERT INTO product_category (merchant_id, name, top_id, parent_id, tree_left, tree_right, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProductCategory :one
SELECT *
FROM product_category
WHERE id = ?
  AND merchant_id = ?;

-- name: UpdateProductCategoryLeftTree :exec
UPDATE product_category
SET tree_right = tree_right + 2,
    updated_at = ?
WHERE tree_left < ?
  and tree_right >= ?
  and top_id = ?
  AND merchant_id = ?;

-- name: UpdateProductCategoryRightTree :exec
UPDATE product_category
SET tree_right = tree_right + 2,
    tree_left  = tree_left + 2,
    updated_at = ?
WHERE tree_left > ?
  and top_id = ?
  AND merchant_id = ?;

-- name: UpdateProductCategory :exec
UPDATE product_category
SET name       = ?,
    updated_at = ?
WHERE id = ?
  AND merchant_id = ?;

-- name: DeleteProductCategory :exec
DELETE
FROM product_category
WHERE top_id = ?
  and tree_left >= ?
  and tree_right <= ?
  AND merchant_id = ?;

-- name: UpdateProductCategoryTopID :exec
UPDATE product_category
SET top_id = ?
where id = ?
  AND merchant_id = ?;

-- name: GetProductCategoryTopCount :one
SELECT COUNT(*)
FROM product_category
WHERE parent_id = 0
  AND merchant_id = ?;

-- name: GetProductCountByCategoryID :one
SELECT COUNT(*)
FROM product
WHERE category_id IN (SELECT a.id
                      FROM product_category AS a,
                           (SELECT id, top_id, tree_left, tree_right
                            FROM product_category
                            WHERE product_category.id = ? AND product_category.merchant_id = ?) AS b
                      WHERE a.tree_left >= b.tree_left
                        AND a.tree_right <= b.tree_right
                        AND a.top_id = b.top_id);

-- name: GetProductCategoryChildrenIDs :many
SELECT a.id
FROM product_category AS a,
     (SELECT * FROM product_category WHERE product_category.id = ? AND product_category.merchant_id = ?) AS b
WHERE a.tree_left >= b.tree_left
  AND a.tree_right <= b.tree_right
  AND a.top_id = b.top_id;