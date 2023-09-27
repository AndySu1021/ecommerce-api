-- name: UpdateProductStockQuantity :exec
UPDATE product_stock
SET quantity   = quantity + ?,
    updated_at = ?
WHERE id = ?
  AND merchant_id = ?;

-- name: GetProductStock :many
SELECT pst.id, spec_1_id, spec_2_id, quantity, `code`, psp1.name AS spec_1_name, IFNULL(psp2.name, '') AS spec_2_name
FROM product_stock pst
         INNER JOIN product_spec psp1 ON pst.spec_1_id = psp1.id
         LEFT JOIN product_spec psp2 ON pst.spec_2_id = psp2.id
WHERE pst.product_id = ?
  AND pst.merchant_id = ?;

-- name: DeleteProductStock :exec
DELETE
FROM product_stock
WHERE id = ?
  AND merchant_id = ?;

-- name: DeleteProductStockByProductID :exec
DELETE
FROM product_stock
WHERE product_id = ?
  AND merchant_id = ?;