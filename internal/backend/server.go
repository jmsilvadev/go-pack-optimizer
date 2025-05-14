package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jmsilvadev/go-pack-optimizer/internal/handler"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/optimizer"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/sizer"
)

type Server struct {
	environment string
	port        string
	dbPath      string
	logger      logger.Logger
}

type ServerOption func(*Server)
type ServerType int

const (
	TypeBackend ServerType = iota
	TypeFrontend
)

func NewServer(options ...ServerOption) *Server {
	svr := &Server{}
	for _, opt := range options {
		opt(svr)
	}
	return svr
}

// Start starts the baceknd server
func (s *Server) Start(ctx context.Context) {
	sz, err := sizer.NewSizer(s.dbPath, s.logger)
	if err != nil {
		log.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer sz.Close()

	op := optimizer.New(sz, s.logger)
	h := handler.New(op)

	r, err := handler.NewRouter(h)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    s.port,
		Handler: r,
	}

	listener := make(chan os.Signal, 1)
	signal.Notify(listener, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.logger.Warn(fmt.Sprint("received a shutdown signal:", <-listener))
		s.logger.Warn("shutdown the server...")
		server.Shutdown(ctx)
		wg.Done()
	}()

	s.logger.Info("server listening at " + s.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to serve: " + err.Error())
	}

	wg.Wait()
	s.logger.Warn("server gracefully stopped")
}
