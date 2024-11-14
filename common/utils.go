package common

import (
	"Accounts/auth"
	"context"
	"errors"
)

const UserContextKey = "userID"
func GetUserIDFromContext(ctx context.Context) (int32, error) {
	userID, ok := ctx.Value(auth.UserContextKey).(int32)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}