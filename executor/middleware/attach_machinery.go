package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/tasks"
)

func AttachMachineryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type ContextKey string
		machinery_server, err := tasks.GetMachineryServer()
		fmt.Println(machinery_server)
		if nil != err {
			http.Error(w, "Could not attach machinery_server", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey("machinery_server"), machinery_server)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
