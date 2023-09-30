-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id serial PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS users;