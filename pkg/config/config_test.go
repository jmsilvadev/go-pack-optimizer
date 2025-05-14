package config

import (
	"context"
	"testing"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewConfig(t *testing.T) {
	l := logger.New(zap.DebugLevel)
	got := New(context.Background(), ":8080", "dev", "test.db", l)
	if got.ServerPort != ":8080" {
		t.Errorf("Got and Expected are not equals. Got: %v, expected: :8080", got.ServerPort)
	}
}

func TestGetDeaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	if config.ServerPort == "" {
		t.Errorf("Got and Expected are not equals. got: '', expected: !''")
	}
}

func TestGetEnv(t *testing.T) {
	v := getEnv("a", "b")
	require.Equal(t, "b", v)
}
