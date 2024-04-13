package main

import (
	"fmt"
	"log"
	"net/http"

	valklyrie_middleware "github.com/deepakdinesh1123/valkyrie/executor/middleware"
	"github.com/deepakdinesh1123/valkyrie/executor/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(valklyrie_middleware.AttachLoggerMiddleware)
	r.Use(valklyrie_middleware.AttachMachineryMiddleware)
	r.Mount("/executor", router.ExecutorRouter())
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
