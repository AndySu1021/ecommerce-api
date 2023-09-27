-- name: CreateProduct :execresult
INSERT INTO product (merchant_id, `name`, category_id, currency_id, price, special_price, special_price_start,
                     special_price_end, single_order_limit, is_single_order_only, temperature, length, width, height,
                     weight, support_delivery_method, is_air_contraband, `description`, pictures, extra, created_at,
                     updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProduct :one
SELECT p.*, pc.name AS category_name
FROM product p
         INNER JOIN product_category pc ON pc.id = p.category_id
WHERE p.id = ?
  AND p.merchant_id = ?;

-- name: UpdateProduct :exec
UPDATE product
SET `name`                  = ?,
    category_id             = ?,
    `description`           = ?,
    price                   = ?,
    special_price           = ?,
    special_price_start     = ?,
    special_price_end       = ?,
    single_order_limit      = ?,
    is_single_order_only    = ?,
    temperature             = ?,
    length                  = ?,
    width                   = ?,
    height                  = ?,
    weight                  = ?,
    support_delivery_method = ?,
    is_air_contraband       = ?,
    pictures                = ?,
    updated_at              = ?
WHERE id = ?
  AND merchant_id = ?;

-- name: DeleteProduct :exec
DELETE
FROM product
WHERE id = ?
  AND merchant_id = ?;

-- name: SwitchProductStatus :exec
UPDATE product
SET is_enabled = !is_enabled
WHERE id = ?
  AND merchant_id = ?;

-- name: UpdateProductSales :exec
UPDATE product
SET sales = sales + ?
WHERE id = ?;