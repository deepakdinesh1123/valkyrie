package middleware

import (
	"context"
	"net/http"
	"regexp"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
)

var eventsPathPattern = regexp.MustCompile(`^/executions/[^/]+/events$`)

func TokenAuth() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if eventsPathPattern.MatchString(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			envConfig, _ := config.GetEnvConfig()

			if envConfig.USER_TOKEN == "" && envConfig.ADMIN_TOKEN == "" {
				r = r.WithContext(context.WithValue(ctx, config.AuthKey, "noauth"))
				h.ServeHTTP(w, r)
				return
			}

			r = r.WithContext(context.WithValue(ctx, config.AuthKey, "auth"))
			headerValue := r.Header.Get("X-Auth-Token")
			switch headerValue {
			case envConfig.USER_TOKEN:
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "user"))
			case envConfig.ADMIN_TOKEN:
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "admin"))
			default:
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
