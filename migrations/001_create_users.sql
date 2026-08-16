-- +goose Up
SELECT 'up SQL query';

CREATE TABLE "users" (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL, 
    passwordhash TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;

SELECT 'down SQL query';
