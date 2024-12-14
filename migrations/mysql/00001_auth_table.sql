-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE TABLE IF NOT EXISTS tokens (
    id BIGINT UNSIGNED PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    access_token VARCHAR(255) NOT NULL, -- hash
    refresh_token VARCHAR(255) NOT NULL, -- hash
    access_expires_at TIMESTAMP(3) NOT NULL,
    refresh_expires_at TIMESTAMP(3) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX tokens_refresh_token_idx ON tokens (refresh_token);

CREATE TABLE IF NOT EXISTS password_resets (
    id BIGINT UNSIGNED PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token VARCHAR(255) NOT NULL, -- hash
    expires_at TIMESTAMP(3) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX password_resets_token_idx ON password_resets (token);

-- +goose Down
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;
