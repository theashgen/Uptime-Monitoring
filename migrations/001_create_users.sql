-- +goose Up
SELECT 'up SQL query';

CREATE TABLE "users" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL, 
    passwordhash TEXT NOT NULL
);

CREATE TABLE urls (
    actualUrl TEXT,
    shortCode VARCHAR(50) UNIQUE,
    user_id UUID,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES "users"(id)
        ON DELETE CASCADE
);


-- +goose Down
DROP TABLE urls;
DROP TABLE users;

SELECT 'down SQL query';
