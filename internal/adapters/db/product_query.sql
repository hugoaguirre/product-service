-- name: GetProduct :one
SELECT * FROM products
WHERE id = ? LIMIT 1;
