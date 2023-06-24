// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"canvas/storage"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	address  string
	database *storage.Database
	log      *zap.Logger
	mux      chi.Router
	server   *http.Server
}

type Options struct {
	Database *storage.Database
	Host     string
	Log      *zap.Logger
	Port     int
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address:  address,
		database: opts.Database,
		log:      opts.Log,
		mux:      mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the Server by setting up routes and listening for HTTP requests on the given address.
func (s *Server) Start() error {

	if err := s.database.Connect(); err != nil {
		s.log.Error("Error while connecting to database", zap.Error(err))
		return fmt.Errorf("error connecting to database: %w", err)
	}
	s.setupRoutes()

	if va, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		s.log.Info("Starting inside lambda", zap.String("AWS_LAMBDA_RUNTIME_API", va))
		lambda.Start(httpadapter.New(s.mux).ProxyWithContext)
	}
	s.log.Info("Starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
