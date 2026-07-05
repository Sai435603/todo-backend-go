package middleware

import (
	"context"
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/response"
)

type contextKey string

const userIDKey contextKey = "userID"

// JWTAuth is a middleware that validates the auth_token cookie and injects userID into context.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				_ = response.Error(w, http.StatusUnauthorized, "authentication required — please sign in")
				return
			}

			claims, err := auth.ParseJWT(cookie.Value, secret)
			if err != nil {
				_ = response.Error(w, http.StatusUnauthorized, "invalid or expired token — please sign in again")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the authenticated user's ID from the request context.
func GetUserID(ctx context.Context) int64 {
	id, _ := ctx.Value(userIDKey).(int64)
	return id
}
