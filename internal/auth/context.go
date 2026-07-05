package auth

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "user_id"

// SetUserID stores the authenticated user's ID in the request context.
func SetUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// GetUserID retrieves the authenticated user's ID from the context.
// Returns an error if the context has no user ID (unauthenticated request).
func GetUserID(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, errors.New("user not authenticated")
	}
	return id, nil
}
