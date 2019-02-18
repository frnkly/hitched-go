package router

import (
	"net/http"
	"time"

	"github.com/frnkly/hitched/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Create - Sets up the router and middleware, and defines all API routes.
func Create() *chi.Mux {
	// Create new router.
	router := chi.NewRouter()

	// Register middlewares.
	router.Use(
		middleware.RedirectSlashes,
		middleware.Logger,
		middleware.Recoverer,
		middleware.DefaultCompress,
		render.SetContentType(render.ContentTypeJSON),
		cors.New(GetCorsOptions()).Handler,
	)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	// Default route.
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hitched API"))
	})

	// Test routes
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})

	router.Get("/panic", func(writer http.ResponseWriter, request *http.Request) {
		panic("test")
	})

	// OAuth for Google Sheets API.

	// RSVP
	router.Get("/rsvp/{code}", handlers.GetRsvp)
	router.Post("/rsvp/{action}/{event}/{code}", handlers.UpdateRsvp)

	return router
}
