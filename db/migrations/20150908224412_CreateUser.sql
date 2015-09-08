-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_email on users (email);

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;

