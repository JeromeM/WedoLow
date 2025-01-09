package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RandomUserAPI  string
	JaegerEndpoint string
}

func NewConfig() *Config {
	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		DatabaseURL:    getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"),
		RandomUserAPI:  getEnvOrDefault("RANDOM_USER_API", "https://randomuser.me/api/"),
		JaegerEndpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
