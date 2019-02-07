package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Router is the main router for the API.
func Router() *chi.Mux {
	// Create new router.
	router := chi.NewRouter()

	// Register middlewares.
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RedirectSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.Recoverer,
	)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	// Default route.
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hitched API"))
	})

	return router
}
