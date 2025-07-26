package middleware

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
)

var eventsPathPattern = regexp.MustCompile(`^/executions/[^/]+/events$`)

func TokenAuth() Middleware {
	return func(h http.Handler) http.Handler {
		log.Println("Token Auth")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if eventsPathPattern.MatchString(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			envConfig, _ := config.GetEnvConfig()

			// If no tokens are configured, skip authentication
			if envConfig.USER_TOKEN == "" && envConfig.ADMIN_TOKEN == "" {
				r = r.WithContext(context.WithValue(ctx, config.AuthKey, "noauth"))
				log.Println(envConfig.USER_TOKEN)
				h.ServeHTTP(w, r)
				return
			}

			// Tokens are configured, so authentication is required
			r = r.WithContext(context.WithValue(ctx, config.AuthKey, "auth"))
			headerValue := r.Header.Get("X-Auth-Token")

			// If header is empty, return unauthorized immediately
			if headerValue == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check against configured tokens (only non-empty tokens)
			if envConfig.USER_TOKEN != "" && headerValue == envConfig.USER_TOKEN {
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "user"))
			} else if envConfig.ADMIN_TOKEN != "" && headerValue == envConfig.ADMIN_TOKEN {
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "admin"))
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
