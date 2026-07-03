package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture the status code and bytes written.
type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteHeader  bool
}

// WriteHeader captures the status code before delegating to the underlying writer.
func (rr *responseRecorder) WriteHeader(code int) {
	if !rr.wroteHeader {
		rr.statusCode = code
		rr.wroteHeader = true
	}
	rr.ResponseWriter.WriteHeader(code)
}

// Write captures bytes written and ensures the header is recorded.
func (rr *responseRecorder) Write(b []byte) (int, error) {
	if !rr.wroteHeader {
		rr.statusCode = http.StatusOK
		rr.wroteHeader = true
	}
	n, err := rr.ResponseWriter.Write(b)
	rr.bytesWritten += n
	return n, err
}

// RequestLogger returns middleware that logs every incoming HTTP request
// with method, path, status, duration, and response size using structured logging.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// Choose log level based on status code.
			logFn := logger.Info
			if rec.statusCode >= http.StatusInternalServerError {
				logFn = logger.Error
			} else if rec.statusCode >= http.StatusBadRequest {
				logFn = logger.Warn
			}

			logFn("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status", rec.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes", rec.bytesWritten,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}
