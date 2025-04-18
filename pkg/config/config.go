package config

import (
	"os"
)

type EnvironmentType string

const (
	Production  EnvironmentType = "prod"
	Development EnvironmentType = "dev"
	Local       EnvironmentType = "local"
)

type Config struct {
	APP         string
	Environment EnvironmentType
	LogLevel    string
	AppURL      string

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}

	Admin struct {
		Username string
		Password string
	}

	Token struct {
		Secret string
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "karavanix-api-server")
	config.Environment = EnvironmentType(getEnv("ENVIRONMENT", "develop"))
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "5m")
	config.AppURL = getEnv("APP_URL", "")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", ":8000")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "karavanix")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.DB.Sslmode = getEnv("POSTGRES_SSLMODE", "disable")

	// admin configuration
	config.Admin.Username = getEnv("ADMIN_USERNAME", "admin")
	config.Admin.Password = getEnv("ADMIN_PASSWORD", "admin")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "secret")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
