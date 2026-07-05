package middleware

import (
	"net/http"
	"strings"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/response"
)

// Auth returns middleware that validates JWT Bearer tokens and injects the
// authenticated user ID into the request context.
func Auth(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				_ = response.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				_ = response.Error(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			claims, err := jwtSvc.ValidateToken(parts[1])
			if err != nil {
				_ = response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := auth.SetUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
