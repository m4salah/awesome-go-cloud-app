package server

import (
	"canvas/handlers"
	"canvas/k"
	"canvas/model"
	"context"

	"github.com/go-chi/chi/v5"
)

type signupperMock struct{}

func (s signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	return "", nil
}

func (s *Server) setupRoutes() {
	s.mux.Route(k.GlobalPrefix, func(r chi.Router) {
		handlers.Health(r)
		handlers.Homepage(r)
		handlers.FrontPage(r)
		handlers.NewsletterSignup(r, &signupperMock{})
		handlers.NewsletterThanks(r)
	})
}
