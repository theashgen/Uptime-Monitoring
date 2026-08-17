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

-- name: GetDueURLs :many
SELECT
    id,
    url,
    interval_seconds,
    next_check_at,
    is_active,
    user_id,
    created_at
FROM urls
WHERE is_active = true
  AND next_check_at <= NOW()
ORDER BY next_check_at
LIMIT $1;

-- name: CreateURLCheck :one
INSERT INTO url_checks (
    url_id,
    is_up,
    status_code,
    response_time_ms,
    error
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;
