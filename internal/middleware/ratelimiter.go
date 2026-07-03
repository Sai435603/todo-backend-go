package middleware

import (
	"net/http"
	"sync"
	"time"
)

// visitor tracks request timestamps for a single client.
type visitor struct {
	tokens    float64
	lastSeen  time.Time
}

// RateLimiter returns middleware that limits requests per client IP
// using a token-bucket algorithm. Each IP gets `rps` requests per second
// with a burst capacity of `burst`.
func RateLimiter(rps float64, burst int) func(http.Handler) http.Handler {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	// Background cleanup: evict stale entries every 3 minutes.
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			v, exists := visitors[ip]
			now := time.Now()

			if !exists {
				v = &visitor{
					tokens:   float64(burst),
					lastSeen: now,
				}
				visitors[ip] = v
			}

			// Refill tokens based on elapsed time.
			elapsed := now.Sub(v.lastSeen).Seconds()
			v.tokens += elapsed * rps
			if v.tokens > float64(burst) {
				v.tokens = float64(burst)
			}
			v.lastSeen = now

			if v.tokens < 1 {
				mu.Unlock()
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"success":false,"error":{"message":"rate limit exceeded"}}`, http.StatusTooManyRequests)
				return
			}

			v.tokens--
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
