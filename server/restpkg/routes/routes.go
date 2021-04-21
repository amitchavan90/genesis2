package routes

import (
	"github.com/go-chi/chi"
)

func MapUrls(r chi.Router) {
	r.Route("/api/", func(r chi.Router) {
		Ping(r)
	})
}
