-- +goose Up
CREATE TYPE todo_status AS ENUM (
    'UNKNOWN', 
    'INITIATE', 
    'IN_PROGRESS', 
    'DROP', 
    'DONE'
);

CREATE TABLE IF NOT EXISTS todos (
    id BIGINT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status todo_status NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS todos;
DROP TYPE IF EXISTS todo_status;
