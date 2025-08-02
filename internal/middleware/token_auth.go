package middleware

import (
	"net/http"
	"strings"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/go-chi/jwtauth/v5"
)

func TokenAuth(ja *jwtauth.JWTAuth) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role := r.Context().Value(config.RoleKey)
			if role == "admin" {
				h.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := ja.Decode(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if token == nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
