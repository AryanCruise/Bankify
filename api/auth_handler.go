package api

import (
	"Accounts/auth"
	database "Accounts/internal/datab"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	EmailID     string `json:"email_id"`
	CountryCode string `json:"country_code"`
	PhoneNo     string `json:"phone_number"`
}
type loginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
type UserResponse struct {
	Message string `json:"message"`
	UserID  int32  `json:"user_id"`
	EmailID string `json:"email_id"`
	PhoneNo string `json:"phone_no"`
}

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description This endpoint allows a user to register by providing a username and password
// @Tags Home Page
// @Accept  json
// @Produce  json
// @Param registerRequest body registerRequest true "User registration request"
// @Success 201 {object} UserResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 409 {string} string "Username already exists"
// @Router /register [post]
func RegisterHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest

		// Decode the registration request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request"+err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Store the new user in the database
		user, err := queries.CreateUser(r.Context(), database.CreateUserParams{
			UserName:     req.UserName,
			EmailID:      req.EmailID,
			PasswordHash: string(hashedPassword),
			CountryCode:  req.CountryCode,
			PhoneNo:      req.PhoneNo,
		})
		if err != nil {
			http.Error(w, "Username/E-mail already exists", http.StatusConflict)
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(UserResponse{
			Message: "User created successfully",
			UserID:  user.UserID,
			EmailID: user.EmailID,
			PhoneNo: user.PhoneNo,
		})
	}
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// loginHandler handles user login
// @Summary Login a user
// @Description This endpoint allows a user to log in by providing a username and password
// @Tags Home Page
// @Accept  json
// @Produce  json
// @Param loginRequest body loginRequest true "User login request"
// @Success 200 {object} UserLoginResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "User not found. Please register first!" // or "Invalid Password"
// @Failure 500 {string} string "Failed to generate tokens"
// @Router /login [post]
func loginHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request"+err.Error(), http.StatusBadRequest)
			return
		}

		users, err := queries.GetUserByUserName(r.Context(), req.UserName)
		if err != nil {
			http.Error(w, "User not found. Please register first!", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(users.PasswordHash), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid Password", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := auth.GenerateTokens(users.UserID, users.UserName)
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserLoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
