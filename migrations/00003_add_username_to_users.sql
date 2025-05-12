-- +goose Up
ALTER TABLE users ADD COLUMN username VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users DROP COLUMN username;