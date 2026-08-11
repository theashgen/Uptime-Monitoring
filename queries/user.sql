-- name: CreateUser :one
INSERT INTO users (
    email,
    passwordhash,
    username
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING email, username;

-- name: GetUserByEmail :one
SELECT id, username, passwordhash FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT id, email, username FROM users
WHERE username = $1;