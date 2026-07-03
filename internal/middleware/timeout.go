package middleware

import (
	"context"
	"net/http"
	"time"
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

const (
	// RequestTimeoutKey is the context key for the request deadline.
	RequestTimeoutKey contextKey = "request_timeout"
)

// Timeout returns middleware that enforces a maximum duration for request processing.
// If the handler does not complete within the timeout, the request context is cancelled.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
