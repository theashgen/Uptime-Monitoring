-- name: CreateURL :one
INSERT INTO urls (
    host,
    interval,
    user_id
)
VALUES (
    $1,
    $2,
    $3
) RETURNING id, host, interval;

-- name: ListURLsByUser :many
SELECT host, interval FROM urls
WHERE user_id = $1;

