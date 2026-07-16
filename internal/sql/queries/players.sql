-- name: CreatePlayer :one
INSERT INTO players (name, user_id)
VALUES ($2, $1)
RETURNING *;

-- name: GetPlayersByUserID :many
SELECT id, name FROM players
WHERE user_id = $1;

-- name: GetPlayersByID :one
SELECT * FROM players
WHERE id = $1;

-- name: DeletePlayersByID :exec
DELETE FROM players
WHERE id = $1;

-- name: UpdatePlayerByID :exec
UPDATE players
SET current_exp = $2, current_level = $3, gold = $4, updated_at = NOW()
WHERE id = $1;