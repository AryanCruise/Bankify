CREATE SEQUENCE user_id_seq START 100;
CREATE TABLE IF NOT EXISTS users (
    user_id INT PRIMARY KEY DEFAULT nextval('user_id_seq'),
    user_name VARCHAR(20) UNIQUE NOT NULL,
    email_id VARCHAR(255) UNIQUE NOT NULL,
    country_code VARCHAR(5) NOT NULL,
    phone_no VARCHAR(15) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  -- Stores the hashed password here
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE account_id_seq START 500;
CREATE TABLE IF NOT EXISTS accounts (
    account_id INT PRIMARY KEY DEFAULT nextval('account_id_seq'),
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    user_name VARCHAR(100) NOT NULL,
    account_balance DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    account_type VARCHAR(15) CHECK (account_type IN ('SALARY', 'SAVINGS', 'FIXED DEPOSIT', 'CURRENT')) NOT NULL,
    create_pin INT NOT NULL CHECK (create_pin >= 1000 AND create_pin <= 9999),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE transaction_id_seq START 1000;
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id INT PRIMARY KEY DEFAULT nextval('transaction_id_seq'),
    account_id INT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type VARCHAR(20) CHECK (transaction_type IN ('DEPOSIT','WITHDRAWAL', 'TRANSFER_IN', 'TRANSFER_OUT')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,  -- Link to user in accounts
    token TEXT NOT NULL UNIQUE,  -- Stores the refresh token as a unique string
    expires_at TIMESTAMP NOT NULL,  -- Expiration timestamp for the token
    revoked BOOLEAN DEFAULT FALSE,  -- Flag to mark as revoked (e.g., upon logout)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp for token creation
);
