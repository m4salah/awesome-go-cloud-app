package server

import (
	"canvas/handlers"

	"github.com/go-chi/chi/v5"
)

const globalPrefix = "/canvas"

func (s *Server) setupRoutes() {
	s.mux.Route(globalPrefix, func(r chi.Router) {
		handlers.Health(r)
		handlers.Homepage(r)
		handlers.FrontPage(s.mux)
	})

}
