-- +goose Up
CREATE SEQUENCE IF NOT EXISTS user_id_seq START 100;

CREATE TABLE IF NOT EXISTS users (
    user_id INT PRIMARY KEY DEFAULT nextval('user_id_seq'),
    user_name VARCHAR(20) UNIQUE NOT NULL,
    email_id VARCHAR(255) UNIQUE NOT NULL,
    country_code VARCHAR(5) NOT NULL,
    phone_no VARCHAR(15) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;