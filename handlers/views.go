package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"canvas/views"
)

func FrontPage(mux chi.Router) {
	mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		_ = views.FrontPage().Render(w)
	})
}
