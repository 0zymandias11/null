-- +goose_up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_posts_user_id ON posts (user_id, username);
CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX idx_post_likes_post_id ON post_likes (post_id);    

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_posts_title ON posts USING gin (title gin_trgm_ops);