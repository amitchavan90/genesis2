package routes

import (
	"github.com/go-chi/chi"
)

func MapUrls(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		Ping(r)
	})
}
