-- name: CreateInventoryItem :exec
INSERT INTO inventory (item_id, player_id, quantity)
VALUES ($1, $2, $3);

-- name: GetInventoryByPlayerID :many
SELECT * FROM inventory
WHERE player_id = $1;

-- name: DeleteInventoryItemByPlayerID :exec
DELETE FROM inventory
WHERE player_id = $1;