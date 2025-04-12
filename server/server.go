package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"astrocyte/server/api"
	mw "astrocyte/server/middleware"
)

type server struct {
	Port    int
	Logger  *slog.Logger
	BaseURL *url.URL
}

type ServerOption func(*server)

// NewServer returns a new server with optional functional configuration
func NewServer(options ...ServerOption) *server {
	server := &server{}

	for _, option := range options {
		option(server)
	}

	return server
}

// Serve starts the astrocyte server
func (s *server) Serve() error {
	mux := http.NewServeMux()

	apis := []api.API{
		api.NewClient(api.WithBaseURL(s.BaseURL)),
		api.NewPushAPI(),
	}

	for _, api := range apis {
		api.RegisterRoutes(mux)
	}

	// set global middleware
	middlewares := []mw.Middleware{
		mw.SetHeader("Content-Type", "application/json"),
		mw.RequestLogging(s.Logger),
	}

	muxWithMiddleware := http.Handler(mux) // cast mux to http.Handler
	for _, middleware := range middlewares {
		muxWithMiddleware = middleware(muxWithMiddleware)
	}

	port := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("Server listening at http://localhost%s - Ctrl+c to quit.\n", port)
	if err := http.ListenAndServe(port, muxWithMiddleware); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}

// WithPort is a helper function that changes the default port
func WithPort(port int) ServerOption {
	return func(s *server) {
		s.Port = port
	}
}

// WithLogger is a helper function that sets up slog
func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *server) {
		s.Logger = logger
	}
}

// WithBaseURL is a helper function that changes the default BaseURL
func WithBaseURL(baseURL *url.URL) ServerOption {
	return func(s *server) {
		s.BaseURL = baseURL
	}
}
