package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	router := Router()

	// List available routes.
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)

		return nil
	}

	log.Println("Available routes:")

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging error: %s\n", err.Error())
	}

	// Launch server
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Serving API on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
