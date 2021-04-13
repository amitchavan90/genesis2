package routes

import (
	"genesis/restpkg/handlers"

	"github.com/go-chi/chi"
)

// Ping Routes function
func Ping(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		r.Get("/king", handlers.King)
	})
}
