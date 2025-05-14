package sizer_test

import (
	"os"
	"testing"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/jmsilvadev/go-pack-optimizer/pkg/sizer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func setupTempDB(t *testing.T) (string, func()) {
	tempDir := t.TempDir()
	return tempDir, func() {
		_ = os.RemoveAll(tempDir)
	}
}

func TestNewSizer_PopulatesDefaults(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	s, err := sizer.NewSizer(dir, logger.New(zapcore.DebugLevel))
	assert.NoError(t, err)
	defer s.Close()

	sizes, err := s.GetAllSizes()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{250, 500, 1000, 2000, 5000}, sizes)
}

func TestAddSize(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	s, err := sizer.NewSizer(dir, logger.New(zapcore.DebugLevel))
	assert.NoError(t, err)
	defer s.Close()

	err = s.AddSize(300)
	assert.NoError(t, err)

	sizes, err := s.GetAllSizes()
	assert.NoError(t, err)
	assert.Contains(t, sizes, 300)
}

func TestRemoveSize(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	s, err := sizer.NewSizer(dir, logger.New(zapcore.DebugLevel))
	assert.NoError(t, err)
	defer s.Close()

	err = s.RemoveSize(250)
	assert.NoError(t, err)

	sizes, err := s.GetAllSizes()
	assert.NoError(t, err)
	assert.NotContains(t, sizes, 250)
}

func TestRemoveSize_NotFound(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	s, err := sizer.NewSizer(dir, logger.New(zapcore.DebugLevel))
	assert.NoError(t, err)
	defer s.Close()

	err = s.RemoveSize(99999)
	assert.Error(t, err)
	assert.Equal(t, "pack size not found", err.Error())
}

func TestClose(t *testing.T) {
	dir, cleanup := setupTempDB(t)
	defer cleanup()

	s, err := sizer.NewSizer(dir, logger.New(zapcore.DebugLevel))
	assert.NoError(t, err)

	err = s.Close()
	assert.NoError(t, err)
}
