package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Global   GlobalConfig
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Name         string
	Password     string
	MaxOpenConns int
}

type GlobalConfig struct {
	SecretPhrase string
	LoggerLevel  slog.Level
	CsvPath      string
}

type ServerConfig struct {
	Port int
	URL  string
}

func getStringEnv(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	slog.Warn("Environment variable is missing", "name", name)
	return defaultVal
}

func getIntEnv(name string, defaultVal int) int {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		slog.Warn("Environment variable is missing", "name", name)
		return defaultVal
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	slog.Error("Environment variable has incorrect value. Expected type - int\n", "name", name, "value", valueStr)
	return defaultVal
}

func getLoggerLevel(name string, defaultVal int) slog.Level {
	intLevel := getIntEnv(name, defaultVal)

	var level slog.Level

	switch intLevel {
	case -4:
		level = slog.LevelDebug
	case 0:
		level = slog.LevelInfo
	case 4:
		level = slog.LevelWarn
	case 8:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return level
}

func Load() *Config {
	slog.Info("Loading config...")

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	return &Config{
		Server: ServerConfig{
			Port: getIntEnv("SERVER_PORT", 8080),
			URL:  getStringEnv("SERVER_URL", ""),
		},
		Database: DatabaseConfig{
			Host:         getStringEnv("DATABASE_HOST", "localhost"),
			Port:         getIntEnv("DATABASE_PORT", 5432),
			User:         getStringEnv("DATABASE_USER", "postgres"),
			Name:         getStringEnv("DATABASE_NAME", "postgres"),
			Password:     getStringEnv("DATABASE_PASSWORD", "postgres"),
			MaxOpenConns: getIntEnv("DATABASE_MAX_OPEN_CONNS", 0),
		},
		Global: GlobalConfig{
			SecretPhrase: getStringEnv("SECRET_PHRASE", "secret"),
			LoggerLevel:  getLoggerLevel("LOGGER_LEVEL", -4),
			CsvPath:      getStringEnv("CSV_PATH", "./pack/"),
		},
	}
}
