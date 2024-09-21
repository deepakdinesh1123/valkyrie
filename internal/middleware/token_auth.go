package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
)

func TokenAuth() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			envConfig, _ := config.GetEnvConfig()
			headerValue := r.Header.Get("X-Auth-Token")
			log.Println(envConfig.ODIN_USER_TOKEN, envConfig.ODIN_ADMIN_TOKEN)
			if envConfig.ODIN_USER_TOKEN == "" && envConfig.ODIN_ADMIN_TOKEN == "" {
				log.Println("ODIN_USER_TOKEN and ODIN_ADMIN_TOKEN are not set")
				h.ServeHTTP(w, r)
				return
			}
			switch headerValue {
			case envConfig.ODIN_USER_TOKEN:
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "user"))
			case envConfig.ODIN_ADMIN_TOKEN:
				r = r.WithContext(context.WithValue(ctx, config.UserKey, "admin"))
			default:
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
