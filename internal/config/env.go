package config

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Environment holds the application configuration.
type Environment struct {
	httpServerPort   int
	devModeEnabled   bool
	streamConnection string
}

type environmentKey string

const (
	port                environmentKey = "HTTP_PORT"
	devMode             environmentKey = "DEV_MODE"
	streamConnectionKey environmentKey = "STREAM_CONN"
)

// GetHttpServerPort returns the configured port for the server.
func (e *Environment) GetHttpServerPort() int {
	return e.httpServerPort
}

// IsDevModeEnabled returns true if the application is in development mode.
func (e *Environment) IsDevModeEnabled() bool {
	return e.devModeEnabled
}

// GetStreamConnection returns the url for the upfluence stream.
func (e *Environment) GetStreamConnection() string {
	return e.streamConnection
}

// LoadFromEnvironment loads the configuration by first reading the .env file.
func LoadFromEnvironment() *Environment {
	// We try to load the .env file before reading environment variables
	if err := loadEnv(".env"); err != nil {
		slog.Warn("no environment file, switching to default variables")
	}

	return &Environment{
		httpServerPort:   getEnvAs(string(port), 8080),
		devModeEnabled:   getEnvAs(string(devMode), true),
		streamConnection: getEnvAs(string(streamConnectionKey), "https://stream.upfluence.co"),
	}
}

// loadEnv reads a basic .env file and sets the environment variables.
func loadEnv(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// We only set the variable if it doesn't already exist in the system
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

// getEnvAs fetches an environment variable by key and parses it to type T.
func getEnvAs[T any](key string, fallback T) T {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}

	var result any
	switch any(fallback).(type) {
	case string:
		result = valStr
	case int:
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return fallback
		}
		result = val
	case bool:
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fallback
		}
		result = val
	default:
		return fallback
	}

	return result.(T)
}
