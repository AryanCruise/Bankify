package api

import (
	"Accounts/common"
	database "Accounts/internal/datab"
	"Accounts/notifications"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

var kafkaProducer sarama.SyncProducer

func Initializekafka(brokerList []string) {
	var err error
	kafkaProducer, err = notifications.InitializeProducer(brokerList)
	if err != nil {
		log.Fatalf("Failed to Initialize the Kafka Producer: %v", err)
	}
}

type DepositWithdrawRequest struct {
	AccountID int32   `json:"account_id"`
	Amount    float64 `json:"amount"`
	PIN       int32   `json:"pin"`
}

// deposit handles the deposit of a specified amount into an account.
// @Summary Deposit money into an account
// @Description Deposits a specified amount into the account with the given account ID.
// @Security BearerAuth
// @Tags Transactions
// @Accept json
// @Produce json
// @Param body body DepositWithdrawRequest true "Account ID and deposit amount"
// @Success 200 {object} map[string]interface{} "Deposit successful"
// @Failure 400 {string} string "Invalid payload or account does not exist"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to update balance or transaction error"
// @Router /transactions/deposit [post]
func deposit(w http.ResponseWriter, r *http.Request, queries *database.Queries) {

	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req DepositWithdrawRequest


	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Payload error"+err.Error(), http.StatusBadRequest)
		return
	}

	account, err := queries.GetAccount(context.Background(), req.AccountID)
	if account.UserID != userID {
		http.Error(w, "Unauthorized: You do not have access to this account", http.StatusUnauthorized)
		return
	}
	if req.PIN!= account.CreatePin{
		http.Error(w, "Incorrect PIN ", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Account does not exist", http.StatusBadRequest)
		return
	}
	newBalance, err := queries.UpdateBalance(context.Background(), database.UpdateBalanceParams{
		AccountID:      req.AccountID,
		AccountBalance: req.Amount,
	})
	if err != nil {
		http.Error(w, "Failed to update balance"+err.Error(), http.StatusInternalServerError)
	}

	transaction, err := queries.CreateTransaction(context.Background(), database.CreateTransactionParams{
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		TransactionType: "DEPOSIT",
	})
	if err != nil {
		http.Error(w, "Transaction failed"+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"AccountHolder Name": account.UserName,
		"PreviousBalance":    account.AccountBalance,
		"Transaction":        transaction,
		"NewBalance":         newBalance,
	})

	message := fmt.Sprintf("Deposit transaction completed for Account: %v. \nAmount: %.2f. \nNew Balance: %.2f.", account.AccountID, req.Amount, newBalance)
	err = notifications.SendNotification(kafkaProducer, userID, "transaction-notifications", message)
	if err != nil {
		log.Printf("Failed to send deposit notifications: %v", err)
	}
}

// withdraw handles the withdrawal of a specified amount from an account.
// @Summary Withdraw money from an account
// @Description Withdraws a specified amount from the account with the given account ID.
// @Security BearerAuth
// @Tags Transactions
// @Accept json
// @Produce json
// @Param body body DepositWithdrawRequest true "Account ID and withdrawal amount"
// @Success 200 {object} map[string]interface{} "Withdrawal successful"
// @Failure 400 {string} string "Invalid payload, insufficient funds, or account does not exist"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to update balance or transaction error"
// @Router /transactions/withdraw [post]
func withdraw(w http.ResponseWriter, r *http.Request, queries *database.Queries) {

	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req DepositWithdrawRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Payload error", http.StatusBadRequest)
		return
	}

	account, err := queries.GetAccount(context.Background(), req.AccountID)
	if account.UserID != userID {
		http.Error(w, "Unauthorized: You do not have access to this account", http.StatusUnauthorized)
		return
	}
	if req.PIN!= account.CreatePin{
		http.Error(w, "Incorrect PIN ", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Invalid Payload error", http.StatusBadRequest)
		return
	}
	if account.AccountBalance < req.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	newBalance, err := queries.UpdateBalance(context.Background(), database.UpdateBalanceParams{
		AccountID:      req.AccountID,
		AccountBalance: -req.Amount,
	})

	if err != nil {
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	transaction, err := queries.CreateTransaction(context.Background(), database.CreateTransactionParams{
		AccountID:       req.AccountID,
		Amount:          -req.Amount,
		TransactionType: "WITHDRAW",
	})
	if err != nil {
		http.Error(w, "Transaction failed"+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"AccountHolder Name": account.UserName,
		"PreviousBalance":    account.AccountBalance,
		"NewBalance":         newBalance,
		"Transaction":        transaction,
	})

	message := fmt.Sprintf("Withdraw transaction completed for Account: %v. \nAmount: %.2f. \nNew Balance: %.2f.", account.AccountID, req.Amount, newBalance)
	err = notifications.SendNotification(kafkaProducer, userID, "transaction-notifications", message)
	if err != nil {
		log.Printf("Failed to send deposit notifications: %v", err)
	}
}

type TransferRequest struct {
	SenderId   int32   `json:"sender_id"`
	RecieverId int32   `json:"reciever_id"`
	Amount     float64 `json:"amount"`
	PIN       int32   `json:"pin"`
}

// transferMoney handles the transfer of a specified amount from one account to another.
// @Summary Transfer money between accounts
// @Description Transfers a specified amount from the sender's account to the receiver's account.
// @Security BearerAuth
// @Tags Transactions
// @Accept json
// @Produce json
// @Param body body TransferRequest true "Sender ID, Receiver ID, and transfer amount"
// @Success 200 {object} map[string]interface{} "Transfer successful"
// @Failure 400 {string} string "Invalid payload, insufficient funds, or account does not exist"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Transaction commit failed or database error"
// @Router /transfer [post]
func transferMoney(w http.ResponseWriter, r *http.Request, queries *database.Queries, db *sql.DB) {

	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req TransferRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//Beginning a new transaction   Context.bg is used when you don't want to ppass any info like deadline
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		http.Error(w, "Could not start transaction", http.StatusInternalServerError)
	}

	queries = queries.WithTx(tx)

	senderAccount, err := queries.GetAccount(context.Background(), req.SenderId)
	if senderAccount.UserID != userID {
		tx.Rollback()
		http.Error(w, "Unauthorized: You do not have access to this account", http.StatusUnauthorized)
		return
	}
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to get account", http.StatusBadRequest)
		return
	}
	if req.PIN!= senderAccount.CreatePin{
		http.Error(w, "Incorrect PIN ", http.StatusBadRequest)
		return
	}
	if senderAccount.AccountBalance < req.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	updatedSenderBalance, err := queries.UpdateBalance(context.Background(), database.UpdateBalanceParams{
		AccountID:      req.SenderId,
		AccountBalance: -req.Amount,
	})
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	recieverAccount, err := queries.GetAccount(context.Background(), req.RecieverId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to get account", http.StatusBadRequest)
		return
	}

	updatedRecieverBalance, err := queries.UpdateBalance(context.Background(), database.UpdateBalanceParams{
		AccountID:      req.RecieverId,
		AccountBalance: req.Amount,
	})
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}
	_, err = queries.CreateTransaction(context.Background(), database.CreateTransactionParams{
		AccountID:       req.SenderId,
		Amount:          -req.Amount,
		TransactionType: "TRANSFER_OUT",
	})
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error recording sender's transaction"+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = queries.CreateTransaction(context.Background(), database.CreateTransactionParams{
		AccountID:       req.RecieverId,
		Amount:          req.Amount,
		TransactionType: "TRANSFER_IN",
	})

	if err != nil {
		tx.Rollback()
		http.Error(w, "Error recording sender's transaction"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Message":                  "Transfer Successful",
		"Sender ID":                senderAccount.AccountID,
		"Sender Name":              senderAccount.UserName,
		"Updated Sender Balance":   updatedSenderBalance,
		"AmountTransfered":         req.Amount,
		"Reciever Id":              recieverAccount.AccountID,
		"Reciever Name":            recieverAccount.UserName,
		"Updated Reciever Balance": updatedRecieverBalance,
	})

	message := fmt.Sprintf("Transfer transaction completed for Account: %v to %v.\nAmount: %.2f. \nPrevious Balance of Sender: %.2f.\nNew balance of Sender: %.2f.", senderAccount.AccountID, recieverAccount.AccountID, req.Amount, senderAccount.AccountBalance, updatedSenderBalance)
	err = notifications.SendNotification(kafkaProducer, userID, "transaction-notifications", message)
	if err != nil {
		log.Printf("Failed to send deposit notifications: %v", err)
	}
}
