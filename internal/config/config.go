package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// AppConfig holds the application configuration loaded from environment variables.
type AppConfig struct {
	Environment string
	LogLevel    string
	TableName   string
	TimeoutSecs int
}

// Load reads environment variables and populates AppConfig, validating required fields.
func Load() (*AppConfig, error) {
	env := getEnv("ENVIRONMENT", "production")
	logLevel := getEnv("LOG_LEVEL", "INFO")
	tableName := getEnv("TABLE_NAME", "")
	
	if tableName == "" {
		return nil, errors.New("missing required environment variable: TABLE_NAME")
	}

	timeoutStr := getEnv("TIMEOUT_SECS", "30")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TIMEOUT_SECS: %w", err)
	}

	return &AppConfig{
		Environment: env,
		LogLevel:    logLevel,
		TableName:   tableName,
		TimeoutSecs: timeout,
	}, nil
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
