package middleware

import (
	"net/http"
)

// SecurityHeaders returns middleware that sets standard security headers
// to protect against common web vulnerabilities (XSS, clickjacking, MIME sniffing, etc.).
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME-type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking.
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS filter in older browsers.
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Control referrer information.
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict browser features/permissions.
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Enforce HTTPS (1 year, include subdomains).
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}
