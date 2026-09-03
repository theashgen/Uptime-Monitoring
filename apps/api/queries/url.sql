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
SELECT id, url, interval_seconds, next_check_at, is_active, created_at FROM urls
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetURLByID :one
SELECT id, url, interval_seconds, next_check_at, is_active, user_id, created_at
FROM urls
WHERE id = $1 AND user_id = $2;

-- name: ListURLChecksByURL :many
SELECT id, url_id, is_up, status_code, response_time_ms, error, checked_at
FROM url_checks
WHERE url_id = $1
ORDER BY checked_at DESC
LIMIT $2;

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

