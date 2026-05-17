package config

import (
	"os"
	"strconv"
	"time"
)

// ConvertConfig contiene configuración para operaciones de conversión.
type ConvertConfig struct {
	DefaultDPI      int
	DefaultFormat   string
	DefaultPoolSize int
	RefreshEvery    int
	RetryOnFailure  bool
}

// ServerConfig contiene configuración para el servidor HTTP.
type ServerConfig struct {
	// HTTP server settings
	Host string
	Port string

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Upload limits
	MaxUploadSize int64

	// Logging
	LogLevel  string
	LogFormat string // "text" o "json"

	// Conversion operations
	Convert ConvertConfig
}

// NewServerConfig crea una configuración del servidor desde variables de entorno.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:         getEnv("HTTP_HOST", "0.0.0.0"),
		Port:         getEnv("HTTP_PORT", "8080"),
		ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 120*time.Second),
		IdleTimeout:  getDuration("HTTP_IDLE_TIMEOUT", 120*time.Second),
		MaxUploadSize: getInt64("HTTP_MAX_UPLOAD_SIZE", 200<<20), // 200MB
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		LogFormat:    getEnv("LOG_FORMAT", "text"),
		Convert: ConvertConfig{
			DefaultDPI:      getInt("CONVERT_DEFAULT_DPI", 150),
			DefaultFormat:   getEnv("CONVERT_DEFAULT_FORMAT", "png"),
			DefaultPoolSize: getInt("CONVERT_DEFAULT_POOL_SIZE", 2),
			RefreshEvery:    getInt("CONVERT_REFRESH_EVERY", 50),
			RetryOnFailure:  getBool("CONVERT_RETRY_ON_FAILURE", false),
		},
	}
}

// Utility functions
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultVal
}

func getDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultVal
}