-- +goose Up
SELECT 'up SQL query';

CREATE TABLE "users" (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL, 
    passwordhash TEXT NOT NULL
);

CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host TEXT NOT NULL UNIQUE,
    interval TEXT NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES "users"(id)
        ON DELETE CASCADE
);


-- +goose Down
DROP TABLE urls;
DROP TABLE users;

SELECT 'down SQL query';
