package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"astrocyte/server/api"
	mw "astrocyte/server/middleware"
)

type ServerOption func(*server)

type server struct {
	Port   int
	Logger *slog.Logger
}

// NewServer returns a server with adjustable defaults
func NewServer(options ...ServerOption) *server {
	server := &server{Port: 8080}

	for _, option := range options {
		option(server)
	}

	return server
}

// Serve starts the astrocyte server
func (s *server) Serve() error {
	mux := http.NewServeMux()

	apis := []api.API{
		api.NewClient(),
		api.NewPushAPI(),
	}

	for _, api := range apis {
		api.Register(mux)
	}

	// set global middleware
	muxWithMiddleware := mw.SetHeader("Content-Type", "application/json")(mux)

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
