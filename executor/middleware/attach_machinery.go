package middleware

import (
	"context"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/constants"
	"github.com/deepakdinesh1123/valkyrie/executor/tasks"
)

// AttachMachineryMiddleware is a middleware function that attaches a machinery server and a user to the request context.
func AttachMachineryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			machineryServerKey constants.ContextKey = "machinery_server"
			userKey            constants.ContextKey = "user"
		)

		// Attempt to get the machinery server.
		machineryServer, err := tasks.GetMachineryServer()
		if err != nil {
			http.Error(w, "Could not attach machinery_server", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), machineryServerKey, machineryServer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
