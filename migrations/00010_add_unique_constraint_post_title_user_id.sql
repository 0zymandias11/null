-- +goose Up
ALTER TABLE posts ADD CONSTRAINT posts_title_user_id_key UNIQUE (title, user_id);

-- +goose Down
ALTER TABLE posts DROP CONSTRAINT posts_title_user_id_key;