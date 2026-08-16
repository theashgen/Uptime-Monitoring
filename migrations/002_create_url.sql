-- +goose Up
SELECT 'up SQL query';
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    url TEXT NOT NULL,

    interval_seconds INT NOT NULL,

    next_check_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    is_active BOOLEAN NOT NULL DEFAULT true,

    user_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_url
        UNIQUE (user_id, url),

    CONSTRAINT positive_interval
        CHECK (interval_seconds > 0)
);
-- +goose Down
DROP TABLE IF EXISTS urls;
DROP TABLE IF EXISTS users;
SELECT 'down SQL query';
