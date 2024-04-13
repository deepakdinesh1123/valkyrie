package router

import (
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/handlers/execution"
	"github.com/go-chi/chi/v5"
)

func ExecutorRouter() http.Handler {
	router := chi.NewRouter()
	router.Post("/execute", execution.Execute)
	return router
}
