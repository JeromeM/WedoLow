package config

import (
	"os"
)

type Config struct {
	Port          string
	RandomUserAPI string
}

func NewConfig() *Config {
	return &Config{
		Port:          getEnvOrDefault("PORT", "8080"),
		RandomUserAPI: getEnvOrDefault("RANDOM_USER_API", "https://randomuser.me/api/"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
