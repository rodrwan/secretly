package config

import (
	"os"
)

// Config contains the server configuration
type Config struct {
	Port     string
	EnvPath  string
	BasePath string
	DBPath   string
}

// New creates a new configuration with default values
func New() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		EnvPath:  getEnv("ENV_PATH", ".env"),
		BasePath: getEnv("BASE_PATH", "/api/v1"),
		DBPath:   getEnv("DB_PATH", "secretly.db"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
