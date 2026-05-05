package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware is a global middleware that catches panics, logs them, and returns 500 JSON errors
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic and stack trace
				log.Printf("PANIC RECOVERED: %v\nStack Trace:\n%s", err, debug.Stack())

				// Return JSON error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal server error",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
