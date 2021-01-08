-- name: CreateUser :one
INSERT INTO users
(
    email,
    phone,
    full_name,
    hashed_password
) VALUES
( $1, $2, $3, $4 )
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;