package middleware

import (
	"context"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
)

func AssignRole() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			envConfig, _ := config.GetEnvConfig()

			headerValue := r.Header.Get("X-Auth-Token")
			switch headerValue {
			case envConfig.ADMIN_TOKEN:
				r = r.WithContext(context.WithValue(ctx, config.RoleKey, "admin"))
			default:
				r = r.WithContext(context.WithValue(ctx, config.RoleKey, "anonymous"))
			}

			h.ServeHTTP(w, r)
		})
	}
}
