package config

import (
	"os"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RandomUserAPI string
}

func NewConfig() *Config {
	return &Config{
		Port:          getEnvOrDefault("PORT", "8080"),
		DatabaseURL:   getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"),
		RandomUserAPI: getEnvOrDefault("RANDOM_USER_API", "https://randomuser.me/api/"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
