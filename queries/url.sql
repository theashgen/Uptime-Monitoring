-- name: CreateURL :one
INSERT INTO urls (
    host,
    interval
)
VALUES (
    $1,
    $2
) RETURNING *;

-- name: ListURLsByUser :many
SELECT host, interval FROM urls
WHERE user_id = $1;