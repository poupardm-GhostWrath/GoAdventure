-- name: GetItems :many
SELECT * FROM items;

-- name: GetItemsByName :one
SELECT * FROM items
WHERE name = $1;

-- name: GetItemsByCategoryID :many
SELECT * FROM items
WHERE category_id = $1;