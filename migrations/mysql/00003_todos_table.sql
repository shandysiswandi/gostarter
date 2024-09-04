-- +goose Up
CREATE TABLE IF NOT EXISTS `todos` (
    `id` BIGINT UNSIGNED PRIMARY KEY,
    `title` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `status` VARCHAR(50) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS todos;
