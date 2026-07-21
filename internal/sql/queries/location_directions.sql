-- name: GetLocationDirectionByID :many
SELECT * FROM location_directions
WHERE location_id = $1;