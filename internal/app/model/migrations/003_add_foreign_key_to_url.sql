-- +goose Up

ALTER TABLE url
    ADD COLUMN user_id INTEGER NOT NULL;

ALTER TABLE url
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(user_id);


-- +goose Down
ALTER TABLE url
DROP COLUMN fk_user_id;





