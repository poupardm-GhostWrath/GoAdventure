-- name: DeleteItems :exec
DELETE FROM items;

-- name: CreateItems :one
INSERT INTO items (name, description, category_id, effect_target, effect_value, value)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetItemsByName :one
SELECT * FROM items
WHERE name = $1;

-- name: GetItemsByCategoryID :many
SELECT * FROM items
WHERE category_id = $1;