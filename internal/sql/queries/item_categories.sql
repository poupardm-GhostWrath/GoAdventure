-- name: GetItemCategoriesByName :one
SELECT * FROM item_categories
WHERE name = $1;