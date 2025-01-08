package main

import (
	"log"
	"os"
	"users-service/api"
	"users-service/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg := config.NewConfig()
	server := api.NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
