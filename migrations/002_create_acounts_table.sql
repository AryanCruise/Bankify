-- +goose Up
CREATE SEQUENCE IF NOT EXISTS account_id_seq START 500;

CREATE TABLE IF NOT EXISTS accounts (
    account_id INT PRIMARY KEY DEFAULT nextval('account_id_seq'),
    user_id INT NOT NULL DEFAULT nextval('user_id_seq'),
    user_name VARCHAR(100) NOT NULL,
    account_balance DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    account_type VARCHAR(15) CHECK (account_type IN ('CURRENT', 'SAVINGS', 'FIXED DEPOSIT', 'SALARY')) NOT NULL,
    create_pin INT NOT NULL CHECK (create_pin >= 1000 AND create_pin <= 9999),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS accounts;
DROP SEQUENCE IF EXISTS account_id_seq;
DROP SEQUENCE IF EXISTS user_id_seq;