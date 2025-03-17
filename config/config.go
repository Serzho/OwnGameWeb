package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
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
	LoggerLevel  int
	CsvPath      string
}

type ServerConfig struct {
	Port int
	Url  string
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

func Load() *Config {
	slog.Info("Loading config...")

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}
	return &Config{
		Server: ServerConfig{
			Port: getIntEnv("SERVER_PORT", 8080),
			Url:  getStringEnv("SERVER_URL", ""),
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
			LoggerLevel:  getIntEnv("LOGGER_LEVEL", -4),
			CsvPath:      getStringEnv("CSV_PATH", "./pack/"),
		}}
}
