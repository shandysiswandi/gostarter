-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT UNSIGNED PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    balance DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE INDEX accounts_user_id_idx ON accounts (user_id);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT UNSIGNED PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    type VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, DEBIT, CREDIT
    status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, PENDING, FAILED, SUCCESS
    remark VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE INDEX transactions_user_id_idx ON transactions (user_id);

CREATE TABLE IF NOT EXISTS topus (
    id BIGINT UNSIGNED PRIMARY KEY,
    transaction_id BIGINT UNSIGNED NOT NULL,
    reference_id VARCHAR(255) NOT NULL UNIQUE, -- from payment gateway for top-up
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE INDEX topus_transaction_id_idx ON topus (transaction_id);

CREATE TABLE IF NOT EXISTS bills (
    id BIGINT UNSIGNED PRIMARY KEY,
    transaction_id BIGINT UNSIGNED NOT NULL,
    reference_id VARCHAR(255) NOT NULL UNIQUE, -- from payment gateway for bills
    type VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, PULSA, LISTRIK, INTERNET
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE INDEX bills_transaction_id_idx ON bills (transaction_id);

CREATE TABLE IF NOT EXISTS transfers ( -- P2P Transfer
    id BIGINT UNSIGNED PRIMARY KEY,
    transaction_id BIGINT UNSIGNED NOT NULL,
    sender_id BIGINT UNSIGNED NOT NULL, -- Sender user ID
    recipient_id BIGINT UNSIGNED NOT NULL, -- Recipient user ID
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE INDEX transfers_sender_id_idx ON transfers (sender_id);
CREATE INDEX transfers_recipient_id_idx ON transfers (recipient_id);
CREATE INDEX transfers_transaction_id_idx ON transfers (transaction_id);

-- +goose Down
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS topus;
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS transfers;
