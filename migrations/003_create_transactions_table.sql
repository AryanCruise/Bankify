-- +goose Up
CREATE SEQUENCE IF NOT EXISTS transaction_id_seq START 1000;

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id INT PRIMARY KEY DEFAULT nextval('transaction_id_seq'),
    account_id INT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type VARCHAR(15) CHECK (transaction_type IN ('DEPOSIT', 'WITHDRAW', 'TRANSFER_IN', 'TRANSFER_OUT')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP SEQUENCE IF EXISTS transaction_id_seq;