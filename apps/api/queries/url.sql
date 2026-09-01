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

-- name: ClaimDueURLs :many
WITH due AS (
    SELECT id
    FROM urls
    WHERE is_active = true
      AND next_check_at <= NOW()
    ORDER BY next_check_at
    FOR UPDATE SKIP LOCKED
    LIMIT $1
)
UPDATE urls u
SET next_check_at = NOW() + (u.interval_seconds * INTERVAL '1 second')
FROM due
WHERE u.id = due.id
RETURNING
    u.id,
    u.url,
    u.interval_seconds,
    u.next_check_at,
    u.is_active,
    u.user_id,
    u.created_at;

-- name: CreateURLChecks :copyfrom
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
);

-- name: UpdateURLNextCheck :exec
UPDATE urls
SET next_check_at = NOW() + (interval_seconds * INTERVAL '1 second')
WHERE id = $1;

