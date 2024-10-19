package middleware

import (
	"log"
	"net/http"
)

// LoggingMiddleware returns a middleware handler that logs the details of incoming HTTP requests.
//
// It takes an http.Handler as a parameter, which is the next handler in the middleware chain.
// Returns an http.Handler that wraps the provided handler with logging functionality.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log the details of the request
		log.Printf("%s %s", r.Method, r.RequestURI)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
