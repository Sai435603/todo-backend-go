package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

// maxBodySize is the maximum allowed request body size (1 MB).
const maxBodySize = 1 << 20 // 1 MB

// ContentType returns middleware that enforces application/json Content-Type
// on mutating requests (POST, PUT, PATCH). GET, DELETE, HEAD, and OPTIONS
// are passed through without checking.
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce on methods that carry a request body.
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := r.Header.Get("Content-Type")
			if ct == "" || !strings.HasPrefix(strings.ToLower(ct), "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"error": map[string]string{
						"message": "Content-Type must be application/json",
					},
				})
				return
			}
		}

		// Limit the request body size to prevent abuse.
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		next.ServeHTTP(w, r)
	})
}
