-- +goose Up
SELECT 'up SQL query';

CREATE TABLE url_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id UUID NOT NULL,
    is_up BOOLEAN NOT NULL,
    status_code INT,
    response_time_ms BIGINT,
    error TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_url
        FOREIGN KEY (url_id)
        REFERENCES urls(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS urls;
DROP TABLE IF EXISTS url_checks;
