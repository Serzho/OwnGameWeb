package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
}

type ServerConfig struct {
}

func Load() *Config {
	return &Config{Server: ServerConfig{}, Database: DatabaseConfig{}}
}
