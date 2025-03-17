package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

func RequestMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Str("method", r.Method).Str("url", r.URL.String()).Msg("request")
			next.ServeHTTP(w, r)
		})
	}
}
