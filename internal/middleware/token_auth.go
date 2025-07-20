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

			if envConfig.USER_TOKEN == "" && envConfig.ADMIN_TOKEN == "" {
				r = r.WithContext(context.WithValue(ctx, config.AuthKey, "noauth"))
				log.Println(envConfig.USER_TOKEN)
				h.ServeHTTP(w, r)
				return
			}

			log.Println(envConfig.USER_TOKEN)
			log.Println(envConfig.ADMIN_TOKEN)

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
