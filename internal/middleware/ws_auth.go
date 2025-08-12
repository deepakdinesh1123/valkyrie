package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func WSAuth(ja *jwtauth.JWTAuth) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jobId := chi.URLParam(r, "jobId")
			if jobId == "" {
				http.Error(w, "Job ID required", http.StatusBadRequest)
				return
			}

			tokenString := r.URL.Query().Get("token")
			if tokenString == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			// Decode and validate the token
			token, err := ja.Decode(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if token == nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Safely extract claims
			claims, err := token.AsMap(context.Background())
			if err != nil {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			tokenJobId, exists := claims["jobId"]
			if !exists {
				http.Error(w, "Token missing jobId claim", http.StatusUnauthorized)
				return
			}

			tokenJobIdStr, ok := tokenJobId.(string)
			if !ok {
				http.Error(w, "Invalid jobId claim format", http.StatusUnauthorized)
				return
			}

			if tokenJobIdStr != jobId {
				http.Error(w, "Token jobId mismatch", http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
