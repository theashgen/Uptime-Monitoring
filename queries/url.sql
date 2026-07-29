-- name: CreateURL :one
INSERT INTO urls (
    actualUrl,
    shortCode,
    user_id
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;