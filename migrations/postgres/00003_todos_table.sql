-- +goose Up
CREATE TABLE IF NOT EXISTS todos (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON todos
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX todos_user_id_idx ON todos (user_id);

-- +goose Down
DROP TABLE IF EXISTS todos;
