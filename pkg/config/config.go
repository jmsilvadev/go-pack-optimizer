package config

import (
	"context"
	"os"

	"github.com/jmsilvadev/go-pack-optimizer/pkg/logger"
)

// Default Values used if no environment variables are set
var (
	dbPath      = "/tmp/packs.db"
	serverPort  = ":8080"
	loggerLevel = "DEBUG"
	environment = "dev"
)

// Config holds application configuration values
type Config struct {
	ServerPort string
	Env        string
	DbPath     string
	Logger     logger.Logger
}

// New creates a new Config instance with provided values
func New(ctx context.Context, port, env, dbPath string, logger logger.Logger) *Config {
	return &Config{
		ServerPort: port,
		Env:        env,
		Logger:     logger,
		DbPath:     dbPath,
	}
}

// GetDefaultConfig loads configuration from environment variables,
// falling back to predefined defaults. It also initializes the logger.
func GetDefaultConfig() *Config {
	environment = getEnv("ENV", environment)
	serverPort = getEnv("SERVER_PORT", serverPort)
	loggerLevel = getEnv("LOG_LEVEL", loggerLevel)
	dbPath = getEnv("DB_PATH", dbPath)

	// Determine log level
	level := logger.LEVEL_ERROR
	if loggerLevel == "INFO" {
		level = logger.LEVEL_INFO
	}
	if loggerLevel == "WARN" {
		level = logger.LEVEL_WARN
	}
	if loggerLevel == "DEBUG" {
		level = logger.LEVEL_DEBUG
	}
	log := logger.New(level)

	ctx := context.Background()
	config := New(ctx, serverPort, environment, dbPath, log)

	return config
}

// getEnv returns the value of the given environment variable,
// or the fallback if the variable is not set.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
