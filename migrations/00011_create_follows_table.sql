-- +goose_up
CREATE TABLE followers (
    user_id BIGINT NOT NULL,
    follower VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose_down
DROP TABLE IF EXISTS followers;