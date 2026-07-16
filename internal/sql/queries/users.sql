-- name: CreateUser :exec
INSERT INTO users (email, password_hash)
VALUES ($1, $2);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUsersByEmail :exec
DELETE FROM users
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE email = $1;