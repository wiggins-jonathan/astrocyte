package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Middleware to add a header
func SetHeader(key, value string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
