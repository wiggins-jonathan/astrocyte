package middleware

import (
	"fmt"
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
			log := fmt.Sprintf("Request received: %s %s", r.Method, r.URL.Path)
			logger.Info(log)

			next.ServeHTTP(w, r)
		})
	}
}
