package server

import (
	"fmt"
	"net/http"
)

type ServerOption func(*server)

type server struct {
	Port  int
	Debug bool
}

// NewServer returns a server with optional defaults
func NewServer(options ...ServerOption) *server {
	server := &server{Port: 8080}

	for _, option := range options {
		option(server)
	}

	return server
}

// Serve starts the astrocyte server
func (s *server) Serve() error {
	port := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("Server listening at http://localhost%s - Ctrl+c to quit.\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}

func WithPort(port int) ServerOption {
	return func(s *server) {
		s.Port = port
	}
}

func WithDebug(debug bool) ServerOption {
	return func(s *server) {
		s.Debug = debug
	}
}
