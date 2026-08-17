-- name: CreateURL :one
INSERT INTO urls (
    url,
    interval_seconds,
    user_id
)
VALUES (
    $1,
    $2,
    $3
) RETURNING id, url, interval_seconds;

-- name: ListURLsByUser :many
SELECT url, interval_seconds FROM urls
WHERE user_id = $1;

