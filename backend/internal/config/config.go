package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Application
	AppName string
	AppEnv  string
	AppPort string
	AppHost string

	// JWT
	JWTSecret            string
	JWTAccessTokenExpiry time.Duration
	JWTRefreshTokenExpiry time.Duration

	// Google Gemini
	GeminiAPIKey    string
	GeminiModel     string
	GeminiTemperature float64
	GeminiMaxTokens int

	// CORS
	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Query Configuration
	QueryTimeoutSeconds int
	MaxResultRows       int
	ExportMaxRows       int

	// Logging
	LogLevel  string
	LogFormat string

	// Security
	BCryptCost          int
	SessionTimeoutMinutes int
}

var AppConfig *Config

func Load() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	AppConfig = &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "mastercard_user"),
		DBPassword: getEnv("DB_PASSWORD", "mastercard_pass"),
		DBName:     getEnv("DB_NAME", "mastercard_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Application
		AppName: getEnv("APP_NAME", "Mastercard NLP-to-SQL Platform"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppHost: getEnv("APP_HOST", "0.0.0.0"),

		// JWT
		JWTSecret:            getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTAccessTokenExpiry:  parseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m")),
		JWTRefreshTokenExpiry: parseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h")),

		// Google Gemini
		GeminiAPIKey:     getEnv("GEMINI_API_KEY", ""),
		GeminiModel:      getEnv("GEMINI_MODEL", "gemini-pro"),
		GeminiTemperature: parseFloat(getEnv("GEMINI_TEMPERATURE", "0.1")),
		GeminiMaxTokens:  parseInt(getEnv("GEMINI_MAX_TOKENS", "2048")),

		// CORS
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

		// Rate Limiting
		RateLimitRequests: parseInt(getEnv("RATE_LIMIT_REQUESTS", "100")),
		RateLimitWindow:   parseDuration(getEnv("RATE_LIMIT_WINDOW", "60s")),

		// Query Configuration
		QueryTimeoutSeconds: parseInt(getEnv("QUERY_TIMEOUT_SECONDS", "30")),
		MaxResultRows:       parseInt(getEnv("MAX_RESULT_ROWS", "10000")),
		ExportMaxRows:       parseInt(getEnv("EXPORT_MAX_ROWS", "100000")),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Security
		BCryptCost:          parseInt(getEnv("BCRYPT_COST", "12")),
		SessionTimeoutMinutes: parseInt(getEnv("SESSION_TIMEOUT_MINUTES", "60")),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(value string) int {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return 0
}

func parseFloat(value string) float64 {
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return 0.0
}

func parseDuration(value string) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return 15 * time.Minute
}

