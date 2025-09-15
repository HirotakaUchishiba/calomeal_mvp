package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DatabaseURL string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPassword  string

	// Service Addresses
	FoodServiceAddr      string
	LogsServiceAddr      string
	AnalyticsServiceAddr string

	// BFF Configuration
	BFFPort string
	BFFHost string

	// Frontend Configuration
	FrontendPort string
	FrontendHost string

	// Security
	JWTSecret  string
	CORSOrigins []string

	// Logging
	LogLevel  string
	LogFormat string

	// Performance
	MaxConnections    int
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration

	// Health Check
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration

	// Monitoring
	MetricsEnabled bool
	MetricsPort    string

	// Error Handling
	RetryMaxAttempts                int
	RetryInitialDelay              time.Duration
	CircuitBreakerFailureThreshold int
	CircuitBreakerResetTimeout     time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Database configuration
	config.DatabaseURL = getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/calomeal?sslmode=disable")
	config.DBHost = getEnv("DB_HOST", "localhost")
	config.DBPort = getEnvAsInt("DB_PORT", 5432)
	config.DBName = getEnv("DB_NAME", "calomeal")
	config.DBUser = getEnv("DB_USER", "postgres")
	config.DBPassword = getEnv("DB_PASSWORD", "password")

	// Service addresses
	config.FoodServiceAddr = getEnv("FOOD_SERVICE_ADDR", "localhost:50051")
	config.LogsServiceAddr = getEnv("LOGS_SERVICE_ADDR", "localhost:50052")
	config.AnalyticsServiceAddr = getEnv("ANALYTICS_SERVICE_ADDR", "localhost:50053")

	// BFF configuration
	config.BFFPort = getEnv("BFF_PORT", "8080")
	config.BFFHost = getEnv("BFF_HOST", "localhost")

	// Frontend configuration
	config.FrontendPort = getEnv("FRONTEND_PORT", "5173")
	config.FrontendHost = getEnv("FRONTEND_HOST", "localhost")

	// Security
	config.JWTSecret = getEnv("JWT_SECRET", "dev_jwt_secret_key")
	config.CORSOrigins = getEnvAsSlice("CORS_ORIGINS", []string{"http://localhost:5173"})

	// Logging
	config.LogLevel = getEnv("LOG_LEVEL", "info")
	config.LogFormat = getEnv("LOG_FORMAT", "text")

	// Performance
	config.MaxConnections = getEnvAsInt("MAX_CONNECTIONS", 10)
	config.ConnectionTimeout = getEnvAsDuration("CONNECTION_TIMEOUT", 10*time.Second)
	config.RequestTimeout = getEnvAsDuration("REQUEST_TIMEOUT", 30*time.Second)

	// Health check
	config.HealthCheckInterval = getEnvAsDuration("HEALTH_CHECK_INTERVAL", 10*time.Second)
	config.HealthCheckTimeout = getEnvAsDuration("HEALTH_CHECK_TIMEOUT", 5*time.Second)

	// Monitoring
	config.MetricsEnabled = getEnvAsBool("METRICS_ENABLED", false)
	config.MetricsPort = getEnv("METRICS_PORT", "9090")

	// Error handling
	config.RetryMaxAttempts = getEnvAsInt("RETRY_MAX_ATTEMPTS", 3)
	config.RetryInitialDelay = getEnvAsDuration("RETRY_INITIAL_DELAY", 100*time.Millisecond)
	config.CircuitBreakerFailureThreshold = getEnvAsInt("CIRCUIT_BREAKER_FAILURE_THRESHOLD", 3)
	config.CircuitBreakerResetTimeout = getEnvAsDuration("CIRCUIT_BREAKER_RESET_TIMEOUT", 10*time.Second)

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as duration with a fallback value
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// getEnvAsSlice gets an environment variable as string slice with a fallback value
func getEnvAsSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		var result []string
		for _, item := range splitAndTrim(value, ",") {
			if item != "" {
				result = append(result, item)
			}
		}
		return result
	}
	return fallback
}

// splitAndTrim splits a string by delimiter and trims spaces from each part
func splitAndTrim(s, delimiter string) []string {
	var result []string
	for _, part := range splitString(s, delimiter) {
		trimmed := trimSpace(part)
		result = append(result, trimmed)
	}
	return result
}

// splitString splits a string by delimiter (simplified version)
func splitString(s, delimiter string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(delimiter) <= len(s) && s[i:i+len(delimiter)] == delimiter {
			result = append(result, s[start:i])
			start = i + len(delimiter)
			i += len(delimiter) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// trimSpace trims leading and trailing whitespace (simplified version)
func trimSpace(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading spaces
	for start < end && s[start] == ' ' {
		start++
	}
	
	// Trim trailing spaces
	for end > start && s[end-1] == ' ' {
		end--
	}
	
	return s[start:end]
}

// ValidateConfig validates the configuration
func (c *Config) ValidateConfig() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.FoodServiceAddr == "" {
		return fmt.Errorf("FOOD_SERVICE_ADDR is required")
	}
	if c.LogsServiceAddr == "" {
		return fmt.Errorf("LOGS_SERVICE_ADDR is required")
	}
	if c.AnalyticsServiceAddr == "" {
		return fmt.Errorf("ANALYTICS_SERVICE_ADDR is required")
	}
	return nil
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.LogLevel == "info" && c.LogFormat == "json"
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return c.DatabaseURL
}
