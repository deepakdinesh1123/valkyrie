package middleware

import (
	"net/http"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
)

func ConditionalRateLimit(lmt *limiter.Limiter, skipHeaders map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for headerName, expectedValue := range skipHeaders {
				if headerValue := r.Header.Get(headerName); headerValue != "" {
					if expectedValue == "" || headerValue == expectedValue {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			// Apply rate limiting
			tollbooth.HTTPMiddleware(lmt)(next).ServeHTTP(w, r)
		})
	}
}
