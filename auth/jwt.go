package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtsecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID   int32  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GenerateTokens(UserID int32, username string) (accessToken string, refreshToken string, err error) {
	accessExpiry := time.Now().Add(10 * time.Minute)
	refreshExpiry := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		UserID:   UserID,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err = CreateToken(claims, jwtsecret, jwt.SigningMethodHS256)
	if err != nil {
		return "", "", err
	}

	claims.ExpiresAt = jwt.NewNumericDate(refreshExpiry)
	refreshToken, err = CreateToken(claims, jwtsecret, jwt.SigningMethodHS256)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func CreateToken(claims *Claims, secret []byte, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(secret)
}

func VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return jwtsecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims from token")
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	// Parse and validate the refresh token
	claims, err := VerifyToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if the refresh token is still valid
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", fmt.Errorf("refresh token has expired")
	}

	// Generate a new access token with a fresh expiry time
	newAccessExpiry := time.Now().Add(5 * time.Minute)
	newAccessTokenClaims := &Claims{
		UserID:   claims.UserID,
		UserName: claims.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(newAccessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newAccessToken, err := CreateToken(newAccessTokenClaims, jwtsecret, jwt.SigningMethodHS256)
	if err != nil {
		return "", fmt.Errorf("failed to create new access token: %w", err)
	}

	return newAccessToken, nil
}
