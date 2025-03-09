package config

import (
	"fmt"
	"github.com/joho/godotenv"
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
	LoggerLevel  string
}

type ServerConfig struct {
	Port int
	Url  string
}

func getStringEnv(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	fmt.Printf("Environment variable %s is missing.\n", name)
	return defaultVal
}

func getIntEnv(name string, defaultVal int) int {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		fmt.Printf("Environment variable %s is missing.\n", name)
		return defaultVal
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	fmt.Printf("Environment variable %s has incorrect value %s. Expected type - int\n", name, valueStr)
	return defaultVal
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
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
			LoggerLevel:  getStringEnv("LOGGER_LEVEL", "info"),
		}}
}
