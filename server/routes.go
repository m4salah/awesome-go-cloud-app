package server

import (
	"canvas/handlers"

	"github.com/go-chi/chi/v5"
)

const globalPrefix = "/hello"

func (s *Server) setupRoutes() {
	s.mux.Route(globalPrefix, func(r chi.Router) {
		handlers.Health(r)
	})
}
