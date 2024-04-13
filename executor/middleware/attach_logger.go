package middleware

import (
	"context"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/logger"
)

func AttachLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type ContextKey string
		log := logger.GetLogger()
		ctx := context.WithValue(r.Context(), ContextKey("logger"), log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
