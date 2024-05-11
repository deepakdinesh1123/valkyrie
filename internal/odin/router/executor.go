package router

import (
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/handlers/execution"
	"github.com/go-chi/chi/v5"
)

func ExecutorRouter() http.Handler {
	router := chi.NewRouter()
	router.Post("/execute/{environment}", execution.Execute)
	router.Get("/execute/{executionId}/state", execution.GetTaskState)
	return router
}
