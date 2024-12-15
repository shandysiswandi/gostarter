-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    balance DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    type VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, DEBIT, CREDIT
    status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, PENDING, FAILED, SUCCESS
    remark VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX transactions_user_id_idx ON transactions (user_id);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS topups (
    id BIGINT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    reference_id VARCHAR(255) NOT NULL UNIQUE, -- from payment gateway for top-up
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX topups_transaction_id_idx ON topups (transaction_id);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON topups
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS bills (
    id BIGINT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    reference_id VARCHAR(255) NOT NULL UNIQUE, -- from payment gateway for bills
    type VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN', -- UNKNOWN, PULSA, LISTRIK, INTERNET
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX bills_transaction_id_idx ON bills (transaction_id);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON bills
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS transfers ( -- P2P Transfer
    id BIGINT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL, -- Sender user ID
    recipient_id BIGINT NOT NULL, -- Recipient user ID
    amount DECIMAL(16, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX transfers_sender_id_idx ON transfers (sender_id);
CREATE INDEX transfers_recipient_id_idx ON transfers (recipient_id);
CREATE INDEX transfers_transaction_id_idx ON transfers (transaction_id);

CREATE TRIGGER update_todos_updated_at
BEFORE UPDATE ON transfers
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS topups;
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS transfers;
