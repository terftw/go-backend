package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: 8080, // or from env
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     5432,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}, nil
}
