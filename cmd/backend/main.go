package main

import (
	"context"

	server "github.com/jmsilvadev/go-pack-optimizer/internal/backend"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/config"
)

func main() {
	c := config.GetDefaultConfig()
	run(c)
}

func run(conf *config.Config) error {
	serverOptions := []server.ServerOption{
		server.WithPort(conf.ServerPort),
		server.WithEnvironment(conf.Env),
		server.WithLogger(conf.Logger),
		server.WithDbPath(conf.DbPath),
	}

	s := server.NewServer(serverOptions...)

	s.Start(context.Background())

	return nil
}
