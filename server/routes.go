package server

import (
	"canvas/handlers"
	"canvas/k"

	"github.com/go-chi/chi/v5"
)

func (s *Server) setupRoutes() {
	s.mux.Route(k.GlobalPrefix, func(r chi.Router) {
		handlers.Health(r)
		handlers.Homepage(r)
		handlers.FrontPage(r)
	})
}
