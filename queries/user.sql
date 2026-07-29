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

-- name: GetUserEmail :one
SELECT username, passwordhash FROM users
WHERE email = $1;

-- name: GetUserUsername :one
SELECT * FROM users
WHERE email = $1
AND passwordhash = $2;
