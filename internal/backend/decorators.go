package server

import (
	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
)

func WithPort(v string) ServerOption {
	return func(s *Server) {
		s.port = v
	}
}

func WithEnvironment(v string) ServerOption {
	return func(s *Server) {
		s.environment = v
	}
}

func WithDbPath(v string) ServerOption {
	return func(s *Server) {
		s.dbPath = v
	}
}

func WithLogger(v logger.Logger) ServerOption {
	return func(s *Server) {
		s.logger = v
	}
}
