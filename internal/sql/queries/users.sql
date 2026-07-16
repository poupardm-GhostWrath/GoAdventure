-- name: CreateUser :exec
INSERT INTO users (email, password_hash)
VALUES ($1, $2);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;