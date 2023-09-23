-- +goose Up
CREATE TABLE IF NOT EXISTS url (
    id serial PRIMARY KEY,
    shorted VARCHAR(20),
    original TEXT UNIQUE,
    deleted BOOLEAN DEFAULT FAlSE
);

-- +goose Down
DROP TABLE IF EXISTS url;