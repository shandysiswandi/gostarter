-- +goose Up
CREATE TABLE IF NOT EXISTS todos (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL
);

CREATE INDEX todos_user_id_idx ON todos (user_id);

-- +goose Down
DROP TABLE IF EXISTS todos;
