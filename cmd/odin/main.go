package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/executor", router.ExecutorRouter())
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
