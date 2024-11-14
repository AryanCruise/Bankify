package api

import (
	"Accounts/common"
	database "Accounts/internal/datab"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreateAccountRequest struct {
	UserName       string  `json:"user_name"`       // User's name
	AccountBalance float64 `json:"account_balance"` // Initial account balance
	AccountType    string  `json:"account_type"`    // Type of the account
	CreatePin      int32   `json:"create_pin"`
}

// createAccount creates a new account
// @Summary Create an account
// @Description Create a new account for the user
// @Security BearerAuth
// @Tags Accounts
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Account details"
// @Success 200 {object} database.Account
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [post]
func createAccount(w http.ResponseWriter, r *http.Request, queries *database.Queries) {

	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req CreateAccountRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Payload error", http.StatusBadRequest)
		return
	}

	account, err := queries.CreateAccount(context.Background(), database.CreateAccountParams{
		UserID:         userID,
		UserName:       req.UserName,
		AccountBalance: req.AccountBalance,
		AccountType:    req.AccountType,
		CreatePin:      req.CreatePin,
	})
	if err != nil {
		http.Error(w, "Failed to create account"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// getAccount retrieves an account
// @Summary Get account details
// @Description Retrieve account details by ID
// @Security BearerAuth
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} database.Account
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [get]
func getAccount(w http.ResponseWriter, r *http.Request, queries *database.Queries) {

	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusInternalServerError)
		return
	}

	getInfo, err := queries.GetAccount(context.Background(), int32(id))
	if err != nil {
		// Check if the error is due to no rows found
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure the account belongs to the user
	if getInfo.UserID != userID {
		http.Error(w, "Unauthorized: You do not have access to this account", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getInfo)
}

// @Summary Get account balance
// @Description Retrieve the balance of a specific account
// @Security BearerAuth
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {number} float64
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id}/balance [get]
func getBalance(w http.ResponseWriter, r *http.Request, queries *database.Queries) {
	userID, err := common.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid accountID", http.StatusInternalServerError)
		return
	}
	account, err := queries.GetAccount(context.Background(), int32(id))
	if err != nil || account.UserID != userID { // Checks -> account belongs to the user
		http.Error(w, "Unauthorized: You do not have access to this account", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account.AccountBalance)
}

// @Summary Delete an account
// @Description Delete an account by ID
// @Security BearerAuth
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [delete]
func deleteAccount(w http.ResponseWriter, r *http.Request, queries *database.Queries) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid accountID", http.StatusInternalServerError)
		return
	}
	err = queries.DeleteAccount(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
}

// derefValue is a generic function that dereferences a pointer and returns a default value if the pointer is nil.
func derefValue[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr // Dereference the pointer if it's not nil
	}
	return defaultVal // Return the default value if the pointer is nil
}

type UpdateAccountRequest struct {
	UserName       *string  `json:"user_name"`
	AccountBalance *float64 `json:"account_balance"`
	AccountType    *string  `json:"account_type"`
}

// @Summary Update an account
// @Description Update account details by ID
// @Security BearerAuth
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param request body UpdateAccountRequest true "Account data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{id} [put]
func updateAccount(w http.ResponseWriter, r *http.Request, queries *database.Queries) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	currentAccount, err := queries.GetAccount(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Failed to retrieve account", http.StatusInternalServerError)
		return
	}

	var req UpdateAccountRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = queries.UpdateAccount(context.Background(), database.UpdateAccountParams{
		AccountID:      int32(id),
		UserName:       derefValue(req.UserName, currentAccount.UserName),
		AccountBalance: derefValue(req.AccountBalance, currentAccount.AccountBalance),
		AccountType:    derefValue(req.AccountType, currentAccount.AccountType),
	})
	if err != nil {
		http.Error(w, "Failed to update account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Your account has been updated"})
}
