// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type Account struct {
	AccountID      int32
	UserID         int32
	UserName       string
	AccountBalance float64
	AccountType    string
	CreatePin      int32
	CreatedAt      time.Time
}

type RefreshToken struct {
	TokenID   int32
	UserID    int32
	Token     string
	ExpiresAt time.Time
	Revoked   sql.NullBool
	CreatedAt sql.NullTime
}

type Transaction struct {
	TransactionID   int32
	AccountID       int32
	Amount          float64
	TransactionType string
	CreatedAt       sql.NullTime
}

type User struct {
	UserID       int32
	UserName     string
	EmailID      string
	CountryCode  string
	PhoneNo      string
	PasswordHash string
	CreatedAt    time.Time
}