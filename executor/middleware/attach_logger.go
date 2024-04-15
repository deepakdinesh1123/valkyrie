package middleware

import (
	"context"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/constants"
	"github.com/deepakdinesh1123/valkyrie/executor/logger"
)

func AttachLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		ctx := context.WithValue(r.Context(), constants.ContextKey("logger"), log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
