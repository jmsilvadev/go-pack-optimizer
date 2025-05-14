package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jmsilvadev/go-pack-optimizer/internal/handler"
	"github.com/jmsilvadev/go-pack-optimizer/internal/handler/mocks"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTempDB(t *testing.T) (string, func()) {
	tempDir := t.TempDir()
	return tempDir, func() {
		_ = os.RemoveAll(tempDir)
	}
}

func TestNewServerBackend(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	mockOptimizer := &mocks.MockOptimizerInterface{}
	mockLogger := logger.New(zap.DebugLevel)
	server := NewServer(
		func(s *Server) {
			s.logger = mockLogger
			s.port = ":8080"
			s.dbPath = dir
		},
	)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.port)

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.New(mockOptimizer).HealthHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}

func TestStart(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dir, cleanup := setupTempDB(t)
	defer cleanup()

	mockLogger := logger.New(zap.DebugLevel)
	server := NewServer(
		func(s *Server) {
			s.logger = mockLogger
			s.port = ":8080"
			s.dbPath = dir
		},
	)

	// Only to cover code
	go server.Start(ctx)
	time.Sleep(time.Microsecond)
}
