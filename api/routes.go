package api

import (
	"Accounts/auth"
	database "Accounts/internal/datab"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin to access your API (adjust if needed)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

		// For preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(queries *database.Queries, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Use(enableCORS)
	router.Handle("/swagger/*", httpSwagger.WrapHandler)

	router.HandleFunc("/register", RegisterHandler(queries))
	router.HandleFunc("/login", loginHandler(queries))

	router.Route("/accounts", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			createAccount(w, r, queries)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			getAccount(w, r, queries)
		})
		r.Get("/{id}/balance", func(w http.ResponseWriter, r *http.Request) {
			getBalance(w, r, queries)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			deleteAccount(w, r, queries)
		})
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			updateAccount(w, r, queries)
		})
	})

	router.Route("/transactions", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/deposit", func(w http.ResponseWriter, r *http.Request) {
			deposit(w, r, queries)
		})
		r.Post("/withdraw", func(w http.ResponseWriter, r *http.Request) {
			withdraw(w, r, queries)
		})
	})

	router.Route("/transfer", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			transferMoney(w, r, queries, db)
		})
	})
	return router
}
