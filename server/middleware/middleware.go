package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// SetHeader is a higher-order function that sets a header
func SetHeader(key, value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)

			next.ServeHTTP(w, r)
		})
	}
}

// RequestLogging is a higher-order function that logs requests to the server
func RequestLogging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})
	}
}

// responseLogger wraps http.ResponseWriter so it can be captured for logging
type responseLogger struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// Write overrides http.Write so we can capture fields for logging
func (rl *responseLogger) Write(b []byte) (int, error) {
	rl.body.Write(b)                  // Copy response body
	return rl.ResponseWriter.Write(b) // Still send to client
}

// WriteHeader overrides http.WriteHeader so we can capture fields for logging
func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
}

// newReponseLogger is the constructor for our wrapped http.ResponseWriter
func newResponseLogger(w http.ResponseWriter) *responseLogger {
	return &responseLogger{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer),
	}
}

// RequestLogging is a higher-order function that logs responses from the server
func ResponseLogging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rl := newResponseLogger(w) // Wrapping the response writer
			next.ServeHTTP(rl, r)      // Proceed with the request

			// Convert the response body to a string
			body := rl.body.Bytes()

			// Log the response with basic details
			logger.Info("Response sent",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rl.statusCode),
				slog.String("body", string(body)),
			)
		})
	}
}
