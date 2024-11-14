-- name: CreateAccount :one
INSERT INTO accounts (user_id, user_name, account_balance, account_type, create_pin) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE account_id = $1::INT;

-- name: GetBalance :one
SELECT account_balance FROM accounts WHERE account_id = $1::INT;

-- name: CreateTransaction :one
INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE account_id = $1;

-- name: UpdateAccount :exec
UPDATE accounts
SET 
    user_name = COALESCE($2, user_name),
    account_balance = COALESCE($3, account_balance),
    account_type = COALESCE($4, account_type)
WHERE account_id = $1;

-- name: UpdateBalance :one
UPDATE accounts
SET account_balance = account_balance + $2
WHERE account_id = $1
RETURNING account_balance;

-- name: CreateUser :one
INSERT INTO users (user_name, email_id, country_code, phone_no, password_hash) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByUserName :one
SELECT user_id, user_name, email_id, phone_no, country_code, password_hash FROM users WHERE user_name = $1;

-- name: GetUserByUserID :one
SELECT user_id, user_name, email_id, phone_no, country_code, password_hash FROM users WHERE user_id = $1;