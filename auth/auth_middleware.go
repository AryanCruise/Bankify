package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("userID")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Headers: %v", r.Header)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing, Login or register first", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token"+err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("User ID from token: %v", claims.UserID)

		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}